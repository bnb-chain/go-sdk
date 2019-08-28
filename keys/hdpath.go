package keys

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcec"
)

const (
	BIPPurpose  = 44
	BIPCoinType = 714
	BIPChange   = false
	BIP44Prefix = "44'/714'/"
	PartialPath = "0'/0/0"
	FullPath    = BIP44Prefix + PartialPath
)

// BIP44Params wraps BIP 44 params (5 level BIP 32 path).
// To receive a canonical string representation ala
// m / purpose' / coin_type' / account' / change / address_index
// call String() on a BIP44Params instance.
type BIP44Params struct {
	purpose    uint32
	coinType   uint32
	account    uint32
	change     bool
	addressIdx uint32
}

// NewParams creates a BIP 44 parameter object from the params:
// m / purpose' / coin_type' / account' / change / address_index
func NewParams(purpose, coinType, account uint32, change bool, addressIdx uint32) *BIP44Params {
	return &BIP44Params{
		purpose:    purpose,
		coinType:   coinType,
		account:    account,
		change:     change,
		addressIdx: addressIdx,
	}
}

// NewBinanceBIP44Params creates a BIP 44 parameter object from the params:
// m / 44' / 714' / account' / 0 / address_index
// The fixed parameters (purpose', coin_type', and change) are determined by what was used in the binance ledger app.
func NewBinanceBIP44Params(account uint32, addressIdx uint32) *BIP44Params {
	return NewParams(BIPPurpose, BIPCoinType, account, BIPChange, addressIdx)
}

// Return the BIP44 fields as an array.
func (p BIP44Params) DerivationPath() []uint32 {
	change := uint32(0)
	if p.change {
		change = 1
	}
	return []uint32{
		p.purpose,
		p.coinType,
		p.account,
		change,
		p.addressIdx,
	}
}

func (p BIP44Params) String() string {
	var changeStr string
	if p.change {
		changeStr = "1"
	} else {
		changeStr = "0"
	}
	// m / purpose' / coin_type' / account' / change / address_index
	return fmt.Sprintf("%d'/%d'/%d'/%s/%d",
		p.purpose,
		p.coinType,
		p.account,
		changeStr,
		p.addressIdx)
}

// ComputeMastersFromSeed returns the master public key, master secret, and chain code in hex.
func ComputeMastersFromSeed(seed []byte) (secret [32]byte, chainCode [32]byte) {
	masterSecret := []byte("Bitcoin seed")
	secret, chainCode = i64(masterSecret, seed)

	return
}

// DerivePrivateKeyForPath derives the private key by following the BIP 32/44 path from privKeyBytes,
// using the given chainCode.
func DerivePrivateKeyForPath(privKeyBytes [32]byte, chainCode [32]byte, path string) ([32]byte, error) {
	data := privKeyBytes
	parts := strings.Split(path, "/")
	for _, part := range parts {
		// do we have an apostrophe?
		harden := part[len(part)-1:] == "'"
		// harden == private derivation, else public derivation:
		if harden {
			part = part[:len(part)-1]
		}
		idx, err := strconv.Atoi(part)
		if err != nil {
			return [32]byte{}, fmt.Errorf("invalid BIP 32 path: %s", err)
		}
		if idx < 0 {
			return [32]byte{}, errors.New("invalid BIP 32 path: index negative ot too large")
		}
		data, chainCode = derivePrivateKey(data, chainCode, uint32(idx), harden)
	}
	var derivedKey [32]byte
	n := copy(derivedKey[:], data[:])
	if n != 32 || len(data) != 32 {
		return [32]byte{}, fmt.Errorf("expected a (secp256k1) key of length 32, got length: %v", len(data))
	}

	return derivedKey, nil
}

// derivePrivateKey derives the private key with index and chainCode.
// If harden is true, the derivation is 'hardened'.
// It returns the new private key and new chain code.
// For more information on hardened keys see:
//  - https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki
func derivePrivateKey(privKeyBytes [32]byte, chainCode [32]byte, index uint32, harden bool) ([32]byte, [32]byte) {
	var data []byte
	if harden {
		index = index | 0x80000000
		data = append([]byte{byte(0)}, privKeyBytes[:]...)
	} else {
		// this can't return an error:
		_, ecPub := btcec.PrivKeyFromBytes(btcec.S256(), privKeyBytes[:])
		pubkeyBytes := ecPub.SerializeCompressed()
		data = pubkeyBytes

		/* By using btcec, we can remove the dependency on tendermint/crypto/secp256k1
		pubkey := secp256k1.PrivKeySecp256k1(privKeyBytes).PubKey()
		public := pubkey.(secp256k1.PubKeySecp256k1)
		data = public[:]
		*/
	}
	data = append(data, uint32ToBytes(index)...)
	data2, chainCode2 := i64(chainCode[:], data)
	x := addScalars(privKeyBytes[:], data2[:])
	return x, chainCode2
}

// modular big endian addition
func addScalars(a []byte, b []byte) [32]byte {
	aInt := new(big.Int).SetBytes(a)
	bInt := new(big.Int).SetBytes(b)
	sInt := new(big.Int).Add(aInt, bInt)
	x := sInt.Mod(sInt, btcec.S256().N).Bytes()
	x2 := [32]byte{}
	copy(x2[32-len(x):], x)
	return x2
}

func uint32ToBytes(i uint32) []byte {
	b := [4]byte{}
	binary.BigEndian.PutUint32(b[:], i)
	return b[:]
}

// i64 returns the two halfs of the SHA512 HMAC of key and data.
func i64(key []byte, data []byte) (IL [32]byte, IR [32]byte) {
	mac := hmac.New(sha512.New, key)
	// sha512 does not err
	_, _ = mac.Write(data)
	I := mac.Sum(nil)
	copy(IL[:], I[:32])
	copy(IR[:], I[32:])
	return
}
