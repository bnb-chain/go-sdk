package sdk

import (
	"encoding/json"
	"fmt"
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
	Code int32  `json:"code"`
	Data string `json:"data"`
	Log  string `json:"log"`
	Hash string `json:"hash"`
}

// GetTx returns transaction details
func (sdk *SDK) GetTx(txHash string) (*TxResult, error) {
	if txHash == "" {
		return nil, fmt.Errorf("Invalid tx hash %s", txHash)
	}

	qp := map[string]string{}
	resp, err := sdk.dexAPI.Get("/tx/"+txHash, qp)
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
func (sdk *SDK) PostTx(hexTx []byte) ([]*TxCommitResult, error) {
	if len(hexTx) == 0 {
		return nil, fmt.Errorf("Invalid tx  %s", hexTx)
	}

	body := map[string]interface{}{"tx": string(hexTx)}
	resp, err := sdk.dexAPI.Post("/tx", body)
	if err != nil {
		return nil, err
	}

	var txResult []*TxCommitResult
	if err := json.Unmarshal(resp, &txResult); err != nil {
		return nil, err
	}

	return txResult, nil
}
