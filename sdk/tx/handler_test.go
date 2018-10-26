package tx

import (
	"fmt"
	"testing"

	"./msgs"
)

func TestSignTx(t *testing.T) {

	acc := []byte(`cosmosaccaddr1hy2e872rqtd675sn72ny87cyyaaanmqeuvwrpc`)
	newOrderMsg := msgs.NewNewOrderMsg(acc, msgs.GenerateOrderID(100, acc), msgs.Side.BUY, "BNB_NNB", 100000000, 5000000000)

	fee := StdFee{
		Amount: Coins{Coin{Denom: "BNB", Amount: 100000000}},
		Gas:    0,
	}

	signMsg := StdSignMsg{
		ChainID:       "bnc-chain-1",
		AccountNumber: 100,
		Sequence:      1,
		Memo:          "",
		Fee:           fee,
		Msgs:          []msgs.Msg{newOrderMsg},
	}

	fmt.Println("signMsg: ", signMsg)

	// if err == nil {
	// 	t.Errorf("GetKlines failed, expected `Error` but got %v", err)
	// }
}
