package rpc

import (
	"errors"
	"fmt"
	"github.com/binance-chain/go-sdk/common/types"
	"strings"
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
	ParamABCIPrefix  = "param"

	StakeTransientStoreName  = "transient_stake"
	ParamsTransientStoreName = "transient_params"
)

func (c *HTTP) ListAllTokens(offset int, limit int) ([]types.Token, error) {
	path := fmt.Sprintf("tokens/list/%d/%d", offset, limit)
	result, err := c.ABCIQuery(path, nil)
	if err != nil {
		return nil, err
	}
	bz := result.Response.GetValue()
	tokens := make([]types.Token, 0)
	err = c.cdc.UnmarshalBinaryLengthPrefixed(bz, &tokens)
	return tokens, err
}

func (c *HTTP) GetTokenInfo(symbol string) (*types.Token, error) {
	path := fmt.Sprintf("tokens/info/%s", symbol)
	result, err := c.ABCIQuery(path, nil)
	if err != nil {
		return nil, err
	}
	bz := result.Response.GetValue()
	token := new(types.Token)
	err = c.cdc.UnmarshalBinaryLengthPrefixed(bz, token)
	return token, err
}

func (c *HTTP) GetAccount(addr types.AccAddress) (acc types.Account, err error) {
	key := append([]byte("account:"), addr.Bytes()...)
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

func (c *HTTP) GetBalances(addr types.AccAddress) ([]types.TokenBalance, error) {
	account, err := c.GetAccount(addr)
	if err != nil {
		return nil, err
	}
	coins := account.GetCoins()
	var denoms map[string]bool
	denoms = map[string]bool{}
	for _, coin := range coins {
		denom := coin.Denom
		exists := c.existsCC(denom)
		// TODO: we probably actually want to show zero balances.
		// if exists && !sdk.Int.IsZero(coins.AmountOf(denom)) {
		if exists {
			denoms[denom] = true
		}
	}

	symbs := make([]string, 0, len(denoms))
	bals := make([]types.TokenBalance, 0, len(denoms))
	for symb := range denoms {
		symbs = append(symbs, symb)
		// count locked and frozen coins
		var locked, frozen int64
		nacc := account.(types.NamedAccount)
		if nacc != nil {
			locked = nacc.GetLockedCoins().AmountOf(symb)
			frozen = nacc.GetFrozenCoins().AmountOf(symb)
		}
		bals = append(bals, types.TokenBalance{
			Symbol: symb,
			Free:   types.Fixed8(coins.AmountOf(symb)),
			Locked: types.Fixed8(locked),
			Frozen: types.Fixed8(frozen),
		})
	}
	return bals, nil
}

func (c *HTTP) GetBalance(addr types.AccAddress, symbol string) (*types.TokenBalance, error) {
	exist := c.existsCC(symbol)
	if !exist {
		return nil, errors.New("symbol not found")
	}
	acc, err := c.GetAccount(addr)
	if err != nil {
		return nil, err
	}
	var locked, frozen int64
	nacc := acc.(types.NamedAccount)
	if nacc != nil {
		locked = nacc.GetLockedCoins().AmountOf(symbol)
		frozen = nacc.GetFrozenCoins().AmountOf(symbol)
	}
	return &types.TokenBalance{
		Symbol: symbol,
		Free:   types.Fixed8(nacc.GetCoins().AmountOf(symbol)),
		Locked: types.Fixed8(locked),
		Frozen: types.Fixed8(frozen),
	}, nil
}

func (c *HTTP) GetFee() ([]types.FeeParam, error) {
	rawFee, err := c.ABCIQuery(fmt.Sprintf("%s/fees", ParamABCIPrefix), nil)
	if err != nil {
		return nil, err
	}
	var fees []types.FeeParam
	err = c.cdc.UnmarshalBinaryLengthPrefixed(rawFee.Response.GetValue(), &fees)
	return fees, err
}

func (c *HTTP) GetOpenOrders(addr types.AccAddress, pair string) ([]types.OpenOrder, error) {
	rawOrders, err := c.ABCIQuery(fmt.Sprintf("dex/openorders/%s/%s", pair, addr), nil)
	if err != nil {
		return nil, err
	}
	bz := rawOrders.Response.GetValue()
	openOrders := make([]types.OpenOrder, 0)
	if bz == nil {
		return openOrders, nil
	}
	if err := c.cdc.UnmarshalBinaryLengthPrefixed(bz, &openOrders); err != nil {
		return nil, err
	} else {
		return openOrders, nil
	}
}

func (c *HTTP) GetTradingPairs(offset int, limit int) ([]types.TradingPair, error) {
	rawTradePairs, err := c.ABCIQuery(fmt.Sprintf("dex/pairs/%d/%d", offset, limit), nil)
	if err != nil {
		return nil, err
	}
	pairs := make([]types.TradingPair, 0)
	if rawTradePairs.Response.GetValue() == nil {
		return pairs, nil
	}
	err = c.cdc.UnmarshalBinaryLengthPrefixed(rawTradePairs.Response.GetValue(), &pairs)
	return pairs, err
}

func (c *HTTP) GetDepth(tradePair string) (*types.OrderBook, error) {
    rawDepth, err := c.ABCIQuery(fmt.Sprintf("dex/orderbook/%s", tradePair), nil)
	if err != nil {
		return nil, err
	}
	var ob types.OrderBook
	err = c.cdc.UnmarshalBinaryLengthPrefixed(rawDepth.Response.GetValue(), &ob)
	if err != nil {
		return nil, err
	}
	return &ob, nil
}



func (c *HTTP) existsCC(symbol string) bool {
	key := []byte(strings.ToUpper(symbol))
	bz, err := c.QueryStore(key, TokenStoreName)
	if err != nil {
		return false
	}
	if bz != nil {
		return true
	}
	return false
}
