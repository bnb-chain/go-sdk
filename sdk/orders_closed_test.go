package sdk

import (
	"reflect"
	"testing"
)

func TestClosedOrdersError(t *testing.T) {
	sdk := &SDK{
		dexAPI: &fDexAPI{},
	}

	_, err := sdk.GetClosedOrders(&ClosedOrdersQuery{})
	if err == nil || err.Error() != "Query.SenderAddress is required" {
		t.Errorf("GetClosedOrders failed, expected `Error: Query.SenderAddress is required` but got %v", err)
	}

	_, err = sdk.GetClosedOrders(&ClosedOrdersQuery{
		SenderAddress: "cosmosaccaddr1q68phq7q6fyup5x25maglf9sxc2t4hy42pa623",
		Side:          "ABC",
	})
	if err == nil || err.Error() != "Invalid `Query.Side` param" {
		t.Errorf("GetClosedOrders failed, expected `Error: Invalid `Query.Side` param` but got %v", err)
	}
}

func TestClosedOrders(t *testing.T) {
	sdk := &SDK{
		dexAPI: &fDexAPI{},
	}

	order, err := sdk.GetClosedOrders(&ClosedOrdersQuery{
		SenderAddress: "cosmosaccaddr1q68phq7q6fyup5x25maglf9sxc2t4hy42pa623",
	})
	if err != nil {
		t.Errorf("GetClosedOrders failed, expected no error but got %v", err)
	}

	expected := []*Order{
		&Order{
			ID:               "cosmosaccaddr1q68phq7q6fyup5x25maglf9sxc2t4hy42pa623-240402",
			Owner:            "cosmosaccaddr1q68phq7q6fyup5x25maglf9sxc2t4hy42pa623",
			Symbol:           "NNB_BNB",
			Price:            1600000000,
			Quantity:         8900000000,
			ExecutedQuantity: 8900000000,
			Side:             OrderSide.SELL,
			Status:           OrderStatus.FULLY_FILLED,
			TimeInForce:      TimeInForce.GTC,
			Type:             OrderType.LIMIT,
		},
	}

	if !reflect.DeepEqual(expected, order) {
		t.Errorf("GetClosedOrders wrong results, expected %v but got %v", expected, order)
	}

}
