package go_sdk

import (
	"fmt"
	"github.com/binance-chain/go-sdk/api"
	"github.com/binance-chain/go-sdk/keys"
)

type Client struct {
	api.IDexAPI
}

func NewBncClient(baseURL, chainId string, keyManager keys.KeyManager) (*Client, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("Invalid base url %s. ", baseURL)
	}
	return &Client{api.NewDefaultDexApi(baseURL, chainId, keyManager)}, nil
}
