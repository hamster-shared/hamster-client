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
	"strconv"
	"strings"
)

type ServiceImpl struct {
	ctx                context.Context
	db                 *gorm.DB
	keyStorageService  keystorage.Client
	accountService     account.Service
	applicationService application.Service
	p2pService         p2p.Service
	deployService      deploy.Service
}

func NewServiceImpl(ctx context.Context, db *gorm.DB, keyStorageService keystorage.Client, accountService account.Service, applicationService application.Service, p2pService p2p.Service, deployService deploy.Service) ServiceImpl {
	return ServiceImpl{ctx, db, keyStorageService, accountService, applicationService, p2pService, deployService}
}

func (c *ServiceImpl) CliLink(applicationId int) (int, error) {
	var port int
	data, err := c.applicationService.QueryCliP2pPort(applicationId)
	if err != nil {
		return data, err
	}
	info, err := c.applicationService.QueryApplicationById(applicationId)
	if err != nil {
		return port, err
	}
	if info.PeerId == "" {
		return port, errors.New("Peerid is empty")
	}
	_, resultErr := c.p2pService.GetSetting()
	if resultErr != nil {
		err := c.p2pService.InitSetting()
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
		err = c.p2pService.LinkByProtocol(protocol, port, info.PeerId)
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
	isActive := c.judgePortIsActive(currentPort, protocol)
	if isActive {
		return currentPort, nil
	}
	err = c.applicationService.UpdateApplicationCliForwardPort(applicationId, nextPort)
	if err != nil {
		fmt.Println("update cli p2p forward port failed,error is: ", err)
		return 0, err
	}
	err = c.p2pService.LinkByProtocol(protocol, nextPort, info.PeerId)
	if err != nil {
		fmt.Println("cli p2p forward link error, error is: ", err)
		return 0, err
	}
	return nextPort, nil
}

func (c *ServiceImpl) judgePortIsActive(port int, protocol string) bool {
	links := c.p2pService.QueryLinks(protocol)
	if len(*links) > 0 {
		for _, value := range *links {
			listenAddress := value.ListenAddress
			if listenAddress != "" {
				lastIndex := strings.LastIndex(listenAddress, "/")
				listenPort := listenAddress[lastIndex+1 : len(listenAddress)-1]
				listenPortInt, err := strconv.Atoi(listenPort)
				if err != nil {
					continue
				}
				if listenPortInt == port {
					if value.Status {
						return true
					}
					return false
				}
			}
		}
	}
	return false
}
