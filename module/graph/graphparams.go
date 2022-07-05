package graph

import (
	"gorm.io/gorm"
	"hamster-client/module/application"
	"time"
)

type GraphParameter struct {
	gorm.Model
	NodeEthereumUrl string                  `json:"nodeEthereumUrl"` //graph-node eth-url
	EthereumUrl     string                  `json:"ethereumUrl"`     //indexer-service、indexer-agent eth-url
	EthereumNetwork string                  `json:"ethereumNetwork"` //eth network
	IndexerAddress  string                  `json:"indexerAddress"`  //indexer address
	Mnemonic        string                  `json:"mnemonic"`        // mnemonic
	Application     application.Application `json:"application"`     //application entity
	ApplicationId   uint                    `json:"applicationId"`   //application id
}

type GraphParameterVo struct {
	NodeEthereumUrl string    `json:"nodeEthereumUrl"` //graph-node eth-url
	EthereumUrl     string    `json:"ethereumUrl"`     //indexer-service、indexer-agent eth-url
	EthereumNetwork string    `json:"ethereumNetwork"` //eth network
	IndexerAddress  string    `json:"indexerAddress"`  //indexer address
	Mnemonic        string    `json:"mnemonic"`        // mnemonic
	ApplicationId   uint      `json:"applicationId"`   //application id
	Name            string    `json:"name"`            //apply name
	Describe        string    `json:"describe"`        //apply describe
	Status          int       `json:"status"`          //apply status 0: not deploy 1:deployed 2:ALL
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type Service interface {
	SaveGraphParameter(data GraphParameter) (bool, error)
	QueryParamByApplyId(applicationId int) (GraphParameterVo, error)
	DeleteGraphAndParams(applicationId int) (bool, error)
	QueryGraphStatus(serviceName string) (int, error)
}
