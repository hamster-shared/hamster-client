package graph

import (
	"context"
	"errors"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/gorm"
	"hamster-client/config"
	"hamster-client/module/application"
	"hamster-client/utils"
)

type ServiceImpl struct {
	ctx      context.Context
	db       *gorm.DB
	httpUtil *utils.HttpUtil
}

func NewServiceImpl(ctx context.Context, db *gorm.DB, httpUtil *utils.HttpUtil) ServiceImpl {
	return ServiceImpl{ctx, db, httpUtil}
}

func (g *ServiceImpl) SaveGraphParameter(data GraphParameter) error {
	var graphParams GraphParameter
	err := g.db.Preload("Application").Where("application_id = ? ", data.ApplicationId).First(&graphParams).Error
	if err != gorm.ErrRecordNotFound {
		g.db.Create(&data)
		return nil
	}
	return errors.New(fmt.Sprintf("graph param -> application :%s already exists", data.Application.Name))
}

func (g *ServiceImpl) QueryParamByApplyId(applicationId int) (GraphParameter, error) {
	var data GraphParameter
	err := g.db.Preload("Application").Where("application_id = ? ", applicationId).First(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

func (g *ServiceImpl) DeleteGraphAndParams(applicationId int) error {
	err := g.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Preload("Application").Where("application_id = ? ", applicationId).First(&GraphParameter{}).Error; err != nil {
			return err
		}
		if err := tx.Preload("Application").Where("application_id = ? ", applicationId).Delete(&GraphParameter{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ? ", applicationId).First(&application.Application{}).Error; err != nil {
			return err
		}
		if err := tx.Debug().Where("id = ?", applicationId).Delete(&application.Application{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (g *ServiceImpl) QueryGraphStatus(serviceName string) (int, error) {
	var status int
	res, err := g.httpUtil.NewRequest().
		SetQueryParam("serviceName", serviceName).
		SetResult(&status).
		Get(config.HttpGraphStatus)
	if err != nil {
		runtime.LogError(g.ctx, "DeployTheGraph http error:"+err.Error())
		return 3, err
	}
	if !res.IsSuccess() {
		runtime.LogError(g.ctx, "DeployTheGraph Response error: "+res.Status())
		return 3, errors.New(fmt.Sprintf("Query status request failed. The request status is:%s", res.Status()))
	}
	return status, nil
}