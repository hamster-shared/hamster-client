package chain

import (
	"gorm.io/gorm"
	"time"
)

type Ethereum struct {
}

type EthereumDeployParam struct {
	ID            uint      `json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Network       string    `json:"network"`       // rinkbey network or mainnet network
	LeaseTerm     int       `json:"leaseTerm"`     //
	ApplicationID uint      `json:"applicationId"` // application id
}

func NewEthereum() Chain {
	return &Ethereum{}
}

func (e *Ethereum) GetDeployParam(appID int, db *gorm.DB) any {
	var deployData EthereumDeployParam
	err := db.Table("ethereum_deploy_params").
		Where("application_id = ?", appID).
		First(&deployData).
		Error
	if err != nil {
		return nil
	} else {
		return deployData
	}
}
