package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/binance-chain/go-sdk/common/crypto/ledger"
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

	msg := `eyJhY2NvdW50X251bWJlciI6IjAiLCJjaGFpbl9pZCI6InRlc3QtY2hhaW4tOFRYd05sIiwiZGF0YSI6bnVsbCwibWVtbyI6IiIsIm1z
Z3MiOlt7ImRlc2NyaXB0aW9uIjoie1wiYmFzZV9hc3NldF9zeW1ib2xcIjpcIkVUSC1EQTNcIixcInF1b3RlX2Fzc2V0X3N5bWJvbFwiOlwiQk5CXCIs
XCJpbml0X3ByaWNlXCI6MTAwMDAwMDAwMCxcImRlc2NyaXB0aW9uXCI6XCJsaXN0IEVUSC1EQTMvQk5CXCIsXCJleHBpcmVfdGltZVwiOlwiMjAxOS0w
OC0xMFQwMDoyMzoyNiswODowMFwifSIsImluaXRpYWxfZGVwb3NpdCI6W3siYW1vdW50IjoiMjAwMDAwMDAwMDAwIiwiZGVub20iOiJCTkIifV0sInBy
b3Bvc2FsX3R5cGUiOiJMaXN0VHJhZGluZ1BhaXIiLCJwcm9wb3NlciI6ImJuYjF3dDVqaHduN3RxeGZzeWNtZ2RmY3E0YTYzaDN1Y2tkYzAzempneCIs
InRpdGxlIjoibGlzdCBFVEgtREEzL0JOQiIsInZvdGluZ19wZXJpb2QiOiI2MDAwMDAwMDAwMCJ9XSwic2VxdWVuY2UiOiI2Iiwic291cmNlIjoiMCJ9`
	decoded, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if ledgerKey, ok := keyManager.GetPrivKey().(*ledger.PrivKeyLedgerSecp256k1); ok {
		// Before sign with ledger key, you must call ShowSignAddr first
		err := ledgerKey.ShowSignAddr()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	signature, err := keyManager.GetPrivKey().Sign(decoded)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(fmt.Sprintf("signature: %s", hex.EncodeToString(signature)))
}
