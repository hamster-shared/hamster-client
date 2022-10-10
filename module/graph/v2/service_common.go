package v2

//
//import (
//	"encoding/json"
//	"fmt"
//	"gorm.io/gorm"
//	"hamster-client/config"
//	"hamster-client/module/account"
//	"hamster-client/module/application"
//	"hamster-client/module/chainmanager"
//	"hamster-client/module/common"
//	"hamster-client/module/deploy"
//	queue2 "hamster-client/module/queue"
//)
//
//type CommonDeploySaveServiceImpl struct {
//	ServiceImpl
//	DeployType int
//}
//
//func (s *CommonDeploySaveServiceImpl) saveDeployParam(appData application.Application, paramData interface{}, tx *gorm.DB) error {
//	var deployData common.EthereumDeployParam
//	deployData.LeaseTerm = appData.LeaseTerm
//	addData := paramData.(AddParam)
//	deployData.Network = addData.SelectNodeType
//	deployData.ApplicationID = appData.ID
//	if err := tx.Create(&deployData).Error; err != nil {
//		return err
//	}
//	return nil
//}
//func (g *CommonDeploySaveServiceImpl) deployJob(addData application.Application) {
//	applicationId := int(addData.ID)
//	var accountInfo account.Account
//	accountInfo, err := g.accountService.GetAccount()
//	if err != nil {
//		accountInfo.WsUrl = config.DefaultPolkadotNode
//	}
//	substrateApi, err := g.accountService.GetSubstrateApi()
//	if err != nil {
//		return
//	}
//	waitResourceJob, _ := NewWaitResourceJob(substrateApi, g.accountService, g.applicationService, g.p2pService, applicationId, g.walletService, g.DeployType)
//
//	var tools chainmanager.Manager
//	tools = chainmanager.NewServiceImpl(g.db, g.applicationService, g.p2pService)
//
//	pullJob := chainmanager.NewPullImageJob(applicationId, tools)
//
//	deployJob := chainmanager.NewStartJob(applicationId, tools)
//
//	queue, err := queue2.NewQueue(applicationId, g.db, waitResourceJob, pullJob, deployJob)
//	if err != nil {
//		fmt.Println("new queue failed,err is: ", err)
//	}
//	channel := make(chan struct{})
//	go queue.Start(channel)
//	<-channel
//}
//
//func (g *CommonDeploySaveServiceImpl) saveJsonParam(id string, paramData interface{}) error {
//	addData := paramData.(AddParam)
//	var deploymentData deploy.DeployParameter
//	pluginDeployInfo := config.PluginMap[addData.SelectNodeType]
//	deploymentData.Initialization.AccountMnemonic = addData.ThegraphIndexer
//	deploymentData.Initialization.LeaseTerm = addData.LeaseTerm
//	deploymentData.Staking.PledgeAmount = addData.StakingAmount
//	deploymentData.Deployment.NodeEthereumUrl = pluginDeployInfo.EthNetwork
//	deploymentData.Deployment.EthereumUrl = pluginDeployInfo.EndpointUrl
//	deploymentData.Deployment.EthereumNetwork = pluginDeployInfo.EthereumNetworkName
//	deploymentData.Staking.NetworkUrl = pluginDeployInfo.EndpointUrl
//	deploymentData.Staking.Address = pluginDeployInfo.TheGraphStakingAddress
//	jsonData, err := json.Marshal(deploymentData)
//	if err != nil {
//		return err
//	}
//	g.keyStorageService.Set(string(application.TYPE_Thegraph)+"_"+id, string(jsonData))
//	return nil
//}
