package transaction

import (
	"github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
)

type HTLTResult struct {
	tx.TxCommitResult
}

func (c *client) HTLT(recipient types.AccAddress, recipientOtherChain []byte, randomNumberHash []byte, timestamp int64,
	outAmount types.Coin, expectedIncome string, heightSpan int64, crossChain bool, sync bool, options ...Option) (*HTLTResult, error) {
	fromAddr := c.keyManager.GetAddr()
	htltMsg := msg.NewHTLTMsg(
		fromAddr,
		recipient,
		recipientOtherChain,
		randomNumberHash,
		timestamp,
		outAmount,
		expectedIncome,
		heightSpan,
		crossChain,
	)
	err := htltMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := c.broadcastMsg(htltMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	return &HTLTResult{*commit}, nil
}

type DepositHTLTResult struct {
	tx.TxCommitResult
}

func (c *client) DepositHTLT(recipient types.AccAddress, randomNumberHash []byte, outAmount types.Coin,
	sync bool, options ...Option) (*DepositHTLTResult, error) {
	fromAddr := c.keyManager.GetAddr()
	depositHTLTMsg := msg.NewDepositHTLTMsg(
		fromAddr,
		recipient,
		outAmount,
		randomNumberHash,
	)
	err := depositHTLTMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := c.broadcastMsg(depositHTLTMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	return &DepositHTLTResult{*commit}, nil
}

type ClaimHTLTResult struct {
	tx.TxCommitResult
}

func (c *client) ClaimHTLT(randomNumberHash []byte, randomNumber []byte, sync bool, options ...Option) (*ClaimHTLTResult, error) {
	fromAddr := c.keyManager.GetAddr()
	claimHTLTMsg := msg.NewClaimHTLTMsg(
		fromAddr,
		randomNumberHash,
		randomNumber,
	)
	err := claimHTLTMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := c.broadcastMsg(claimHTLTMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	return &ClaimHTLTResult{*commit}, nil
}

type RefundHTLTResult struct {
	tx.TxCommitResult
}

func (c *client) RefundHTLT(randomNumberHash []byte, sync bool, options ...Option) (*RefundHTLTResult, error) {
	fromAddr := c.keyManager.GetAddr()
	refundHTLTMsg := msg.NewRefundHTLTMsg(
		fromAddr,
		randomNumberHash,
	)
	err := refundHTLTMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := c.broadcastMsg(refundHTLTMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	return &RefundHTLTResult{*commit}, nil
}
