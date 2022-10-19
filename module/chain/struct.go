package chain

type DeployParam struct {
	Name            string `json:"name"` //apply name
	ServiceType     string `json:"serviceType"`
	SelectNodeType  string `json:"selectNodeType"` //apply plugin
	LeaseTerm       int    `json:"leaseTerm"`
	ThegraphIndexer string `json:"thegraphIndexer"` // mnemonic
	StakingAmount   int    `json:"stakingAmount"`
	Blockchain      string `json:"blockchain"`
}

type DeployResult struct {
	ID     uint `json:"id"`
	Result bool `json:"result"`
}

//type base struct {
//	ctx     context.Context
//	db      *gorm.DB
//	ks      keystorage.Service
//	account account.Service
//	app     application.Service
//	p2p     p2p.Service
//	wallet  wallet.Service
//	queue   queue.Service
//}
