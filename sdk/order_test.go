package sdk

import (
	"reflect"
	"testing"
)

func TestOrderError(t *testing.T) {
	sdk := &SDK{
		dexAPI: &fDexAPI{},
	}

	_, err := sdk.GetOrder("")
	if err == nil {
		t.Errorf("GetOrder failed, expected `Error` but got %v", err)
	}
}

func TestOrder(t *testing.T) {
	sdk := &SDK{
		dexAPI: &fDexAPI{},
	}

	order, err := sdk.GetOrder("cosmosaccaddr1q68phq7q6fyup5x25maglf9sxc2t4hy42pa623")
	if err != nil {
		t.Errorf("GetOrder failed, expected no error but got %v", err)
	}

	expected := &Order{
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
	}

	if !reflect.DeepEqual(expected, order) {
		t.Errorf("GetOrder wrong results, expected %v but got %v", expected, order)
	}

}
