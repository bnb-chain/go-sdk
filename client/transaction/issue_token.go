package transaction

import (
	"encoding/json"
	"fmt"

	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
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

func (c *client) IssueToken(name, symbol string, supply int64, sync bool, mintable bool) (*IssueTokenResult, error) {
	if symbol == "" {
		return nil, fmt.Errorf("Freeze token symbol can'c be empty ")
	}
	fromAddr := c.keyManager.GetAddr()

	issueMsg := msg.NewTokenIssueMsg(
		fromAddr,
		name,
		symbol,
		supply,
		mintable,
	)
	err := issueMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := c.broadcastMsg(issueMsg, sync)
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
