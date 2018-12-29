package sdk

import (
	"fmt"
	"github.com/binance-chain/go-sdk/sdk/api"
	"github.com/binance-chain/go-sdk/sdk/keys"
)

// Client wrapper
type Client struct {
	api.IDexAPI
}

// NewBncCLient init
func NewBncCLient(baseURL, chainid string, keyManager keys.KeyManager) (*Client, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("Invalid base url %s. ", baseURL)
	}
	return &Client{api.NewDefaultDexApi(baseURL, chainid, keyManager)}, nil
}
