package p2p

import (
	"gorm.io/gorm"
	//"link-server/repository"
)

type P2pConfig struct {
	gorm.Model
	PrivateKey string `json:"privateKey"`
	Port       int    `json:"port"`
	PeerId     string `json:"peerId"`
}

type LinkInfo struct {
	Protocol      string
	ListenAddress string
	TargetAddress string
	Status        bool
}

type Service interface {
	Link(port int, peerId string) error
	LinkByProtocol(protocol string, localPort int, peerId string) error
	Close(target string) (int, error)
	Destroy() error
	GetLinks() *[]LinkInfo
	InitSetting() error
	GetSetting() (P2pConfig, error)
	ProLink(peerId string) error
	GetProviderLinks() *[]LinkInfo
	JudgeP2pReconnection() bool
	ReconnectionProLink() (bool, error)
	QueryLinks(protocol string) *[]LinkInfo
}
