package sdk

import (
	"reflect"
	"testing"
)

func TestTrades(t *testing.T) {
	sdk := &SDK{
		dexAPI: &fDexAPI{},
	}

	trades, err := sdk.GetTrades(&TradesQuery{})
	if err != nil {
		t.Errorf("GetTrades failed, expected no error but got %v", err)
	}

	expected := []*Trade{
		&Trade{
			BuyerOrderID:  "order-buy-1",
			BuyFee:        "0.50000000",
			Price:         "0.75000000",
			Quantity:      "1.00000000",
			SellerOrderID: "order-sell-1",
			SellFee:       "0.50000000",
			Symbol:        "BNB_NNB",
			Time:          1000000000,
			TradeID:       "trade-1",
		},
	}

	if !reflect.DeepEqual(expected, trades) {
		t.Errorf("GetTrades wrong results, expected %v but got %v", expected, trades)
	}

}
