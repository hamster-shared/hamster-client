package deploy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/gorm"
	"hamster-client/config"
	"hamster-client/module/application"
	"hamster-client/module/keystorage"
	"hamster-client/module/p2p"
	"hamster-client/module/wallet"
	"hamster-client/utils"
	"net/url"
	"strconv"
	"time"
)

type ServiceImpl struct {
	ctx                context.Context
	httpUtil           *utils.HttpUtil
	db                 *gorm.DB
	keyStorageService  keystorage.Client
	p2pService         p2p.Service
	walletService      wallet.Service
	applicationService application.Service
}

func NewServiceImpl(ctx context.Context, httpUtil *utils.HttpUtil, db *gorm.DB, keyStorageService *keystorage.Client, p2pService p2p.Service, walletService wallet.Service, applicationService application.Service) ServiceImpl {
	return ServiceImpl{ctx, httpUtil, db, *keyStorageService, p2pService, walletService, applicationService}
}

func (s *ServiceImpl) DeployTheGraph(id int, jsonData string) (bool, error) {
	//Judge whether the account has peerId
	info, err := s.applicationService.QueryApplicationById(id)
	if err != nil {
		return false, err
	}
	if info.PeerId == "" {
		//Modify the status of the application to wait for resources
		result := s.db.Model(application.Application{}).Where("id = ?", id).Update("status", application.Deploying).Error
		if result != nil {
			return false, result
		}
		return true, nil
	}
	err = s.setupP2p(info.PeerId)
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

	var data application.Application
	queryResult := s.db.Where("id = ? ", id).First(&data)
	if queryResult.Error != nil {
		return false, queryResult.Error
	}
	providerUrl := fmt.Sprintf(config.Httpprovider, data.P2pForwardPort)

	err = s.deployApi(sendData, providerUrl)
	if err != nil {
		return false, err
	}
	//Modification status is in deployment
	result := s.db.Model(application.Application{}).Where("id = ?", id).Update("status", application.Deploying).Error
	if result != nil {
		return false, result
	}
	go s.queryDeployStatus(id)
	return true, nil
}

func (s *ServiceImpl) deployApi(sendData DeployParams, url string) error {
	runtime.LogInfo(s.ctx, "start Deploy the graph")
	pair, err := s.walletService.GetWalletKeypair()
	if err != nil {
		return err
	}
	res, err := s.httpUtil.NewRequest().SetHeader("SS58AuthData", utils.GetSS58AuthDataWithKeyringPair(pair)).SetBody(sendData).Post(url)
	if err != nil {
		runtime.LogError(s.ctx, "DeployTheGraph http error:"+err.Error())
		return err
	}
	if !res.IsSuccess() {
		runtime.LogError(s.ctx, "DeployTheGraph Response error: "+res.Status())
		return errors.New(res.Status())
	}
	return nil
}

func (s *ServiceImpl) DeployGraph(id int, sendData DeployParams) (bool, error) {
	var data application.Application
	queryResult := s.db.Where("id = ? ", id).First(&data)
	if queryResult.Error != nil {
		return false, queryResult.Error
	}
	providerUrl := fmt.Sprintf(config.Httpprovider, data.P2pForwardPort)
	fmt.Println("start Deploy the graph:")
	err := s.deployApi(sendData, providerUrl)
	if err != nil {
		fmt.Println("DeployTheGraph http error:", err.Error())
		return false, err
	}
	//Modification status is in deployment
	result := s.db.Model(application.Application{}).Where("id = ?", id).Update("status", application.Deploying).Error
	if result != nil {
		return false, result
	}
	err = s.queryDeployStatus(id)
	if err != nil {
		return false, err
	}
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

	var application application.Application
	queryResult := g.db.Where("id = ? ", id).First(&application)
	if queryResult.Error != nil {
		return config.RequestStatusFailed, queryResult.Error
	}
	providerUrl := fmt.Sprintf(config.HttpGraphStatus, application.P2pForwardPort)
	data, err := g.graphStatusApi(providerUrl, serviceName...)
	if err != nil {
		return config.RequestStatusFailed, err
	}
	return data.Result, nil
}

func (g *ServiceImpl) graphStatusApi(providerUrl string, serviceName ...string) (*Result, error) {
	// SetHeader("SS58AuthData", utils.GetSS58AuthDataWithKeyringPair(pair))).
	pair, err := g.walletService.GetWalletKeypair()
	if err != nil {
		return nil, err
	}
	var data Result
	res, err := g.httpUtil.NewRequest().
		SetQueryParamsFromValues(url.Values{"serviceName": serviceName}).
		SetHeader("SS58AuthData", utils.GetSS58AuthDataWithKeyringPair(pair)).
		SetResult(&data).
		Get(providerUrl)
	if err != nil {
		runtime.LogError(g.ctx, "DeployTheGraph http error:"+err.Error())
		return nil, err
	}
	if !res.IsSuccess() {
		runtime.LogError(g.ctx, "DeployTheGraph Response error: "+res.Status())
		return nil, errors.New(fmt.Sprintf("Query status request failed. The request status is:%s", res.Status()))
	}
	return &data, nil
}

// query deploy graph status
func (s *ServiceImpl) queryDeployStatus(id int) error {
	containerIds := []string{"graph-node", "postgres", "index-service", "index-agent", "index-cli"}
	numbers := 0
	for {
		time.Sleep(time.Duration(40) * time.Second)
		res, err := s.QueryGraphStatus(id, containerIds...)
		if err != nil {
			fmt.Println("docker status :", res)
			if numbers >= 4 {
				return err
			}
			numbers = numbers + 1
			continue
		}
		fmt.Println("docker status :", res)
		if res == 1 {
			result := s.db.Model(application.Application{}).Where("id = ?", id).Update("status", application.Running).Error
			if result != nil {
				fmt.Println("save  status :", result.Error())
				if numbers >= 4 {
					return result
				}
				numbers = numbers + 1
				continue
			}
			return nil
		} else if res == config.RequestStatusFailed {
			numbers = numbers + 1
			continue
		} else {
			if numbers >= 4 {
				s.db.Model(application.Application{}).Where("id = ?", id).Update("status", application.DeploymentFailed)
				return errors.New("deploy failed")
			}
		}
		numbers = numbers + 1
	}
}

func (s *ServiceImpl) closeP2p() {
	data := s.p2pService.GetProviderLinks()
	res := *data
	if len(res) > 0 {
		for _, value := range res {
			s.p2pService.Close(value.TargetAddress)
		}
	}
}

func (s *ServiceImpl) setupP2p(peerId string) error {
	_, resultErr := s.p2pService.GetSetting()
	if resultErr != nil {
		err := s.p2pService.InitSetting()
		if err != nil {
			return err
		}
	}
	//close p2p link
	s.closeP2p()
	fmt.Println("p2p start")
	fmt.Println(peerId)
	proErr := s.p2pService.ProLink(peerId)
	if proErr != nil {
		runtime.LogError(s.ctx, "provider link error:"+proErr.Error())
		return proErr
	}
	fmt.Println("p2p end")
	return nil
}
