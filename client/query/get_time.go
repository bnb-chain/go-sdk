package query

import (
	"encoding/json"
)

type Time struct {
	ApTime    string `json:"ap_time"`
	BlockTime string `json:"block_time"`
}

// GetTime returns market depth records
func (c *client) GetTime() (*Time, error) {
	qp := map[string]string{}
	resp, err := c.baseClient.Get("/time", qp)
	if err != nil {
		return nil, err
	}

	var t Time
	if err := json.Unmarshal(resp, &t); err != nil {
		return nil, err
	}

	return &t, nil
}
