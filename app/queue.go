package app

import "hamster-client/module/queue"

type Queue struct {
	service queue.Client
}

func NewQueueApp(service queue.Client) Queue {
	return Queue{
		service: service,
	}
}

type QueueInfo struct {
	Info []queue.StatusInfo `json:"info"`
}

func (q *Queue) GetQueueInfo(id int) (QueueInfo, error) {
	info, err := q.service.GetStatusInfo(id)
	if err != nil {
		return QueueInfo{}, err
	}
	return QueueInfo{Info: info}, nil
}
