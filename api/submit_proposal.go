package api

import (
	"encoding/json"

	"github.com/binance-chain/go-sdk/tx"
	"github.com/binance-chain/go-sdk/tx/txmsg"
)

type SubmitProposalResult struct {
	TxCommitResult
	ProposalId int64 `json:"proposal_id"`
}

func (dex *dexAPI) SubmitListPairProposal(title string, param txmsg.ListTradingPairParams, initialDeposit int64, sync bool) (*SubmitProposalResult, error) {
	bz, err := json.Marshal(&param)
	if err != nil {
		return nil, err
	}
	return dex.SubmitProposal(title, string(bz), txmsg.ProposalTypeListTradingPair, initialDeposit, sync)
}

func (dex *dexAPI) SubmitProposal(title string, description string, proposalType txmsg.ProposalKind, initialDeposit int64, sync bool) (*SubmitProposalResult, error) {
	fromAddr := dex.keyManager.GetAddr()
	coins := txmsg.Coins{txmsg.Coin{Denom: NativeSymbol, Amount: initialDeposit}}
	proposalMsg := txmsg.NewMsgSubmitProposal(title, description, proposalType, fromAddr, coins)
	err := proposalMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := dex.broadcastMsg(proposalMsg, sync)
	if err != nil {
		return nil, err
	}
	var proposalId int64
	if commit.Ok && sync{
		// Todo since ap do not return proposal id now, do not return err
		tx.Cdc.UnmarshalBinaryBare([]byte(commit.Data), &proposalId)
	}
	return &SubmitProposalResult{*commit, proposalId}, err

}
