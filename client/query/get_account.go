package query

import (
	"encoding/json"
	"net/http"

	"github.com/bnb-chain/go-sdk/common/types"
)

// GetAccount returns list of trading pairs
func (c *client) GetAccount(address string) (*types.BalanceAccount, error) {
	if address == "" {
		return nil, types.AddressMissingError
	}

	qp := map[string]string{}
	resp, code, err := c.baseClient.Get("/account/"+address, qp)
	if err != nil {
		if code == http.StatusNotFound {
			return &types.BalanceAccount{}, nil
		}
		return nil, err
	}
	var account types.BalanceAccount
	if err := json.Unmarshal(resp, &account); err != nil {
		return nil, err
	}
	return &account, nil
}
