package main

import (
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
	fmt.Println(keyManager.GetAddr())
}
