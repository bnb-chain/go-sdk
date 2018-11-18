package sdk

import (
	"reflect"
	"testing"
)

func TestTokens(t *testing.T) {
	sdk := &SDK{
		dexAPI: &fDexAPI{},
	}

	pair, err := sdk.GetTokens()
	if err != nil {
		t.Errorf("GetTokens failed, expected no error but got %v", err)
	}

	expected := []*Token{
		&Token{
			Name:        "ABC Token",
			Symbol:      "BNB",
			TotalSupply: "100000000",
			Owner:       "cosmosaccaddr1hy2e872rqtd675sn72ny87cyyaaanmqeuvwrpc",
		},
	}

	if !reflect.DeepEqual(expected, pair) {
		t.Errorf("GetTokens wrong results, expected %s but got %s", expected, pair)
	}

}
