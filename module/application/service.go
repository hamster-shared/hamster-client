package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"hamster-client/config"
)

type ServiceImpl struct {
	ctx context.Context
	db  *gorm.DB
}

func NewServiceImpl(ctx context.Context, db *gorm.DB) ServiceImpl {
	return ServiceImpl{ctx, db}
}

// AddApplication add application data
func (a *ServiceImpl) AddApplication(application *AddApplicationParam) (bool, error) {
	var applyData Application
	err := a.db.Where("name=?", application.Name).First(&applyData).Error
	if err == gorm.ErrRecordNotFound {
		applyData.SelectNodeType = application.SelectNodeType
		applyData.Name = application.Name
		a.db.Create(&applyData)
		return true, nil
	}
	return false, errors.New(fmt.Sprintf("application:%s already exists", application.Name))
}

// UpdateApplication update application field
func (a *ServiceImpl) UpdateApplication(id int, name string, plugin string) (bool, error) {
	var applyData Application
	result := a.db.Model(applyData).Where("id = ?", id).Updates(Application{Name: name, SelectNodeType: plugin})
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

// DeleteApplication delete application by id
func (a *ServiceImpl) DeleteApplication(id int) (bool, error) {
	result := a.db.Debug().Where("id = ?", id).Delete(&Application{})
	fmt.Println(result.Error)
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

// QueryApplicationById query application data by id
func (a *ServiceImpl) QueryApplicationById(id int) (ApplyVo, error) {
	var data Application
	var resultData ApplyVo
	result := a.db.Where("id = ? ", id).First(&data)
	if result.Error != nil {
		return resultData, result.Error
	}
	resultData.ID = data.ID
	resultData.Name = data.Name
	resultData.SelectNodeType = data.SelectNodeType
	resultData.CreatedAt = data.CreatedAt
	resultData.LeaseTerm = data.LeaseTerm
	resultData.Status = data.Status
	resultData.P2pForwardPort = data.P2pForwardPort
	return resultData, nil
}

// ApplicationList Paging query application data
func (a *ServiceImpl) ApplicationList(page, pageSize int, name string, status int) (PageApplicationVo, error) {
	var total int64
	var list []Application
	var listVo []ListVo
	var data PageApplicationVo
	tx := a.db.Model(Application{})
	if name != "" {
		tx = tx.Where("name like ? ", "%"+name+"%")
	}
	if status != config.ALL {
		tx = tx.Where("status = ?", status)
	}
	result := tx.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Count(&total)
	if result.Error != nil {
		return data, result.Error
	}
	if len(list) > 0 {
		copier.Copy(&listVo, &list)
	}
	data.Items = listVo
	data.Total = total
	return data, nil
}

// UpdateApplicationStatus update deploy status
func (a *ServiceImpl) UpdateApplicationStatus(id, status int) error {
	var applyData Application
	result := a.db.Model(applyData).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (a *ServiceImpl) UpdateApplicationP2pForwardPort(id, port int) error {
	var applyData Application
	result := a.db.Model(applyData).Where("id = ?", id).Update("p2p_forward_port", port)
	return result.Error
}

func (a *ServiceImpl) QueryNextP2pPort() int {
	var data Application
	tx := a.db.Model(Application{})
	result := tx.Order("p2p_forward_port desc").Limit(1).First(&data)
	if result.Error != nil {
		return 34000
	}
	if data.P2pForwardPort < 34000 {
		return 34000
	}
	return data.P2pForwardPort + 1
}

func (a *ServiceImpl) QueryCliP2pPort(id int) (int, error) {
	var data Application
	result := a.db.Where("id = ? ", id).First(&data)
	if result.Error != nil {
		return data.CliForwardPort, result.Error
	}
	return data.CliForwardPort, nil
}

func (a *ServiceImpl) UpdateApplicationCliForwardPort(id, port int) error {
	var applyData Application
	result := a.db.Model(applyData).Where("id = ?", id).Update("cli_forward_port", port)
	return result.Error
}

func (a *ServiceImpl) QueryNextCliP2pPort() int {
	var data Application
	tx := a.db.Model(Application{})
	result := tx.Order("cli_forward_port desc").Limit(1).First(&data)
	if result.Error != nil {
		return 44000
	}
	if data.P2pForwardPort < 44000 {
		return 44000
	}
	return data.P2pForwardPort + 1
}
