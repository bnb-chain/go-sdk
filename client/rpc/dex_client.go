package rpc

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/binance-chain/go-sdk/common"
	"github.com/binance-chain/go-sdk/common/types"
	sdk "github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/keys"
	gtypes "github.com/binance-chain/go-sdk/types"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
)

type SyncType int

const (
	Async SyncType = iota
	Sync
	Commit
)

const (
	AccountStoreName    = "acc"
	OracleStoreName     = "oracle"
	SideChainStoreName  = "sc"
	BridgeStoreName     = "bridge"
	ParamABCIPrefix     = "param"
	TimeLockMsgRoute    = "timelock"
	AtomicSwapStoreName = "atomic_swap"

	TimeLockrcNotFoundErrorCode = 458760
)

type DexClient interface {
	Broadcast(m msg.Msg, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error)
	TxInfoSearch(query string, prove bool, page, perPage int) ([]Info, error)
	ListAllTokens(offset int, limit int) ([]types.Token, error)
	GetTokenInfo(symbol string) (*types.Token, error)
	GetAccount(addr types.AccAddress) (acc types.Account, err error)
	GetCommitAccount(addr types.AccAddress) (acc types.Account, err error)

	GetBalances(addr types.AccAddress) ([]types.TokenBalance, error)
	GetBalance(addr types.AccAddress, symbol string) (*types.TokenBalance, error)
	GetFee() ([]types.FeeParam, error)
	GetOpenOrders(addr types.AccAddress, pair string) ([]types.OpenOrder, error)
	GetTradingPairs(offset int, limit int) ([]types.TradingPair, error)
	GetDepth(tradePair string, level int) (*types.OrderBook, error)
	GetProposals(status types.ProposalStatus, numLatest int64) ([]types.Proposal, error)
	GetSideChainProposals(status types.ProposalStatus, numLatest int64, sideChainId string) ([]types.Proposal, error)
	GetSideChainProposal(proposalId int64, sideChainId string) (types.Proposal, error)
	GetProposal(proposalId int64) (types.Proposal, error)
	GetTimelocks(addr types.AccAddress) ([]types.TimeLockRecord, error)
	GetTimelock(addr types.AccAddress, recordID int64) (*types.TimeLockRecord, error)
	GetSwapByID(swapID types.SwapBytes) (types.AtomicSwap, error)
	GetSwapByCreator(creatorAddr string, offset int64, limit int64) ([]types.SwapBytes, error)
	GetSwapByRecipient(recipientAddr string, offset int64, limit int64) ([]types.SwapBytes, error)
	GetSideChainParams(sideChainId string) ([]msg.SCParam, error)

	ListAllMiniTokens(offset int, limit int) ([]types.MiniToken, error)
	GetMiniTokenInfo(symbol string) (*types.MiniToken, error)
	GetMiniTradingPairs(offset int, limit int) ([]types.TradingPair, error)

	SetKeyManager(k keys.KeyManager)
	SendToken(transfers []msg.Transfer, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error)
	CreateOrder(baseAssetSymbol, quoteAssetSymbol string, op int8, price, quantity int64, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error)
	CancelOrder(baseAssetSymbol, quoteAssetSymbol, refId string, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error)
	HTLT(recipient types.AccAddress, recipientOtherChain, senderOtherChain string, randomNumberHash []byte, timestamp int64,
		amount types.Coins, expectedIncome string, heightSpan int64, crossChain bool, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error)
	DepositHTLT(recipient types.AccAddress, swapID []byte, amount types.Coins,
		syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error)
	ClaimHTLT(swapID []byte, randomNumber []byte, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error)
	RefundHTLT(swapID []byte, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error)
	TransferTokenOwnership(symbol string, newOwner types.AccAddress, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error)

	Bind(symbol string, amount int64, contractAddress msg.SmartChainAddress, contractDecimals int8, expireTime int64, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error)
	Unbind(symbol string, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error)
	TransferOut(to msg.SmartChainAddress, amount types.Coin, expireTime int64, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error)

	Claim(chainId sdk.IbcChainID, sequence uint64, payload []byte, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error)
	GetProphecy(chainId sdk.IbcChainID, sequence int64) (*msg.Prophecy, error)
	GetCurrentOracleSequence(chainId sdk.IbcChainID) (int64, error)

	SideChainVote(proposalID int64, option msg.VoteOption, sideChainId string, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error)
	SideChainDeposit(proposalID int64, amount types.Coins, sideChainId string, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error)
	SideChainSubmitSCParamsProposal(title string, scParam msg.SCChangeParams, initialDeposit types.Coins, votingPeriod time.Duration, sideChainId string, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error)
	SideChainSubmitCSCParamsProposal(title string, cscParam msg.CSCParamChange, initialDeposit types.Coins, votingPeriod time.Duration, sideChainId string, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error)
	SideChainSubmitProposal(title string, description string, proposalType msg.ProposalKind, initialDeposit types.Coins, votingPeriod time.Duration, sideChainId string, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error)
	SubmitListProposal(title string, param msg.ListTradingPairParams, proposalType msg.ProposalKind, initialDeposit types.Coins, votingPeriod time.Duration, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error)

	SubmitProposal(title string, description string, proposalType msg.ProposalKind, initialDeposit types.Coins, votingPeriod time.Duration, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error)
	Deposit(proposalID int64, amount types.Coins, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error)
	Vote(proposalID int64, option msg.VoteOption, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error)
}

func (c *HTTP) TxInfoSearch(query string, prove bool, page, perPage int) ([]Info, error) {
	if err := ValidateTxSearchQueryStr(query); err != nil {
		return nil, err
	}
	return c.WSEvents.TxInfoSearch(query, prove, page, perPage)
}

func (c *HTTP) ListAllTokens(offset int, limit int) ([]types.Token, error) {
	if err := ValidateOffset(offset); err != nil {
		return nil, err
	}
	if err := ValidateLimit(limit); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("tokens/list/%d/%d", offset, limit)
	result, err := c.ABCIQuery(path, nil)
	if err != nil {
		return nil, err
	}
	if !result.Response.IsOK() {
		return nil, fmt.Errorf(result.Response.Log)
	}
	bz := result.Response.GetValue()
	tokens := make([]types.Token, 0)
	err = c.cdc.UnmarshalBinaryLengthPrefixed(bz, &tokens)
	return tokens, err
}

func (c *HTTP) GetTokenInfo(symbol string) (*types.Token, error) {
	if err := ValidateSymbol(symbol); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("tokens/info/%s", symbol)
	result, err := c.ABCIQuery(path, nil)
	if err != nil {
		return nil, err
	}
	if !result.Response.IsOK() {
		return nil, fmt.Errorf(result.Response.Log)
	}
	bz := result.Response.GetValue()
	token := new(types.Token)
	err = c.cdc.UnmarshalBinaryLengthPrefixed(bz, token)
	return token, err
}

// Always fetch the account from the commit store at (currentHeight-1) in node.
// example:
// 1. currentCommitHeight: 1000, accountA(balance: 10BNB, sequence: 10)
// 2. height: 999, accountA(balance: 11BNB, sequence: 9)
// 3. GetCommitAccount will return accountA(balance: 11BNB, sequence: 9).
// If the (currentHeight-1) do not exist, will return  account at currentHeight.
// 1. currentCommitHeight: 1000, accountA(balance: 10BNB, sequence: 10)
// 2. height: 999. the state do not exist
// 3. GetCommitAccount will return accountA(balance: 10BNB, sequence: 10).
func (c *HTTP) GetCommitAccount(addr types.AccAddress) (acc types.Account, err error) {
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

// Always fetch the latest account from the cache in node even the previous transaction from
// this account is not included in block yet.
// example:
// 1. AccountA(Balance: 10BNB, sequence: 1), AccountB(Balance: 5BNB, sequence: 1)
// 2. Node receive Tx(AccountA --> AccountB 2BNB) and check have passed, but not included in block yet.
// 3. GetAccount will return AccountA(Balance: 8BNB, sequence: 2), AccountB(Balance: 7BNB, sequence: 1)
func (c *HTTP) GetAccount(addr types.AccAddress) (acc types.Account, err error) {
	result, err := c.ABCIQuery(fmt.Sprintf("/account/%s", addr.String()), nil)
	if err != nil {
		return nil, err
	}
	resp := result.Response
	if !resp.IsOK() {
		return nil, errors.New(resp.Log)
	}
	value := result.Response.GetValue()
	if len(value) == 0 {
		return nil, nil
	}
	err = c.cdc.UnmarshalBinaryBare(value, &acc)
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
	if account == nil {
		return []types.TokenBalance{}, nil
	}
	coins := account.GetCoins()

	symbs := make([]string, 0, len(coins))
	bals := make([]types.TokenBalance, 0, len(coins))
	for _, coin := range coins {
		symbs = append(symbs, coin.Denom)
		// count locked and frozen coins
		var locked, frozen int64
		nacc := account.(types.NamedAccount)
		if nacc != nil {
			locked = nacc.GetLockedCoins().AmountOf(coin.Denom)
			frozen = nacc.GetFrozenCoins().AmountOf(coin.Denom)
		}
		bals = append(bals, types.TokenBalance{
			Symbol: coin.Denom,
			Free:   types.Fixed8(coins.AmountOf(coin.Denom)),
			Locked: types.Fixed8(locked),
			Frozen: types.Fixed8(frozen),
		})
	}
	return bals, nil
}

func (c *HTTP) GetBalance(addr types.AccAddress, symbol string) (*types.TokenBalance, error) {
	if err := ValidateSymbol(symbol); err != nil {
		return nil, err
	}
	exist := c.existsCC(symbol)
	if !exist {
		return nil, errors.New("symbol not found")
	}
	acc, err := c.GetAccount(addr)
	if err != nil {
		return nil, err
	}
	if acc == nil {
		return &types.TokenBalance{
			Symbol: symbol,
			Free:   types.Fixed8Zero,
			Locked: types.Fixed8Zero,
			Frozen: types.Fixed8Zero,
		}, nil
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
	if !rawFee.Response.IsOK() {
		return nil, fmt.Errorf(rawFee.Response.Log)
	}
	var fees []types.FeeParam
	err = c.cdc.UnmarshalBinaryLengthPrefixed(rawFee.Response.GetValue(), &fees)
	return fees, err
}

func (c *HTTP) GetOpenOrders(addr types.AccAddress, pair string) ([]types.OpenOrder, error) {
	if err := ValidatePair(pair); err != nil {
		return nil, err
	}
	rawOrders, err := c.ABCIQuery(fmt.Sprintf("dex/openorders/%s/%s", pair, addr), nil)
	if err != nil {
		return nil, err
	}
	if !rawOrders.Response.IsOK() {
		return nil, fmt.Errorf(rawOrders.Response.Log)
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
	if err := ValidateLimit(limit); err != nil {
		return nil, err
	}
	if err := ValidateOffset(offset); err != nil {
		return nil, err
	}
	rawTradePairs, err := c.ABCIQuery(fmt.Sprintf("dex/pairs/%d/%d", offset, limit), nil)
	if err != nil {
		return nil, err
	}
	if !rawTradePairs.Response.IsOK() {
		return nil, fmt.Errorf(rawTradePairs.Response.Log)
	}
	pairs := make([]types.TradingPair, 0)
	if rawTradePairs.Response.GetValue() == nil {
		return pairs, nil
	}
	err = c.cdc.UnmarshalBinaryLengthPrefixed(rawTradePairs.Response.GetValue(), &pairs)
	return pairs, err
}

func (c *HTTP) GetDepth(tradePair string, level int) (*types.OrderBook, error) {
	if err := ValidatePair(tradePair); err != nil {
		return nil, err
	}
	if err := ValidateDepthLevel(level); err != nil {
		return nil, err
	}
	rawDepth, err := c.ABCIQuery(fmt.Sprintf("dex/orderbook/%s/%d", tradePair, level), nil)
	if err != nil {
		return nil, err
	}
	if !rawDepth.Response.IsOK() {
		return nil, fmt.Errorf(rawDepth.Response.Log)
	}
	var ob types.OrderBook
	err = c.cdc.UnmarshalBinaryLengthPrefixed(rawDepth.Response.GetValue(), &ob)
	if err != nil {
		return nil, err
	}
	return &ob, nil
}

func (c *HTTP) GetTimelocks(addr types.AccAddress) ([]types.TimeLockRecord, error) {

	params := types.QueryTimeLocksParams{
		Account: addr,
	}

	bz, err := c.cdc.MarshalJSON(params)

	if err != nil {
		fmt.Errorf("marshal params failed %v", err)
	}

	rawRecords, err := c.ABCIQuery(fmt.Sprintf("custom/%s/%s", TimeLockMsgRoute, "timelocks"), bz)

	if err != nil {
		return nil, err
	}
	if rawRecords == nil {
		return nil, fmt.Errorf("zero records")
	}
	if !rawRecords.Response.IsOK() {
		return nil, fmt.Errorf(rawRecords.Response.Log)
	}
	records := make([]types.TimeLockRecord, 0)

	if err = c.cdc.UnmarshalJSON(rawRecords.Response.GetValue(), &records); err != nil {
		return nil, err
	} else {
		return records, nil
	}

}

func (c *HTTP) GetTimelock(addr types.AccAddress, recordID int64) (*types.TimeLockRecord, error) {

	params := types.QueryTimeLockParams{
		Account: addr,
		Id:      recordID,
	}

	bz, err := c.cdc.MarshalJSON(params)

	if err != nil {
		return nil, fmt.Errorf("incorrectly formatted request data %s", err.Error())
	}

	rawRecord, err := c.ABCIQuery(fmt.Sprintf("custom/%s/%s", TimeLockMsgRoute, "timelock"), bz)

	if err != nil {
		return nil, fmt.Errorf("error query %s", err.Error())
	}
	if rawRecord.Response.Code == TimeLockrcNotFoundErrorCode {
		return nil, nil
	}
	var record types.TimeLockRecord

	err = c.cdc.UnmarshalJSON(rawRecord.Response.GetValue(), &record)
	if err != nil {
		return nil, err
	}

	return &record, nil

}

func (c *HTTP) GetProposals(status types.ProposalStatus, numLatest int64) ([]types.Proposal, error) {
	return c.getProposals(status, "", numLatest)
}

func (c *HTTP) GetSideChainProposals(status types.ProposalStatus, numLatest int64, sideChainId string) ([]types.Proposal, error) {
	return c.getProposals(status, sideChainId, numLatest)
}

func (c *HTTP) getProposals(status types.ProposalStatus, sideChainId string, numLatest int64) ([]types.Proposal, error) {
	params := types.QueryProposalsParams{}
	if status != types.StatusNil {
		params.ProposalStatus = status
	}
	if numLatest > 0 {
		params.NumLatestProposals = numLatest
	}
	params.SideChainId = sideChainId

	bz, err := c.cdc.MarshalJSON(&params)
	if err != nil {
		return nil, err
	}
	rawProposals, err := c.ABCIQuery("custom/gov/proposals", bz)
	if err != nil {
		return nil, err
	}
	if !rawProposals.Response.IsOK() {
		return nil, fmt.Errorf(rawProposals.Response.Log)
	}
	proposals := make([]types.Proposal, 0)

	err = c.cdc.UnmarshalJSON(rawProposals.Response.GetValue(), &proposals)
	return proposals, err
}

func (c *HTTP) GetProposal(proposalId int64) (types.Proposal, error) {
	return c.getProposal(proposalId, "")
}

func (c *HTTP) GetSideChainProposal(proposalId int64, sideChainId string) (types.Proposal, error) {
	return c.getProposal(proposalId, sideChainId)
}

func (c *HTTP) getProposal(proposalId int64, sideChainId string) (types.Proposal, error) {
	params := types.QueryProposalParams{
		ProposalID: proposalId,
	}
	params.SideChainId = sideChainId
	bz, err := c.cdc.MarshalJSON(params)
	if err != nil {
		return nil, err
	}
	rawProposal, err := c.ABCIQuery("custom/gov/proposal", bz)
	if err != nil {
		return nil, err
	}
	if !rawProposal.Response.IsOK() {
		return nil, fmt.Errorf(rawProposal.Response.Log)
	}
	var proposal types.Proposal

	err = c.cdc.UnmarshalJSON(rawProposal.Response.GetValue(), &proposal)
	return proposal, err
}

func (c *HTTP) GetSideChainParams(sideChainId string) ([]msg.SCParam, error) {
	data, err := c.cdc.MarshalJSON(sideChainId)
	if err != nil {
		return nil, err
	}
	rawParams, err := c.ABCIQuery("param/sideParams", data)
	if err != nil {
		return nil, err
	}
	if !rawParams.Response.IsOK() {
		return nil, fmt.Errorf(rawParams.Response.Log)
	}
	var params []msg.SCParam
	err = c.cdc.UnmarshalJSON(rawParams.Response.GetValue(), &params)
	return params, err
}

func (c *HTTP) existsCC(symbol string) bool {
	resp, err := c.ABCIQuery(fmt.Sprintf("tokens/info/%s", symbol), nil)
	if err != nil {
		return false
	}
	if !resp.Response.IsOK() {
		return false
	}
	if len(resp.Response.GetValue()) == 0 {
		return false
	}
	var token types.Token
	err = c.cdc.UnmarshalBinaryLengthPrefixed(resp.Response.GetValue(), &token)
	if err != nil {
		return false
	}
	return true
}

func (c *HTTP) GetSwapByID(swapID types.SwapBytes) (types.AtomicSwap, error) {
	params := types.QuerySwapByID{
		SwapID: swapID,
	}
	bz, err := c.cdc.MarshalJSON(params)
	if err != nil {
		return types.AtomicSwap{}, err
	}

	resp, err := c.ABCIQuery(fmt.Sprintf("custom/%s/%s", msg.AtomicSwapRoute, "swapid"), bz)
	if err != nil {
		return types.AtomicSwap{}, err
	}
	if !resp.Response.IsOK() {
		return types.AtomicSwap{}, fmt.Errorf(resp.Response.Log)
	}
	if len(resp.Response.GetValue()) == 0 {
		return types.AtomicSwap{}, fmt.Errorf("zero records")
	}
	var result types.AtomicSwap
	err = c.cdc.UnmarshalJSON(resp.Response.GetValue(), &result)
	if err != nil {
		return types.AtomicSwap{}, err
	}
	return result, nil
}

func (c *HTTP) GetSwapByCreator(creatorAddr string, offset int64, limit int64) ([]types.SwapBytes, error) {
	addr, err := types.AccAddressFromBech32(creatorAddr)
	if err != nil {
		return nil, err
	}

	params := types.QuerySwapByCreatorParams{
		Creator: addr,
		Limit:   limit,
		Offset:  offset,
	}

	bz, err := c.cdc.MarshalJSON(params)
	if err != nil {
		return nil, err
	}

	resp, err := c.ABCIQuery(fmt.Sprintf("custom/%s/%s", msg.AtomicSwapRoute, "swapcreator"), bz)
	if err != nil {
		return nil, err
	}
	if !resp.Response.IsOK() {
		return nil, fmt.Errorf(resp.Response.Log)
	}
	if len(resp.Response.GetValue()) == 0 {
		return nil, fmt.Errorf("zero records")
	}
	var swapIDList []types.SwapBytes
	err = c.cdc.UnmarshalJSON(resp.Response.GetValue(), &swapIDList)
	if err != nil {
		return nil, err
	}
	return swapIDList, nil
}

func (c *HTTP) GetSwapByRecipient(recipientAddr string, offset int64, limit int64) ([]types.SwapBytes, error) {
	recipient, err := types.AccAddressFromBech32(recipientAddr)
	if err != nil {
		return nil, err
	}
	params := types.QuerySwapByRecipientParams{
		Recipient: recipient,
		Limit:     limit,
		Offset:    offset,
	}

	bz, err := c.cdc.MarshalJSON(params)
	if err != nil {
		return nil, err
	}

	resp, err := c.ABCIQuery(fmt.Sprintf("custom/%s/%s", msg.AtomicSwapRoute, "swaprecipient"), bz)
	if err != nil {
		return nil, err
	}
	if !resp.Response.IsOK() {
		return nil, fmt.Errorf(resp.Response.Log)
	}
	if len(resp.Response.GetValue()) == 0 {
		return nil, fmt.Errorf("zero records")
	}
	var swapIDList []types.SwapBytes
	err = c.cdc.UnmarshalJSON(resp.Response.GetValue(), &swapIDList)
	if err != nil {
		return nil, err
	}
	return swapIDList, nil
}

func (c *HTTP) ListAllMiniTokens(offset int, limit int) ([]types.MiniToken, error) {
	if err := ValidateOffset(offset); err != nil {
		return nil, err
	}
	if err := ValidateLimit(limit); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("mini-tokens/list/%d/%d", offset, limit)
	result, err := c.ABCIQuery(path, nil)
	if err != nil {
		return nil, err
	}
	if !result.Response.IsOK() {
		return nil, fmt.Errorf(result.Response.Log)
	}
	bz := result.Response.GetValue()
	tokens := make([]types.MiniToken, 0)
	err = c.cdc.UnmarshalBinaryLengthPrefixed(bz, &tokens)
	return tokens, err
}

func (c *HTTP) GetMiniTokenInfo(symbol string) (*types.MiniToken, error) {
	if err := ValidateSymbol(symbol); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("mini-tokens/info/%s", symbol)
	result, err := c.ABCIQuery(path, nil)
	if err != nil {
		return nil, err
	}
	if !result.Response.IsOK() {
		return nil, fmt.Errorf(result.Response.Log)
	}
	bz := result.Response.GetValue()
	token := new(types.MiniToken)
	err = c.cdc.UnmarshalBinaryLengthPrefixed(bz, token)
	return token, err
}

func (c *HTTP) GetMiniTradingPairs(offset int, limit int) ([]types.TradingPair, error) {
	if err := ValidateLimit(limit); err != nil {
		return nil, err
	}
	if err := ValidateOffset(offset); err != nil {
		return nil, err
	}
	rawTradePairs, err := c.ABCIQuery(fmt.Sprintf("dex-mini/pairs/%d/%d", offset, limit), nil)
	if err != nil {
		return nil, err
	}
	if !rawTradePairs.Response.IsOK() {
		return nil, fmt.Errorf(rawTradePairs.Response.Log)
	}
	pairs := make([]types.TradingPair, 0)
	if rawTradePairs.Response.GetValue() == nil {
		return pairs, nil
	}
	err = c.cdc.UnmarshalBinaryLengthPrefixed(rawTradePairs.Response.GetValue(), &pairs)
	return pairs, err
}

func (c *HTTP) SendToken(transfers []msg.Transfer, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyMissingError
	}
	fromAddr := c.key.GetAddr()
	fromCoins := types.Coins{}
	for _, t := range transfers {
		t.Coins = t.Coins.Sort()
		fromCoins = fromCoins.Plus(t.Coins)
	}
	sendMsg := msg.CreateSendMsg(fromAddr, fromCoins, transfers)
	return c.Broadcast(sendMsg, syncType, options...)

}

func (c *HTTP) CreateOrder(baseAssetSymbol, quoteAssetSymbol string, op int8, price, quantity int64, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyMissingError
	}
	if baseAssetSymbol == "" || quoteAssetSymbol == "" {
		return nil, fmt.Errorf("BaseAssetSymbol or QuoteAssetSymbol is missing. ")
	}
	fromAddr := c.key.GetAddr()
	newOrderMsg := msg.NewCreateOrderMsg(
		fromAddr,
		"",
		op,
		common.CombineSymbol(baseAssetSymbol, quoteAssetSymbol),
		price,
		quantity,
	)
	return c.Broadcast(newOrderMsg, syncType, options...)
}

func (c *HTTP) HTLT(recipient types.AccAddress, recipientOtherChain, senderOtherChain string, randomNumberHash []byte, timestamp int64,
	amount types.Coins, expectedIncome string, heightSpan int64, crossChain bool, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyMissingError
	}
	fromAddr := c.key.GetAddr()
	htltMsg := msg.NewHTLTMsg(
		fromAddr,
		recipient,
		recipientOtherChain,
		senderOtherChain,
		randomNumberHash,
		timestamp,
		amount,
		expectedIncome,
		heightSpan,
		crossChain,
	)
	return c.Broadcast(htltMsg, syncType, options...)
}

func (c *HTTP) DepositHTLT(recipient types.AccAddress, swapID []byte, amount types.Coins,
	syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyMissingError
	}
	fromAddr := c.key.GetAddr()
	depositHTLTMsg := msg.NewDepositHTLTMsg(
		fromAddr,
		swapID,
		amount,
	)
	return c.Broadcast(depositHTLTMsg, syncType, options...)
}

func (c *HTTP) ClaimHTLT(swapID []byte, randomNumber []byte, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyMissingError
	}
	fromAddr := c.key.GetAddr()
	claimHTLTMsg := msg.NewClaimHTLTMsg(
		fromAddr,
		swapID,
		randomNumber,
	)
	return c.Broadcast(claimHTLTMsg, syncType, options...)
}

func (c *HTTP) RefundHTLT(swapID []byte, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyMissingError
	}
	fromAddr := c.key.GetAddr()
	refundHTLTMsg := msg.NewRefundHTLTMsg(
		fromAddr,
		swapID,
	)
	return c.Broadcast(refundHTLTMsg, syncType, options...)
}

func (c *HTTP) TransferTokenOwnership(symbol string, newOwner types.AccAddress, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyMissingError
	}
	fromAddr := c.key.GetAddr()
	transferOwnershipMsg := msg.NewTransferOwnershipMsg(fromAddr, symbol, newOwner)
	return c.Broadcast(transferOwnershipMsg, syncType, options...)
}

func (c *HTTP) SideChainVote(proposalID int64, option msg.VoteOption, sideChainId string, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyMissingError
	}
	fromAddr := c.key.GetAddr()
	msg := msg.NewSideChainVoteMsg(fromAddr, proposalID, option, sideChainId)
	return c.Broadcast(msg, syncType, options...)
}

func (c *HTTP) SideChainDeposit(proposalID int64, amount types.Coins, sideChainId string, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyMissingError
	}
	fromAddr := c.key.GetAddr()
	msg := msg.NewSideChainDepositMsg(fromAddr, proposalID, amount, sideChainId)
	return c.Broadcast(msg, syncType, options...)
}

func (c *HTTP) SideChainSubmitSCParamsProposal(title string, scParam msg.SCChangeParams, initialDeposit types.Coins, votingPeriod time.Duration, sideChainId string, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyMissingError
	}
	err := scParam.Check()
	if err != nil {
		return nil, err
	}
	scParamsBz, err := c.cdc.MarshalJSON(scParam)
	if err != nil {
		return nil, err
	}
	fromAddr := c.key.GetAddr()
	msg := msg.NewSideChainSubmitProposalMsg(title, string(scParamsBz), msg.ProposalTypeSCParamsChange, fromAddr, initialDeposit, votingPeriod, sideChainId)
	return c.Broadcast(msg, syncType, options...)
}

func (c *HTTP) SideChainSubmitCSCParamsProposal(title string, cscParam msg.CSCParamChange, initialDeposit types.Coins, votingPeriod time.Duration, sideChainId string, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyMissingError
	}
	err := cscParam.Check()
	if err != nil {
		return nil, err
	}
	// cscParam get interface field, use amino
	cscParamsBz, err := c.cdc.MarshalJSON(cscParam)
	if err != nil {
		return nil, err
	}
	fromAddr := c.key.GetAddr()
	msg := msg.NewSideChainSubmitProposalMsg(title, string(cscParamsBz), msg.ProposalTypeCSCParamsChange, fromAddr, initialDeposit, votingPeriod, sideChainId)
	return c.Broadcast(msg, syncType, options...)
}

func (c *HTTP) SideChainSubmitProposal(title string, description string, proposalType msg.ProposalKind, initialDeposit types.Coins, votingPeriod time.Duration, sideChainId string, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyMissingError
	}
	fromAddr := c.key.GetAddr()
	msg := msg.NewSideChainSubmitProposalMsg(title, description, proposalType, fromAddr, initialDeposit, votingPeriod, sideChainId)
	return c.Broadcast(msg, syncType, options...)
}

func (c *HTTP) SubmitListProposal(title string, param msg.ListTradingPairParams, proposalType msg.ProposalKind, initialDeposit types.Coins, votingPeriod time.Duration, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyMissingError
	}
	bz, err := json.Marshal(&param)
	if err != nil {
		return nil, err
	}
	fromAddr := c.key.GetAddr()
	msg := msg.NewMsgSubmitProposal(title, string(bz), proposalType, fromAddr, initialDeposit, votingPeriod)
	return c.Broadcast(msg, syncType, options...)
}

func (c *HTTP) SubmitProposal(title string, description string, proposalType msg.ProposalKind, initialDeposit types.Coins, votingPeriod time.Duration, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyMissingError
	}
	fromAddr := c.key.GetAddr()
	msg := msg.NewMsgSubmitProposal(title, description, proposalType, fromAddr, initialDeposit, votingPeriod)
	return c.Broadcast(msg, syncType, options...)
}

func (c *HTTP) Deposit(proposalID int64, amount types.Coins, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyMissingError
	}
	fromAddr := c.key.GetAddr()
	msg := msg.NewDepositMsg(fromAddr, proposalID, amount)
	return c.Broadcast(msg, syncType, options...)
}

func (c *HTTP) Vote(proposalID int64, option msg.VoteOption, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyMissingError
	}
	fromAddr := c.key.GetAddr()
	msg := msg.NewMsgVote(fromAddr, proposalID, option)
	return c.Broadcast(msg, syncType, options...)
}

func (c *HTTP) CancelOrder(baseAssetSymbol, quoteAssetSymbol, refId string, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyMissingError
	}
	if baseAssetSymbol == "" || quoteAssetSymbol == "" {
		return nil, fmt.Errorf("BaseAssetSymbol or QuoteAssetSymbol is missing. ")
	}
	if refId == "" {
		return nil, fmt.Errorf("OrderId or Order RefId is missing. ")
	}

	fromAddr := c.key.GetAddr()

	cancelOrderMsg := msg.NewCancelOrderMsg(fromAddr, common.CombineSymbol(baseAssetSymbol, quoteAssetSymbol), refId)
	return c.Broadcast(cancelOrderMsg, syncType, options...)
}

func (c *HTTP) TransferOut(to msg.SmartChainAddress, amount types.Coin, expireTime int64, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyMissingError
	}

	fromAddr := c.key.GetAddr()

	transferOutMsg := msg.NewTransferOutMsg(fromAddr, to, amount, expireTime)

	return c.Broadcast(transferOutMsg, syncType, options...)
}

func (c *HTTP) Bind(symbol string, amount int64, contractAddress msg.SmartChainAddress, contractDecimals int8, expireTime int64, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyMissingError
	}

	fromAddr := c.key.GetAddr()

	bindMsg := msg.NewBindMsg(fromAddr, symbol, amount, contractAddress, contractDecimals, expireTime)

	return c.Broadcast(bindMsg, syncType, options...)
}

func (c *HTTP) Unbind(symbol string, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyMissingError
	}

	fromAddr := c.key.GetAddr()

	unbindMsg := msg.NewUnbindMsg(fromAddr, symbol)
	return c.Broadcast(unbindMsg, syncType, options...)
}

func (c *HTTP) Broadcast(m msg.Msg, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error) {
	signBz, err := c.sign(m, options...)
	if err != nil {
		return nil, err
	}
	switch syncType {
	case Async:
		return c.BroadcastTxAsync(signBz)
	case Sync:
		return c.BroadcastTxSync(signBz)
	case Commit:
		commitRes, err := c.BroadcastTxCommit(signBz)
		if err != nil {
			return nil, err
		}
		if commitRes.CheckTx.IsErr() {
			return &core_types.ResultBroadcastTx{
				Code: commitRes.CheckTx.Code,
				Log:  commitRes.CheckTx.Log,
				Hash: commitRes.Hash,
				Data: commitRes.CheckTx.Data,
			}, nil
		}
		return &core_types.ResultBroadcastTx{
			Code: commitRes.DeliverTx.Code,
			Log:  commitRes.DeliverTx.Log,
			Hash: commitRes.Hash,
			Data: commitRes.DeliverTx.Data,
		}, nil
	default:
		return nil, fmt.Errorf("unknown synctype")
	}
}

func (c *HTTP) Claim(chainId sdk.IbcChainID, sequence uint64, payload []byte, syncType SyncType, options ...tx.Option) (*core_types.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyMissingError
	}

	fromAddr := c.key.GetAddr()

	claimMsg := msg.NewClaimMsg(chainId, sequence, payload, fromAddr)

	return c.Broadcast(claimMsg, syncType, options...)
}

func (c *HTTP) GetProphecy(chainId sdk.IbcChainID, sequence int64) (*msg.Prophecy, error) {
	key := []byte(msg.GetClaimId(chainId, msg.OracleChannelId, sequence))
	bz, err := c.QueryStore(key, OracleStoreName)
	if err != nil {
		return nil, err
	}
	if bz == nil {
		return nil, nil
	}

	dbProphecy := new(msg.DBProphecy)
	err = c.cdc.UnmarshalBinaryBare(bz, &dbProphecy)
	if err != nil {
		return nil, err
	}

	prophecy, err := dbProphecy.DeserializeFromDB()
	if err != nil {
		return nil, err
	}

	return &prophecy, err
}

func (c *HTTP) GetCurrentOracleSequence(chainId sdk.IbcChainID) (int64, error) {
	key := types.GetReceiveSequenceKey(chainId, msg.OracleChannelId)
	bz, err := c.QueryStore(key, SideChainStoreName)
	if err != nil {
		return 0, err
	}
	if bz == nil {
		return 0, nil
	}

	sequence := binary.BigEndian.Uint64(bz)
	return int64(sequence), err
}

func (c *HTTP) sign(m msg.Msg, options ...tx.Option) ([]byte, error) {
	if c.key == nil {
		return nil, fmt.Errorf("keymanager is missing, use SetKeyManager to set key")
	}
	// prepare message to sign
	chainID := gtypes.ProdChainID
	if types.Network == types.TestNetwork {
		chainID = gtypes.TestnetChainID
	} else if types.Network == types.TmpTestNetwork {
		chainID = gtypes.KongoChainId
	} else if types.Network == types.GangesNetwork {
		chainID = gtypes.GangesChainId
	}
	signMsg := &tx.StdSignMsg{
		ChainID:       chainID,
		AccountNumber: -1,
		Sequence:      -1,
		Memo:          "",
		Msgs:          []msg.Msg{m},
		Source:        tx.Source,
	}

	for _, op := range options {
		signMsg = op(signMsg)
	}

	if signMsg.Sequence == -1 || signMsg.AccountNumber == -1 {
		fromAddr := c.key.GetAddr()
		acc, err := c.GetAccount(fromAddr)
		if err != nil {
			return nil, err
		}
		if acc == nil {
			return nil, fmt.Errorf("the signer account do not exist in the chain")
		}
		signMsg.Sequence = acc.GetSequence()
		signMsg.AccountNumber = acc.GetAccountNumber()
	}

	// special logic for createOrder, to save account query
	if orderMsg, ok := m.(msg.CreateOrderMsg); ok {
		orderMsg.ID = msg.GenerateOrderID(signMsg.Sequence+1, c.key.GetAddr())
		signMsg.Msgs[0] = orderMsg
	}

	for _, m := range signMsg.Msgs {
		if err := m.ValidateBasic(); err != nil {
			return nil, err
		}
	}
	return c.key.Sign(*signMsg)
}
