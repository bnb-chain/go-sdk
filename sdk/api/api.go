package api

import (
	"fmt"
	"github.com/binance-chain/go-sdk/sdk/keys"
	"github.com/binance-chain/go-sdk/sdk/tx"
	"github.com/binance-chain/go-sdk/sdk/tx/txmsg"
	resty "gopkg.in/resty.v1"
	"net/http"
)

const (
	defaultAPIVersionPrefix = "/api/v1"
	NativeSymbol            = "BNB"

	GoSdkSource = 2
	SideBuy     = "BUY"
	SideSell    = "SELL"
)

// dexAPI wrapper
type dexAPI struct {
	keyManager keys.KeyManager
	apiUrl     string
	chainId    string
}

// IDexAPI methods
type IDexAPI interface {
	// Query api
	GetClosedOrders(query *ClosedOrdersQuery) (*CloseOrders, error)
	GetDepth(query *DepthQuery) (*MarketDepth, error)
	GetKlines(query *KlineQuery) ([]Kline, error)
	GetMarkets(query *MarketsQuery) ([]SymbolPair, error)
	GetOrder(orderID string) (*Order, error)
	GetOpenOrders(query *OpenOrdersQuery) (*OpenOrders, error)
	GetTicker24h(query *Ticker24hQuery) ([]Ticker24h, error)
	GetTrades(query *TradesQuery) (*Trades, error)
	GetAccount(string) (*Account, error)
	GetTime() (*Time, error)
	GetTokens() ([]Token, error)
	GetTx(txHash string) (*TxResult, error)

	//Transaction api
	CreateOrder(baseAssetSymbol, quoteAssetSymbol string, op int8, price, quantity int64, sync bool) (*CreateOrderResult, error)
	CancelOrder(baseAssetSymbol, quoteAssetSymbol, id, refId string, sync bool) (*CancelOrderResult, error)
	BurnToken(symbol string, amount int64, sync bool) (*BurnTokenResult, error)
	ListPair(proposalId int64, baseAssetSymbol string, quoteAssetSymbol string, initPrice int64, sync bool) (*ListPairResult, error)
	FreezeToken(symbol string, amount int64, sync bool) (*FreezeTokenResult, error)
	UnfreezeToken(symbol string, amount int64, sync bool) (*UnfreezeTokenResult, error)
	IssueToken(name, symbol string, supply int64, sync bool, mintable bool) (*IssueTokenResult, error)
	SendToken(dest txmsg.AccAddress, symbol string, quantity int64, sync bool) (*SendTokenResult, error)
	MintToken(symbol string, amount int64, sync bool) (*MintTokenResult, error)

	//Gov api
	SubmitListPairProposal(title string, param txmsg.ListTradingPairParams, initialDeposit int64, sync bool) (*SubmitProposalResult, error)
	SubmitProposal(title string, description string, proposalType txmsg.ProposalKind, initialDeposit int64, sync bool) (*SubmitProposalResult, error)
	DepositProposal(proposalID int64, amount int64, sync bool) (*DepositProposalResult, error)
	VoteProposal(proposalID int64, option txmsg.VoteOption, sync bool) (*VoteProposalResult, error)
}

func init() {
	resty.DefaultClient.SetRedirectPolicy(resty.FlexibleRedirectPolicy(10))
}
func NewDefaultDexApi(baseUrl string, chainId string, keyManager keys.KeyManager) IDexAPI {
	return &dexAPI{apiUrl: baseUrl + defaultAPIVersionPrefix, chainId: chainId, keyManager: keyManager}
}

// Get generic method
func (dex *dexAPI) Get(path string, qp map[string]string) ([]byte, error) {
	resp, err := resty.R().SetQueryParams(qp).Get(dex.apiUrl + path)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() >= http.StatusMultipleChoices || resp.StatusCode() < http.StatusOK {
		err = fmt.Errorf("bad response, status code %d, response: %s", resp.StatusCode(), string(resp.Body()))
	}
	return resp.Body(), err
}

// Post generic method
func (dex *dexAPI) Post(path string, body interface{}, param map[string]string) ([]byte, error) {
	resp, err := resty.R().
		SetHeader("Content-Type", "text/plain").
		SetBody(body).
		SetQueryParams(param).
		Post(dex.apiUrl + path)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode() >= http.StatusMultipleChoices {
		err = fmt.Errorf("bad response, status code %d, response: %s", resp.StatusCode(), string(resp.Body()))
	}
	return resp.Body(), err
}

func (dex *dexAPI) broadcastMsg(msg txmsg.Msg, sync bool) (*TxCommitResult, error) {
	fromAddr := dex.keyManager.GetAddr()
	acc, err := dex.GetAccount(fromAddr.String())
	if err != nil {
		return nil, err
	}
	sequence := acc.Sequence
	// prepare message to sign
	signMsg := tx.StdSignMsg{
		ChainID:       dex.chainId,
		AccountNumber: acc.Number,
		Sequence:      sequence,
		Memo:          "",
		Msgs:          []txmsg.Msg{msg},
		Source:        GoSdkSource,
	}

	// Hex encoded signed transaction, ready to be posted to BncChain API
	hexTx, err := dex.keyManager.Sign(signMsg)
	if err != nil {
		return nil, err
	}
	param := map[string]string{}
	if sync {
		param["sync"] = "true"
	}
	commits, err := dex.PostTx(hexTx, param)
	if err != nil {
		return nil, err
	}
	if len(commits) < 1 {
		return nil, fmt.Errorf("Len of tx Commit result is less than 1 ")
	}
	return &commits[0], nil
}
