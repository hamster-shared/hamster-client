package pallet

import (
	"errors"
	"fmt"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	log "github.com/sirupsen/logrus"
	"math/big"
)

func GetResourceList(meta *types.Metadata, api *gsrpc.SubstrateAPI, filter func(resource *ComputingResource) bool) (map[types.U64]*ComputingResource, error) {

	key, err := types.CreateStorageKey(meta, "Provider", "ProviderOnlineList")

	if err != nil {
		return nil, err
	}

	var onlineList []types.AccountID

	ok, err := api.RPC.State.GetStorageLatest(key, &onlineList)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("query Substrate error")
	}

	var mapData = make(map[types.U64]*ComputingResource)

	for _, account := range onlineList {
		publicKey, _ := AddressToPublicKey(AccountIdToAddress(account))
		key, err := types.CreateStorageKey(meta, "Provider", "Providers", publicKey)
		if err != nil {
			return nil, err
		}
		var providerList []types.U64
		_, err = api.RPC.State.GetStorageLatest(key, &providerList)
		if err != nil {
			return nil, err
		}

		for _, providerIndex := range providerList {
			r, err := GetResource(providerIndex, meta, api)

			if err != nil {
				return nil, err
			}
			if filter != nil {
				if !filter(r) {
					continue
				}
			}
			mapData[r.Index] = r
		}
	}

	return mapData, nil
}

func GetResource(index types.U64, meta *types.Metadata, api *gsrpc.SubstrateAPI) (*ComputingResource, error) {
	bytes, err := types.Encode(index)
	if err != nil {
		return nil, err
	}
	key, err := types.CreateStorageKey(meta, "Provider", "Resources", bytes)

	var computeResource ComputingResource
	ok, err := api.RPC.State.GetStorageLatest(key, &computeResource)
	if err != nil {
		return nil, err
	}

	if ok {
		return &computeResource, nil
	} else {
		return nil, errors.New("no resource")
	}
}

func CallAndWatch(api *gsrpc.SubstrateAPI, c types.Call, meta *types.Metadata, hook func(header *types.Header) error, pair signature.KeyringPair) error {

	// Create the extrinsic
	ext := types.NewExtrinsic(c)
	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return err
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return err
	}

	// Get the nonce for Account
	key, err := types.CreateStorageKey(meta, "System", "Account", pair.PublicKey)
	if err != nil {
		return err
	}

	var accountInfo types.AccountInfo
	ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("GetStorageLatest fail")
	}

	nonce := uint32(accountInfo.Nonce)
	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	// Sign the transaction using User's default account
	err = ext.Sign(pair, o)
	if err != nil {
		return err
	}

	sub, err := api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer sub.Unsubscribe()

	for {
		status := <-sub.Chan()
		fmt.Printf("Transaction status: %#v\n", status)

		if status.IsInBlock {
			hh, err := api.RPC.Chain.GetHeader(status.AsInBlock)
			fmt.Printf("Completed at block hash: %#x\n", hh.Number)

			if hook != nil {
				if err != nil {
					return err
				}
				return hook(hh)
			}

			return nil
		}

		if status.IsDropped || status.IsInvalid {
			fmt.Println("submit tx fail , blocknumber is : ")
			return errors.New("submit tx fail")
		}
	}
}

func Bond(api *gsrpc.SubstrateAPI, meta *types.Metadata, amount int64, pair signature.KeyringPair) error {
	log.Info("bonding...")
	c, err := types.NewCall(meta, "Market.bond", types.NewU128(*big.NewInt(amount)))
	if err != nil {
		fmt.Println(err.Error())
		log.Errorf("NewCall error: %s", err.Error())
		return err
	}
	log.Infof("call succeed")
	return CallAndWatch(api, c, meta, func(header *types.Header) error {
		fmt.Println(header.Digest)
		return err
	}, pair)
}

func GetBlockNumber(api *gsrpc.SubstrateAPI) (uint64, error) {
	header, err := api.RPC.Chain.GetHeaderLatest()
	if err != nil {
		return 0, err
	}

	return uint64(header.Number), nil
}

func GetEvent(api *gsrpc.SubstrateAPI, meta *types.Metadata, blockNumber uint64) (*MyEventRecords, error) {
	log.Info("get event: ", blockNumber)
	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, err
	}
	bh, err := api.RPC.Chain.GetBlockHash(blockNumber)
	if err != nil {
		return nil, err
	}
	key, err := types.CreateStorageKey(meta, "System", "Events", nil)
	if err != nil {
		return nil, err
	}
	raw, err := api.RPC.State.GetStorageRaw(key, bh)
	if err != nil {
		return nil, err
	}
	// Decode the event records
	events := MyEventRecords{}
	err = types.EventRecordsRaw(*raw).DecodeEventRecords(meta, &events)
	return &events, err
}

// CancelOrder  cancel order or cancel agreement
func CancelOrder(api *gsrpc.SubstrateAPI, meta *types.Metadata, keypair signature.KeyringPair, orderIndex int) error {
	order, err := GetOrder(types.NewU64(uint64(orderIndex)), meta, api)
	if err != nil {
		return err
	}

	if order.Status.IsCanceled {
		return nil
	}

	if order.Status.IsPending {
		return cancelOrder(api, meta, keypair, orderIndex)
	}

	return cancelAgreement(api, meta, keypair, order)

}

func cancelOrder(api *gsrpc.SubstrateAPI, meta *types.Metadata, keypair signature.KeyringPair, orderIndex int) error {
	c, err := types.NewCall(meta, "ResourceOrder.cancel_order", types.NewU64(uint64(orderIndex)))

	if err != nil {
		return err
	}

	return CallAndWatch(api, c, meta, nil, keypair)
}

func cancelAgreement(api *gsrpc.SubstrateAPI, meta *types.Metadata, keypair signature.KeyringPair, order *ResourceOrder) error {
	if order.AgreementIndex.IsNone() {
		return nil
	}

	ok, agreementId := order.AgreementIndex.Unwrap()

	if !ok {
		return nil
	}

	agreement, err := GetRentalAgreement(agreementId, meta, api)
	if err != nil {
		return err
	}

	// check agreement status

	if agreement.Status.IsUsing {
		c, err := types.NewCall(meta, "ResourceOrder.cancel_agreement", agreementId)

		if err != nil {
			return err
		}

		return CallAndWatch(api, c, meta, nil, keypair)
	}
	return nil
}

func GetOrder(index types.U64, meta *types.Metadata, api *gsrpc.SubstrateAPI) (*ResourceOrder, error) {
	bytes, err := types.Encode(index)
	if err != nil {
		return nil, err
	}
	key, err := types.CreateStorageKey(meta, "ResourceOrder", "ResourceOrders", bytes)

	var order ResourceOrder
	ok, err := api.RPC.State.GetStorageLatest(key, &order)
	if err != nil {
		return nil, err
	}

	if ok {
		return &order, nil
	} else {
		return nil, errors.New("no resource")
	}
}

func GetRentalAgreement(index types.U64, meta *types.Metadata, api *gsrpc.SubstrateAPI) (*RentalAgreement, error) {
	bytes, err := types.Encode(index)
	if err != nil {
		return nil, err
	}
	key, err := types.CreateStorageKey(meta, "ResourceOrder", "RentalAgreements", bytes)

	var agreement RentalAgreement
	ok, err := api.RPC.State.GetStorageLatest(key, &agreement)
	if err != nil {
		return nil, err
	}

	if ok {
		return &agreement, nil
	} else {
		return nil, errors.New("no resource")
	}
}
