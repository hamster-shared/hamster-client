package application

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
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
		applyData.Describe = application.Describe
		applyData.Name = application.Name
		a.db.Create(&applyData)
		return true, nil
	}
	return false, errors.New(fmt.Sprintf("application:%s already exists", application.Name))
}

// UpdateApplication update application field
func (a *ServiceImpl) UpdateApplication(id int, name string, des string) (bool, error) {
	var applyData Application
	result := a.db.Model(applyData).Where("id = ?", id).Updates(Application{Name: name, Describe: des})
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
	resultData.Describe = data.Describe
	resultData.CreatedAt = data.CreatedAt
	resultData.UpdatedAt = data.UpdatedAt
	return resultData, nil
}

// ApplicationList Paging query application data
func (a *ServiceImpl) ApplicationList(page, pageSize int, name string, status int) (PageApplicationVo, error) {
	var total int64
	var list []Application
	var data PageApplicationVo
	tx := a.db.Model(Application{})
	if name != "" {
		tx = tx.Where("name like ? ", "%"+name+"%")
	}
	if status != 2 {
		tx = tx.Where("status = ?", status)
	}
	result := tx.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Count(&total)
	if result.Error != nil {
		return data, result.Error
	}
	data.Items = list
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
