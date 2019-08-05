package transaction

import (
	"github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
)

type HashTimerLockTransferResult struct {
	tx.TxCommitResult
}

func (c *client) HashTimerLockTransfer(to types.AccAddress, toOnOtherChain []byte, randomNumberHash []byte, timestamp int64,
	outAmount types.Coin, inAmount int64, heightSpan int64, sync bool, options ...Option) (*HashTimerLockTransferResult, error) {
	fromAddr := c.keyManager.GetAddr()
	hashTimerLockTransferMsg := msg.NewHashTimerLockTransferMsg(
		fromAddr,
		to,
		toOnOtherChain,
		randomNumberHash,
		timestamp,
		outAmount,
		inAmount,
		heightSpan,
	)
	err := hashTimerLockTransferMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := c.broadcastMsg(hashTimerLockTransferMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	return &HashTimerLockTransferResult{*commit}, nil
}

type ClaimHashTimerLockResult struct {
	tx.TxCommitResult
}

func (c *client) ClaimHashTimerLock(randomNumberHash []byte, randomNumber []byte, sync bool, options ...Option) (*ClaimHashTimerLockResult, error) {
	fromAddr := c.keyManager.GetAddr()
	claimHashTimerLockMsg := msg.NewClaimHashTimerLockMsg(
		fromAddr,
		randomNumberHash,
		randomNumber,
	)
	err := claimHashTimerLockMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := c.broadcastMsg(claimHashTimerLockMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	return &ClaimHashTimerLockResult{*commit}, nil
}

type RefundHashTimerLockResult struct {
	tx.TxCommitResult
}

func (c *client) RefundHashTimerLock(randomNumberHash []byte, sync bool, options ...Option) (*RefundHashTimerLockResult, error) {
	fromAddr := c.keyManager.GetAddr()
	refundHashTimerLockMsg := msg.NewRefundLockedAssetMsg(
		fromAddr,
		randomNumberHash,
	)
	err := refundHashTimerLockMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := c.broadcastMsg(refundHashTimerLockMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	return &RefundHashTimerLockResult{*commit}, nil
}
