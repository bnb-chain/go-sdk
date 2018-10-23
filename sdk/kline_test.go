package sdk

import (
	"reflect"
	"testing"
)

func TestKlineError(t *testing.T) {
	sdk := &SDK{
		dexAPI: &fDexAPI{},
	}

	_, err := sdk.GetKline(&KlineQuery{})
	if err == nil {
		t.Errorf("GetKline failed, expected `Error` but got %v", err)
	}
}

func TestKline(t *testing.T) {
	sdk := &SDK{
		dexAPI: &fDexAPI{},
	}

	kline, err := sdk.GetKline(&KlineQuery{
		Symbol:   "BNB_NNB",
		Interval: "1h",
	})
	if err != nil {
		t.Errorf("GetKline failed, expected no error but got %v", err)
	}

	expected := &Kline{
		Close:            50000000,
		CloseTime:        90000000,
		High:             150000000,
		Low:              500000000,
		NumberOfTrades:   150,
		Open:             500000000,
		OpenTime:         1000000,
		QuoteAssetVolume: 800000000,
		Volume:           2000000000,
	}

	if !reflect.DeepEqual(expected, kline) {
		t.Errorf("GetKline wrong results, expected %v but got %v", expected, kline)
	}

}
