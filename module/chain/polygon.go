package chain

import (
	"hamster-client/module/application"
	"hamster-client/module/queue"

	"gorm.io/gorm"
)

type Polygon struct {
	common     *Common
	DeployType int
}

type PolygonDeployParam struct{}

func NewPolygon(common *Common, deployType int) *Polygon {
	return &Polygon{
		common:     common,
		DeployType: deployType,
	}
}

func (p *Polygon) DeployJob(appData application.Application) error {
	return p.common.deployJob(appData, p.DeployType)
}

func (p *Polygon) CreateQueue(appData application.Application) (queue.Queue, error) {
	return p.common.createQueue(appData, p.DeployType)
}

func (p *Polygon) SaveDeployParam(
	appInfo application.Application,
	deployParam DeployParam,
	db *gorm.DB,
) error {
	panic("not implemented") // TODO: Implement
}

func (p *Polygon) SaveJsonParam(id string, deployParam DeployParam) error {
	panic("not implemented") // TODO: Implement
}

func (p *Polygon) GetDeployParam(appID int, db *gorm.DB) any {
	var deployData PolygonDeployParam
	err := db.Table("polygon_deploy_params").
		Where("application_id = ?", appID).
		First(&deployData).
		Error
	if err != nil {
		return nil
	} else {
		return deployData
	}
}
