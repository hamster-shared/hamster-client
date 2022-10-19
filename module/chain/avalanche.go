package chain

import (
	"hamster-client/module/application"
	"hamster-client/module/queue"

	"gorm.io/gorm"
)

type Avalanche struct {
	common     *Common
	DeployType int
}

type AvalancheDeployParam struct{}

func NewAvalanche(common *Common, deployType int) *Avalanche {
	return &Avalanche{
		common:     common,
		DeployType: deployType,
	}
}

func (a *Avalanche) DeployJob(appData application.Application) error {
	return a.common.deployJob(appData, a.DeployType)
}

func (a *Avalanche) CreateQueue(appData application.Application) (queue.Queue, error) {
	return a.common.createQueue(appData, a.DeployType)
}

func (a *Avalanche) SaveDeployParam(
	appInfo application.Application,
	deployParam DeployParam,
	db *gorm.DB,
) error {
	panic("not implemented") // TODO: Implement
}

func (a *Avalanche) SaveJsonParam(id string, deployParam DeployParam) error {
	panic("not implemented") // TODO: Implement
}

func (a *Avalanche) GetDeployParam(appID int, db *gorm.DB) any {
	var deployData AvalancheDeployParam
	err := db.Table("avalanche_deploy_params").
		Where("application_id = ?", appID).
		First(&deployData).
		Error
	if err != nil {
		return nil
	} else {
		return deployData
	}
}
