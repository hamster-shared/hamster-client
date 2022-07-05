package app

import (
	"context"
	"hamster-client/module/application"
	"hamster-client/module/graph"
)

type Application struct {
	ctx                context.Context
	applicationService application.Service
	graphService       graph.Service
}

func NewApplicationApp(service application.Service, graphService graph.Service) Application {
	return Application{
		applicationService: service,
		graphService:       graphService,
	}
}

func (a *Application) WailsInit(ctx context.Context) error {
	a.ctx = ctx
	return nil
}

// AddApplication add application
func (a *Application) AddApplication(application application.AddApplicationParam) (bool, error) {
	return a.applicationService.AddApplication(&application)
}

// UpdateApplication edit application
func (a *Application) UpdateApplication(application application.UpdateApplicationParam) (bool, error) {
	return a.applicationService.UpdateApplication(int(application.ID), application.Name, application.Describe)
}

func (a *Application) DeleteApplication(id int) (bool, error) {
	return a.applicationService.DeleteApplication(id)
}

// ApplicationList Paging query application list
func (a *Application) ApplicationList(page, pageSize int, name string, status int) (application.PageApplicationVo, error) {
	return a.applicationService.ApplicationList(page, pageSize, name, status)
}

// QueryApplicationById query application by applicationId
func (a *Application) QueryApplicationById(id int) (application.ApplyVo, error) {
	return a.applicationService.QueryApplicationById(id)
}

// DeleteGraphAndParams delete application
func (a *Application) DeleteGraphAndParams(applicationId int) (bool, error) {
	return a.graphService.DeleteGraphAndParams(applicationId)
}

// QueryGraphStatus query graph status
func (a *Application) QueryGraphStatus(serviceName string) (int, error) {
	return a.graphService.QueryGraphStatus(serviceName)
}
