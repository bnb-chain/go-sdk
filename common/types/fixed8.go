package types

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/bnb-chain/node/common/utils"
)

var (
	Fixed8Decimals = utils.Fixed8Decimals
	Fixed8One      = utils.Fixed8One
	Fixed8Zero     = utils.NewFixed8(0)
)

type Fixed8 = utils.Fixed8

var (
	NewFixed8          = utils.NewFixed8
	Fixed8DecodeString = utils.Fixed8DecodeString
)

type Double float64

func (n *Double) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		if p, err := strconv.ParseFloat(s, 64); err == nil {
			*n = Double(p)
		} else {
			return err
		}
	} else {
		return err
	}
	return nil
}

func (n *Double) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%.8f", float64(*n)))
}
