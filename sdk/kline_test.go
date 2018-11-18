package sdk

import (
	"reflect"
	"testing"
)

func TestKlineError(t *testing.T) {
	sdk := &SDK{
		dexAPI: &fDexAPI{},
	}

	_, err := sdk.GetKlines(&KlineQuery{})
	if err == nil {
		t.Errorf("GetKlines failed, expected `Error` but got %v", err)
	}
}

func TestKline(t *testing.T) {
	sdk := &SDK{
		dexAPI: &fDexAPI{},
	}

	kline, err := sdk.GetKlines(&KlineQuery{
		Symbol:   "BNB_NNB",
		Interval: "1h",
	})
	if err != nil {
		t.Errorf("GetKlines failed, expected no error but got %v", err)
	}

	expected := []*Kline{
		&Kline{
			Close:            50000000,
			CloseTime:        90000000,
			High:             150000000,
			Low:              500000000,
			NumberOfTrades:   150,
			Open:             500000000,
			OpenTime:         1000000,
			QuoteAssetVolume: 800000000,
			Volume:           2000000000,
		},
	}

	if !reflect.DeepEqual(expected, kline) {
		t.Errorf("GetKlines wrong results, expected %v but got %v", expected, kline)
	}

}
