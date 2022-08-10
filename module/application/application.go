package application

import (
	"gorm.io/gorm"
	"time"
)

type Application struct {
	ID             uint           `json:"id" gorm:"primarykey"`
	Name           string         `json:"name"`           //apply name
	SelectNodeType string         `json:"selectNodeType"` //apply plugin
	Status         int            `json:"status"`         //apply status 0: not deploy 1:deployed 2:ALL 3:wait resource 4:In deployment 5:deploy failed
	P2pForwardPort int            `json:"p2pForwardPort"`
	CliForwardPort int            `json:"cliForwardPort"`
	GrtIncome      int64          `json:"grtIncome"`
	LeaseTerm      int            `json:"leaseTerm"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

type AddApplicationParam struct {
	Name           string `json:"name"`           //apply name
	SelectNodeType string `json:"selectNodeType"` //apply plugin
}

type UpdateApplicationParam struct {
	ID             uint   `json:"id"`             //application ID
	Name           string `json:"name"`           //apply name
	SelectNodeType string `json:"selectNodeType"` //apply plugin
}

type PageApplicationVo struct {
	Items []ListVo `json:"items"`
	Total int64    `json:"total"`
}

type ApplyVo struct {
	ID             uint      `json:"id"`
	CreatedAt      time.Time `json:"createdAt"`
	Name           string    `json:"name"`           //apply name
	SelectNodeType string    `json:"selectNodeType"` //apply plugin
	Status         int       `json:"status"`
	LeaseTerm      int       `json:"leaseTerm"`
	P2pForwardPort int       `json:"p2pForwardPort"`
}

type ListVo struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`           //apply name
	SelectNodeType string `json:"selectNodeType"` //apply plugin
	Status         int    `json:"status"`
	GrtIncome      int64  `json:"grtIncome"`
}

type Service interface {
	AddApplication(application *AddApplicationParam) (bool, error)
	UpdateApplication(id int, name string, plugin string) (bool, error)
	DeleteApplication(id int) (bool, error)
	QueryApplicationById(id int) (ApplyVo, error)
	ApplicationList(page, pageSize int, name string, status int) (PageApplicationVo, error)
	UpdateApplicationStatus(id, status int) error
	UpdateApplicationP2pForwardPort(id, port int) error
	QueryNextP2pPort() int
	QueryCliP2pPort(id int) (int, error)
	QueryNextCliP2pPort() int
	UpdateApplicationCliForwardPort(id, port int) error
}
