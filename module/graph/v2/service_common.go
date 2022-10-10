package v2

import (
	"encoding/json"
	"hamster-client/config"
	"hamster-client/module/application"
	"hamster-client/module/common"
	"hamster-client/module/deploy"

	"gorm.io/gorm"
)

type CommonDeploySaveServiceImpl struct {
	ServiceImpl
	DeployType int
}

func (s *CommonDeploySaveServiceImpl) saveDeployParam(
	appData application.Application,
	paramData interface{},
	tx *gorm.DB,
) error {
	var deployData common.EthereumDeployParam
	deployData.LeaseTerm = appData.LeaseTerm
	addData := paramData.(AddParam)
	deployData.Network = addData.SelectNodeType
	deployData.ApplicationID = appData.ID
	if err := tx.Create(&deployData).Error; err != nil {
		return err
	}
	return nil
}

func (g *CommonDeploySaveServiceImpl) saveJsonParam(id string, paramData interface{}) error {
	addData := paramData.(AddParam)
	var deploymentData deploy.DeployParameter
	pluginDeployInfo := config.PluginMap[addData.SelectNodeType]
	deploymentData.Initialization.AccountMnemonic = addData.ThegraphIndexer
	deploymentData.Initialization.LeaseTerm = addData.LeaseTerm
	deploymentData.Staking.PledgeAmount = addData.StakingAmount
	deploymentData.Deployment.NodeEthereumUrl = pluginDeployInfo.EthNetwork
	deploymentData.Deployment.EthereumUrl = pluginDeployInfo.EndpointUrl
	deploymentData.Deployment.EthereumNetwork = pluginDeployInfo.EthereumNetworkName
	deploymentData.Staking.NetworkUrl = pluginDeployInfo.EndpointUrl
	deploymentData.Staking.Address = pluginDeployInfo.TheGraphStakingAddress
	jsonData, err := json.Marshal(deploymentData)
	if err != nil {
		return err
	}
	g.keyStorageService.Set(string(application.TYPE_Thegraph)+"_"+id, string(jsonData))
	return nil
}
