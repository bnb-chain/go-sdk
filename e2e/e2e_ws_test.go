package e2e

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	sdk "github.com/binance-chain/go-sdk/client"
	"github.com/binance-chain/go-sdk/client/websocket"
	ctypes "github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/keys"
	"github.com/binance-chain/go-sdk/types"
)

func NewClient(t *testing.T) sdk.DexClient {
	mnemonic := "test mnemonic"
	keyManager, err := keys.NewMnemonicKeyManager(mnemonic)
	assert.NoError(t, err)

	client, err := sdk.NewDexClient("testnet-dex.binance.org", ctypes.TestNetwork, keyManager)
	assert.NoError(t, err)
	return client
}

func TestSubscribeAllTickerEvent(t *testing.T) {
	client := NewClient(t)

	quit := make(chan struct{})
	err := client.SubscribeAllTickerEvent(quit, func(events []*websocket.TickerEvent) {
		bz, _ := json.Marshal(events)
		fmt.Println(string(bz))
	}, func(err error) {
		assert.NoError(t, err)
	}, nil)
	assert.NoError(t, err)
	time.Sleep(10 * time.Second)
	close(quit)
	time.Sleep(1 * time.Second)
}

func TestSubscribeTickerEvent(t *testing.T) {
	client := NewClient(t)

	markets, err := client.GetMarkets(ctypes.NewMarketsQuery().WithLimit(1))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(markets))
	tradeSymbol := markets[0].BaseAssetSymbol
	if markets[0].QuoteAssetSymbol != types.NativeSymbol {
		tradeSymbol = markets[0].QuoteAssetSymbol
	}
	quit := make(chan struct{})
	err = client.SubscribeTickerEvent(tradeSymbol, types.NativeSymbol, quit, func(event *websocket.TickerEvent) {
		bz, _ := json.Marshal(event)
		fmt.Println(string(bz))
	}, func(err error) {
		assert.NoError(t, err)
	}, nil)
	assert.NoError(t, err)
	time.Sleep(10 * time.Second)
	close(quit)
	time.Sleep(1 * time.Second)
}

func TestSubscribeAllMiniTickersEvent(t *testing.T) {
	client := NewClient(t)

	quit := make(chan struct{})
	err := client.SubscribeAllMiniTickersEvent(quit, func(events []*websocket.MiniTickerEvent) {
		bz, _ := json.Marshal(events)
		fmt.Println(string(bz))
	}, func(err error) {
		assert.NoError(t, err)
	}, nil)
	assert.NoError(t, err)
	time.Sleep(10 * time.Second)
	close(quit)
	time.Sleep(1 * time.Second)
}

func TestSubscribeMiniTickersEvent(t *testing.T) {
	client := NewClient(t)

	markets, err := client.GetMarkets(ctypes.NewMarketsQuery().WithLimit(1))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(markets))
	tradeSymbol := markets[0].BaseAssetSymbol
	if markets[0].QuoteAssetSymbol != types.NativeSymbol {
		tradeSymbol = markets[0].QuoteAssetSymbol
	}
	quit := make(chan struct{})
	err = client.SubscribeMiniTickerEvent(tradeSymbol, types.NativeSymbol, quit, func(event *websocket.MiniTickerEvent) {
		bz, err := json.Marshal(event)
		assert.NoError(t, err)
		fmt.Println(string(bz))
	}, func(err error) {
		assert.NoError(t, err)
	}, nil)
	assert.NoError(t, err)
	time.Sleep(10 * time.Second)
	close(quit)
	time.Sleep(1 * time.Second)
}

func TestSubscribeTradeEvent(t *testing.T) {
	client := NewClient(t)

	markets, err := client.GetMarkets(ctypes.NewMarketsQuery().WithLimit(1))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(markets))
	tradeSymbol := markets[0].BaseAssetSymbol
	if markets[0].QuoteAssetSymbol != types.NativeSymbol {
		tradeSymbol = markets[0].QuoteAssetSymbol
	}
	quit := make(chan struct{})
	err = client.SubscribeTradeEvent(tradeSymbol, types.NativeSymbol, quit, func(events []*websocket.TradeEvent) {
		bz, err := json.Marshal(events)
		assert.NoError(t, err)
		fmt.Println(string(bz))
	}, func(err error) {
		assert.NoError(t, err)
	}, nil)
	assert.NoError(t, err)
	time.Sleep(10 * time.Second)
	close(quit)
	time.Sleep(1 * time.Second)
}

func TestSubscribeOrderEvent(t *testing.T) {
	client := NewClient(t)

	quit := make(chan struct{})
	err := client.SubscribeOrderEvent(client.GetKeyManager().GetAddr().String(), quit, func(event []*websocket.OrderEvent) {
		bz, _ := json.Marshal(event)
		fmt.Println(string(bz))
	}, func(err error) {
		assert.NoError(t, err)
	}, nil)
	assert.NoError(t, err)
	time.Sleep(10 * time.Second)
	close(quit)
	time.Sleep(1 * time.Second)
}

func TestSubscribeAccountEvent(t *testing.T) {
	client := NewClient(t)

	quit := make(chan struct{})
	err := client.SubscribeAccountEvent(client.GetKeyManager().GetAddr().String(), quit, func(event *websocket.AccountEvent) {
		bz, _ := json.Marshal(event)
		fmt.Println(string(bz))
	}, func(err error) {
		assert.NoError(t, err)
	}, nil)
	assert.NoError(t, err)
	time.Sleep(10 * time.Second)
	close(quit)
	time.Sleep(1 * time.Second)
}

func TestSubscribeBlockHeightEvent(t *testing.T) {
	client := NewClient(t)

	quit := make(chan struct{})
	err := client.SubscribeBlockHeightEvent(quit, func(event *websocket.BlockHeightEvent) {
		bz, _ := json.Marshal(event)
		fmt.Println(string(bz))
	}, func(err error) {
		assert.NoError(t, err)
	}, nil)
	assert.NoError(t, err)
	time.Sleep(10 * time.Second)
	close(quit)
	time.Sleep(1 * time.Second)
}

func TestSubscribeKlineEvent(t *testing.T) {
	client := NewClient(t)

	markets, err := client.GetMarkets(ctypes.NewMarketsQuery().WithLimit(1))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(markets))
	tradeSymbol := markets[0].BaseAssetSymbol
	if markets[0].QuoteAssetSymbol != types.NativeSymbol {
		tradeSymbol = markets[0].QuoteAssetSymbol
	}
	quit := make(chan struct{})
	err = client.SubscribeKlineEvent(tradeSymbol, types.NativeSymbol, websocket.OneMinuteInterval, quit, func(event *websocket.KlineEvent) {
		bz, err := json.Marshal(event)
		assert.NoError(t, err)
		fmt.Println(string(bz))
	}, func(err error) {
		assert.NoError(t, err)
	}, nil)
	assert.NoError(t, err)
	time.Sleep(10 * time.Second)
	close(quit)
	time.Sleep(1 * time.Second)
}

func TestSubscribeMarketDiffEvent(t *testing.T) {
	client := NewClient(t)

	markets, err := client.GetMarkets(ctypes.NewMarketsQuery().WithLimit(1))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(markets))
	tradeSymbol := markets[0].BaseAssetSymbol
	if markets[0].QuoteAssetSymbol != types.NativeSymbol {
		tradeSymbol = markets[0].QuoteAssetSymbol
	}
	quit := make(chan struct{})
	err = client.SubscribeMarketDiffEvent(tradeSymbol, types.NativeSymbol, quit, func(event *websocket.MarketDeltaEvent) {
		bz, err := json.Marshal(event)
		assert.NoError(t, err)
		fmt.Println(string(bz))
	}, func(err error) {
		assert.NoError(t, err)
	}, nil)
	assert.NoError(t, err)
	time.Sleep(10 * time.Second)
	close(quit)
	time.Sleep(1 * time.Second)
}

func TestSubscribeMarketDepthEvent(t *testing.T) {
	client := NewClient(t)

	markets, err := client.GetMarkets(ctypes.NewMarketsQuery().WithLimit(1))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(markets))
	tradeSymbol := markets[0].BaseAssetSymbol
	if markets[0].QuoteAssetSymbol != types.NativeSymbol {
		tradeSymbol = markets[0].BaseAssetSymbol
	}
	quit := make(chan struct{})
	err = client.SubscribeMarketDepthEvent(tradeSymbol, types.NativeSymbol, quit, func(event *websocket.MarketDepthEvent) {
		bz, err := json.Marshal(event)
		assert.NoError(t, err)
		fmt.Println(string(bz))
	}, func(err error) {
		assert.NoError(t, err)
	}, nil)
	assert.NoError(t, err)
	time.Sleep(10 * time.Second)
	close(quit)
	time.Sleep(1 * time.Second)
}
