package chain

import (
	"hamster-client/module/application"
	"hamster-client/module/queue"
	"time"

	"gorm.io/gorm"
)

type Graph struct {
	common     *Common
	DeployType int
}

type GraphDeployParam struct {
	GraphDeployParameter
}

type GraphDeployParameter struct {
	ID              uint                    `json:"id"`
	CreatedAt       time.Time               `json:"createdAt"`
	UpdatedAt       time.Time               `json:"updatedAt"`
	Network         string                  `json:"network"`         // rinkbey network or mainnet network
	LeaseTerm       int                     `json:"leaseTerm"`       // indexer-service、indexer-agent eth-url
	ThegraphIndexer string                  `json:"thegraphIndexer"` // mnemonic
	StakingAmount   int                     `json:"stakingAmount"`
	Application     application.Application `json:"application"`   // application entity
	ApplicationID   uint                    `json:"applicationId"` // application id
}

type QueueInfo struct {
	Info []queue.StatusInfo `json:"info"`
}

func (g *Graph) DeployJob(appData application.Application) error {
	return g.common.deployJob(appData, g.DeployType)
}

func (g *Graph) CreateQueue(appData application.Application) (queue.Queue, error) {
	return g.common.createQueue(appData, g.DeployType)
}

func (g *Graph) SaveDeployParam(
	appInfo application.Application,
	deployParam DeployParam,
	db *gorm.DB,
) error {
	panic("not implemented") // TODO: Implement
}

func (g *Graph) SaveJsonParam(id string, deployParam DeployParam) error {
	panic("not implemented") // TODO: Implement
}

func (g *Graph) GetDeployParam(appID int, db *gorm.DB) any {
	var deployData GraphDeployParam
	err := db.Table("graph_deploy_params").
		Where("application_id = ?", appID).
		First(&deployData).
		Error
	if err != nil {
		return nil
	} else {
		return deployData
	}
}
