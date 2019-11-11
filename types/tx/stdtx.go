package tx

import (
	"github.com/binance-chain/go-sdk/types/msg"
)

const Source int64 = 0

type Tx interface {

	// Gets the Msg.
	GetMsgs() []msg.Msg
}

// StdTx def
type StdTx struct {
	Msgs       []msg.Msg      `json:"msg"`
	Signatures []StdSignature `json:"signatures"`
	Memo       string         `json:"memo"`
	Source     int64          `json:"source"`
	Data       []byte         `json:"data"`
}

// NewStdTx to instantiate an instance
func NewStdTx(msgs []msg.Msg, sigs []StdSignature, memo string, source int64, data []byte) StdTx {
	return StdTx{
		Msgs:       msgs,
		Signatures: sigs,
		Memo:       memo,
		Source:     source,
		Data:       data,
	}
}

// GetMsgs def
func (tx StdTx) GetMsgs() []msg.Msg { return tx.Msgs }
