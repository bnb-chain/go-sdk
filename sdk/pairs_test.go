package sdk

import (
	"reflect"
	"testing"
)

func TestPairs(t *testing.T) {
	sdk := &SDK{
		dexAPI: &fDexAPI{},
	}

	pair, err := sdk.GetPairs(1)
	if err != nil {
		t.Errorf("GetPairs failed, expected no error but got %v", err)
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
		t.Errorf("GetPairs wrong results, expected %s but got %s", expected, pair)
	}

}
