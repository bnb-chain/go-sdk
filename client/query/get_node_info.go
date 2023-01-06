package query

import (
	"encoding/json"

	"github.com/bnb-chain/go-sdk/common/types"
)

func (c *client) GetNodeInfo() (*types.ResultStatus, error) {
	qp := map[string]string{}
	resp, _, err := c.baseClient.Get("/node-info", qp)
	if err != nil {
		return nil, err
	}
	var resultStatus types.ResultStatus
	if err := json.Unmarshal(resp, &resultStatus); err != nil {
		return nil, err
	}

	return &resultStatus, nil
}
