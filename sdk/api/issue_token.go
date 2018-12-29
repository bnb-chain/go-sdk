package api

import (
	"encoding/json"
	"fmt"
	"github.com/binance-chain/go-sdk/sdk/tx/txmsg"
	"strings"
)

type IssueTokenResult struct {
	TxCommitResult
	Symbol string `json:"symbol"`
}

type IssueTokenValue struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	OrigSymbol  string `json:"original_symbol"`
	TotalSupply string `json:"total_supply"`
	Owner       string `json:"owner"`
}

func (dex *dexAPI) IssueToken(name, symbol string, supply int64, sync bool, mintable bool) (*IssueTokenResult, error) {
	if symbol == "" {
		return nil, fmt.Errorf("Freeze token symbol can't be empty ")
	}
	fromAddr := dex.keyManager.GetAddr()

	issueMsg := txmsg.NewTokenIssueMsg(
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
	commit, err := dex.broadcastMsg(issueMsg, sync)
	if err != nil {
		return nil, err
	}
	var issueTokenValue IssueTokenValue
	issueSymbol := symbol
	if commit.Ok && sync {
		// Do no return error here, we get no data for now.
		err = json.Unmarshal([]byte(commit.Data), &issueTokenValue)
		if err == nil {
			issueSymbol = issueTokenValue.Symbol
		} else {
			// Todo Try to find symbol in log, delete this code when chain fix
			resLogs := strings.Split(commit.Log, " ")
			if len(resLogs) > 0 {
				issueSymbol = resLogs[len(resLogs)-1]
			}
		}
	}

	return &IssueTokenResult{*commit, issueSymbol}, nil

}
