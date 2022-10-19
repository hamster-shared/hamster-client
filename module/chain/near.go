package chain

import (
	"hamster-client/module/application"
	"hamster-client/module/queue"

	"gorm.io/gorm"
)

type Near struct {
	common     *Common
	DeployType int
}

type NearDeployParam struct{}

func NewNear(common *Common, deployType int) *Near {
	return &Near{
		common:     common,
		DeployType: deployType,
	}
}

func (n *Near) DeployJob(appData application.Application) error {
	return n.common.deployJob(appData, n.DeployType)
}

func (n *Near) CreateQueue(appData application.Application) (queue.Queue, error) {
	return n.common.createQueue(appData, n.DeployType)
}

func (n *Near) SaveDeployParam(
	appInfo application.Application,
	deployParam DeployParam,
	db *gorm.DB,
) error {
	panic("not implemented") // TODO: Implement
}

func (n *Near) SaveJsonParam(id string, deployParam DeployParam) error {
	panic("not implemented") // TODO: Implement
}

func (n *Near) GetDeployParam(appID int, db *gorm.DB) any {
	var deployData NearDeployParam
	err := db.Table("near_deploy_params").
		Where("application_id = ?", appID).
		First(&deployData).
		Error
	if err != nil {
		return nil
	} else {
		return deployData
	}
}
