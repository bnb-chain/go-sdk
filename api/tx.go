package api

import (
	"encoding/json"
	"fmt"
)

const (
	CodeOk int32 = 0
)

// TxResult def
type TxResult struct {
	Hash string `json:"hash"`
	Log  string `json:"log"`
	Data string `json:"data"`
	Code int32  `json:"code"`
}

// TxCommitResult for POST tx results
type TxCommitResult struct {
	Ok   bool   `json:"ok"`
	Log  string `json:"log"`
	Hash string `json:"hash"`
	Code int32  `json:"code"`
	Data string `json:"data"`
}

// GetTx returns transaction details
func (dex *dexAPI) GetTx(txHash string) (*TxResult, error) {
	if txHash == "" {
		return nil, fmt.Errorf("Invalid tx hash %s ", txHash)
	}

	qp := map[string]string{}
	resp, err := dex.Get("/tx/"+txHash, qp)
	if err != nil {
		return nil, err
	}

	var txResult TxResult
	if err := json.Unmarshal(resp, &txResult); err != nil {
		return nil, err
	}

	return &txResult, nil
}

// PostTx returns transaction details
func (dex *dexAPI) PostTx(hexTx []byte, param map[string]string) ([]TxCommitResult, error) {
	if len(hexTx) == 0 {
		return nil, fmt.Errorf("Invalid tx  %s", hexTx)
	}

	body := hexTx
	resp, err := dex.Post("/broadcast", body, param)
	if err != nil {
		return nil, err
	}
	txResult := make([]TxCommitResult, 0)
	if err := json.Unmarshal(resp, &txResult); err != nil {
		return nil, err
	}

	return txResult, nil
}
