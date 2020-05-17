package types

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
)

// number of decimal places
const (
	Precision = 8

	// bytes required to represent the above precision
	// ceil(log2(9999999999))
	DecimalPrecisionBits = 34
)

var (
	precisionReuse       = new(big.Int).Exp(big.NewInt(10), big.NewInt(Precision), nil).Int64()
	precisionMultipliers []int64
	zeroInt              = big.NewInt(0)
	oneInt               = big.NewInt(1)
	tenInt               = big.NewInt(10)
)

// Set precision multipliers
func init() {
	precisionMultipliers = make([]int64, Precision+1)
	for i := 0; i <= Precision; i++ {
		precisionMultipliers[i] = calcPrecisionMultiplier(int64(i))
	}
}

func precisionInt() int64 {
	return precisionReuse
}

type Dec struct {
	int64 `json:"int"`
}

func (d Dec) String() string {
	return strconv.FormatInt(d.int64, 10)
}

func (d Dec) MarshalText() ([]byte, error) {
	return []byte(strconv.FormatInt(d.int64, 10)), nil
}

func (d *Dec) UnmarshalText(text []byte) error {
	v, err := strconv.ParseInt(string(text), 10, 64)
	d.int64 = v
	return err
}

// requires a valid JSON string - strings quotes and calls UnmarshalText
func (d *Dec) UnmarshalAmino(v int64) (err error) {
	d.int64 = v
	return nil
}
func (d Dec) MarshalAmino() (int64, error) {
	return d.int64, nil
}

// MarshalJSON marshals the decimal
func (d Dec) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON defines custom decoding scheme
func (d *Dec) UnmarshalJSON(bz []byte) error {
	var text string
	err := json.Unmarshal(bz, &text)
	if err != nil {
		return err
	}
	// TODO: Reuse dec allocation
	newDec, err := NewDecFromStr(text)
	if err != nil {
		return err
	}
	d.int64 = newDec.int64
	return nil
}

func NewDecFromStr(str string) (d Dec, err error) {
	value, parseErr := strconv.ParseInt(str, 10, 64)
	if parseErr != nil {
		return d, fmt.Errorf("bad string to integer conversion, input string: %v, error: %v", str, parseErr)
	}
	return Dec{value}, nil
}

//nolint
func (d Dec) IsNil() bool       { return false }               // is decimal nil
func (d Dec) IsZero() bool      { return d.int64 == 0 }        // is equal to zero
func (d Dec) Equal(d2 Dec) bool { return d.int64 == d2.int64 } // equal decimals
func (d Dec) GT(d2 Dec) bool    { return d.int64 > d2.int64 }  // greater than
func (d Dec) GTE(d2 Dec) bool   { return d.int64 >= d2.int64 } // greater than or equal
func (d Dec) LT(d2 Dec) bool    { return d.int64 < d2.int64 }  // less than
func (d Dec) LTE(d2 Dec) bool   { return d.int64 <= d2.int64 } // less than or equal
func (d Dec) Neg() Dec          { return Dec{-d.int64} }       // reverse the decimal sign
func (d Dec) Abs() Dec {
	if d.int64 < 0 {
		return d.Neg()
	}
	return d
}

// subtraction
func (d Dec) Sub(d2 Dec) Dec {
	c := d.int64 - d2.int64
	if (c < d.int64) != (d2.int64 > 0) {
		panic("Int overflow")
	}
	return Dec{c}
}

// nolint - common values
func ZeroDec() Dec { return Dec{0} }
func OneDec() Dec  { return Dec{precisionInt()} }
func NewDecWithPrec(i, prec int64) Dec {
	if i == 0 {
		return Dec{0}
	}
	c := i * precisionMultiplier(prec)
	if c/i != precisionMultiplier(prec) {
		panic("Int overflow")
	}
	return Dec{c}
}

// get the precision multiplier, do not mutate result
func precisionMultiplier(prec int64) int64 {
	if prec > Precision {
		panic(fmt.Sprintf("too much precision, maximum %v, provided %v", Precision, prec))
	}
	return precisionMultipliers[prec]
}

// calculate the precision multiplier
func calcPrecisionMultiplier(prec int64) int64 {
	if prec > Precision {
		panic(fmt.Sprintf("too much precision, maximum %v, provided %v", Precision, prec))
	}
	zerosToAdd := Precision - prec
	multiplier := new(big.Int).Exp(tenInt, big.NewInt(zerosToAdd), nil).Int64()
	return multiplier
}
