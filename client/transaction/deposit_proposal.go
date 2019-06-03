package transaction

import (
	ctypes "github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/types"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
)

type DepositProposalResult struct {
	tx.TxCommitResult
}

func (c *client) DepositProposal(proposalID int64, amount int64, sync bool, memo string, source int64) (*DepositProposalResult, error) {
	fromAddr := c.keyManager.GetAddr()
	coins := ctypes.Coins{ctypes.Coin{Denom: types.NativeSymbol, Amount: amount}}
	depositMsg := msg.NewDepositMsg(fromAddr, proposalID, coins)
	err := depositMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := c.broadcastMsg(depositMsg, sync, memo, source)
	if err != nil {
		return nil, err
	}

	return &DepositProposalResult{*commit}, err

}
