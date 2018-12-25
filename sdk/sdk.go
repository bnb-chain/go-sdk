package sdk

import (
	"fmt"
	"github.com/binance-chain/go-sdk/sdk/api"
	"github.com/binance-chain/go-sdk/sdk/keys"
)

// SDK wrapper
type SDK struct {
	api.IDexAPI
}

// NewBncSDK init
func NewBncSDK(baseURL, chainid string, keyManager keys.KeyManager) (*SDK, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("Invalid base url %s. ", baseURL)
	}
	return &SDK{api.NewDefaultDexApi(baseURL, chainid, keyManager)}, nil
}
