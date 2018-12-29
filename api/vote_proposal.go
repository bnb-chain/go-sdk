package api

import (
	"github.com/binance-chain/go-sdk/tx/txmsg"
)

type VoteProposalResult struct {
	TxCommitResult
}

func (dex *dexAPI) VoteProposal(proposalID int64, option txmsg.VoteOption,sync bool) (*VoteProposalResult, error) {
	fromAddr := dex.keyManager.GetAddr()
	voteMsg := txmsg.NewMsgVote(fromAddr, proposalID, option)
	err := voteMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := dex.broadcastMsg(voteMsg, sync)
	if err != nil {
		return nil, err
	}

	return &VoteProposalResult{*commit}, err

}
