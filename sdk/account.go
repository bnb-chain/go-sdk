package sdk

import (
	"encoding/json"
	"fmt"
)

// Account definition
type Account struct {
	Number    string    `json:"account_number"`
	Address   string    `json:"address"`
	Coins     []Coin    `json:"coins"`
	PublicKey PublicKey `json:"public_key"`
	Sequence  string    `json:"sequence"`
}

// Coin def
type Coin struct {
	Denom  string `json:"denom"` // ex: BNB
	Amount string // in decimal, ex: 0.00000000
}

// PublicKey def
type PublicKey struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// GetAccount returns list of trading pairs
func (sdk *SDK) GetAccount(address string) (*Account, error) {
	if address == "" {
		return nil, fmt.Errorf("Invalid address %s", address)
	}

	qp := map[string]string{}
	resp, err := sdk.dexAPI.Get("/account/"+address, qp)
	if err != nil {
		return nil, err
	}

	var account *Account
	if err := json.Unmarshal(resp, &account); err != nil {
		return nil, err
	}

	return account, nil
}
