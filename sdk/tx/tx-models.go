package tx

import (
	"encoding/json"

	"./msgs"
	amino "github.com/tendermint/go-amino"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
)

type Codec = amino.Codec

var Cdc *Codec

func init() {
	cdc := amino.NewCodec()
	cryptoAmino.RegisterAmino(cdc)
	Cdc = cdc.Seal()
}

type StdSignDoc struct {
	ChainID       string            `json:"chain_id"`
	AccountNumber int64             `json:"account_number"`
	Sequence      int64             `json:"sequence"`
	Memo          string            `json:"memo"`
	Fee           json.RawMessage   `json:"fee"`
	Msgs          []json.RawMessage `json:"msgs"`
}

type StdSignMsg struct {
	ChainID       string
	AccountNumber int64
	Sequence      int64
	Memo          string
	Fee           StdFee
	Msgs          []msgs.Msg
}

// Standard Signature
type StdSignature struct {
	tmcrypto.PubKey `json:"pub_key"` // optional
	Signature       []byte           `json:"signature"`
	AccountNumber   int64            `json:"account_number"`
	Sequence        int64            `json:"sequence"`
}

// Bytes gets message bytes
func (msg StdSignMsg) Bytes() []byte {
	return StdSignBytes(msg.ChainID, msg.AccountNumber, msg.Sequence, msg.Fee, msg.Msgs, msg.Memo)
}

// StdSignBytes returns the bytes to sign for a transaction.
func StdSignBytes(chainID string, accnum int64, sequence int64, fee StdFee, msgs []msgs.Msg, memo string) []byte {
	var msgsBytes []json.RawMessage
	for _, msg := range msgs {
		msgsBytes = append(msgsBytes, json.RawMessage(msg.GetSignBytes()))
	}
	bz, err := Cdc.MarshalJSON(StdSignDoc{
		ChainID:       chainID,
		AccountNumber: accnum,
		Sequence:      sequence,
		Memo:          memo,
		Fee:           json.RawMessage(fee.Bytes()),
		Msgs:          msgsBytes,
	})
	if err != nil {
		panic(err)
	}
	return bz //sdk.MustSortJSON(bz)
}

type StdFee struct {
	Amount Coins `json:"amount"`
	Gas    int64 `json:"gas"`
}

// fee bytes for signing later
func (fee StdFee) Bytes() []byte {
	if len(fee.Amount) == 0 {
		fee.Amount = Coins{}
	}
	bz, err := Cdc.MarshalJSON(fee)
	if err != nil {
		panic(err)
	}
	return bz
}

type Coin struct {
	Denom  string `json:"denom"`
	Amount int64  `json:"amount"`
}

type Coins []Coin

type Tx struct{}
