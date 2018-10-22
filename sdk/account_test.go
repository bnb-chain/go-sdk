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
		Coins: []Coin{
			Coin{
				Denom:  "BNB",
				Amount: "18975020177895000",
			},
			Coin{
				Denom:  "NNB",
				Amount: "1737120240518526",
			},
			Coin{
				Denom:  "ZCB",
				Amount: "1887578172962946",
			},
		},
		PublicKey: PublicKey{
			Type:  "tendermint/PubKeySecp256k1",
			Value: "A58TeSbC3MRQ1ig5heN/XPinu9kjZrK4gp60DD7czU8J",
		},
		Sequence: "298113",
	}

	if !reflect.DeepEqual(expected, depth) {
		t.Errorf("GetAccount wrong results, expected %v but got %v", expected, depth)
	}

}
