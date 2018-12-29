package tx

import (
	"github.com/binance-chain/go-sdk/tx/txmsg"
)

const Source int64 = 2

type Tx interface {

	// Gets the Msg.
	GetMsgs() []txmsg.Msg
}

// StdTx def
type StdTx struct {
	Msgs       []txmsg.Msg    `json:"msg"`
	Signatures []StdSignature `json:"signatures"`
	Memo       string         `json:"memo"`
	Source     int64          `json:"source"`
	Data       []byte         `json:"data"`
}

// NewStdTx to instantiate an instance
func NewStdTx(msgs []txmsg.Msg, sigs []StdSignature, memo string, source int64, data []byte) StdTx {
	return StdTx{
		Msgs:       msgs,
		Signatures: sigs,
		Memo:       memo,
		Source:     source,
		Data:       data,
	}
}

// GetMsgs def
func (tx StdTx) GetMsgs() []txmsg.Msg { return tx.Msgs }
