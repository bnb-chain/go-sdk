package sdk

import (
	"reflect"
	"testing"
)

func TestMarkets(t *testing.T) {
	sdk := &SDK{
		dexAPI: &fDexAPI{},
	}

	pair, err := sdk.GetMarkets(1)
	if err != nil {
		t.Errorf("GetMarkets failed, expected no error but got %v", err)
	}

	expected := []*SymbolPair{
		&SymbolPair{
			TradeAsset: "BNB",
			QuoteAsset: "NNB",
			Price:      "100000000",
			TickSize:   "0.00005000",
			LotSize:    "0.10000000",
		},
	}

	if !reflect.DeepEqual(expected, pair) {
		t.Errorf("GetMarkets wrong results, expected %s but got %s", expected, pair)
	}

}
