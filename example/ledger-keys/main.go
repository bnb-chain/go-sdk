package main

import (
	"encoding/hex"
	"fmt"

	"github.com/binance-chain/go-sdk/common/crypto/ledger"
	"github.com/binance-chain/go-sdk/keys"
	"github.com/binance-chain/go-sdk/types"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
)

func main() {
	//Check whether there are variable ledger devices
	ledgerDevice, err := ledger.DiscoverLedger()
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to find ledger device: %s", err.Error()))
		return
	}
	err = ledgerDevice.Close()
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to find ledger device: %s", err.Error()))
		return
	}

	bip44Params := keys.NewBinanceBIP44Params(0, 0)
	keyManager1, err := keys.NewLedgerKeyManager(bip44Params.DerivationPath())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(fmt.Sprintf("address: %s", keyManager1.GetAddr()))
	fmt.Println(fmt.Sprintf("pubkey: %s", hex.EncodeToString(keyManager1.GetPrivKey().PubKey().Bytes())))

	receiverAddr, err := types.AccAddressFromBech32("bnb1pkppf2ar3wj38cu7y8khyg0clmhvf2f2nzt4w6")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	sendMsg := msg.CreateSendMsg(keyManager1.GetAddr(), types.Coins{types.Coin{Denom: "BNB", Amount: 100000000000000}}, []msg.Transfer{{receiverAddr, types.Coins{types.Coin{Denom: "BNB", Amount: 100000000000000}}}})
	stdTx := tx.StdSignMsg{
		ChainID:       "binance-chain",
		AccountNumber: 0,
		Sequence:      0,
		Msgs:          []msg.Msg{sendMsg},
		Memo:          "test ledger sign",
		Source:        0,
	}

	if ledgerKey, ok := keyManager1.GetPrivKey().(*ledger.PrivKeyLedgerSecp256k1); ok {
		fmt.Println(fmt.Sprintf("Please verify if the address displayed on your ledger screen is identical to %s", types.AccAddress(keyManager1.GetAddr()).String()))
		fmt.Println("If so, please click confirm button on your ledger device")
		// Before sign with ledger key, you must call ShowSignAddr first
		err := ledgerKey.ShowSignAddr()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		fmt.Println("Invalid ledger keyManager")
		return
	}
	fmt.Println("Please verify transaction data")
	signature, err := keyManager1.GetPrivKey().Sign(stdTx.Bytes())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(fmt.Sprintf("signature: %s", hex.EncodeToString(signature)))
}
