package chainjob

import (
	"encoding/json"
	"errors"
	"fmt"
	"hamster-client/config"
	"hamster-client/module/chainhelper"
	"hamster-client/module/pallet"
	"hamster-client/module/queue"
	"hamster-client/utils"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type VersionVo struct {
	Version string `json:"version"`
}

type WaitResourceJob struct {
	appID      int
	si         queue.StatusInfo
	helper     chainhelper.Helper
	bond       bool
	meta       *types.Metadata
}

func NewWaitResourceJob(appID int, helper chainhelper.Helper) queue.Job {
	return &WaitResourceJob{
		appID:  appID,
		helper: helper,
	}
}

func (j *WaitResourceJob) InitStatus() {
	j.si.Name = "Pull Image"
	j.si.Status = queue.None
}

func (j *WaitResourceJob) Run(si chan queue.StatusInfo) (queue.StatusInfo, error) {
	j.si.Status = queue.Running
	si <- j.si

	deployType, err := j.helper.DeployType(j.appID)
	if err != nil {
		j.si.Status = queue.Failed
		j.si.Error = err.Error()
		si <- j.si
		return j.si, err
	}

	api, err := j.helper.Account().GetSubstrateApi()
	if err != nil {
		j.si.Status = queue.Failed
		j.si.Error = err.Error()
		si <- j.si
		return j.si, err
	}

	pair, err := j.helper.Wallet().GetWalletKeypair()
	if err != nil {
		j.si.Status = queue.Failed
		j.si.Error = "WALLET_LOAD_ERROR"
		si <- j.si
		return j.si, err
	}
	if pair.Address == "" {
		j.si.Status = queue.Failed
		j.si.Error = errors.New("get Keypair failed").Error()
		si <- j.si
		return j.si, err
	}

	if !j.bond {
		// 100 unit
		err = pallet.Bond(api, j.meta, 100000000000000, pair)
		if err != nil {
			j.si.Status = queue.Failed
			j.si.Error = "WALLET_LOAD_ERROR"
			si <- j.si
			return j.si, err
		}
	}

	for i := 0; i < 60; i++ {

		mapData, err := pallet.GetResourceList(
			j.meta,
			api,
			func(resource *pallet.ComputingResource) bool {
				return resource.Status.IsUnused
			},
		)
		if err != nil {
			j.si.Status = queue.Failed
			j.si.Error = err.Error()
			si <- j.si
			return j.si, err
		}

		fmt.Println("可用资源数：", len(mapData))

		failSet := make(map[int]string)

		for _, val := range mapData {

			if _, isMapContainsKey := failSet[int(val.Index)]; isMapContainsKey {
				continue
			}

			if val.Status.IsUnused {
				resourceIndex := val.Index
				fmt.Println("发现未使用资源，占用。资源号：", resourceIndex)
				data, err := j.helper.App().QueryApplicationById(j.appID)
				if err != nil {
					fmt.Println("get application failed,err is: ", err)
					continue
				}

				// check p2p connection can be connected
				port := j.helper.App().QueryNextP2pPort()
				err = j.helper.P2p().
					LinkByProtocol(config.ProviderProtocol, port, string(val.PeerId))

				if err != nil {
					fmt.Println("create p2p network forward fail")
					failSet[int(val.Index)] = "fail"
					_, _ = j.helper.P2p().Close(fmt.Sprintf("/p2p/%s", string(val.PeerId)))
					continue
				}

				// check http is Ok

				url := fmt.Sprintf("http://localhost:%d/version", port)
				fmt.Println("测试 api 连通性:", url)
				req := utils.NewHttp().NewRequest()
				resp, httperr := req.Get(url)
				fmt.Println("连通性结果： ", httperr == nil)
				if httperr != nil {
					_, _ = j.helper.P2p().Close(fmt.Sprintf("/p2p/%s", string(val.PeerId)))
					failSet[int(val.Index)] = "fail"
					continue
				}

				var version VersionVo
				bandErr := json.Unmarshal(resp.Body(), &version)
				if bandErr != nil {
					_, _ = j.helper.P2p().Close(fmt.Sprintf("/p2p/%s", string(val.PeerId)))
					failSet[int(val.Index)] = "fail"
					continue
				}
				fmt.Println("provider version: ", version.Version)

				c, err := types.NewCall(
					j.meta,
					"ResourceOrder.create_order_info",
					resourceIndex,
					types.NewU32(uint32(data.LeaseTerm)),
					"",
					types.NewU32(uint32(deployType)),
				)
				if err != nil {
					fmt.Println(err.Error())
					_, _ = j.helper.P2p().Close(fmt.Sprintf("/p2p/%s", string(val.PeerId)))
					failSet[int(val.Index)] = "fail"
					continue
				}
				var events pallet.MyEventRecords
				err = pallet.CallAndWatch(api, c, j.meta, func(header *types.Header) error {
					fmt.Printf("资源占用成功，资源号： %d, 交易序号： %d", resourceIndex, header.Number)
					// get order index
					e, err := pallet.GetEvent(api, j.meta, uint64(header.Number))
					events = *e
					return err
				}, pair)

				fmt.Println("watch event error : ", err)
				if err != nil {
					_, _ = j.helper.P2p().Close(fmt.Sprintf("/p2p/%s", string(val.PeerId)))
					failSet[int(val.Index)] = "fail"
					continue
				}

				var orderIndex types.U64
				for _, e := range events.ResourceOrder_CreateOrderSuccess {
					if e.ResourceIndex == resourceIndex {
						orderIndex = e.OrderIndex
					}
				}

				fmt.Println("query orderIndex: ", orderIndex)

				// check provider exec order
				number, err := pallet.GetBlockNumber(api)
				if err != nil {
					_ = pallet.CancelOrder(api, j.meta, pair, int(orderIndex))
					_, _ = j.helper.P2p().Close(fmt.Sprintf("/p2p/%s", string(val.PeerId)))
					continue
				}

				accept := false

				for i := number; i < number+10; i++ {
					events, err := pallet.GetEvent(api, j.meta, i)
					if err != nil {
						continue
					}
					if len(events.ResourceOrder_OrderExecSuccess) > 0 {
						for _, e := range events.ResourceOrder_OrderExecSuccess {
							if e.OrderIndex == orderIndex {
								accept = true
								fmt.Println("confirm order : ", orderIndex)
								break
							}
						}
						if accept {
							break
						}
					}
					time.Sleep(time.Second * 6)
				}

				if !accept {
					_ = pallet.CancelOrder(api, j.meta, pair, int(orderIndex))
					_, _ = j.helper.P2p().Close(fmt.Sprintf("/p2p/%s", string(val.PeerId)))
					failSet[int(val.Index)] = "fail"
					continue
				}

				_ = j.helper.App().UpdatePeerIdAndOrderIndex(
					j.appID,
					int(orderIndex),
					int(resourceIndex),
					string(val.PeerId),
				)

				err = j.helper.App().UpdateApplicationP2pForwardPort(j.appID, port)
				if err != nil {
					fmt.Println("query p2p max port fail")
				}

				j.si.Status = queue.Succeeded
				j.si.Error = ""
				si <- j.si
				return j.si, nil
			}
		}

		time.Sleep(time.Second * 30)

	}
	j.si.Status = queue.Failed
	j.si.Error = "NO_RESOURCE_TO_USE"
	si <- j.si
	return j.si, errors.New("NO_RESOURCE_TO_USE")
}

func (j *WaitResourceJob) Status() queue.StatusInfo {
	return j.si
}
