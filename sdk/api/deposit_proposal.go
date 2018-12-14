package api

import (
	"github.com/BiJie/bnc-go-sdk/sdk/tx/txmsg"
)

type DepositProposalResult struct {
	TxCommitResult
}

func (dex *dexAPI) DepositProposal(proposalID int64, amount int64, sync bool) (*DepositProposalResult, error) {
	fromAddr := dex.keyManager.GetAddr()
	coins := txmsg.Coins{txmsg.Coin{Denom: NativeSymbol, Amount: amount}}
	depositMsg := txmsg.NewDepositMsg(fromAddr, proposalID, coins)
	err := depositMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := dex.broadcastMsg(depositMsg, sync)
	if err != nil {
		return nil, err
	}

	return &DepositProposalResult{*commit}, err

}
