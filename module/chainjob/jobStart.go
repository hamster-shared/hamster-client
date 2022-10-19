package chainjob

import (
	"fmt"
	"hamster-client/module/application"
	"hamster-client/module/chainhelper"
	"hamster-client/module/queue"
	"hamster-client/utils"
	"time"

	log "github.com/sirupsen/logrus"
)

// StartJob container start job
type StartJob struct {
	appID  int
	helper chainhelper.Helper
	si     queue.StatusInfo
	Param  any
}

func NewStartJob(appID int, helper chainhelper.Helper, Param any) *StartJob {
	return &StartJob{
		appID:  appID,
		helper: helper,
		Param:  Param,
	}
}

func (j *StartJob) InitStatus() {
	j.si.Name = "Start"
	j.si.Status = queue.None
}

func (j *StartJob) Run(si chan queue.StatusInfo) (queue.StatusInfo, error) {
	j.si.Status = queue.Running
	si <- j.si

	// get p2p forward port
	p2pForwardPort, err := j.helper.GetP2pForwardPort(j.appID)
	if err != nil {
		j.si.Error = err.Error()
		j.si.Status = queue.Failed
		si <- j.si
		return j.si, err
	}

	// send pull image request
	url := fmt.Sprintf("http://localhost:%d/api/v1/chains/start", p2pForwardPort)

	for i := 0; i < 3; i++ {

		req := utils.NewHttp().NewRequest()
		req.SetBody(j.Param)
		response, err := req.Post(url)
		if err != nil {
			j.si.Error = err.Error()
			fmt.Println(string(response.Body()))
			continue
		}
		if response.IsSuccess() {
			err := j.helper.DB().
				Model(application.Application{}).
				Where("id = ?", j.appID).
				Update("status", application.Running).
				Error
			if err != nil {
				log.Errorf("update db application status error: %s", err.Error())
			}
			j.si.Status = queue.Succeeded
			j.si.Error = ""
			si <- j.si
			fmt.Println("========================== chain start success ==========================")
			return j.si, nil
		} else {
			time.Sleep(3 * time.Second)
			continue
		}
	}

	j.si.Status = queue.Failed
	j.si.Error = "chain start failed"
	si <- j.si
	return j.si, fmt.Errorf(j.si.Error)
}

func (j *StartJob) Status() queue.StatusInfo {
	return j.si
}
