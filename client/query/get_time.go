package query

import (
	"encoding/json"

	"github.com/bnb-chain/go-sdk/common/types"
)

// GetTime returns market depth records
func (c *client) GetTime() (*types.Time, error) {
	qp := map[string]string{}
	resp, _, err := c.baseClient.Get("/time", qp)
	if err != nil {
		return nil, err
	}

	var t types.Time
	if err := json.Unmarshal(resp, &t); err != nil {
		return nil, err
	}

	return &t, nil
}
