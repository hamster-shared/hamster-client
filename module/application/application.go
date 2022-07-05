package application

import (
	"gorm.io/gorm"
	"time"
)

type Application struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Name      string         `json:"name"`     //apply name
	Describe  string         `json:"describe"` //apply describe
	Status    int            `json:"status"`   //apply status 0: not deploy 1:deployed 2:ALL
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

type AddApplicationParam struct {
	Name     string `json:"name"`     //apply name
	Describe string `json:"describe"` //apply describe
}

type UpdateApplicationParam struct {
	ID       uint   `json:"id"`       //application ID
	Name     string `json:"name"`     //apply name
	Describe string `json:"describe"` //apply describe
}

type PageApplicationVo struct {
	Items []Application `json:"items"`
	Total int64         `json:"total"`
}

type ApplyVo struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name"`     //apply name
	Describe  string    `json:"describe"` //apply describe
	Status    int       `json:"status"`
}

type Service interface {
	AddApplication(application *AddApplicationParam) (bool, error)
	UpdateApplication(id int, name string, des string) (bool, error)
	DeleteApplication(id int) (bool, error)
	QueryApplicationById(id int) (ApplyVo, error)
	ApplicationList(page, pageSize int, name string, status int) (PageApplicationVo, error)
	UpdateApplicationStatus(id, status int) error
}
