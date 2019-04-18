package types

import (
	"encoding/json"
	"fmt"
	"strconv"
)

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