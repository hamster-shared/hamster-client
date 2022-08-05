package v2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"hamster-client/module/application"
	"hamster-client/module/deploy"
	"hamster-client/module/keystorage"
	"strconv"
)

type ServiceImpl struct {
	ctx               context.Context
	db                *gorm.DB
	keyStorageService keystorage.Service
}

func NewServiceImpl(ctx context.Context, db *gorm.DB, keyStorageService keystorage.Service) ServiceImpl {
	return ServiceImpl{ctx, db, keyStorageService}
}

func (g *ServiceImpl) SaveGraphDeployParameterAndApply(addData AddParam) (bool, error) {
	var applyData application.Application
	var deployData GraphDeployParameter
	err := g.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("name=?", addData.Name).First(&applyData).Error; err != gorm.ErrRecordNotFound {
			return errors.New(fmt.Sprintf("application:%s already exists", addData.Name))
		}
		applyData.Name = addData.Name
		applyData.Plugin = addData.Plugin
		if err := tx.Create(&applyData).Error; err != nil {
			return err
		}
		deployData.Application = applyData
		deployData.LeaseTerm = addData.LeaseTerm
		deployData.Network = addData.Network
		deployData.Mnemonic = addData.Mnemonic
		deployData.PledgeAmount = addData.PledgeAmount
		deployData.ApplicationID = applyData.ID
		if err := tx.Create(&deployData).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	var deploymentData deploy.DeployParameter
	deploymentData.Initialization.AccountMnemonic = addData.Mnemonic
	deploymentData.Initialization.LeaseTerm = addData.LeaseTerm
	deploymentData.Staking.NetworkUrl = addData.Network
	deploymentData.Staking.PledgeAmount = addData.PledgeAmount
	deploymentData.Deployment.EthereumUrl = addData.Network
	jsonData, err := json.Marshal(deploymentData)
	if err != nil {
		return false, err
	}
	g.keyStorageService.Set("graph_"+strconv.Itoa(int(applyData.ID)), string(jsonData))
	return true, nil
}

func (g *ServiceImpl) DeleteGraphDeployParameterAndApply(id int) (bool, error) {
	err := g.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Debug().Where("id = ?", id).Delete(&application.Application{}).Error; err != nil {
			return err
		}
		if err := tx.Debug().Where("applicationId = ?", id).Delete(GraphDeployParameter{}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
