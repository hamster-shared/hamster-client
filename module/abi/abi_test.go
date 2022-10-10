package abi

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"math/big"
	"testing"
)

const (
	senderAddress = "0xe9C29208012f52230866b6Cd367A789fEF6B189D"
	endpointUrl   = "https://rinkeby.infura.io/v3/62d7b5f33ae443e784919f1c2afe24a3"
	privateKey    = "46d5cce3f60bf3557f7f999a834850579708c4b0999f82376b629028dac44395"
)

var (
	sender  = GetEthAddress(senderAddress)
	privKey = GetPrivateKey(privateKey)
	client  = getClient()
)

func getClient() *ethclient.Client {
	client, err := GetEthClient(endpointUrl)
	if err != nil {
		log.Fatalln(err)
	}
	return client
}

func TestAll(t *testing.T) {
	var stakingAddress common.Address
	var err error
	stakingAmount := new(big.Int)
	stakingAmount.SetString("100000000000000000000000", 10)

	// 获取质押合约地址
	stakingAddress, err = StakeProxyFactoryAbiGetStakingAddress(context.Background(), sender, client)
	if err != nil {
		// 如果获取质押合约地址失败，则创建质押合约
		err = StakeProxyFactoryAbiCreateStakingContract(sender, client, big.NewInt(4), privKey, context.Background(), client)
		if err != nil {
			t.Fatal(err)
		}
		// 再次获取质押合约地址
		stakingAddress, err = StakeProxyFactoryAbiGetStakingAddress(context.Background(), sender, client)
		if err != nil {
			t.Fatal(err)
		}
	}
	fmt.Println("get staking stakingAddress:", stakingAddress)

	// 批准
	err = Ecr20AbiApprove(stakingAddress, client, big.NewInt(4), stakingAmount, privKey, context.Background(), client)
	if err != nil {
		t.Fatal(err)
	}

	// 从质押地址获取质押金额
	amount, err := StakeDistributionProxyAbiGetStakingAmount(context.Background(), stakingAddress, client)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("get staking amount:", amount)
	if amount.Cmp(big.NewInt(0)) == 0 {
		// 质押
		err = StakeDistributionProxyAbiStaking(stakingAddress, client, big.NewInt(4), stakingAmount, privKey, context.Background(), client)
		if err != nil {
			t.Fatal(err)
		}
	}
}
