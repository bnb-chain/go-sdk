package tx

import (
	"github.com/BiJie/bnc-go-sdk/sdk/tx/txmsg"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
)

// Coin def
type Coin struct {
	Denom  string `json:"denom"`
	Amount int64  `json:"amount"`
}

// Coins def
type Coins []Coin

// Codec def
type Codec = amino.Codec

// Cdc global variable
var Cdc *Codec

func init() {
	cdc := amino.NewCodec()

	cryptoAmino.RegisterAmino(cdc)

	cdc.RegisterInterface((*txmsg.Msg)(nil), nil)

	cdc.RegisterConcrete(txmsg.CreateOrderMsg{}, "dex/NewOrder", nil)
	cdc.RegisterConcrete(txmsg.CancelOrderMsg{}, "dex/CancelOrder", nil)
	cdc.RegisterConcrete(txmsg.TokenIssueMsg{}, "tokens/IssueMsg", nil)
	cdc.RegisterConcrete(txmsg.TokenBurnMsg{}, "tokens/BurnMsg", nil)
	cdc.RegisterConcrete(txmsg.TokenFreezeMsg{}, "tokens/FreezeMsg", nil)
	cdc.RegisterConcrete(txmsg.TokenUnfreezeMsg{}, "tokens/UnfreezeMsg", nil)
	cdc.RegisterConcrete(txmsg.DexListMsg{}, "dex/ListMsg", nil)

	Cdc = cdc.Seal()
}
