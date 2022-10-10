package chainmanager

import (
	"fmt"
	"hamster-client/module/account"
	"hamster-client/module/application"
	"hamster-client/module/chainhelper"
	"hamster-client/module/chainjob"
	"hamster-client/module/p2p"
	"hamster-client/module/queue"
	"hamster-client/module/wallet"
	"sync"

	"gorm.io/gorm"
)

type VersionVo struct {
	Version string `json:"version"`
}

type ChainManager struct {
	db      *gorm.DB
	queues  sync.Map
	app     application.Service
	p2p     p2p.Service
	account account.Service
	wallet  wallet.Service
}

func NewManager(
	db *gorm.DB,
	app application.Service,
	p2p p2p.Service,
	account account.Service,
	wallet wallet.Service,
) Manager {
	return &ChainManager{
		db:      db,
		app:     app,
		p2p:     p2p,
		account: account,
		wallet:  wallet,
	}
}

// addQueue add queue use app id
func (c *ChainManager) addQueue(appID int, q queue.Queue) {
	c.queues.Store(appID, q)
}

// rmQueue remove queue
func (c *ChainManager) rmQueue(id int) {
	c.queues.Delete(id)
}

func (c *ChainManager) CreateAndStartQueue(appID int) error {
	helper := chainhelper.NewHelper(c.db, c.app, c.p2p, c.account, c.wallet)

	waitResourceJob := chainjob.NewWaitResourceJob(appID, helper)
	pullImageJob := chainjob.NewPullImageJob(appID, helper)
	startJob := chainjob.NewStartJob(appID, helper)

	q, err := queue.NewQueue(appID, c.db, waitResourceJob, pullImageJob, startJob)
	if err != nil {
		return err
	}
	c.addQueue(appID, q)

	// start queue, in a new goroutine
	done := make(chan struct{})
	go q.Start(done)

	// wait
	<-done

	return nil
}

func (c *ChainManager) GetQueueInfo(appID int) (QueueInfo, error) {
	q, ok := c.queues.Load(appID)
	if !ok {
		return QueueInfo{}, fmt.Errorf("queue %d not found", appID)
	}
	info, err := q.(queue.Queue).GetStatus()
	if err != nil {
		return QueueInfo{}, err
	}
	return QueueInfo{
		Info: info,
	}, nil
}

func (c *ChainManager) RetryStartQueue(appID int, runNow bool) error {
	if runNow {
		return c.CreateAndStartQueue(appID)
	}
	helper := chainhelper.NewHelper(c.db, c.app, c.p2p, c.account, c.wallet)

	waitResourceJob := chainjob.NewWaitResourceJob(appID, helper)
	pullImageJob := chainjob.NewPullImageJob(appID, helper)
	startJob := chainjob.NewStartJob(appID, helper)

	q, err := queue.NewQueue(appID, c.db, waitResourceJob, pullImageJob, startJob)
	if err != nil {
		return err
	}
	c.addQueue(appID, q)
	return nil
}
