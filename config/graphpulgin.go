package config

type PluginDeployInfo struct {
	EndpointUrl            string `json:"endpointUrl"`
	EthNetwork             string `json:"ethNetwork"`
	EthereumNetworkName    string `json:"ethereumNetworkName"`
	TheGraphStakingAddress string `json:"theGraphStakingAddress"`
	TheGraphTokenAddress   string `json:"theGraphTokenAddress"`
}

var PluginMap = map[string]PluginDeployInfo{
	"thegraph_mainnet": {
		EndpointUrl:            "mainnet:https://cloudflare-eth.com",
		EthNetwork:             "mainnet:https://cloudflare-eth.com",
		EthereumNetworkName:    "mainnet",
		TheGraphStakingAddress: "0xF55041E37E12cD407ad00CE2910B8269B01263b9",
		TheGraphTokenAddress:   "0xc944E90C64B2c07662A292be6244BDf05Cda44a7",
	},
	"thegraph_rinkeby": {
		EndpointUrl:            "https://rinkeby.infura.io/v3/62d7b5f33ae443e784919f1c2afe24a3",
		EthNetwork:             "mainnet:https://cloudflare-eth.com",
		EthereumNetworkName:    "rinkeby",
		TheGraphStakingAddress: "0x2d44C0e097F6cD0f514edAC633d82E01280B4A5c",
		TheGraphTokenAddress:   "0x54Fe55d5d255b8460fB3Bc52D5D676F9AE5697CD",
	},
	"thegraph_gorli": {
		EndpointUrl:            "https://goerli.infura.io/v3/62d7b5f33ae443e784919f1c2afe24a3",
		EthNetwork:             "mainnet:https://cloudflare-eth.com",
		EthereumNetworkName:    "goerli",
		TheGraphStakingAddress: "0x35e3Cb6B317690d662160d5d02A5b364578F62c9",
		TheGraphTokenAddress:   "0x5c946740441C12510a167B447B7dE565C20b9E3C",
	},
}