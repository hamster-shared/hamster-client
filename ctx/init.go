package ctx

import (
	"context"
	_ "embed"
	"fmt"
	"hamster-client/app"
	"hamster-client/module/account"
	"hamster-client/module/application"
	"hamster-client/module/chainmanager"
	"hamster-client/module/common"
	"hamster-client/module/graph"
	param "hamster-client/module/graph/v2"
	"hamster-client/module/keystorage"
	"hamster-client/module/logger"
	"hamster-client/module/p2p"
	"hamster-client/module/resource"
	"hamster-client/module/wallet"
	"hamster-client/utils"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type App struct {
	gormDB   *gorm.DB
	httpUtil *utils.HttpUtil
	ctx      context.Context

	AccountService  account.Service
	P2pService      p2p.Service
	ResourceService resource.Service
	WalletService   wallet.Service
	// DeployService           deploy.Service
	ApplicationService application.Service
	// GraphParamsService      graph.Service
	KeyStorageService *keystorage.Client
	// GraphDeployParamService param.Client
	// CliService              cli.Service
	ChainManagerService chainmanager.Manager

	AccountApp  app.Account
	P2pApp      app.P2p
	ResourceApp app.Resource
	SettingApp  app.Setting
	WalletApp   app.Wallet
	// DeployApp       app.Deploy
	ApplicationApp app.Application
	// GraphApp        app.Graph
	KeyStorageApp   app.KeyStorage
	ChainManagerApp app.ChainManager
	QueueApp        app.Queue
}

func NewApp() *App {
	a := &App{}
	a.init()
	return a
}

func (a *App) init() {
	// initialize the database
	a.initDB()
	// tired of initializing http tools
	a.initHttp()
}

func (a *App) initDB() {
	configPath := initConfigPath()
	db, err := gorm.Open(sqlite.Open(filepath.Join(configPath, "link.db")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(
		&account.Account{},
		&p2p.P2pConfig{},
		&resource.Resource{},
		&wallet.Wallet{},
		//&application.Application{},
		&graph.GraphParameter{},
		&param.GraphDeployParameter{},
		&common.EthereumDeployParam{},
	)
	var user account.Account
	result := db.First(&user)
	if result.Error != nil {
		println("------------------------------------")
		user := account.Account{WsUrl: "ws://127.0.0.1:9944"}
		db.Create(&user)
	}
	if err != nil {
		panic("failed to AutoMigrate Account")
	}
	a.gormDB = db
}

func (a *App) initHttp() {
	a.httpUtil = utils.NewHttp()
}

func (a *App) initService() {
	// graphParamServiceImpl := graph.NewClient(a.ctx, a.gormDB, a.httpUtil)
	// a.GraphParamsService = &graphParamServiceImpl
	accountServiceImpl := account.NewServiceImpl(a.ctx, a.gormDB, a.httpUtil)
	a.AccountService = &accountServiceImpl
	applicationServiceImpl := application.NewServiceImpl(a.ctx, a.gormDB)
	a.ApplicationService = &applicationServiceImpl
	p2pServiceImpl := p2p.NewServiceImpl(a.ctx, a.gormDB, a.AccountService)
	a.P2pService = &p2pServiceImpl
	resourceServiceImpl := resource.NewServiceImpl(a.ctx, a.gormDB, a.httpUtil)
	a.ResourceService = &resourceServiceImpl
	walletServiceImpl := wallet.NewServiceImpl(a.ctx, a.gormDB)
	a.WalletService = &walletServiceImpl
	keyStorageServiceImpl := keystorage.NewClient(a.gormDB)
	a.KeyStorageService = &keyStorageServiceImpl
	// deployServiceImpl := deploy.NewClient(a.ctx, a.httpUtil, a.gormDB, a.KeyStorageService, a.P2pService, a.WalletService, a.ApplicationService)
	// a.DeployService = &deployServiceImpl
	// graphDeployParamServiceImpl := param.NewClient(a.ctx, a.gormDB, *a.KeyStorageService, a.AccountService, a.ApplicationService, a.P2pService, a.DeployService, a.WalletService, a.QueueService)
	// a.GraphDeployParamService = &graphDeployParamServiceImpl
	// cliServiceImpl := cli.NewClient(a.ctx, a.gormDB, *a.KeyStorageService, a.AccountService, a.ApplicationService, a.P2pService, a.DeployService)
	// a.CliService = &cliServiceImpl
	a.ChainManagerService = chainmanager.NewManager(
		a.gormDB,
		a.ApplicationService,
		a.P2pService,
		a.AccountService,
		a.WalletService,
	)
}

func (a *App) initApp() {
	a.AccountApp = app.NewAccountApp(a.AccountService)
	a.P2pApp = app.NewP2pApp(a.P2pService)
	a.ResourceApp = app.NewResourceApp(a.ResourceService, a.AccountService)
	// a.SettingApp = app.NewSettingApp(a.P2pService, a.AccountService, a.gormDB, *a.KeyStorageService)
	a.WalletApp = app.NewWalletApp(a.WalletService)
	// a.DeployApp = app.NewDeployApp(a.DeployService, a.AccountService, a.P2pService)
	// a.ApplicationApp = app.NewApplicationApp(a.ApplicationService, a.P2pService, a.DeployService)
	// a.GraphApp = app.NewGraphApp(a.GraphParamsService, a.CliService, a.GraphDeployParamService)
	a.KeyStorageApp = app.NewKeyStorageApp(a.KeyStorageService)
	a.ChainManagerApp = app.NewChainManagerApp(a.ChainManagerService)
}

func initConfigPath() string {
	// initialize the configuration file
	dir := "~/.link/"
	linkConfig, err := homedir.Expand(dir)
	if err != nil {
		panic("failed to homedir Expand")
	}
	_, err = os.Stat(linkConfig)
	if err == nil {
		return linkConfig
	}
	err = os.MkdirAll(linkConfig, os.ModePerm)
	if err != nil {
		fmt.Printf("failed to config Mkdir err%s\n", err)
		panic("failed to config Mkdir err")
	}
	return linkConfig
}

// Startup is called at application startup
func (a *App) Startup(context context.Context) {
	// Perform your setup here
	a.ctx = context
	// initialize service
	a.initService()
	// initialize app
	a.initApp()
	// init chain queue
	a.initChainQueue()
}

// DomReady is called after the front-end dom has been loaded
func (a *App) DomReady(ctx context.Context) {
	// Add your action here
}

// Shutdown is called at application termination
func (a *App) Shutdown(ctx context.Context) {
	// Perform your teardown here
}

func (a *App) initChainQueue() {
	logger.Infof("start init chain queue")
	list, err := a.ApplicationService.ApplicationList(0, 1000, "", application.All)
	if err != nil {
		logger.Errorf("get ApplicationList error: %s", err)
		return
	}
	for _, item := range list.Items {
		err := a.ChainManagerApp.RetryStartQueue(int(item.ID), false)
		if err != nil {
			logger.Errorf("init queue error: %s, app id: %d", err, item.ID)
			continue
		}
	}
}
