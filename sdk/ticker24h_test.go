package sdk

import (
	"reflect"
	"testing"
)

func TestTicker24hError(t *testing.T) {
	sdk := &SDK{
		dexAPI: &fDexAPI{},
	}

	_, err := sdk.GetTicker24h("")
	if err == nil || err.Error() != "Symbol is required" {
		t.Errorf("GetTicker24h failed, expected `Error Symbol is required` but got %v", err)
	}
}

func TestTicker24h(t *testing.T) {
	sdk := &SDK{
		dexAPI: &fDexAPI{},
	}

	ticker24h, err := sdk.GetTicker24h("BNB_NNB")
	if err != nil {
		t.Errorf("GetTicker24h failed, expected no error but got %v", err)
	}

	expected := &Ticker24h{
		Symbol:             "BNB",
		AskPrice:           "1.0000000",
		AskQuantity:        "1.0000000",
		BidPrice:           "1.0000000",
		BidQuantity:        "1.0000000",
		CloseTime:          10000000,
		Count:              100,
		FirstID:            "order-1",
		HighPrice:          "1.0000000",
		LastID:             "order-100",
		LastPrice:          "1.0000000",
		LastQuantity:       "1.0000000",
		LowPrice:           "1.0000000",
		OpenPrice:          "1.0000000",
		OpenTime:           2000000,
		PrevClosePrice:     "1.0000000",
		PriceChange:        "1.0000000",
		PriceChangePercent: "15",
		QuoteVolume:        "1.0000000",
		Volume:             "1.0000000",
		WeightedAvgPrice:   "1.0000000",
	}

	if !reflect.DeepEqual(expected, ticker24h) {
		t.Errorf("GetTicker24h wrong results, expected \n%v but got \n%v", expected, ticker24h)
	}

}
