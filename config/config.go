package config

const (
	DefaultPolkadotNode = "ws://127.0.0.1:9944"
	// HttpGetResource get my resources
	HttpGetResource = "http://42.192.53.10:8888/api/resources/use"
	// HttpGetCode get verification code
	HttpGetCode = "https://ttchain.tntlinking.com/api/authorization/verification/code"
	// HttpLogin login
	HttpLogin       = "https://ttchain.tntlinking.com/api/authorization/oauth/token"
	Httpprovider    = "http://localhost:%d/api/v1/thegraph/deploy"
	HttpGraphStatus = "http://localhost:%d/api/v1/thegraph/status"

	//Port p2p port configuration
	Port = 4001

	// p2p swarm key
	SwarmKey = "/key/swarm/psk/1.0.0/\n/base16/\n55158d9b6b7e5a8e41aa8b34dd057ff1880e38348613d27ae194ad7c5b9670d7"

	// rinkbey network
	EndpointUrl = "https://rinkeby.infura.io/v3/62d7b5f33ae443e784919f1c2afe24a3"
	// eth main net
	EthMainNetwork             = "mainnet:https://cloudflare-eth.com"
	EthereumRinkebyNetworkName = "rinkeby"
	EthereumMainNetworkName    = "mainnet"
	TheGraphStakingAddress     = "0x2d44C0e097F6cD0f514edAC633d82E01280B4A5c"
)

// deploy status
const (
	NOT_DEPLOYED  int = iota //not deployed 0
	DEPLOYED                 //deployed   1
	ALL                      // all      2
	WAIT_RESOURCE            //waiting for resource  3
	IN_DEPLOYMENT            // in deployment  4
	DEPLOY_FAILED            //deploy failed 5
)

//HTTP request docker status Failure Return
const RequestStatusFailed = 5
