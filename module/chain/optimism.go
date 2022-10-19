package chain

import (
	"hamster-client/module/application"
	"hamster-client/module/queue"

	"gorm.io/gorm"
)

type Optimism struct {
	common     *Common
	DeployType int
}

type OptimismDeployParam struct{}

func NewOptimism(common *Common, deployType int) *Optimism {
	return &Optimism{
		common:     common,
		DeployType: deployType,
	}
}

func (o *Optimism) DeployJob(appData application.Application) error {
	return o.common.deployJob(appData, o.DeployType)
}

func (o *Optimism) CreateQueue(appData application.Application) (queue.Queue, error) {
	return o.common.createQueue(appData, o.DeployType)
}

func (o *Optimism) SaveDeployParam(
	appInfo application.Application,
	deployParam DeployParam,
	db *gorm.DB,
) error {
	panic("not implemented") // TODO: Implement
}

func (o *Optimism) SaveJsonParam(id string, deployParam DeployParam) error {
	panic("not implemented") // TODO: Implement
}

func (o *Optimism) GetDeployParam(appID int, db *gorm.DB) any {
	var deployData OptimismDeployParam
	err := db.Table("optimism_deploy_params").
		Where("application_id = ?", appID).
		First(&deployData).
		Error
	if err != nil {
		return nil
	} else {
		return deployData
	}
}
