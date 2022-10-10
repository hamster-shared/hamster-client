package v2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hamster-client/config"
	"hamster-client/module/account"
	"hamster-client/module/application"
	"hamster-client/module/deploy"
	"hamster-client/module/keystorage"
	"hamster-client/module/p2p"
	"hamster-client/module/pallet"
	queue2 "hamster-client/module/queue"
	"hamster-client/module/wallet"
	"hamster-client/utils"
	"strconv"
)

type ServiceImpl struct {
	ctx                context.Context
	db                 *gorm.DB
	keyStorageService  keystorage.Service
	accountService     account.Service
	applicationService application.Service
	p2pService         p2p.Service
	deployService      deploy.Service
	walletService      wallet.Service
	queueService       queue2.Service
}

type ServiceDeploySaveService interface {
	saveDeployParam(appData application.Application, paramData interface{}, tx *gorm.DB) error
	saveJsonParam(id string, paramData interface{}) error
	deployJob(appData application.Application)
}

func NewServiceImpl(ctx context.Context, db *gorm.DB, keyStorageService keystorage.Service, accountService account.Service, applicationService application.Service, p2pService p2p.Service, deployService deploy.Service, walletService wallet.Service, queueService queue2.Service) ServiceImpl {
	return ServiceImpl{ctx, db, keyStorageService, accountService, applicationService, p2pService, deployService, walletService, queueService}
}

func (g *ServiceImpl) SaveGraphDeployParameterAndApply(addParam AddParam) (AddApplicationVo, error) {

	var saveService ServiceDeploySaveService

	if addParam.ServiceType == application.TYPE_Thegraph {
		saveService = &ThegraphDeploySaveServiceImpl{*g}
		//} else if val, isKeyExists := application.GetDeployEnumMap()[addParam.ServiceType]; isKeyExists {
		//saveService = &CommonDeploySaveServiceImpl{*g, val}
	} else {
		return AddApplicationVo{}, errors.New("Unsupport deploy type!")
	}

	var applyData application.Application
	var applicationVo AddApplicationVo
	err := g.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("name=?", addParam.Name).First(&applyData).Error; err != gorm.ErrRecordNotFound {
			return errors.New(fmt.Sprintf("application:%s already exists", addParam.Name))
		}
		applyData.Name = addParam.Name
		applyData.ServiceType = addParam.ServiceType
		applyData.SelectNodeType = addParam.SelectNodeType
		applyData.LeaseTerm = addParam.LeaseTerm
		applyData.Status = application.Deploying
		if err := tx.Create(&applyData).Error; err != nil {
			return err
		}
		if err := saveService.saveDeployParam(applyData, addParam, tx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		applicationVo.Result = false
		return applicationVo, err
	}
	err = saveService.saveJsonParam(strconv.Itoa(int(applyData.ID)), addParam)
	if err != nil {
		applicationVo.Result = false
		return applicationVo, err
	}

	go saveService.deployJob(applyData)
	applicationVo.Result = true
	applicationVo.ID = applyData.ID
	return applicationVo, nil
}

func (g *ServiceImpl) DeleteGraphDeployParameterAndApply(id int) (bool, error) {
	_ = g.queueService.StopQueue(id)
	app, err := g.applicationService.QueryApplicationById(id)
	// close p2p port
	_, _ = g.p2pService.Close(fmt.Sprintf("/p2p/%s", app.PeerId))
	if err != nil {
		return false, err
	}
	err = g.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Debug().Where("id = ?", id).Delete(&application.Application{}).Error; err != nil {
			return err
		}
		if err := tx.Debug().Where("application_id = ?", id).Delete(GraphDeployParameter{}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	//delete key storage
	//stop docker
	keypair, err := g.walletService.GetWalletKeypair()
	if err != nil {
		return true, nil
	}
	if api, err := g.accountService.GetSubstrateApi(); err == nil {
		if meta, err := api.RPC.State.GetMetadataLatest(); err == nil {
			_ = pallet.CancelOrder(api, meta, keypair, app.OrderIndex)
		}
	}

	return true, nil
}

func (g *ServiceImpl) RetryDeployGraphJob(applicationId int, runNow bool) error {
	var queue queue2.Queue
	queue, err := queue2.GetQueue(applicationId)
	if err != nil {
		var data application.Application
		result := g.db.Where("id = ? ", applicationId).First(&data)
		if result.Error != nil {
			fmt.Println("query application failed,error is: ", result.Error)
			return result.Error
		}
		pluginDeployInfo := config.PluginMap[data.SelectNodeType]
		stakingJob := NewGraphStakingJob(g.keyStorageService, applicationId, pluginDeployInfo.EndpointUrl, pluginDeployInfo.ChainId)
		substrateApi, err := g.accountService.GetSubstrateApi()
		if err != nil {
			return err
		}
		const deployType = 0
		waitResourceJob, _ := NewWaitResourceJob(substrateApi, g.accountService, g.applicationService, g.p2pService, applicationId, g.walletService, deployType)

		pullJob := NewPullImageJob(g.applicationService, applicationId, g.p2pService, g.accountService, g.walletService)

		deployJob := NewServiceDeployJob(g.keyStorageService, g.deployService, applicationId, g.p2pService, g.accountService, g.applicationService, g.walletService)

		queue, err = queue2.NewQueue(applicationId, g.db, &stakingJob, waitResourceJob, &pullJob, &deployJob)
		if err != nil {
			fmt.Println("new queue failed, err is: ", err)
			return err
		}
	}

	if runNow {
		channel := make(chan struct{})
		go queue.Start(channel)
		<-channel
	} else {
		statusInfo, err := queue.GetStatus()
		if err != nil {
			log.Errorf("get status failed, error: %v", err)
			return nil
		}
		for _, job := range statusInfo {
			if job.Status == queue2.Running {
				log.Infof("job %s status is running, set to failed", job.Name)
				queue.SetJobStatus(job.Name, queue2.StatusInfo{
					Name:   job.Name,
					Status: queue2.Failed,
					Error:  "Abnormal exit",
				})
			}
		}
	}

	return nil
}

func (g *ServiceImpl) GetQueueInfo(applicationId int) (QueueInfo, error) {
	var data application.Application
	result := g.db.Where("id = ? ", applicationId).First(&data)
	if result.Error != nil {
		fmt.Println("query application failed,error is: ", result.Error)
		return QueueInfo{}, result.Error
	}
	pluginDeployInfo := config.PluginMap[data.SelectNodeType]
	stakingJob := NewGraphStakingJob(g.keyStorageService, applicationId, pluginDeployInfo.EndpointUrl, pluginDeployInfo.ChainId)
	substrateApi, err := g.accountService.GetSubstrateApi()
	if err != nil {
		return QueueInfo{}, err
	}
	const deployType = 0
	waitResourceJob, _ := NewWaitResourceJob(substrateApi, g.accountService, g.applicationService, g.p2pService, applicationId, g.walletService, deployType)

	pullJob := NewPullImageJob(g.applicationService, applicationId, g.p2pService, g.accountService, g.walletService)

	deployJob := NewServiceDeployJob(g.keyStorageService, g.deployService, applicationId, g.p2pService, g.accountService, g.applicationService, g.walletService)

	_, err = queue2.NewQueue(applicationId, g.db, &stakingJob, waitResourceJob, &pullJob, &deployJob)
	if err != nil {
		fmt.Println("new queue failed,err is: ", err)
		return QueueInfo{}, err
	}
	if err != nil {
		return QueueInfo{}, err
	}
	info, err := g.queueService.GetStatusInfo(applicationId)
	if err != nil {
		return QueueInfo{Info: info}, err
	}
	return QueueInfo{Info: info}, nil
}

func (g *ServiceImpl) GraphStart(appID int, deploymentID string) error {
	port, peerId, err := g.getP2pPort(appID)
	_ = g.p2pService.LinkByProtocol(config.ProviderProtocol, port, peerId)
	if err != nil {
		return err
	}
	keyringPair, err := g.walletService.GetWalletKeypair()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("http://localhost:%d/api/v1/thegraph/start", port)
	req := utils.NewHttp().NewRequest()
	req.SetHeader("SS58AuthData", utils.GetSS58AuthDataWithKeyringPair(keyringPair))
	resp, err := req.Get(fmt.Sprintf("%s?deploymentID=%s", url, deploymentID))
	if err != nil {
		return err
	}
	if resp.IsSuccess() {
		return nil
	}
	return parseResponseError(resp)
}

func (g *ServiceImpl) GraphStop(appID int, deploymentID string) error {
	port, peerId, err := g.getP2pPort(appID)
	_ = g.p2pService.LinkByProtocol(config.ProviderProtocol, port, peerId)
	if err != nil {
		return err
	}
	keyringPair, err := g.walletService.GetWalletKeypair()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("http://localhost:%d/api/v1/thegraph/stop", port)
	req := utils.NewHttp().NewRequest()
	req.SetHeader("SS58AuthData", utils.GetSS58AuthDataWithKeyringPair(keyringPair))
	resp, err := req.Get(fmt.Sprintf("%s?deploymentID=%s", url, deploymentID))
	if err != nil {
		return err
	}
	if resp.IsSuccess() {
		return nil
	}
	return parseResponseError(resp)
}

func (g *ServiceImpl) GraphRules(appID int) ([]GraphRule, error) {
	port, peerId, err := g.getP2pPort(appID)
	_ = g.p2pService.LinkByProtocol(config.ProviderProtocol, port, peerId)

	fmt.Println("#### p2p port : ", port)
	if err != nil {
		return nil, err
	}
	keyringPair, err := g.walletService.GetWalletKeypair()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("http://localhost:%d/api/v1/thegraph/rules", port)
	req := utils.NewHttp().NewRequest()
	req.SetHeader("SS58AuthData", utils.GetSS58AuthDataWithKeyringPair(keyringPair))
	resp, err := req.Get(url)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, parseResponseError(resp)
	}
	var result Result
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	filter := func(in []GraphRule) []GraphRule {
		n := 0
		for _, val := range in {
			if val.DecisionBasis == "always" {
				in[n] = val
				n++
			}
		}
		return in[:n]
	}

	return filter(result.Result), nil
}

func parseResponseError(resp *resty.Response) error {
	var result Result
	err := json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return fmt.Errorf("%s", resp.Body())
	}
	return fmt.Errorf("%s", result.Message)
}

type Result struct {
	Code    uint64      `json:"code"`
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Result  []GraphRule `json:"result"`
}

type GraphRule struct {
	AllocationAmount        string `json:"allocationAmount"`
	AllocationLifetime      string `json:"allocationLifetime"`
	AutoRenewal             bool   `json:"autoRenewal"`
	DecisionBasis           string `json:"decisionBasis"`
	Identifier              string `json:"identifier"`
	IdentifierType          string `json:"identifierType"`
	MaxAllocationPercentage string `json:"maxAllocationPercentage"`
	MaxSignal               string `json:"maxSignal"`
	MinAverageQueryFees     string `json:"minAverageQueryFees"`
	MinSignal               string `json:"minSignal"`
	MinStake                string `json:"minStake"`
	ParallelAllocations     string `json:"parallelAllocations"`
	RequireSupported        bool   `json:"requireSupported"`
}

func (g *ServiceImpl) getP2pPort(appID int) (int, string, error) {
	vo, err := g.applicationService.QueryApplicationById(appID)
	if err != nil {
		return 0, "", err
	}
	return vo.P2pForwardPort, vo.PeerId, nil
}
