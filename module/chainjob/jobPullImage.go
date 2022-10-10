package chainjob

import (
	"encoding/json"
	"fmt"
	"hamster-client/module/chainhelper"
	"hamster-client/module/queue"
	"hamster-client/utils"
	"time"
)

// PullImageJob pull image job
type PullImageJob struct {
	appID  int
	si     queue.StatusInfo
	helper chainhelper.Helper
}

func NewPullImageJob(appID int, helper chainhelper.Helper) queue.Job {
	return &PullImageJob{
		appID:  appID,
		helper: helper,
	}
}

func (j *PullImageJob) InitStatus() {
	j.si.Name = "Pull Image"
	j.si.Status = queue.None
}

func (j *PullImageJob) Run(si chan queue.StatusInfo) (queue.StatusInfo, error) {
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
	url := fmt.Sprintf("http://localhost:%d/api/v1/chains/pullImage", p2pForwardPort)

	for i := 0; i < 3; i++ {

		deployType, err := j.helper.DeployType(j.appID)
		if err != nil {
			fmt.Printf("get deploy type error: %s", err.Error())
		}

		data, err := j.helper.GetChain(deployType)
		if err != nil {
			fmt.Printf("get deploy type error: %s", err.Error())
		}

		jsonStr, err := json.Marshal(data)
		if err != nil {
			continue
		}
		fmt.Println("param: ", string(jsonStr))

		req := utils.NewHttp().NewRequest()
		req.SetBody(data)
		response, err := req.Post(url)
		if err != nil {
			j.si.Error = err.Error()
			fmt.Println(string(response.Body()))
			continue
		}
		if response.IsSuccess() {
			j.si.Status = queue.Succeeded
			j.si.Error = ""
			si <- j.si
			fmt.Println(
				"========================== chain pull image success ==========================",
			)
			return j.si, nil
		} else {
			time.Sleep(3 * time.Second)
			continue
		}
	}

	j.si.Status = queue.Failed
	j.si.Error = "chain pull image failed"
	si <- j.si
	return j.si, fmt.Errorf(j.si.Error)
}

func (j *PullImageJob) Status() queue.StatusInfo {
	return j.si
}
