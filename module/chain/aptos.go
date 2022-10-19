package chain

import (
	"hamster-client/module/application"
	"hamster-client/module/queue"

	"gorm.io/gorm"
)

type Aptos struct {
	common     *Common
	DeployType int
}

type AptosDeployParam struct{}

func NewAptos(common *Common, deployType int) Chain {
	return &Aptos{
		common:     common,
		DeployType: deployType,
	}
}

func (a *Aptos) DeployJob(appData application.Application) error {
	return a.common.deployJob(appData, a.DeployType)
}

func (a *Aptos) CreateQueue(appData application.Application) (queue.Queue, error) {
	return a.common.createQueue(appData, a.DeployType)
}

func (a *Aptos) SaveDeployParam(
	appInfo application.Application,
	deployParam DeployParam,
	db *gorm.DB,
) error {
	panic("not implemented") // TODO: Implement
}

func (a *Aptos) SaveJsonParam(id string, deployParam DeployParam) error {
	panic("not implemented") // TODO: Implement
}

func (a *Aptos) GetDeployParam(appID int, db *gorm.DB) any {
	var deployData AptosDeployParam
	err := db.Table("aptos_deploy_params").
		Where("application_id = ?", appID).
		First(&deployData).
		Error
	if err != nil {
		return nil
	} else {
		return deployData
	}
}
