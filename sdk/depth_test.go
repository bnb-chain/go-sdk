package sdk

import (
	"reflect"
	"testing"
)

func TestDepthError(t *testing.T) {
	sdk := &SDK{
		dexAPI: &fDexAPI{},
	}

	_, err := sdk.GetDepth(&DepthQuery{})
	if err == nil || err.Error() != "Query.Symbol is required" {
		t.Errorf("GetDepth failed, expected `Error Query.Symbol is required` but got %v", err)
	}
}

func TestDepth(t *testing.T) {
	sdk := &SDK{
		dexAPI: &fDexAPI{},
	}

	depth, err := sdk.GetDepth(&DepthQuery{Symbol: "BNB_NNB"})
	if err != nil {
		t.Errorf("GetDepth failed, expected no error but got %v", err)
	}

	expected := &MarketDepth{
		LastUpdateID: 1000,
		Symbol:       "BNB_NNB",
		Bids:         [][]string{[]string{"0.00240000", "50"}, []string{"0.00230000", "100"}},
		Asks:         [][]string{[]string{"0.00250000", "90"}, []string{"0.00260000", "120"}, []string{"0.0030000", "350"}},
	}

	if !reflect.DeepEqual(expected, depth) {
		t.Errorf("GetDepth wrong results, expected %v but got %v", expected, depth)
	}

}
