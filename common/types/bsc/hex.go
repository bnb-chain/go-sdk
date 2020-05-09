package bsc

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
)

var (
	ErrEmptyString   = &DecError{"empty hex string"}
	ErrSyntax        = &DecError{"invalid hex string"}
	ErrMissingPrefix = &DecError{"hex string without 0x prefix"}
	ErrOddLength     = &DecError{"hex string of odd length"}
	ErrEmptyNumber   = &DecError{"hex string \"0x\""}
	ErrLeadingZero   = &DecError{"hex number with leading zero digits"}
	ErrUint64Range   = &DecError{"hex number > 64 bits"}
	ErrBig256Range   = &DecError{"hex number > 256 bits"}
)

type DecError struct{ msg string }

func (err DecError) Error() string { return err.msg }

const UintBits = 32 << (uint64(^uint(0)) >> 63)

var BigWordNibbles int

func init() {
	// This is a weird way to compute the number of nibbles required for big.Word.
	// The usual way would be to use constant arithmetic but go vet can't handle that.
	b, _ := new(big.Int).SetString("FFFFFFFFFF", 16)
	switch len(b.Bits()) {
	case 1:
		BigWordNibbles = 16
	case 2:
		BigWordNibbles = 8
	default:
		panic("weird big.Word size")
	}
}

const BadNibble = ^uint64(0)

func HexDecodeNibble(in byte) uint64 {
	switch {
	case in >= '0' && in <= '9':
		return uint64(in - '0')
	case in >= 'A' && in <= 'F':
		return uint64(in - 'A' + 10)
	case in >= 'a' && in <= 'f':
		return uint64(in - 'a' + 10)
	default:
		return BadNibble
	}
}

// DecodeUint64 decodes a hex string with 0x prefix as a quantity.
func HexDecodeUint64(input string) (uint64, error) {
	raw, err := checkNumber(input)
	if err != nil {
		return 0, err
	}
	dec, err := strconv.ParseUint(raw, 16, 64)
	if err != nil {
		err = mapError(err)
	}
	return dec, err
}

// EncodeUint64 encodes i as a hex string with 0x prefix.
func HexEncodeUint64(i uint64) string {
	enc := make([]byte, 2, 10)
	copy(enc, "0x")
	return string(strconv.AppendUint(enc, i, 16))
}

// EncodeBig encodes bigint as a hex string with 0x prefix.
// The sign of the integer is ignored.
func HexEncodeBig(bigint *big.Int) string {
	nbits := bigint.BitLen()
	if nbits == 0 {
		return "0x0"
	}
	return fmt.Sprintf("%#x", bigint)
}

func checkNumber(input string) (raw string, err error) {
	if len(input) == 0 {
		return "", ErrEmptyString
	}
	if !has0xPrefix(input) {
		return "", ErrMissingPrefix
	}
	input = input[2:]
	if len(input) == 0 {
		return "", ErrEmptyNumber
	}
	if len(input) > 1 && input[0] == '0' {
		return "", ErrLeadingZero
	}
	return input, nil
}

func mapError(err error) error {
	if err, ok := err.(*strconv.NumError); ok {
		switch err.Err {
		case strconv.ErrRange:
			return ErrUint64Range
		case strconv.ErrSyntax:
			return ErrSyntax
		}
	}
	if _, ok := err.(hex.InvalidByteError); ok {
		return ErrSyntax
	}
	if err == hex.ErrLength {
		return ErrOddLength
	}
	return err
}

func hexEncode(b []byte) string {
	enc := make([]byte, len(b)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], b)
	return string(enc)
}

// Decode decodes a hex string with 0x prefix.
func hexDecode(input string) ([]byte, error) {
	if !has0xPrefix(input) {
		return nil, fmt.Errorf("hex string without 0x prefix")
	}
	return hex.DecodeString(input[2:])
}

func has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}
