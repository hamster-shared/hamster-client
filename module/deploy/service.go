package deploy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/gorm"
	"hamster-client/config"
	"hamster-client/module/account"
	"hamster-client/module/application"
	"hamster-client/module/keystorage"
	"hamster-client/module/p2p"
	"hamster-client/utils"
	"net/url"
	"strconv"
	"time"
)

type ServiceImpl struct {
	ctx               context.Context
	httpUtil          *utils.HttpUtil
	db                *gorm.DB
	keyStorageService keystorage.Service
	accountService    account.Service
	p2pServer         p2p.Service
}

func NewServiceImpl(ctx context.Context, httpUtil *utils.HttpUtil, db *gorm.DB, keyStorageService *keystorage.Service, accountService account.Service, p2pServer p2p.Service) ServiceImpl {
	return ServiceImpl{ctx, httpUtil, db, *keyStorageService, accountService, p2pServer}
}

func (s *ServiceImpl) DeployTheGraph(id int, jsonData string) (bool, error) {
	//Judge whether the account has peerId
	info, err := s.accountService.GetAccount()
	if err != nil {
		return false, err
	}
	if info.PeerId == "" {
		//Modify the status of the application to wait for resources
		result := s.db.Model(application.Application{}).Where("id = ?", id).Update("status", config.WAIT_RESOURCE).Error
		if result != nil {
			return false, result
		}
		return true, nil
	}
	err = s.setupP2p()
	//Determine whether to initialize configuration
	if err != nil {
		return false, err
	}
	var param DeployParameter
	jsonParam := s.keyStorageService.Get("graph_" + strconv.Itoa(id))
	if err := json.Unmarshal([]byte(jsonParam), &param); err != nil {
		return false, err
	}
	var sendData DeployParams
	sendData.Mnemonic = param.Initialization.AccountMnemonic
	sendData.Id = id
	sendData.EthereumUrl = param.Deployment.EthereumUrl
	sendData.IndexerAddress = param.Deployment.IndexerAddress
	sendData.NodeEthereumUrl = param.Deployment.NodeEthereumUrl
	sendData.EthereumNetwork = param.Deployment.EthereumNetwork
	runtime.LogInfo(s.ctx, "start Deploy the graph")
	res, err := s.httpUtil.NewRequest().SetBody(sendData).Post(config.Httpprovider)
	if err != nil {
		runtime.LogError(s.ctx, "DeployTheGraph http error:"+err.Error())
		return false, err
	}
	if !res.IsSuccess() {
		runtime.LogError(s.ctx, "DeployTheGraph Response error: "+res.Status())
		return false, errors.New(res.Status())
	}
	//Modification status is in deployment
	result := s.db.Model(application.Application{}).Where("id = ?", id).Update("status", config.IN_DEPLOYMENT).Error
	if result != nil {
		return false, result
	}
	go s.queryDeployStatus(id)
	return true, nil
}
func (s *ServiceImpl) DeployGraph(id int, sendData DeployParams) (bool, error) {
	var data application.Application
	queryResult := s.db.Where("id = ? ", id).First(&data)
	if queryResult.Error != nil {
		return false, queryResult.Error
	}
	providerUrl := fmt.Sprintf(config.Httpprovider, data.P2pForwardPort)
	fmt.Println("start Deploy the graph:")
	res, err := s.httpUtil.NewRequest().SetBody(sendData).Post(providerUrl)
	if err != nil {
		fmt.Println("DeployTheGraph http error:", err.Error())
		return false, err
	}
	if !res.IsSuccess() {
		fmt.Println("DeployTheGraph Response error: ", res.Status())
		return false, errors.New(res.Status())
	}
	//Modification status is in deployment
	result := s.db.Model(application.Application{}).Where("id = ?", id).Update("status", config.IN_DEPLOYMENT).Error
	if result != nil {
		return false, result
	}
	go s.queryDeployStatus(id)
	return true, nil
}
func (s *ServiceImpl) GetDeployInfo(id int) (DeployParameter, error) {
	data := s.keyStorageService.Get("graph_" + strconv.Itoa(id))
	var param DeployParameter
	if err := json.Unmarshal([]byte(data), &param); err != nil {
		return param, err
	}
	return param, nil
}

func (s *ServiceImpl) SaveDeployInfo(id int, json string) (bool, error) {
	s.keyStorageService.Set("graph_"+strconv.Itoa(id), json)
	err := s.keyStorageService.Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (g *ServiceImpl) QueryGraphStatus(id int, serviceName ...string) (int, error) {
	var data Result
	var application application.Application
	queryResult := g.db.Where("id = ? ", id).First(&application)
	if queryResult.Error != nil {
		return config.RequestStatusFailed, queryResult.Error
	}
	providerUrl := fmt.Sprintf(config.HttpGraphStatus, application.P2pForwardPort)
	res, err := g.httpUtil.NewRequest().
		SetQueryParamsFromValues(url.Values{"serviceName": serviceName}).
		SetResult(&data).
		Get(providerUrl)
	if err != nil {
		runtime.LogError(g.ctx, "DeployTheGraph http error:"+err.Error())
		return config.RequestStatusFailed, err
	}
	if !res.IsSuccess() {
		runtime.LogError(g.ctx, "DeployTheGraph Response error: "+res.Status())
		return config.RequestStatusFailed, errors.New(fmt.Sprintf("Query status request failed. The request status is:%s", res.Status()))
	}
	return data.Result, nil
}

// query deploy graph status
func (s *ServiceImpl) queryDeployStatus(id int) {
	containerIds := []string{"graph-node", "postgres", "index-service", "index-agent", "index-cli"}
	numbers := 0
	for {
		time.Sleep(time.Duration(10) * time.Second)
		res, _ := s.QueryGraphStatus(id, containerIds...)
		fmt.Println("docker status :", res)
		if res == 1 {
			result := s.db.Model(application.Application{}).Where("status = ?", config.IN_DEPLOYMENT).Update("status", config.DEPLOYED).Error
			if result == nil {
				return
			}
		} else if res == config.RequestStatusFailed {
			continue
		} else {
			if numbers >= 3 {
				s.db.Model(application.Application{}).Where("status = ?", config.IN_DEPLOYMENT).Update("status", config.DEPLOY_FAILED)
				return
			}
		}
		numbers = numbers + 1
	}
}

func (s *ServiceImpl) closeP2p() {
	data := s.p2pServer.GetProviderLinks()
	res := *data
	if len(res) > 0 {
		for _, value := range res {
			s.p2pServer.Close(value.TargetAddress)
		}
	}
}

func (s *ServiceImpl) setupP2p() error {

	info, err := s.accountService.GetAccount()
	if err != nil {
		return err
	}

	peerId := info.PeerId

	_, resultErr := s.p2pServer.GetSetting()
	if resultErr != nil {
		err := s.p2pServer.InitSetting()
		if err != nil {
			return err
		}
	}
	//close p2p link
	s.closeP2p()
	fmt.Println("p2p start")
	fmt.Println(peerId)
	proErr := s.p2pServer.ProLink(peerId)
	if proErr != nil {
		runtime.LogError(s.ctx, "provider link error:"+proErr.Error())
		return proErr
	}
	fmt.Println("p2p end")
	return nil
}
