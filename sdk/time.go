package sdk

import (
	"encoding/json"
)

// Time to be broadcasted to the user
type Time struct {
	Time string `json:"time"`
}

// GetTime returns market depth records
func (sdk *SDK) GetTime() (string, error) {
	qp := map[string]string{}
	resp, err := sdk.dexAPI.Get("/time", qp)
	if err != nil {
		return "", err
	}

	var t Time
	if err := json.Unmarshal(resp, &t); err != nil {
		return "", err
	}

	return t.Time, nil
}
