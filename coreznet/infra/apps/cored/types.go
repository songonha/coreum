package cored

import (
	"fmt"
	"math/big"
)

// Ports defines ports used by cored application
type Ports struct {
	RPC     int `json:"rpc"`
	P2P     int `json:"p2p"`
	GRPC    int `json:"grpc"`
	GRPCWeb int `json:"grpcWeb"`
	PProf   int `json:"pprof"`
}

// Wallet stores information related to wallet
type Wallet struct {
	// Name is the name of the key stored in keystore
	Name string

	// Address is the address of the wallet
	Address string
}

func (w Wallet) String() string {
	return fmt.Sprintf("%s@%s", w.Name, w.Address)
}

// Balance stores balance of denom
type Balance struct {
	// Amount is stored amount
	Amount *big.Int `json:"amount"`

	// Denom is a token symbol
	Denom string `json:"denom"`
}

func (b Balance) String() string {
	return b.Amount.String() + b.Denom
}