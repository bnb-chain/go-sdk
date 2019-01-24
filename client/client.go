package client

import (
	"gopkg.in/resty.v1"

	"github.com/binance-chain/go-sdk/client/basic"
	"github.com/binance-chain/go-sdk/client/query"
	"github.com/binance-chain/go-sdk/client/transaction"
	"github.com/binance-chain/go-sdk/keys"
	"github.com/binance-chain/go-sdk/types"
)

// dexClient wrapper
type dexClient struct {
	query.QueryClient
	transaction.TransactionClient
	basic.BasicClient
}

// DexClient methods
type DexClient interface {
	basic.BasicClient
	query.QueryClient
	transaction.TransactionClient
}

func init() {
	resty.DefaultClient.SetRedirectPolicy(resty.FlexibleRedirectPolicy(10))
}

func NewDexClient(baseUrl string, network types.ChainNetwork, keyManager keys.KeyManager) (DexClient, error) {
	types.Network = network
	c := basic.NewClient(baseUrl + types.DefaultAPIVersionPrefix)
	q := query.NewClient(c)
	n, err := q.GetNodeInfo()
	if err != nil {
		return nil, err
	}
	t := transaction.NewClient(n.NodeInfo.Network, keyManager, q, c)
	return &dexClient{BasicClient: c, QueryClient: q, TransactionClient: t}, nil
}
