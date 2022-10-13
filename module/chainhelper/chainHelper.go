package chainhelper

import (
	"fmt"
	"hamster-client/config"
	"hamster-client/module/account"
	"hamster-client/module/application"
	"hamster-client/module/keystorage"
	"hamster-client/module/p2p"
	"hamster-client/module/queue"
	"hamster-client/module/wallet"

	"gorm.io/gorm"
)

type ChainHelper struct {
	db      *gorm.DB
	app     application.Service
	p2p     p2p.Service
	account account.Service
	wallet  wallet.Service
	queue   queue.Service
	ks      keystorage.Service
}

func NewHelper(db *gorm.DB, app application.Service, p2p p2p.Service, account account.Service, wallet wallet.Service, queue queue.Service, ks keystorage.Service) Helper {
	return &ChainHelper{
		db:      db,
		app:     app,
		p2p:     p2p,
		account: account,
		wallet:  wallet,
		queue:   queue,
		ks:      ks,
	}
}

func (c *ChainHelper) GetP2pForwardPort(appID int) (int, error) {
	vo, err := c.app.QueryApplicationById(appID)
	if err != nil {
		fmt.Println("query application by id failed, err: ", err)
		return 0, err
	}

	fmt.Println("pull before: reForwardLink:", vo.PeerId)
	if _, err = c.p2p.GetSetting(); err != nil {
		_ = c.p2p.InitSetting()
	}

	err = reForwardLink(c.p2p, vo.P2pForwardPort, vo.PeerId)
	if err != nil {
		fmt.Println("reForwardLink failed, err: ", err)
		return 0, err
	}
	return vo.P2pForwardPort, nil
}

func reForwardLink(p2pService p2p.Service, port int, peerID string) error {
	protocol := config.ProviderProtocol
	err := p2pService.LinkByProtocol(protocol, port, peerID)
	return err
}

func (c *ChainHelper) DB() *gorm.DB {
	return c.db
}

func (c *ChainHelper) App() application.Service {
	return c.app
}

func (c *ChainHelper) P2p() p2p.Service {
	return c.p2p
}

func (c *ChainHelper) Account() account.Service {
	return c.account
}

func (c *ChainHelper) Wallet() wallet.Service {
	return c.wallet
}

func (c *ChainHelper) Queue() queue.Service {
	return c.queue
}

func (c *ChainHelper) KS() keystorage.Service {
	return c.ks
}
