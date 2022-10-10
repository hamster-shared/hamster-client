package app

import "hamster-client/module/chainmanager"

type ChainManager struct {
	manager chainmanager.Manager
}

func NewChainManagerApp(manager chainmanager.Manager) ChainManager {
	return ChainManager{
		manager: manager,
	}
}

func (c *ChainManager) StartQueue(appID int) error {
	return c.manager.CreateAndStartQueue(appID)
}

func (c *ChainManager) GetQueueInfo(appID int) (chainmanager.QueueInfo, error) {
	return c.manager.GetQueueInfo(appID)
}

func (c *ChainManager) RetryStartQueue(appID int, runNow bool) error {
	return c.manager.RetryStartQueue(appID, runNow)
}
