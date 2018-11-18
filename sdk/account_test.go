package sdk

import (
	"reflect"
	"testing"
)

func TestAccountError(t *testing.T) {
	sdk := &SDK{
		dexAPI: &fDexAPI{},
	}

	_, err := sdk.GetAccount("")
	if err == nil {
		t.Errorf("GetAccount failed, expected `Error` but got %v", err)
	}
}

func TestAccount(t *testing.T) {
	sdk := &SDK{
		dexAPI: &fDexAPI{},
	}

	depth, err := sdk.GetAccount("cosmosaccaddr1hy2e872rqtd675sn72ny87cyyaaanmqeuvwrpc")
	if err != nil {
		t.Errorf("GetAccount failed, expected no error but got %v", err)
	}

	expected := &Account{
		Number:  "0",
		Address: "cosmosaccaddr1hy2e872rqtd675sn72ny87cyyaaanmqeuvwrpc",
		Balances: []Coin{
			Coin{
				Symbol: "BNB",
				Free:   "18975020177895000",
				Locked: "000000000",
				Frozen: "000000000",
			},
			Coin{
				Symbol: "NNB",
				Free:   "828912912928291",
				Locked: "000000000",
				Frozen: "75000000",
			},
		},
		PublicKey: []int8{1, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		Sequence:  "298113",
	}

	if !reflect.DeepEqual(expected, depth) {
		t.Errorf("GetAccount wrong results, expected \n%v but got \n%v", expected, depth)
	}

}
