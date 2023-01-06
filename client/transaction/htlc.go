package transaction

import (
	"github.com/bnb-chain/go-sdk/common/types"
	"github.com/bnb-chain/go-sdk/types/msg"
	"github.com/bnb-chain/go-sdk/types/tx"
)

type HTLTResult struct {
	tx.TxCommitResult
}

func (c *client) HTLT(recipient types.AccAddress, recipientOtherChain, senderOtherChain string, randomNumberHash []byte, timestamp int64,
	amount types.Coins, expectedIncome string, heightSpan int64, crossChain bool, sync bool, options ...Option) (*HTLTResult, error) {
	fromAddr := c.keyManager.GetAddr()
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
	commit, err := c.broadcastMsg(htltMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	return &HTLTResult{*commit}, nil
}

type DepositHTLTResult struct {
	tx.TxCommitResult
}

func (c *client) DepositHTLT(swapID []byte, amount types.Coins,
	sync bool, options ...Option) (*DepositHTLTResult, error) {
	fromAddr := c.keyManager.GetAddr()
	depositHTLTMsg := msg.NewDepositHTLTMsg(
		fromAddr,
		amount,
		swapID,
	)
	commit, err := c.broadcastMsg(depositHTLTMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	return &DepositHTLTResult{*commit}, nil
}

type ClaimHTLTResult struct {
	tx.TxCommitResult
}

func (c *client) ClaimHTLT(swapID []byte, randomNumber []byte, sync bool, options ...Option) (*ClaimHTLTResult, error) {
	fromAddr := c.keyManager.GetAddr()
	claimHTLTMsg := msg.NewClaimHTLTMsg(
		fromAddr,
		swapID,
		randomNumber,
	)
	commit, err := c.broadcastMsg(claimHTLTMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	return &ClaimHTLTResult{*commit}, nil
}

type RefundHTLTResult struct {
	tx.TxCommitResult
}

func (c *client) RefundHTLT(swapID []byte, sync bool, options ...Option) (*RefundHTLTResult, error) {
	fromAddr := c.keyManager.GetAddr()
	refundHTLTMsg := msg.NewRefundHTLTMsg(
		fromAddr,
		swapID,
	)
	commit, err := c.broadcastMsg(refundHTLTMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	return &RefundHTLTResult{*commit}, nil
}
