package transaction

import (
	"encoding/json"
	"strings"

	sdk "github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
)

type TransferInResult struct {
	tx.TxCommitResult
}

func (c *client) TransferIn(sequence int64, contractAddr msg.EthereumAddress,
	refundAddresses []msg.EthereumAddress, receiverAddresses []sdk.AccAddress, amounts []int64, symbol string,
	relayFee sdk.Coin, expireTime int64, sync bool, options ...Option) (*TransferInResult, error) {
	fromAddr := c.keyManager.GetAddr()
	claim := msg.TransferInClaim{
		ContractAddress:   contractAddr,
		RefundAddresses:   refundAddresses,
		ReceiverAddresses: receiverAddresses,
		Amounts:           amounts,
		Symbol:            symbol,
		RelayFee:          relayFee,
		ExpireTime:        expireTime,
	}

	claimBz, err := json.Marshal(claim)
	if err != nil {
		return nil, err
	}

	claimMsg := msg.NewClaimMsg(msg.ClaimTypeTransferIn, sequence, string(claimBz), fromAddr)

	commit, err := c.broadcastMsg(claimMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	return &TransferInResult{*commit}, nil
}

type BindResult struct {
	tx.TxCommitResult
}

func (c *client) Bind(symbol string, amount int64, contractAddress msg.EthereumAddress, contractDecimals int8, expireTime int64, sync bool, options ...Option) (*BindResult, error) {
	fromAddr := c.keyManager.GetAddr()
	bindMsg := msg.NewBindMsg(fromAddr, symbol, amount, contractAddress, contractDecimals, expireTime)
	commit, err := c.broadcastMsg(bindMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	return &BindResult{*commit}, nil
}

type TransferOutResult struct {
	tx.TxCommitResult
}

func (c *client) TransferOut(to msg.EthereumAddress, amount sdk.Coin, expireTime int64, sync bool, options ...Option) (*TransferOutResult, error) {
	fromAddr := c.keyManager.GetAddr()
	transferOutMsg := msg.NewTransferOutMsg(fromAddr, to, amount, expireTime)
	commit, err := c.broadcastMsg(transferOutMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	return &TransferOutResult{*commit}, nil
}

type TransferOutRefundResult struct {
	tx.TxCommitResult
}

func (c *client) TransferOutRefund(sequence int64, refundAddr sdk.AccAddress, amount sdk.Coin, refundReason msg.RefundReason, sync bool, options ...Option) (*TransferOutRefundResult, error) {
	fromAddr := c.keyManager.GetAddr()
	claim := msg.TransferOutRefundClaim{
		RefundAddress: refundAddr,
		Amount:        amount,
		RefundReason:  refundReason,
	}

	claimBz, err := json.Marshal(claim)
	if err != nil {
		return nil, err
	}

	claimMsg := msg.NewClaimMsg(msg.ClaimTypeTransferOutRefund, sequence, string(claimBz), fromAddr)

	commit, err := c.broadcastMsg(claimMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	return &TransferOutRefundResult{*commit}, nil
}

type UpdateBindResult struct {
	tx.TxCommitResult
}

func (c *client) UpdateBind(sequence int64, symbol string, contractAddress msg.EthereumAddress, status msg.BindStatus, sync bool, options ...Option) (*UpdateBindResult, error) {
	fromAddr := c.keyManager.GetAddr()
	claim := msg.UpdateBindClaim{
		Status:          status,
		Symbol:          strings.ToUpper(symbol),
		ContractAddress: contractAddress,
	}

	claimBz, err := json.Marshal(claim)
	if err != nil {
		return nil, err
	}

	claimMsg := msg.NewClaimMsg(msg.ClaimTypeUpdateBind, sequence, string(claimBz), fromAddr)

	commit, err := c.broadcastMsg(claimMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	return &UpdateBindResult{*commit}, nil
}
