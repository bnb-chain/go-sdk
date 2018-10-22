package sdk

import "fmt"

// SDK wrapper
type SDK struct {
	dexAPI IDexAPI
}

// NewSDK init
func NewSDK(baseURL string) (*SDK, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("Invalid baseURL %s", baseURL)
	}

	return &SDK{
		dexAPI: &DexAPI{baseURL},
	}, nil
}
