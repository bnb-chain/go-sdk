package bsc

import (
	"encoding/json"
	"errors"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"golang.org/x/crypto/sha3"

	"github.com/binance-chain/go-sdk/common/types/bsc/rlp"
)

type Header struct {
	ParentHash  Hash       `json:"parentHash"       gencodec:"required"`
	UncleHash   Hash       `json:"sha3Uncles"       gencodec:"required"`
	Coinbase    Address    `json:"miner"            gencodec:"required"`
	Root        Hash       `json:"stateRoot"        gencodec:"required"`
	TxHash      Hash       `json:"transactionsRoot" gencodec:"required"`
	ReceiptHash Hash       `json:"receiptsRoot"     gencodec:"required"`
	Bloom       Bloom      `json:"logsBloom"        gencodec:"required"`
	Difficulty  int64      `json:"difficulty"       gencodec:"required"`
	Number      int64      `json:"number"           gencodec:"required"`
	GasLimit    uint64     `json:"gasLimit"         gencodec:"required"`
	GasUsed     uint64     `json:"gasUsed"          gencodec:"required"`
	Time        uint64     `json:"timestamp"        gencodec:"required"`
	Extra       []byte     `json:"extraData"        gencodec:"required"`
	MixDigest   Hash       `json:"mixHash"`
	Nonce       BlockNonce `json:"nonce"`
}

// MarshalJSON marshals as JSON.
func (h Header) MarshalJSON() ([]byte, error) {
	type Header struct {
		ParentHash  Hash       `json:"parentHash"       gencodec:"required"`
		UncleHash   Hash       `json:"sha3Uncles"       gencodec:"required"`
		Coinbase    Address    `json:"miner"            gencodec:"required"`
		Root        Hash       `json:"stateRoot"        gencodec:"required"`
		TxHash      Hash       `json:"transactionsRoot" gencodec:"required"`
		ReceiptHash Hash       `json:"receiptsRoot"     gencodec:"required"`
		Bloom       Bloom      `json:"logsBloom"        gencodec:"required"`
		Difficulty  *Big       `json:"difficulty"       gencodec:"required"`
		Number      *Big       `json:"number"           gencodec:"required"`
		GasLimit    Uint64     `json:"gasLimit"         gencodec:"required"`
		GasUsed     Uint64     `json:"gasUsed"          gencodec:"required"`
		Time        Uint64     `json:"timestamp"        gencodec:"required"`
		Extra       Bytes      `json:"extraData"        gencodec:"required"`
		MixDigest   Hash       `json:"mixHash"`
		Nonce       BlockNonce `json:"nonce"`
	}
	var enc Header
	enc.ParentHash = h.ParentHash
	enc.UncleHash = h.UncleHash
	enc.Coinbase = h.Coinbase
	enc.Root = h.Root
	enc.TxHash = h.TxHash
	enc.ReceiptHash = h.ReceiptHash
	enc.Bloom = h.Bloom
	enc.Difficulty = (*Big)(big.NewInt(h.Difficulty))
	enc.Number = (*Big)(big.NewInt(h.Number))
	enc.GasLimit = Uint64(h.GasLimit)
	enc.GasUsed = Uint64(h.GasUsed)
	enc.Time = Uint64(h.Time)
	enc.Extra = h.Extra
	enc.MixDigest = h.MixDigest
	enc.Nonce = h.Nonce
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
func (h *Header) UnmarshalJSON(input []byte) error {
	type Header struct {
		ParentHash  *Hash       `json:"parentHash"       gencodec:"required"`
		UncleHash   *Hash       `json:"sha3Uncles"       gencodec:"required"`
		Coinbase    *Address    `json:"miner"            gencodec:"required"`
		Root        *Hash       `json:"stateRoot"        gencodec:"required"`
		TxHash      *Hash       `json:"transactionsRoot" gencodec:"required"`
		ReceiptHash *Hash       `json:"receiptsRoot"     gencodec:"required"`
		Bloom       *Bloom      `json:"logsBloom"        gencodec:"required"`
		Difficulty  *Big        `json:"difficulty"       gencodec:"required"`
		Number      *Big        `json:"number"           gencodec:"required"`
		GasLimit    *Uint64     `json:"gasLimit"         gencodec:"required"`
		GasUsed     *Uint64     `json:"gasUsed"          gencodec:"required"`
		Time        *Uint64     `json:"timestamp"        gencodec:"required"`
		Extra       *Bytes      `json:"extraData"        gencodec:"required"`
		MixDigest   *Hash       `json:"mixHash"`
		Nonce       *BlockNonce `json:"nonce"`
	}
	var dec Header
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.ParentHash == nil {
		return errors.New("missing required field 'parentHash' for Header")
	}
	h.ParentHash = *dec.ParentHash
	if dec.UncleHash == nil {
		return errors.New("missing required field 'sha3Uncles' for Header")
	}
	h.UncleHash = *dec.UncleHash
	if dec.Coinbase == nil {
		return errors.New("missing required field 'miner' for Header")
	}
	h.Coinbase = *dec.Coinbase
	if dec.Root == nil {
		return errors.New("missing required field 'stateRoot' for Header")
	}
	h.Root = *dec.Root
	if dec.TxHash == nil {
		return errors.New("missing required field 'transactionsRoot' for Header")
	}
	h.TxHash = *dec.TxHash
	if dec.ReceiptHash == nil {
		return errors.New("missing required field 'receiptsRoot' for Header")
	}
	h.ReceiptHash = *dec.ReceiptHash
	if dec.Bloom == nil {
		return errors.New("missing required field 'logsBloom' for Header")
	}
	h.Bloom = *dec.Bloom
	if dec.Difficulty == nil {
		return errors.New("missing required field 'difficulty' for Header")
	}
	h.Difficulty = dec.Difficulty.ToInt().Int64()
	if dec.Number == nil {
		return errors.New("missing required field 'number' for Header")
	}
	h.Number = dec.Number.ToInt().Int64()
	if dec.GasLimit == nil {
		return errors.New("missing required field 'gasLimit' for Header")
	}
	h.GasLimit = uint64(*dec.GasLimit)
	if dec.GasUsed == nil {
		return errors.New("missing required field 'gasUsed' for Header")
	}
	h.GasUsed = uint64(*dec.GasUsed)
	if dec.Time == nil {
		return errors.New("missing required field 'timestamp' for Header")
	}
	h.Time = uint64(*dec.Time)
	if dec.Extra == nil {
		return errors.New("missing required field 'extraData' for Header")
	}
	h.Extra = *dec.Extra
	if dec.MixDigest != nil {
		h.MixDigest = *dec.MixDigest
	}
	if dec.Nonce != nil {
		h.Nonce = *dec.Nonce
	}
	return nil
}

const extraSeal = 65

func (h *Header) GetSignature() ([]byte, error) {
	if len(h.Extra) < extraSeal {
		return nil, errors.New("extra-data 65 byte signature suffix missing")
	}
	signature := h.Extra[len(h.Extra)-extraSeal:]
	return signature, nil
}

func (h *Header) ExtractSignerFromHeader() (signer Address, err error) {
	signature, err := h.GetSignature()
	if err != nil {
		return
	}
	pubKey, err := secp256k1.RecoverPubkey(SealHash(h).Bytes(), signature)
	if err != nil {
		return
	}
	copy(signer[:], Keccak256(pubKey[1:])[12:])
	return
}

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256(data ...[]byte) []byte {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

// SealHash returns the hash of a block prior to it being sealed.
func SealHash(header *Header) (hash Hash) {
	hasher := sha3.NewLegacyKeccak256()
	encodeSigHeader(hasher, header)
	hasher.Sum(hash[:0])
	return hash
}

func encodeSigHeader(w io.Writer, header *Header) {
	err := rlp.Encode(w, []interface{}{
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		big.NewInt(header.Difficulty),
		big.NewInt(header.Number),
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra[:len(header.Extra)-65], // this will panic if extra is too short, should check before calling encodeSigHeader
		header.MixDigest,
		header.Nonce,
	})
	if err != nil {
		panic("can't encode: " + err.Error())
	}
}
