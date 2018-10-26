package tx

import (
	"fmt"

	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
)

type Coin struct {
	Denom  string `json:"denom"`
	Amount int64  `json:"amount"`
}

type Coins []Coin

type Codec = amino.Codec

var Cdc *Codec

func init() {
	fmt.Println("init cdc")
	cdc := amino.NewCodec()
	cryptoAmino.RegisterAmino(cdc)
	Cdc = cdc.Seal()
}

type Tx struct{}
