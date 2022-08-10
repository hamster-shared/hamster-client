package v2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"hamster-client/config"
	ethAbi "hamster-client/module/abi"
	"hamster-client/module/account"
	"hamster-client/module/application"
	"hamster-client/module/deploy"
	"hamster-client/module/keystorage"
	"hamster-client/module/p2p"
	"hamster-client/module/pallet"
	"hamster-client/module/queue"
	"hamster-client/utils"
	"math/big"
	"strconv"
	"time"
)

type PullImageJob struct {
	statusInfo         queue.StatusInfo
	applicationId      int
	applicationService application.Service
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

	url := fmt.Sprintf("http://localhost:%d/api/v1/thegraph/pullImage", vo.P2pForwardPort)
	for i := 0; i < 3; i++ {
		response, err := utils.NewHttp().NewRequest().Post(url)
		if err != nil {
			j.statusInfo.Error = err.Error()
			continue
		}
		if response.IsSuccess() {
			j.statusInfo.Status = queue.Succeeded
			j.statusInfo.Error = ""
			sc <- j.statusInfo
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

func NewPullImageJob(service application.Service, applicationId int) PullImageJob {
	return PullImageJob{
		applicationId:      applicationId,
		applicationService: service,
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
}

func (j *WaitResourceJob) InitStatus() {
	j.statusInfo.Name = "Resource Waiting"
	j.statusInfo.Status = queue.None
}

func (j *WaitResourceJob) Run(sc chan queue.StatusInfo) (queue.StatusInfo, error) {
	j.statusInfo.Status = queue.Running
	sc <- j.statusInfo

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

		for _, val := range mapData {

			if val.Status.IsUnused {
				fmt.Println("发现未使用资源，占用。资源号：", val.Index)
				data, err := j.applicationService.QueryApplicationById(j.applicationId)
				if err != nil {
					fmt.Println("get application failed,err is: ", err)
					continue
				}
				c, err := types.NewCall(j.meta, "ResourceOrder.create_order_info", val.Index, types.NewU32(uint32(data.LeaseTerm)), "")
				if err != nil {
					continue
				}
				err = pallet.CallAndWatch(j.api, c, j.meta, func(header *types.Header) error {
					fmt.Println("资源占用成功，资源号：", val.Index)
					return nil
				})
				if err != nil {
					continue
				}

				ac, _ := j.accountService.GetAccount()
				ac.PeerId = string(val.PeerId)
				j.accountService.SaveAccount(&ac)

				port := j.applicationService.QueryNextP2pPort()

				err = j.applicationService.UpdateApplicationP2pForwardPort(j.applicationId, port)
				if err != nil {
					fmt.Println("query p2p max port fail")
				}
				err = j.p2pService.LinkByProtocol("/x/provider", port, ac.PeerId)

				if err != nil {
					fmt.Println("create p2p network forward fail")
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

func NewWaitResourceJob(api *gsrpc.SubstrateAPI, accountService account.Service, applicationService application.Service, p2pService p2p.Service, applicationId int) (*WaitResourceJob, error) {

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
	}, nil
}

type GraphStakingJob struct {
	statusInfo        queue.StatusInfo
	id                int
	keyStorageService keystorage.Service
}

func (g *GraphStakingJob) InitStatus() {
	g.statusInfo.Name = "TheGraph Staking"
	g.statusInfo.Status = queue.None
}

func NewGraphStakingJob(service keystorage.Service, applicationId int) GraphStakingJob {
	return GraphStakingJob{
		id:                applicationId,
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
	client, err := ethAbi.GetEthClient(config.EndpointUrl)
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
	if stakingAddress == ethAbi.GetEthAddress("0") {
		//Create agent pledge address
		err = ethAbi.StakeProxyFactoryAbiCreateStakingContract(address, client, big.NewInt(4), privateKey)
		if err != nil {
			fmt.Println("Create agent pledge address failed, err is :", err)
			g.statusInfo.Status = queue.Failed
			g.statusInfo.Error = err.Error()
			sc <- g.statusInfo
			return g.statusInfo, err
		}
		// Get the agent pledge address again
		stakingAddress, err = ethAbi.StakeProxyFactoryAbiGetStakingAddress(context.Background(), address, client)
		if err != nil {
			fmt.Println("Get the agent pledge address again failed, err is :", err)
			g.statusInfo.Status = queue.Failed
			g.statusInfo.Error = err.Error()
			sc <- g.statusInfo
			return g.statusInfo, err
		}
		//Convert the pledged amount into Wei
		stakingAmount := utils.ToWei18(int64(param.Staking.PledgeAmount))
		// Authorize the agency pledge address
		err = ethAbi.Ecr20AbiApprove(stakingAddress, client, big.NewInt(4), stakingAmount, privateKey)
		if err != nil {
			fmt.Println("approve failed, err is :", err)
			g.statusInfo.Status = queue.Failed
			g.statusInfo.Error = err.Error()
			sc <- g.statusInfo
			return g.statusInfo, err
		}
		//GRT pledge
		err = ethAbi.StakeDistributionProxyAbiStaking(stakingAddress, client, big.NewInt(4), stakingAmount, privateKey)
		if err != nil {
			fmt.Println("staking failed, err is :", err)
			g.statusInfo.Status = queue.Failed
			g.statusInfo.Error = err.Error()
			sc <- g.statusInfo
			return g.statusInfo, err
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
	statusInfo        queue.StatusInfo
	id                int
	deployService     deploy.Service
	keyStorageService keystorage.Service
}

func (s *ServiceDeployJob) InitStatus() {
	s.statusInfo.Name = "Service Deploy"
	s.statusInfo.Status = queue.None
}

func NewServiceDeployJob(service keystorage.Service, deployService deploy.Service, applicationId int) ServiceDeployJob {
	return ServiceDeployJob{
		id:                applicationId,
		keyStorageService: service,
		deployService:     deployService,
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
	_, err := s.deployService.DeployGraph(s.id, sendData)
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
