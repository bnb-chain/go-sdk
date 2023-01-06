package main

import (
	"fmt"
	"strconv"

	"github.com/bnb-chain/go-sdk/client"
	"github.com/bnb-chain/go-sdk/common/ledger"
	"github.com/bnb-chain/go-sdk/common/types"
	"github.com/bnb-chain/go-sdk/keys"
	"github.com/bnb-chain/go-sdk/types/msg"
)

// To run this example, please make sure your key address have more than 1:BNB on testnet
func main() {
	types.SetNetwork(types.TestNetwork)

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
	keyManager, err := keys.NewLedgerKeyManager(bip44Params.DerivationPath())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	receiverAddr, err := types.AccAddressFromBech32("tbnb15339dcwlq5nza4atfmqxfx6mhamywz35he2cvv")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	dexClient, err := client.NewDexClient("testnet-dex.binance.org:443", types.TestNetwork, keyManager)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	account, err := dexClient.GetAccount(keyManager.GetAddr().String())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	floatAmount := 0.0
	for _, coin := range account.Balances {
		if coin.Symbol == "BNB" {
			fmt.Println(fmt.Sprintf("Your account has %s:BNB", coin.Free))
			floatAmount, err = strconv.ParseFloat(coin.Free.String(), 64)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			break
		}
	}
	if floatAmount <= 1.0 {
		fmt.Println("Your account doesn't have enough bnb")
	}

	fmt.Println(fmt.Sprintf("Please verify sign key address (%s) and transaction data", types.AccAddress(keyManager.GetAddr()).String()))
	sendResult, err := dexClient.SendToken([]msg.Transfer{{receiverAddr, types.Coins{types.Coin{Denom: "BNB", Amount: 10000000}}}}, true)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(fmt.Sprintf("Send result: %t", sendResult.Ok))
}
