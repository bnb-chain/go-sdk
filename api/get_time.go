package api

import (
	"encoding/json"
)

type Time struct {
	ApTime    string `json:"ap_time"`
	BlockTime string `json:"block_time"`
}

// GetTime returns market depth records
func (dex *dexAPI) GetTime() (*Time, error) {
	qp := map[string]string{}
	resp, err := dex.Get("/time", qp)
	if err != nil {
		return nil, err
	}

	var t Time
	if err := json.Unmarshal(resp, &t); err != nil {
		return nil, err
	}

	return &t, nil
}
