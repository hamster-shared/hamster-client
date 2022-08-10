package cli

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"hamster-client/module/account"
	"hamster-client/module/application"
	"hamster-client/module/deploy"
	"hamster-client/module/keystorage"
	"hamster-client/module/p2p"
)

type ServiceImpl struct {
	ctx                context.Context
	db                 *gorm.DB
	keyStorageService  keystorage.Service
	accountService     account.Service
	applicationService application.Service
	p2pServer          p2p.Service
	deployService      deploy.Service
}

func NewServiceImpl(ctx context.Context, db *gorm.DB, keyStorageService keystorage.Service, accountService account.Service, applicationService application.Service, p2pServer p2p.Service, deployService deploy.Service) ServiceImpl {
	return ServiceImpl{ctx, db, keyStorageService, accountService, applicationService, p2pServer, deployService}
}

func (c *ServiceImpl) CliLink(applicationId int) (int, error) {
	var port int
	data, err := c.applicationService.QueryCliP2pPort(applicationId)
	if err != nil {
		return data, err
	}
	info, err := c.accountService.GetAccount()
	if err != nil {
		return port, err
	}
	if info.PeerId == "" {
		return port, errors.New("Peerid is empty")
	}
	_, resultErr := c.p2pServer.GetSetting()
	if resultErr != nil {
		err := c.p2pServer.InitSetting()
		if err != nil {
			return port, err
		}
	}
	protocol := "/x/graph-cli"
	nextPort := c.applicationService.QueryNextCliP2pPort()
	if nextPort == 44000 {
		port = nextPort
		err = c.applicationService.UpdateApplicationCliForwardPort(applicationId, port)
		if err != nil {
			fmt.Println("update cli p2p forward port failed,error is: ", err)
			return 0, err
		}
		err = c.p2pServer.LinkByProtocol(protocol, port, info.PeerId)
		if err != nil {
			fmt.Println("cli p2p forward link error, error is: ", err)
			return 0, err
		}
		return port, nil
	}
	currentPort, err := c.applicationService.QueryCliP2pPort(applicationId)
	if err != nil {
		return port, err
	}

	return currentPort, nil
}

func (c *ServiceImpl) judgePortIsActive(port int, protocol string) (bool, error) {
	links := c.p2pServer.QueryLinks(protocol)
	if len(*links) > 0 {
		for _, value := range *links {
			if value.Status {

			}
		}
		return false, nil
	}
	return true, nil
}
