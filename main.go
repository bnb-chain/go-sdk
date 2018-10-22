package main

import (
	"fmt"

	"./sdk"
)

func main() {
	sdk, _ := sdk.NewSDK("http://localhost:8080/api/v1")
	pairs, _ := sdk.GetPairs(100)
	fmt.Println("pairs: ", pairs)
}
