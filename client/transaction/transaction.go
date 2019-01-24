package transaction

import (
	"fmt"

	"github.com/binance-chain/go-sdk/client/basic"
	"github.com/binance-chain/go-sdk/client/query"
	"github.com/binance-chain/go-sdk/keys"
	"github.com/binance-chain/go-sdk/types"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
)

type TransactionClient interface {
	CreateOrder(baseAssetSymbol, quoteAssetSymbol string, op int8, price, quantity int64, sync bool) (*CreateOrderResult, error)
	CancelOrder(baseAssetSymbol, quoteAssetSymbol, id, refId string, sync bool) (*CancelOrderResult, error)
	BurnToken(symbol string, amount int64, sync bool) (*BurnTokenResult, error)
	ListPair(proposalId int64, baseAssetSymbol string, quoteAssetSymbol string, initPrice int64, sync bool) (*ListPairResult, error)
	FreezeToken(symbol string, amount int64, sync bool) (*FreezeTokenResult, error)
	UnfreezeToken(symbol string, amount int64, sync bool) (*UnfreezeTokenResult, error)
	IssueToken(name, symbol string, supply int64, sync bool, mintable bool) (*IssueTokenResult, error)
	SendToken(dest types.AccAddress, symbol string, quantity int64, sync bool) (*SendTokenResult, error)
	MintToken(symbol string, amount int64, sync bool) (*MintTokenResult, error)

	SubmitListPairProposal(title string, param msg.ListTradingPairParams, initialDeposit int64, sync bool) (*SubmitProposalResult, error)
	SubmitProposal(title string, description string, proposalType msg.ProposalKind, initialDeposit int64, sync bool) (*SubmitProposalResult, error)
	DepositProposal(proposalID int64, amount int64, sync bool) (*DepositProposalResult, error)
	VoteProposal(proposalID int64, option msg.VoteOption, sync bool) (*VoteProposalResult, error)
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

func (c *client) broadcastMsg(m msg.Msg, sync bool) (*tx.TxCommitResult, error) {
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
		Memo:          "",
		Msgs:          []msg.Msg{m},
		Source:        types.GoSdkSource,
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
