package main

import (
	"encoding/hex"
	"fmt"

	"github.com/binance-chain/go-sdk/keys"
)

func main() {
	bip44Params := keys.NewFundraiserParams(0, 0)
	keyManager, err := keys.NewLedgerKeyManager(bip44Params.DerivationPath())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(fmt.Sprintf("address: %s", keyManager.GetAddr()))
	fmt.Println(fmt.Sprintf("pubkey: %s", hex.EncodeToString(keyManager.GetPrivKey().PubKey().Bytes())))
}
