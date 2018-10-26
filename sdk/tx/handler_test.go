package tx

import (
	"fmt"
	"testing"

	"./txmsg"
)

func TestSignTx(t *testing.T) {

	priv, acc := PrivAndAddr()
	newOrderMsg := txmsg.NewNewOrderMsg(acc, txmsg.GenerateOrderID(1, acc), txmsg.OrderSide.BUY, "BNB_NNB", 100000000, 5000000000)

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
		Msgs:          []txmsg.Msg{newOrderMsg},
	}

	fmt.Println("signMsg: ", signMsg)

	tx := &Tx{}
	privKey := priv.Bytes()
	stdTx, err := tx.Sign(privKey, signMsg)

	fmt.Println("stdTx: ", stdTx)

	if err != nil {
		t.Errorf("tx.Sign() failed, expected signed tx but got error: %v", err)
	}

	if len(stdTx) == 0 {
		t.Errorf("tx.Sign() failed, expected signed tx but got empty data: %v", stdTx)
	}
}
