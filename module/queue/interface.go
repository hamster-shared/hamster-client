package queue

type Job interface {
	InitStatus()
	Run(si chan StatusInfo) (StatusInfo, error)
	Status() StatusInfo
}

type Queue interface {
	ID() int
	GetStatus() (info []StatusInfo, err error)
	Start(done chan struct{})
	Stop() error
	Reset()
	saveStatus() error
	loadStatus() error
	InitStatus()
	SetJobStatus(jobName string, statusInfo StatusInfo)
}
