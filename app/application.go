package app

import (
	"context"
	"fmt"
	"hamster-client/config"
	"hamster-client/module/application"
	"hamster-client/module/deploy"
	param "hamster-client/module/graph/v2"
	"hamster-client/module/p2p"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Application struct {
	ctx                context.Context
	applicationService application.Service
	p2pService         p2p.Service
	deployService      deploy.Service
}

func NewApplicationApp(
	service application.Service,
	p2pService p2p.Service,
	deployService deploy.Service,
) Application {
	return Application{
		applicationService: service,
		p2pService:         p2pService,
		deployService:      deployService,
	}
}

func (a *Application) WailsInit(ctx context.Context) error {
	a.ctx = ctx
	return nil
}

// AddApplication add application
func (a *Application) AddApplication(
	applicationData param.AddParam,
) (param.AddApplicationVo, error) {
	// TODO: add application
	return param.AddApplicationVo{}, nil
}

// UpdateApplication edit application
func (a *Application) UpdateApplication(
	application application.UpdateApplicationParam,
) (bool, error) {
	return a.applicationService.UpdateApplication(
		int(application.ID),
		application.Name,
		application.SelectNodeType,
	)
}

func (a *Application) DeleteApplication(id int) (bool, error) {
	fmt.Println("DeleteApplication: ", id)
	// TODO: delete application
	return true, nil
}

// ApplicationList Paging query application list
func (a *Application) ApplicationList(
	page, pageSize int,
	name string,
	status int,
) (application.PageApplicationVo, error) {
	return a.applicationService.ApplicationList(page, pageSize, name, status)
}

// QueryApplicationById query application by applicationId
func (a *Application) QueryApplicationById(id int) (application.ApplyVo, error) {
	vo, err := a.applicationService.QueryApplicationById(id)

	if vo.PeerId == "" {
		return vo, err
	}

	if vo.Status == application.Running || vo.Status == application.Offline {
		_ = a.p2pService.LinkByProtocol(config.ProviderProtocol, vo.P2pForwardPort, vo.PeerId)
		containerIds := []string{
			"graph-node",
			"postgres",
			"index-client",
			"index-agent",
			"index-cli",
		}
		status, err := a.deployService.QueryGraphStatus(int(vo.ID), containerIds...)
		fmt.Println("status:", status, "error: ", err)
		if err != nil || status != 1 {
			_ = a.applicationService.UpdateApplicationStatus(int(vo.ID), application.Offline)
		} else if status == 1 {
			_ = a.applicationService.UpdateApplicationStatus(int(vo.ID), application.Running)
		}
	}

	return vo, err
}

func (a *Application) ReconnectionProLink(applicationId int) (bool, error) {
	applicationInfo, err := a.applicationService.QueryApplicationById(applicationId)
	if err != nil {
		runtime.LogError(a.ctx, "Get application error")
		return false, err
	}
	if applicationInfo.PeerId != "" {
		err := a.p2pService.ProLink(applicationInfo.PeerId)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (a *Application) UpdateApplicationIncome(id int, income float64) (bool, error) {
	return a.applicationService.UpdateApplicationIncome(id, income)
}

func (a *Application) UpdateThinkingTime(id, time int) (bool, error) {
	return a.applicationService.UpdateThinkingTime(id, time)
}
