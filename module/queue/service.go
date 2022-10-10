package queue

type Client interface {
	GetStatusInfo(id int) ([]StatusInfo, error)
	StopQueue(id int) error
}

type client struct{}

func NewServiceImpl() Client {
	return &client{}
}

func (c *client) GetStatusInfo(id int) ([]StatusInfo, error) {
	q, err := GetQueue(id)
	if err != nil {
		return nil, err
	}
	return q.(Queue).GetStatus()
}

func (c *client) StopQueue(id int) error {
	q, err := GetQueue(id)
	if err != nil {
		return nil
	}

	err = q.(Queue).Stop()
	if err != nil {
		return err
	}
	queues.Delete(id)
	return nil
}
