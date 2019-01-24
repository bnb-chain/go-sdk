package transaction

import (
	"encoding/json"
	"strconv"

	"github.com/binance-chain/go-sdk/types"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
)

type SubmitProposalResult struct {
	tx.TxCommitResult
	ProposalId int64 `json:"proposal_id"`
}

func (c *client) SubmitListPairProposal(title string, param msg.ListTradingPairParams, initialDeposit int64, sync bool) (*SubmitProposalResult, error) {
	bz, err := json.Marshal(&param)
	if err != nil {
		return nil, err
	}
	return c.SubmitProposal(title, string(bz), msg.ProposalTypeListTradingPair, initialDeposit, sync)
}

func (c *client) SubmitProposal(title string, description string, proposalType msg.ProposalKind, initialDeposit int64, sync bool) (*SubmitProposalResult, error) {
	fromAddr := c.keyManager.GetAddr()
	coins := types.Coins{types.Coin{Denom: types.NativeSymbol, Amount: initialDeposit}}
	proposalMsg := msg.NewMsgSubmitProposal(title, description, proposalType, fromAddr, coins)
	err := proposalMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := c.broadcastMsg(proposalMsg, sync)
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
