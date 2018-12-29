package api

import (
	"encoding/json"
)

// Account definition
type Account struct {
	Number    int64 `json:"account_number"`
	Address   string `json:"address"`
	Balances  []Coin `json:"balances"`
	PublicKey []uint8 `json:"public_key"`
	Sequence  int64 `json:"sequence"`
}

// Coin def
type Coin struct {
	Symbol string `json:"symbol"` // ex: BNB
	Free   string `json:"free"`   // in decimal, ex: 0.00000000
	Locked string `json:"locked"` // in decimal, ex: 0.00000000
	Frozen string `json:"frozen"` // in decimal, ex: 0.00000000
}

// GetAccount returns list of trading pairs
func (dex *dexAPI) GetAccount(address string) (*Account, error) {
	if address == "" {
		return nil, AddressMissingError
	}

	qp := map[string]string{}
	resp, err := dex.Get("/account/"+address, qp)
	if err != nil {
		return nil, err
	}
	var account Account
	if err := json.Unmarshal(resp, &account); err != nil {
		return nil, err
	}

	return &account, nil
}
