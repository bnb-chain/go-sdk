package sdk

import (
	"reflect"
	"testing"
)

func TestTickerError(t *testing.T) {
	sdk := &SDK{
		dexAPI: &fDexAPI{},
	}

	_, err := sdk.GetTicker("")
	if err == nil || err.Error() != "Symbol is required" {
		t.Errorf("GetTicker failed, expected `Error Symbol is required` but got %v", err)
	}
}

func TestTicker(t *testing.T) {
	sdk := &SDK{
		dexAPI: &fDexAPI{},
	}

	ticker, err := sdk.GetTicker("BNB_NNB")
	if err != nil {
		t.Errorf("GetTicker failed, expected no error but got %v", err)
	}

	expected := &Ticker{
		Symbol:   "BNB",
		AskPrice: "1.0000000",
		AskQty:   "1.0000000",
		BidPrice: "1.0000000",
		BidQty:   "1.0000000",
	}

	if !reflect.DeepEqual(expected, ticker) {
		t.Errorf("GetTicker wrong results, expected \n%v but got \n%v", expected, ticker)
	}

}
