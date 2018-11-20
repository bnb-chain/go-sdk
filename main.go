package main

import (
	"fmt"

	"github.com/BiJie/bnc-go-sdk/sdk"
)

func main() {
	sdk, _ := sdk.NewBncSDK("http://localhost:8080/api/v1")
	markets, _ := sdk.GetMarkets(100)
	fmt.Println("markets: ", markets)
}
