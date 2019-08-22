package transaction

import (
	"fmt"
	"time"

	"github.com/binance-chain/go-sdk/client/basic"
	"github.com/binance-chain/go-sdk/client/query"
	"github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/keys"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
)

type TransactionClient interface {
	CreateOrder(baseAssetSymbol, quoteAssetSymbol string, op int8, price, quantity int64, sync bool, options ...Option) (*CreateOrderResult, error)
	CancelOrder(baseAssetSymbol, quoteAssetSymbol, refId string, sync bool, options ...Option) (*CancelOrderResult, error)
	BurnToken(symbol string, amount int64, sync bool, options ...Option) (*BurnTokenResult, error)
	ListPair(proposalId int64, baseAssetSymbol string, quoteAssetSymbol string, initPrice int64, sync bool, options ...Option) (*ListPairResult, error)
	FreezeToken(symbol string, amount int64, sync bool, options ...Option) (*FreezeTokenResult, error)
	UnfreezeToken(symbol string, amount int64, sync bool, options ...Option) (*UnfreezeTokenResult, error)
	IssueToken(name, symbol string, supply int64, sync bool, mintable bool, options ...Option) (*IssueTokenResult, error)
	SendToken(transfers []msg.Transfer, sync bool, options ...Option) (*SendTokenResult, error)
	MintToken(symbol string, amount int64, sync bool, options ...Option) (*MintTokenResult, error)
	TimeLock(description string, amount types.Coins, lockTime int64, sync bool, options ...Option) (*TimeLockResult, error)
	TimeUnLock(id int64, sync bool, options ...Option) (*TimeUnLockResult, error)
	TimeReLock(id int64, description string, amount types.Coins, lockTime int64, sync bool, options ...Option) (*TimeReLockResult, error)
	SetAccountFlags(flags uint64, sync bool, options ...Option) (*SetAccountFlagsResult, error)
	AddAccountFlags(flagOptions []types.FlagOption, sync bool, options ...Option) (*SetAccountFlagsResult, error)
	HTLT(recipient types.AccAddress, recipientOtherChain []byte, randomNumberHash []byte, timestamp int64, outAmount types.Coin, expectedIncome string, heightSpan int64, crossChain bool, sync bool, options ...Option) (*HTLTResult, error)
	DepositHTLT(recipient types.AccAddress, randomNumberHash []byte, outAmount types.Coin, sync bool, options ...Option) (*DepositHTLTResult, error)
	ClaimHTLT(randomNumberHash []byte, randomNumber []byte, sync bool, options ...Option) (*ClaimHTLTResult, error)
	RefundHTLT(randomNumberHash []byte, sync bool, options ...Option) (*RefundHTLTResult, error)

	SubmitListPairProposal(title string, param msg.ListTradingPairParams, initialDeposit int64, votingPeriod time.Duration, sync bool, options ...Option) (*SubmitProposalResult, error)
	SubmitProposal(title string, description string, proposalType msg.ProposalKind, initialDeposit int64, votingPeriod time.Duration, sync bool, options ...Option) (*SubmitProposalResult, error)
	DepositProposal(proposalID int64, amount int64, sync bool, options ...Option) (*DepositProposalResult, error)
	VoteProposal(proposalID int64, option msg.VoteOption, sync bool, options ...Option) (*VoteProposalResult, error)

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

type Option func(*tx.StdSignMsg) *tx.StdSignMsg

func WithSource(source int64) Option {
	return func(txMsg *tx.StdSignMsg) *tx.StdSignMsg {
		txMsg.Source = source
		return txMsg
	}
}

func WithMemo(memo string) Option {
	return func(txMsg *tx.StdSignMsg) *tx.StdSignMsg {
		txMsg.Memo = memo
		return txMsg
	}
}

func WithAcNumAndSequence(accountNum, seq int64) Option {
	return func(txMsg *tx.StdSignMsg) *tx.StdSignMsg {
		txMsg.Sequence = seq
		txMsg.AccountNumber = accountNum
		return txMsg
	}
}

func (c *client) broadcastMsg(m msg.Msg, sync bool, options ...Option) (*tx.TxCommitResult, error) {
	// prepare message to sign
	signMsg := &tx.StdSignMsg{
		ChainID:       c.chainId,
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
		fromAddr := c.keyManager.GetAddr()
		acc, err := c.queryClient.GetAccount(fromAddr.String())
		if err != nil {
			return nil, err
		}
		signMsg.Sequence = acc.Sequence
		signMsg.AccountNumber = acc.Number
	}

	// special logic for createOrder, to save account query
	if orderMsg, ok := m.(msg.CreateOrderMsg); ok {
		orderMsg.ID = msg.GenerateOrderID(signMsg.Sequence+1, c.keyManager.GetAddr())
		signMsg.Msgs[0] = orderMsg
	}

	for _, m := range signMsg.Msgs {
		if err := m.ValidateBasic(); err != nil {
			return nil, err
		}
	}

	// Hex encoded signed transaction, ready to be posted to BncChain API
	hexTx, err := c.keyManager.Sign(*signMsg)
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
