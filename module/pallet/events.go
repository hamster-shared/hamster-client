package pallet

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type EventFreeResourceApplied struct {
	Phase      types.Phase
	OrderIndex types.U64
	PeerId     string
	Topics     []types.Hash
}

type EventBalance struct {
	Phase     types.Phase
	AccountId types.AccountID
	Amount    types.U128
	Topics    []types.Hash
}

type EventResourceOrderCreateOrderSuccess struct {
	Phase         types.Phase
	AccountId     types.AccountID
	OrderIndex    types.U64
	ResourceIndex types.U64
	Duration      types.U32
	DeployType types.U32
	PublicKey     string
	Topics        []types.Hash
}

type EventResourceOrderExecSuccess struct {
	Phase          types.Phase
	AccountId      types.AccountID
	OrderIndex     types.U64
	ResourceIndex  types.U64
	AgreementIndex types.U64
	Topics         []types.Hash
}

type MyEventRecords struct {
	types.EventRecords
	ResourceOrder_CreateOrderSuccess []EventResourceOrderCreateOrderSuccess //nolint:stylecheck,golint
	ResourceOrder_OrderExecSuccess   []EventResourceOrderExecSuccess
	Balances_Withdraw                []EventBalance //org approve event
}
