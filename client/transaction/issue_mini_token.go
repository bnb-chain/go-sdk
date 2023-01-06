package transaction

import (
	"encoding/json"
	"fmt"

	"github.com/bnb-chain/go-sdk/types/msg"
	"github.com/bnb-chain/go-sdk/types/tx"
)

type IssueMiniTokenResult struct {
	tx.TxCommitResult
	Symbol string `json:"symbol"`
}

type IssueMiniTokenValue struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	OrigSymbol  string `json:"original_symbol"`
	TotalSupply string `json:"total_supply"`
	TokenURI    string `json:"token_uri"`
	Owner       string `json:"owner"`
}

func (c *client) IssueMiniToken(name, symbol string, supply int64, sync bool, mintable bool, tokenURI string, options ...Option) (*IssueMiniTokenResult, error) {
	if symbol == "" {
		return nil, fmt.Errorf("Issue mini token symbol can't be empty ")
	}
	fromAddr := c.keyManager.GetAddr()

	issueMsg := msg.NewMiniTokenIssueMsg(
		fromAddr,
		name,
		symbol,
		supply,
		mintable,
		tokenURI,
	)
	commit, err := c.broadcastMsg(issueMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	var issueTokenValue IssueMiniTokenValue
	issueSymbol := symbol
	if commit.Ok && sync {
		err = json.Unmarshal([]byte(commit.Data), &issueTokenValue)
		if err != nil {
			return nil, err
		}
		issueSymbol = issueTokenValue.Symbol
	}

	return &IssueMiniTokenResult{*commit, issueSymbol}, nil
}
