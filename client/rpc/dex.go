package rpc

import (
	"fmt"
	"github.com/binance-chain/go-sdk/common/types"
)
const (
	MainStoreName    = "main"
	AccountStoreName = "acc"
	ValAddrStoreName = "val"
	TokenStoreName   = "tokens"
	DexStoreName     = "dex"
	PairStoreName    = "pairs"
	StakeStoreName   = "stake"
	ParamsStoreName  = "params"
	GovStoreName     = "gov"

	StakeTransientStoreName  = "transient_stake"
	ParamsTransientStoreName = "transient_params"
)

func (c *HTTP) ListAllTokens(offset int, limit int) ([]types.Token, error) {
	path:=fmt.Sprintf("tokens/list/%d/%d", offset, limit)
	result, err := c.ABCIQuery(path, nil)
	if err != nil {
		return nil, err
	}
	bz:=result.Response.GetValue()
	tokens := make([]types.Token, 0)
	err = c.cdc.UnmarshalBinaryLengthPrefixed(bz, &tokens)
	return tokens, err
}

func (c *HTTP) GetTokenInfo(symbol string) (*types.Token, error) {
	path:=fmt.Sprintf("tokens/info/%s", symbol)
	result, err := c.ABCIQuery(path, nil)
	if err != nil {
		return nil, err
	}
	bz:=result.Response.GetValue()
	token := new(types.Token)
	err = c.cdc.UnmarshalBinaryLengthPrefixed(bz, token)
	return token, err
}

func (c *HTTP) GetAccount(addr types.AccAddress)(acc types.Account, err error)  {
	key:=append([]byte("account:"), addr.Bytes()...)
	bz, err := c.QueryStore(key, AccountStoreName)
	if err != nil {
		return nil, err
	}
	if bz == nil {
		return nil, nil
	}
	err = c.cdc.UnmarshalBinaryBare(bz, &acc)
	if err != nil {
		return nil, err
	}
	return acc, err
}


