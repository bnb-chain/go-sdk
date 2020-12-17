package transaction

import (
	"fmt"

	"github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
)

type TransferTokenOwnershipResult struct {
	tx.TxCommitResult
}

func (c *client) TransferTokenOwnership(symbol string, newOwner types.AccAddress, sync bool, options ...Option) (*TransferTokenOwnershipResult, error) {
	if symbol == "" {
		return nil, fmt.Errorf("Transfer token ownership, symbol can't be empty ")
	}
	if newOwner == nil {
		return nil, fmt.Errorf("Transfer token ownership, new owner can't be nil ")
	}
	fromAddr := c.keyManager.GetAddr()
	transferOwnershipMsg := msg.NewTransferOwnershipMsg(fromAddr, symbol, newOwner)
	commit, err := c.broadcastMsg(transferOwnershipMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	return &TransferTokenOwnershipResult{*commit}, nil
}
