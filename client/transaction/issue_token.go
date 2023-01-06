package transaction

import (
	"encoding/json"
	"fmt"

	"github.com/bnb-chain/go-sdk/types/msg"
	"github.com/bnb-chain/go-sdk/types/tx"
)

type IssueTokenResult struct {
	tx.TxCommitResult
	Symbol string `json:"symbol"`
}

type IssueTokenValue struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	OrigSymbol  string `json:"original_symbol"`
	TotalSupply string `json:"total_supply"`
	Owner       string `json:"owner"`
}

func (c *client) IssueToken(name, symbol string, supply int64, sync bool, mintable bool, options ...Option) (*IssueTokenResult, error) {
	if symbol == "" {
		return nil, fmt.Errorf("Issue token symbol can't be empty ")
	}
	fromAddr := c.keyManager.GetAddr()

	issueMsg := msg.NewTokenIssueMsg(
		fromAddr,
		name,
		symbol,
		supply,
		mintable,
	)
	commit, err := c.broadcastMsg(issueMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	var issueTokenValue IssueTokenValue
	issueSymbol := symbol
	if commit.Ok && sync {
		err = json.Unmarshal([]byte(commit.Data), &issueTokenValue)
		if err != nil {
			return nil, err
		}
		issueSymbol = issueTokenValue.Symbol
	}

	return &IssueTokenResult{*commit, issueSymbol}, nil

}
