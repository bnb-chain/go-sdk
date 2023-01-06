package transaction

import (
	"encoding/json"
	"strconv"
	"time"

	ctypes "github.com/bnb-chain/go-sdk/common/types"
	"github.com/bnb-chain/go-sdk/types"
	"github.com/bnb-chain/go-sdk/types/msg"
	"github.com/bnb-chain/go-sdk/types/tx"
)

type SubmitProposalResult struct {
	tx.TxCommitResult
	ProposalId int64 `json:"proposal_id"`
}

func (c *client) SubmitListPairProposal(title string, param msg.ListTradingPairParams, initialDeposit int64, votingPeriod time.Duration, sync bool, options ...Option) (*SubmitProposalResult, error) {
	bz, err := json.Marshal(&param)
	if err != nil {
		return nil, err
	}
	return c.SubmitProposal(title, string(bz), msg.ProposalTypeListTradingPair, initialDeposit, votingPeriod, sync, options...)
}

func (c *client) SubmitProposal(title string, description string, proposalType msg.ProposalKind, initialDeposit int64, votingPeriod time.Duration, sync bool, options ...Option) (*SubmitProposalResult, error) {
	fromAddr := c.keyManager.GetAddr()
	coins := ctypes.Coins{ctypes.Coin{Denom: types.NativeSymbol, Amount: initialDeposit}}
	proposalMsg := msg.NewMsgSubmitProposal(title, description, proposalType, fromAddr, coins, votingPeriod)
	commit, err := c.broadcastMsg(proposalMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	var proposalId int64
	if commit.Ok && sync {
		// Todo since ap do not return proposal id now, do not return err
		proposalId, err = strconv.ParseInt(string(commit.Data), 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return &SubmitProposalResult{*commit, proposalId}, err

}
