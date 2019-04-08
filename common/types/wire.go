package types

import "github.com/tendermint/go-amino"

func RegisterWire(cdc *amino.Codec) {

	cdc.RegisterConcrete(Token{}, "bnbchain/Token", nil)
	cdc.RegisterInterface((*Account)(nil), nil)
	cdc.RegisterInterface((*NamedAccount)(nil), nil)
	cdc.RegisterConcrete(&AppAccount{}, "bnbchain/Account", nil)
}
