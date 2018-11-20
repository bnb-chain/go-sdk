package tx

import (
	"encoding/json"

	"github.com/BiJie/bnc-go-sdk/sdk/tx/txmsg"
	tmcrypto "github.com/tendermint/tendermint/crypto"
)

// StdSignDoc def
type StdSignDoc struct {
	ChainID       string            `json:"chain_id"`
	AccountNumber int64             `json:"account_number"`
	Sequence      int64             `json:"sequence"`
	Memo          string            `json:"memo"`
	Fee           json.RawMessage   `json:"fee"`
	Msgs          []json.RawMessage `json:"msgs"`
}

// StdSignMsg def
type StdSignMsg struct {
	AccountNumber int64
	ChainID       string
	Fee           StdFee
	Memo          string
	Msgs          []txmsg.Msg
	Sequence      int64
}

// StdSignature Signature
type StdSignature struct {
	tmcrypto.PubKey    `json:"pub_key"` // optional
	tmcrypto.Signature `json:"signature"`
	AccountNumber      int64 `json:"account_number"`
	Sequence           int64 `json:"sequence"`
}

// Bytes gets message bytes
func (msg StdSignMsg) Bytes() []byte {
	return StdSignBytes(msg.ChainID, msg.AccountNumber, msg.Sequence, msg.Fee, msg.Msgs, msg.Memo)
}

// StdSignBytes returns the bytes to sign for a transaction.
func StdSignBytes(chainID string, accnum int64, sequence int64, fee StdFee, msgs []txmsg.Msg, memo string) []byte {
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
	return MustSortJSON(bz)
}
