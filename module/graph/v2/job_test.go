package v2

import (
	"context"
	"fmt"
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
	ctx := context.Background()
	httpUtil := utils.NewHttp()
	accountService := account.NewServiceImpl(ctx, db, httpUtil)
	applicationService := application.NewServiceImpl(ctx, db)
	p2pService := p2p.NewServiceImpl(ctx, db, &accountService)
	keyStorageService := keystorage.NewClient(db)
	walletService := wallet.NewServiceImpl(ctx, db)
	deployService := deploy.NewServiceImpl(ctx, httpUtil, db, &keyStorageService, &p2pService, &walletService, &applicationService)
	queueService := queue2.NewServiceImpl()
	graphParamService := NewServiceImpl(ctx, db, keyStorageService, &accountService, &applicationService, &p2pService, &deployService, &walletService, queueService)
	//create application
	var addParam AddParam
	addParam.Name = "Client one12"
	addParam.ThegraphIndexer = "chef moon high razor hockey steak better version myself large purchase cave"
	addParam.SelectNodeType = "thegraph_rinkeby"
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
	queue, err := queue2.GetQueue(int(res.ID))
	if err != nil {
		fmt.Println("new queue failed,err is: ", err)
		t.Error(err)
	}
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
	// view status
	info, _ := queue.GetStatus()

	fmt.Println(info)
}

func TestGraphRules(t *testing.T) {
	g := getGraphParamService()
	rules, err := g.GraphRules(4)
	if err != nil {
		t.Error(err)
	}
	for _, v := range rules {
		fmt.Println(v)
	}
}

func TestGraphStart(t *testing.T) {
	g := getGraphParamService()
	err := g.GraphStart(4, "QmVqMeQUwvQ3XjzCYiMhRvQjRiQLGpVt8C3oHgvDi3agJ2")
	if err != nil {
		t.Error(err)
	}
}

func getGraphParamService() ServiceImpl {
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
	ctx := context.Background()
	httpUtil := utils.NewHttp()
	accountService := account.NewServiceImpl(ctx, db, httpUtil)
	applicationService := application.NewServiceImpl(ctx, db)
	p2pService := p2p.NewServiceImpl(ctx, db, &accountService)
	keyStorageService := keystorage.NewClient(db)
	walletService := wallet.NewServiceImpl(ctx, db)
	deployService := deploy.NewServiceImpl(ctx, httpUtil, db, &keyStorageService, &p2pService, &walletService, &applicationService)
	queueService := queue2.NewServiceImpl()
	graphParamService := NewServiceImpl(ctx, db, keyStorageService, &accountService, &applicationService, &p2pService, &deployService, &walletService, queueService)
	return graphParamService
}

func TestHttp(t *testing.T) {
	url := fmt.Sprintf("http://localhost:%d/version", 8888)
	fmt.Println("测试 api 连通性:", url)
	req := utils.NewHttp().NewRequest()
	resp, err := req.Get(url)

	fmt.Println("连通性结果： ", err != nil)

	fmt.Println(string(resp.Body()))

	if err != nil {
		panic(err)
	}
}
