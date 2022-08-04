package v2

import (
	"hamster-client/module/application"
	"time"
)

type GraphDeployParameter struct {
	ID            uint                    `json:"id"`
	CreatedAt     time.Time               `json:"createdAt"`
	UpdatedAt     time.Time               `json:"updatedAt"`
	Network       string                  `json:"network"`   //rinkbey network or mainnet network
	LeaseTerm     int                     `json:"leaseTerm"` //indexer-service、indexer-agent eth-url
	Mnemonic      string                  `json:"mnemonic"`  // mnemonic
	PledgeAmount  int                     `json:"pledgeAmount"`
	Application   application.Application `json:"application"`   //application entity
	ApplicationID uint                    `json:"applicationId"` //application id
}

type AddParam struct {
	Name         string `json:"name"`   //apply name
	Plugin       string `json:"plugin"` //apply plugin
	LeaseTerm    int    `json:"leaseTerm"`
	Network      string `json:"network"`
	Mnemonic     string `json:"mnemonic"` // mnemonic
	PledgeAmount int    `json:"pledgeAmount"`
}

type Service interface {
	SaveGraphDeployParameterAndApply(data AddParam) (bool, error)
}
