package chainmanager

import (
	"hamster-client/module/queue"
)

type QueueInfo struct {
	Info []queue.StatusInfo `json:"info"`
}

type Manager interface {
	CreateAndStartQueue(appID int) error
	RetryStartQueue(appID int, runNow bool) error
	GetQueueInfo(appID int) (QueueInfo, error)
}
