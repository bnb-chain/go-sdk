package types

import (
	"github.com/cosmos/cosmos-sdk/types"
)

type Dec = types.Dec

var (
	NewDecFromStr  = types.NewDecFromStr
	ZeroDec        = types.ZeroDec
	OneDec         = types.OneDec
	NewDecWithPrec = types.NewDecWithPrec
	NewDec         = types.NewDec
)
