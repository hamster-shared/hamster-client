package v2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hamster-client/config"
	ethAbi "hamster-client/module/abi"
	"hamster-client/module/account"
	"hamster-client/module/application"
	"hamster-client/module/deploy"
	"hamster-client/module/keystorage"
	"hamster-client/module/p2p"
	"hamster-client/module/pallet"
	"hamster-client/module/queue"
	"hamster-client/module/wallet"
	"hamster-client/utils"
	"math/big"
	"strconv"
	"sync"
	"time"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/sirupsen/logrus"
)

var (
	mutex sync.RWMutex
)

func init() {
	mutex = sync.RWMutex{}
}

type PullImageJob struct {
	statusInfo         queue.StatusInfo
	applicationId      int
	applicationService application.Service
	walletService      wallet.Service
	p2pService         p2p.Service
	accountService     account.Service
}

func (j *PullImageJob) InitStatus() {
	j.statusInfo.Name = "Service Pull"
	j.statusInfo.Status = queue.None
}

func (j *PullImageJob) Run(sc chan queue.StatusInfo) (queue.StatusInfo, error) {
	j.statusInfo.Status = queue.Running
	sc <- j.statusInfo
	vo, err := j.applicationService.QueryApplicationById(j.applicationId)
	if err != nil {
		fmt.Println("query application fail, err is :", err)
		j.statusInfo.Status = queue.Failed
		j.statusInfo.Error = err.Error()
		sc <- j.statusInfo
		return j.statusInfo, err
	}

	pair, err := j.walletService.GetWalletKeypair()
	if err != nil {
		fmt.Println("query wallet fail, err is :", err)
		j.statusInfo.Status = queue.Failed
		j.statusInfo.Error = err.Error()
		sc <- j.statusInfo
		return j.statusInfo, err
	}
	fmt.Println("pull before: reForwardLink: ", vo.PeerId)
	if _, err := j.p2pService.GetSetting(); err != nil {
		_ = j.p2pService.InitSetting()
	}
	err = reForwardLink(j.p2pService, vo.P2pForwardPort, vo.PeerId)
	if err != nil {
		fmt.Println("reconnect fail, err is :", err)
		j.statusInfo.Status = queue.Failed
		j.statusInfo.Error = err.Error()
		sc <- j.statusInfo
		return j.statusInfo, err
	}
	url := fmt.Sprintf("http://localhost:%d/api/v1/thegraph/pullImage", vo.P2pForwardPort)
	for i := 0; i < 3; i++ {
		req := utils.NewHttp().NewRequest()
		req.SetHeader("SS58AuthData", utils.GetSS58AuthDataWithKeyringPair(pair))
		response, err := req.Post(url)
		if err != nil {
			j.statusInfo.Error = err.Error()
			fmt.Println(string(response.Body()))
			continue
		}
		if response.IsSuccess() {
			j.statusInfo.Status = queue.Succeeded
			j.statusInfo.Error = ""
			sc <- j.statusInfo
			fmt.Println("========  pull image success ===============")
			return j.statusInfo, nil
		} else {
			time.Sleep(time.Second * 3)
			continue
		}
	}
	j.statusInfo.Status = queue.Failed
	j.statusInfo.Error = "api response fail"
	sc <- j.statusInfo

	return j.statusInfo, errors.New("api response fail")

}

func (j *PullImageJob) Status() queue.StatusInfo {
	return j.statusInfo
}

func NewPullImageJob(service application.Service, applicationId int, p2pService p2p.Service, accountService account.Service, walletService wallet.Service) PullImageJob {
	return PullImageJob{
		applicationId:      applicationId,
		applicationService: service,
		walletService:      walletService,
		p2pService:         p2pService,
		accountService:     accountService,
	}
}

type WaitResourceJob struct {
	statusInfo         queue.StatusInfo
	api                *gsrpc.SubstrateAPI
	meta               *types.Metadata
	accountService     account.Service
	applicationService application.Service
	p2pService         p2p.Service
	applicationId      int
	walletService      wallet.Service
	bond               bool
}

func (j *WaitResourceJob) InitStatus() {
	j.statusInfo.Name = "Resource Waiting"
	j.statusInfo.Status = queue.None
	j.bond = false
}

func (j *WaitResourceJob) Run(sc chan queue.StatusInfo) (queue.StatusInfo, error) {
	j.statusInfo.Status = queue.Running
	sc <- j.statusInfo

	pair, err := j.walletService.GetWalletKeypair()
	if err != nil {
		j.statusInfo.Status = queue.Failed
		j.statusInfo.Error = "WALLET_LOAD_ERROR"
		sc <- j.statusInfo
		return j.statusInfo, err
	}
	if pair.Address == "" {
		j.statusInfo.Status = queue.Failed
		j.statusInfo.Error = errors.New("get Keypair failed").Error()
		sc <- j.statusInfo
		return j.statusInfo, err
	}

	if !j.bond {
		// 100 unit
		err = pallet.Bond(j.api, j.meta, 100000000000000, pair)
		if err != nil {
			j.statusInfo.Status = queue.Failed
			j.statusInfo.Error = "WALLET_LOAD_ERROR"
			sc <- j.statusInfo
			return j.statusInfo, err
		}
	}

	for i := 0; i < 60; i++ {

		mapData, err := pallet.GetResourceList(j.meta, j.api, func(resource *pallet.ComputingResource) bool {
			return resource.Status.IsUnused
		})
		if err != nil {
			j.statusInfo.Status = queue.Failed
			j.statusInfo.Error = err.Error()
			sc <- j.statusInfo
			return j.statusInfo, err
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
				data, err := j.applicationService.QueryApplicationById(j.applicationId)
				if err != nil {
					fmt.Println("get application failed,err is: ", err)
					continue
				}

				// check p2p connection can be connected
				port := j.applicationService.QueryNextP2pPort()
				err = j.p2pService.LinkByProtocol(config.ProviderProtocol, port, string(val.PeerId))

				if err != nil {
					fmt.Println("create p2p network forward fail")
					failSet[int(val.Index)] = "fail"
					_, _ = j.p2pService.Close(fmt.Sprintf("/p2p/%s", string(val.PeerId)))
					continue
				}

				// check http is Ok
				url := fmt.Sprintf("http://localhost:%d/version", port)
				req := utils.NewHttp().NewRequest()
				resp, err := req.Get(url)
				if err != nil {
					_, _ = j.p2pService.Close(fmt.Sprintf("/p2p/%s", string(val.PeerId)))
					failSet[int(val.Index)] = "fail"
					continue
				}

				var version VersionVo
				err = json.Unmarshal(resp.Body(), &version)
				if err != nil {
					_, _ = j.p2pService.Close(fmt.Sprintf("/p2p/%s", string(val.PeerId)))
					failSet[int(val.Index)] = "fail"
					continue
				}
				fmt.Println("provider version: ", version.Version)

				logrus.Infof("调用 ResourceOrder.create_order_info, 参数：%v, %v", resourceIndex, data.LeaseTerm)
				c, err := types.NewCall(j.meta, "ResourceOrder.create_order_info", resourceIndex, types.NewU32(uint32(data.LeaseTerm)), " ", types.NewU32(0))
				if err != nil {
					fmt.Println(err.Error())
					logrus.Error("调用 ResourceOrder.create_order_info 失败： ",err.Error())
					_, _ = j.p2pService.Close(fmt.Sprintf("/p2p/%s", string(val.PeerId)))
					failSet[int(val.Index)] = "fail"
					continue
				}
				var events pallet.MyEventRecords
				err = pallet.CallAndWatch(j.api, c, j.meta, func(header *types.Header) error {
					fmt.Printf("资源占用成功，资源号： %d, 交易序号： %d", resourceIndex, header.Number)
					// get order index
					logrus.Debug("调用 pallet.GetEvent")
					e, err := pallet.GetEvent(j.api, j.meta, uint64(header.Number))
					events = *e
					if err != nil {
						logrus.Error("获取事件失败： ", err.Error())
					}
					return err
				}, pair)
				if err != nil {
					logrus.Errorf("pallet.CallAndWatch 失败： %v", err.Error())
					_, _ = j.p2pService.Close(fmt.Sprintf("/p2p/%s", string(val.PeerId)))
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
				number, err := pallet.GetBlockNumber(j.api)
				if err != nil {
					_ = pallet.CancelOrder(j.api, j.meta, pair, int(orderIndex))
					_, _ = j.p2pService.Close(fmt.Sprintf("/p2p/%s", string(val.PeerId)))
					continue
				}

				accept := false

				for i := number; i < number+10; i++ {
					events, err := pallet.GetEvent(j.api, j.meta, i)
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
					_ = pallet.CancelOrder(j.api, j.meta, pair, int(orderIndex))
					_, _ = j.p2pService.Close(fmt.Sprintf("/p2p/%s", string(val.PeerId)))
					failSet[int(val.Index)] = "fail"
					continue
				}

				_ = j.applicationService.UpdatePeerIdAndOrderIndex(j.applicationId, int(orderIndex), int(resourceIndex), string(val.PeerId))

				err = j.applicationService.UpdateApplicationP2pForwardPort(j.applicationId, port)
				if err != nil {
					fmt.Println("query p2p max port fail")
				}

				j.statusInfo.Status = queue.Succeeded
				j.statusInfo.Error = ""
				sc <- j.statusInfo
				return j.statusInfo, nil
			}
		}

		time.Sleep(time.Second * 30)

	}
	j.statusInfo.Status = queue.Failed
	j.statusInfo.Error = "NO_RESOURCE_TO_USE"
	sc <- j.statusInfo
	return j.statusInfo, errors.New("NO_RESOURCE_TO_USE")
}

func (j *WaitResourceJob) Status() queue.StatusInfo {
	return j.statusInfo
}

func NewWaitResourceJob(api *gsrpc.SubstrateAPI, accountService account.Service, applicationService application.Service, p2pService p2p.Service, applicationId int, walletService wallet.Service) (*WaitResourceJob, error) {

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, err
	}

	return &WaitResourceJob{
		api:                api,
		meta:               meta,
		accountService:     accountService,
		applicationService: applicationService,
		p2pService:         p2pService,
		applicationId:      applicationId,
		walletService:      walletService,
	}, nil
}

type GraphStakingJob struct {
	statusInfo        queue.StatusInfo
	id                int
	chainId           int64
	network           string
	keyStorageService keystorage.Service
}

func (g *GraphStakingJob) InitStatus() {
	g.statusInfo.Name = "TheGraph Staking"
	g.statusInfo.Status = queue.None
}

func NewGraphStakingJob(service keystorage.Service, applicationId int, network string, chainId int64) GraphStakingJob {
	return GraphStakingJob{
		id:                applicationId,
		chainId:           chainId,
		network:           network,
		keyStorageService: service,
	}
}

func (g *GraphStakingJob) Run(sc chan queue.StatusInfo) (queue.StatusInfo, error) {
	g.statusInfo.Status = queue.Running
	sc <- g.statusInfo
	//Obtain pledge amount and mnemonic words
	var param deploy.DeployParameter
	jsonParam := g.keyStorageService.Get("graph_" + strconv.Itoa(g.id))
	if err := json.Unmarshal([]byte(jsonParam), &param); err != nil {
		g.statusInfo.Status = queue.Failed
		g.statusInfo.Error = err.Error()
		sc <- g.statusInfo
		return g.statusInfo, err
	}
	if param.Initialization.AccountMnemonic == "" {
		err := errors.New("Saved mnemonic is empty")
		g.statusInfo.Status = queue.Failed
		g.statusInfo.Error = err.Error()
		sc <- g.statusInfo
		return g.statusInfo, err
	}
	//Get address from mnemonic
	addr, err := ethAbi.GetPrivateKeyHexStringWithMnemonic(param.Initialization.AccountMnemonic)
	if err != nil {
		fmt.Println("Get address from mnemonic failed, err is :", err)
		g.statusInfo.Status = queue.Failed
		g.statusInfo.Error = err.Error()
		sc <- g.statusInfo
		return g.statusInfo, err
	}
	address := ethAbi.GetEthAddress(addr)
	//Get eth client
	client, err := ethAbi.GetEthClient(g.network)
	if err != nil {
		fmt.Println("Get eth client failed, err is :", err)
		g.statusInfo.Status = queue.Failed
		g.statusInfo.Error = err.Error()
		sc <- g.statusInfo
		return g.statusInfo, err
	}
	// Obtain the agent pledge address
	stakingAddress, err := ethAbi.StakeProxyFactoryAbiGetStakingAddress(context.Background(), address, client)
	if err != nil {
		fmt.Println("Get agent proxy address failed, err is :", err)
		g.statusInfo.Status = queue.Failed
		g.statusInfo.Error = err.Error()
		sc <- g.statusInfo
		return g.statusInfo, err
	}
	//Get private key from mnemonic
	privateKey, err := ethAbi.GetPrivateKeyWithMnemonic(param.Initialization.AccountMnemonic)
	if err != nil {
		fmt.Println("Get private key from mnemonic failed, err is :", err)
		g.statusInfo.Status = queue.Failed
		g.statusInfo.Error = err.Error()
		sc <- g.statusInfo
		return g.statusInfo, err
	}
	//Convert the pledged amount into Wei
	stakingAmount := utils.ToWei18(int64(param.Staking.PledgeAmount))
	if stakingAddress == ethAbi.GetEthAddress("0") {
		//Create agent pledge address
		err = ethAbi.StakeProxyFactoryAbiCreateStakingContract(address, client, big.NewInt(g.chainId), privateKey, context.Background(), client)
		if err != nil {
			fmt.Println("Create agent pledge address failed, err is :", err)
			g.statusInfo.Status = queue.Failed
			g.statusInfo.Error = err.Error()
			sc <- g.statusInfo
			return g.statusInfo, err
		}
		time.Sleep(time.Second * 2)
		// Get the agent pledge address again
		stakingAddress, err = ethAbi.StakeProxyFactoryAbiGetStakingAddress(context.Background(), address, client)
		if err != nil {
			fmt.Println("Get the agent pledge address again failed, err is :", err)
			g.statusInfo.Status = queue.Failed
			g.statusInfo.Error = err.Error()
			sc <- g.statusInfo
			return g.statusInfo, err
		}
		// Authorize the agency pledge address
		err = ethAbi.Ecr20AbiApprove(stakingAddress, client, big.NewInt(g.chainId), stakingAmount, privateKey, context.Background(), client)
		if err != nil {
			fmt.Println("approve failed, err is :", err)
			g.statusInfo.Status = queue.Failed
			g.statusInfo.Error = err.Error()
			sc <- g.statusInfo
			return g.statusInfo, err
		}
		time.Sleep(time.Second * 3)
		//GRT pledge
		err = ethAbi.StakeDistributionProxyAbiStaking(stakingAddress, client, big.NewInt(g.chainId), stakingAmount, privateKey, context.Background(), client)
		if err != nil {
			fmt.Println("staking failed, err is :", err)
			g.statusInfo.Status = queue.Failed
			g.statusInfo.Error = err.Error()
			sc <- g.statusInfo
			return g.statusInfo, err
		}
	} else {
		// get staking amount
		amount, err := ethAbi.StakeDistributionProxyAbiGetStakingAmount(context.Background(), stakingAddress, client)
		if err != nil {
			fmt.Println("get stake amount failed, err is :", err)
			g.statusInfo.Status = queue.Failed
			g.statusInfo.Error = err.Error()
			sc <- g.statusInfo
			return g.statusInfo, err
		}
		if amount.Cmp(big.NewInt(0)) == 0 {
			// Authorize the agency pledge address
			err = ethAbi.Ecr20AbiApprove(stakingAddress, client, big.NewInt(g.chainId), stakingAmount, privateKey, context.Background(), client)
			if err != nil {
				fmt.Println("approve failed, err is :", err)
				g.statusInfo.Status = queue.Failed
				g.statusInfo.Error = err.Error()
				sc <- g.statusInfo
				return g.statusInfo, err
			}
			time.Sleep(time.Second * 3)
			//GRT pledge
			err = ethAbi.StakeDistributionProxyAbiStaking(stakingAddress, client, big.NewInt(g.chainId), stakingAmount, privateKey, context.Background(), client)
			if err != nil {
				fmt.Println("staking failed, err is :", err)
				g.statusInfo.Status = queue.Failed
				g.statusInfo.Error = err.Error()
				sc <- g.statusInfo
				return g.statusInfo, err
			}
		}
	}
	param.Deployment.IndexerAddress = addr
	param.Staking.AgentAddress = stakingAddress.Hex()
	jsonData, err := json.Marshal(param)
	if err != nil {
		g.statusInfo.Status = queue.Failed
		g.statusInfo.Error = err.Error()
		sc <- g.statusInfo
		return g.statusInfo, err
	}
	g.keyStorageService.Set("graph_"+strconv.Itoa(int(g.id)), string(jsonData))
	g.statusInfo.Status = queue.Succeeded
	g.statusInfo.Error = ""
	sc <- g.statusInfo
	return g.statusInfo, nil
}

func (g *GraphStakingJob) Status() queue.StatusInfo {
	return g.statusInfo
}

type ServiceDeployJob struct {
	statusInfo         queue.StatusInfo
	id                 int
	deployService      deploy.Service
	keyStorageService  keystorage.Service
	p2pService         p2p.Service
	accountService     account.Service
	applicationService application.Service
	walletService      wallet.Service
}

func (s *ServiceDeployJob) InitStatus() {
	s.statusInfo.Name = "Service Deploy"
	s.statusInfo.Status = queue.None
}

func NewServiceDeployJob(service keystorage.Service, deployService deploy.Service, applicationId int, p2pService p2p.Service, accountService account.Service, applicationService application.Service, walletService wallet.Service) ServiceDeployJob {
	return ServiceDeployJob{
		id:                 applicationId,
		keyStorageService:  service,
		deployService:      deployService,
		p2pService:         p2pService,
		accountService:     accountService,
		applicationService: applicationService,
		walletService:      walletService,
	}
}

func (s *ServiceDeployJob) Run(sc chan queue.StatusInfo) (queue.StatusInfo, error) {
	s.statusInfo.Status = queue.Running
	sc <- s.statusInfo
	var param deploy.DeployParameter
	var sendData deploy.DeployParams
	jsonParam := s.keyStorageService.Get("graph_" + strconv.Itoa(s.id))
	if err := json.Unmarshal([]byte(jsonParam), &param); err != nil {
		s.statusInfo.Status = queue.Failed
		s.statusInfo.Error = err.Error()
		sc <- s.statusInfo
		return s.statusInfo, err
	}
	sendData.EthereumUrl = param.Deployment.EthereumUrl
	sendData.EthereumNetwork = param.Deployment.EthereumNetwork
	sendData.NodeEthereumUrl = param.Deployment.NodeEthereumUrl
	sendData.IndexerAddress = param.Staking.AgentAddress
	sendData.Mnemonic = param.Initialization.AccountMnemonic
	vo, err := s.applicationService.QueryApplicationById(s.id)
	if err != nil {
		fmt.Println("query application fail, err is :", err)
		s.statusInfo.Status = queue.Failed
		s.statusInfo.Error = err.Error()
		sc <- s.statusInfo
		return s.statusInfo, err
	}
	fmt.Println("deploy before: reForwardLink")
	err = reForwardLink(s.p2pService, vo.P2pForwardPort, vo.PeerId)
	if err != nil {
		fmt.Println("reconnect fail, err is :", err)
		s.statusInfo.Status = queue.Failed
		s.statusInfo.Error = err.Error()
		sc <- s.statusInfo
		return s.statusInfo, err
	}
	_, err = s.deployService.DeployGraph(s.id, sendData)
	if err != nil {
		fmt.Println("deploy service failed, err is :", err)
		s.statusInfo.Status = queue.Failed
		s.statusInfo.Error = err.Error()
		sc <- s.statusInfo
		return s.statusInfo, err
	}
	s.statusInfo.Status = queue.Succeeded
	s.statusInfo.Error = ""
	sc <- s.statusInfo
	return s.statusInfo, nil
}

func (s *ServiceDeployJob) Status() queue.StatusInfo {
	return s.statusInfo
}

func reForwardLink(p2pService p2p.Service, port int, peerId string) error {
	protocol := config.ProviderProtocol
	err := p2pService.LinkByProtocol(protocol, port, peerId)
	return err
}
