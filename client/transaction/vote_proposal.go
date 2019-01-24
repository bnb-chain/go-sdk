package transaction

import (
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
)

type VoteProposalResult struct {
	tx.TxCommitResult
}

func (c *client) VoteProposal(proposalID int64, option msg.VoteOption, sync bool) (*VoteProposalResult, error) {
	fromAddr := c.keyManager.GetAddr()
	voteMsg := msg.NewMsgVote(fromAddr, proposalID, option)
	err := voteMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := c.broadcastMsg(voteMsg, sync)
	if err != nil {
		return nil, err
	}

	return &VoteProposalResult{*commit}, err

}
