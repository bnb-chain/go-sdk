package transaction

import (
	"fmt"
	"time"

	"github.com/binance-chain/go-sdk/client/basic"
	"github.com/binance-chain/go-sdk/client/query"
	"github.com/binance-chain/go-sdk/keys"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
)

type TransactionClient interface {
	CreateOrder(baseAssetSymbol, quoteAssetSymbol string, op int8, price, quantity int64, sync bool, memo string, source int64) (*CreateOrderResult, error)
	CancelOrder(baseAssetSymbol, quoteAssetSymbol, refId string, sync bool, memo string, source int64) (*CancelOrderResult, error)
	BurnToken(symbol string, amount int64, sync bool, memo string, source int64) (*BurnTokenResult, error)
	ListPair(proposalId int64, baseAssetSymbol string, quoteAssetSymbol string, initPrice int64, sync bool, memo string, source int64) (*ListPairResult, error)
	FreezeToken(symbol string, amount int64, sync bool, memo string, source int64) (*FreezeTokenResult, error)
	UnfreezeToken(symbol string, amount int64, sync bool, memo string, source int64) (*UnfreezeTokenResult, error)
	IssueToken(name, symbol string, supply int64, sync bool, mintable bool, memo string, source int64) (*IssueTokenResult, error)
	SendToken(transfers []msg.Transfer, sync bool, memo string, source int64) (*SendTokenResult, error)
	MintToken(symbol string, amount int64, sync bool, memo string, source int64) (*MintTokenResult, error)

	SubmitListPairProposal(title string, param msg.ListTradingPairParams, initialDeposit int64, votingPeriod time.Duration, sync bool, memo string, source int64) (*SubmitProposalResult, error)
	SubmitProposal(title string, description string, proposalType msg.ProposalKind, initialDeposit int64, votingPeriod time.Duration, sync bool, memo string, source int64) (*SubmitProposalResult, error)
	DepositProposal(proposalID int64, amount int64, sync bool, memo string, source int64) (*DepositProposalResult, error)
	VoteProposal(proposalID int64, option msg.VoteOption, sync bool, memo string, source int64) (*VoteProposalResult, error)

	GetKeyManager() keys.KeyManager
}

type client struct {
	basicClient basic.BasicClient
	queryClient query.QueryClient
	keyManager  keys.KeyManager
	chainId     string
}

func NewClient(chainId string, keyManager keys.KeyManager, queryClient query.QueryClient, basicClient basic.BasicClient) TransactionClient {
	return &client{basicClient, queryClient, keyManager, chainId}
}

func (c *client) GetKeyManager() keys.KeyManager {
	return c.keyManager
}

func (c *client) broadcastMsg(m msg.Msg, sync bool, memo string, source int64) (*tx.TxCommitResult, error) {
	fromAddr := c.keyManager.GetAddr()
	acc, err := c.queryClient.GetAccount(fromAddr.String())
	if err != nil {
		return nil, err
	}
	sequence := acc.Sequence
	// prepare message to sign
	signMsg := tx.StdSignMsg{
		ChainID:       c.chainId,
		AccountNumber: acc.Number,
		Sequence:      sequence,
		Memo:          memo,
		Msgs:          []msg.Msg{m},
		Source:        source,
	}

	// Hex encoded signed transaction, ready to be posted to BncChain API
	hexTx, err := c.keyManager.Sign(signMsg)
	if err != nil {
		return nil, err
	}
	param := map[string]string{}
	if sync {
		param["sync"] = "true"
	}
	commits, err := c.basicClient.PostTx(hexTx, param)
	if err != nil {
		return nil, err
	}
	if len(commits) < 1 {
		return nil, fmt.Errorf("Len of tx Commit result is less than 1 ")
	}
	return &commits[0], nil
}
