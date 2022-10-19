package chain

import (
	"hamster-client/module/application"
	"hamster-client/module/queue"

	"gorm.io/gorm"
)

type Sui struct {
	common     *Common
	DeployType int
}

type SuiDeployParam struct{}

func NewSui(common *Common, deployType int) *Sui {
	return &Sui{
		common:     common,
		DeployType: deployType,
	}
}

func (s *Sui) DeployJob(appData application.Application) error {
	return s.common.deployJob(appData, s.DeployType)
}

func (s *Sui) CreateQueue(appData application.Application) (queue.Queue, error) {
	return s.common.createQueue(appData, s.DeployType)
}

func (s *Sui) SaveDeployParam(
	appInfo application.Application,
	deployParam DeployParam,
	db *gorm.DB,
) error {
	panic("not implemented") // TODO: Implement
}

func (s *Sui) SaveJsonParam(id string, deployParam DeployParam) error {
	panic("not implemented") // TODO: Implement
}

func (s *Sui) GetDeployParam(appID int, db *gorm.DB) any {
	var deployData SuiDeployParam
	err := db.Table("sui_deploy_params").
		Where("application_id = ?", appID).
		First(&deployData).
		Error
	if err != nil {
		return nil
	} else {
		return deployData
	}
}
