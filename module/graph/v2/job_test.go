package v2

import (
	"context"
	"fmt"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/mitchellh/go-homedir"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"hamster-client/module/account"
	"hamster-client/module/application"
	"hamster-client/module/deploy"
	"hamster-client/module/graph"
	"hamster-client/module/keystorage"
	"hamster-client/module/p2p"
	queue2 "hamster-client/module/queue"
	"hamster-client/module/resource"
	"hamster-client/module/wallet"
	"hamster-client/utils"
	"path/filepath"
	"testing"
	"time"
)

func TestDeploy(t *testing.T) {
	configPath, _ := homedir.Expand("~/.link/")
	db, err := gorm.Open(sqlite.Open(filepath.Join(configPath, "link.db")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(
		&account.Account{},
		&p2p.P2pConfig{},
		&resource.Resource{},
		&wallet.Wallet{},
		&application.Application{},
		&graph.GraphParameter{},
	)

	//applicationId := 11
	var applicationId int
	ctx := context.Background()
	httpUtil := utils.NewHttp()
	accountService := account.NewServiceImpl(ctx, db, httpUtil)
	applicationService := application.NewServiceImpl(ctx, db)
	p2pService := p2p.NewServiceImpl(ctx, db)
	keyStorageService := keystorage.NewServiceImpl(ctx, db)
	deployService := deploy.NewServiceImpl(ctx, httpUtil, db, &keyStorageService, &accountService, &p2pService)
	graphParamService := NewServiceImpl(ctx, db, keyStorageService, &accountService, &applicationService, &p2pService, &deployService)
	//create application
	var addParam AddParam
	addParam.Name = "Service one12"
	addParam.ThegraphIndexer = "chef moon high razor hockey steak better version myself large purchase cave"
	addParam.SelectNodeType = "The Graph"
	addParam.StakingAmount = 100000
	addParam.LeaseTerm = 1
	res, err := graphParamService.SaveGraphDeployParameterAndApply(addParam)
	if err != nil {
		fmt.Println("create deploy service failed,err is: ", err)
		return
	}
	if !res.Result {
		fmt.Println("create deploy service failed,err is: ", err)
		return
	}
	applicationId = int(res.ID)
	stakingJob := NewGraphStakingJob(keyStorageService, applicationId)
	pullJob := NewPullImageJob(&applicationService, applicationId)
	substrateApi, _ := gsrpc.NewSubstrateAPI("ws://183.66.65.207:49944")

	waitResourceJob, _ := NewWaitResourceJob(substrateApi, &accountService, &applicationService, &p2pService, applicationId)
	deployJob := NewServiceDeployJob(keyStorageService, &deployService, applicationId)
	queue, err := queue2.NewQueue("1", &stakingJob, waitResourceJob, &pullJob, &deployJob)
	if err != nil {
		fmt.Println("new queue failed,err is: ", err)
		t.Error(err)
	}
	channel := make(chan struct{})
	go queue.Start(channel)
	go func() {
		for {
			time.Sleep(time.Second)
			info, err := queue.GetStatus()
			if err != nil {
				t.Error(err)
			}
			for _, v := range info {
				fmt.Print(v, "; ")
			}
			fmt.Println()
		}
	}()
	// wait
	<-channel
	// view status
	info, _ := queue.GetStatus()

	fmt.Println(info)
}
