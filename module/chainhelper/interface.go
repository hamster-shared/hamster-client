package chainhelper

import (
	"hamster-client/module/account"
	"hamster-client/module/application"
	"hamster-client/module/chain"
	"hamster-client/module/p2p"
	"hamster-client/module/wallet"

	"gorm.io/gorm"
)

type Helper interface {
	GetChain(deployType int) (chain.Chain, error)
	DeployType(appID int) (int, error)
	GetP2pForwardPort(appID int) (int, error)

	DB() *gorm.DB
	App() application.Service
	Account() account.Service
	P2p() p2p.Service
	Wallet() wallet.Service
}
