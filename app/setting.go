package app

import (
	"context"
	"gorm.io/gorm"
	"hamster-client/module/account"
	"hamster-client/module/deploy"
	"hamster-client/module/keystorage"
	"hamster-client/module/p2p"
	"hamster-client/module/pallet"
)

type Config struct {
	PublicKey string `json:"publicKey"`
	Port      int    `json:"port"`
	PeerId    string `json:"peerId"`
	WsUrl     string `json:"wsUrl"`
}

type Setting struct {
	ctx               context.Context
	db                *gorm.DB
	keyStorageService keystorage.Service
	deployService     deploy.Service
	p2pService        p2p.Service
	accountService    account.Service
	chainListener     pallet.ChainListener
}

func NewSettingApp(service p2p.Service, accountService account.Service, db *gorm.DB, keyStorageService keystorage.Service, deployService deploy.Service) Setting {
	return Setting{
		p2pService:        service,
		accountService:    accountService,
		db:                db,
		keyStorageService: keyStorageService,
		deployService:     deployService,
	}
}

func (s *Setting) WailsInit(ctx context.Context) error {
	s.ctx = ctx
	return nil
}

// GetSetting view configuration information
func (s *Setting) GetSetting() (*Config, error) {
	config := &Config{}
	//query the public key in the user information
	info, err := s.accountService.GetAccount()
	if err != nil {
		return config, err
	}
	config.PublicKey = info.PublicKey
	config.WsUrl = info.WsUrl
	//query the setting information in the p2p setting
	p2pConfig, err := s.p2pService.GetSetting()
	if err != nil {
		return config, nil
	}
	config.Port = p2pConfig.Port
	config.PeerId = p2pConfig.PeerId
	return config, nil
}

// Setting set public key information
func (s *Setting) Setting(publicKey string, wsUrl string) (bool, error) {
	accountInfo, _ := s.accountService.GetAccount()

	//set user public key
	accountInfo.PublicKey = publicKey
	accountInfo.WsUrl = wsUrl
	s.accountService.SaveAccount(&accountInfo)
	// close go func
	//s.chainListener.CancelListen()
	//start go func
	//s.chainListener.StartListen(s.db, s.keyStorageService, s.deployService)
	return true, nil
}

//InitP2pSetting initialize p2p settings
func (s *Setting) InitP2pSetting() (bool, error) {
	//initial configuration
	err := s.p2pService.InitSetting()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Setting) SettingWsUrl(wsUrl string) (bool, error) {
	accountInfo, err := s.accountService.GetAccount()
	if err != nil {
		return false, err
	}
	accountInfo.WsUrl = wsUrl
	s.accountService.SaveAccount(&accountInfo)
	// close go func
	s.chainListener.CancelListen()
	//start go func
	s.chainListener.StartListen(s.db, s.keyStorageService, s.deployService)
	return true, nil
}

func (s *Setting) SettingPublicKey(publicKey string) (bool, error) {
	accountInfo, err := s.accountService.GetAccount()
	if err != nil {
		return false, err
	}
	accountInfo.PublicKey = publicKey
	s.accountService.SaveAccount(&accountInfo)
	return true, nil
}
