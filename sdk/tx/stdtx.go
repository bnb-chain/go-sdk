package tx

import (
	"./txmsg"
)

// StdTx def
type StdTx struct {
	Msgs       []txmsg.Msg    `json:"msg"`
	Fee        StdFee         `json:"fee"`
	Signatures []StdSignature `json:"signatures"`
	Memo       string         `json:"memo"`
}

// NewStdTx to instantiate an instance
func NewStdTx(msgs []txmsg.Msg, fee StdFee, sigs []StdSignature, memo string) StdTx {
	return StdTx{
		Msgs:       msgs,
		Fee:        fee,
		Signatures: sigs,
		Memo:       memo,
	}
}

// GetMemo def
func (tx StdTx) GetMemo() string { return tx.Memo }

// GetMsgs def
func (tx StdTx) GetMsgs() []txmsg.Msg { return tx.Msgs }

// GetSignatures def
func (tx StdTx) GetSignatures() []StdSignature { return tx.Signatures }

// GetSigners def
func (tx StdTx) GetSigners() []txmsg.AccAddress {
	seen := map[string]bool{}
	var signers []txmsg.AccAddress
	for _, msg := range tx.GetMsgs() {
		for _, addr := range msg.GetSigners() {
			if !seen[addr.String()] {
				signers = append(signers, addr)
				seen[addr.String()] = true
			}
		}
	}
	return signers
}
