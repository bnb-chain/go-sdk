package sdk

import (
	"fmt"
)

// SDK wrapper
type SDK struct {
	dexAPI IDexAPI
}

// NewBncSDK init
func NewBncSDK(baseURL string) (*SDK, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("Invalid baseURL %s", baseURL)
	}

	return &SDK{
		dexAPI: &DexAPI{baseURL},
	}, nil
}
