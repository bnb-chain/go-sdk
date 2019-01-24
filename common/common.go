package common

import (
	"encoding/json"
	"fmt"
)

func QueryParamToMap(qp interface{}) (map[string]string, error) {
	queryMap := make(map[string]string, 0)
	bz, err := json.Marshal(qp)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bz, &queryMap)
	if err != nil {
		return nil, err
	}
	return queryMap, nil
}

func CombineSymbol(baseAssetSymbol, quoteAssetSymbol string) string {
	return fmt.Sprintf("%s_%s", baseAssetSymbol, quoteAssetSymbol)
}
