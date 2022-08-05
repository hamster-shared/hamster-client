package app

import (
	"context"
	"hamster-client/module/application"
	param "hamster-client/module/graph/v2"
)

type Application struct {
	ctx                     context.Context
	applicationService      application.Service
	graphDeployParamService param.Service
}

func NewApplicationApp(service application.Service, graphDeployParamService param.Service) Application {
	return Application{
		applicationService:      service,
		graphDeployParamService: graphDeployParamService,
	}
}

func (a *Application) WailsInit(ctx context.Context) error {
	a.ctx = ctx
	return nil
}

// AddApplication add application
func (a *Application) AddApplication(applicationData param.AddParam) (bool, error) {
	return a.graphDeployParamService.SaveGraphDeployParameterAndApply(applicationData)
}

// UpdateApplication edit application
func (a *Application) UpdateApplication(application application.UpdateApplicationParam) (bool, error) {
	return a.applicationService.UpdateApplication(int(application.ID), application.Name, application.Plugin)
}

func (a *Application) DeleteApplication(id int) (bool, error) {
	return a.graphDeployParamService.DeleteGraphDeployParameterAndApply(id)
}

// ApplicationList Paging query application list
func (a *Application) ApplicationList(page, pageSize int, name string, status int) (application.PageApplicationVo, error) {
	return a.applicationService.ApplicationList(page, pageSize, name, status)
}

// QueryApplicationById query application by applicationId
func (a *Application) QueryApplicationById(id int) (application.ApplyVo, error) {
	return a.applicationService.QueryApplicationById(id)
}
