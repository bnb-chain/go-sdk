package tx

import (
	"./msgs"
)

type StdTx struct {
	Msgs       []msgs.Msg     `json:"msg"`
	Fee        StdFee         `json:"fee"`
	Signatures []StdSignature `json:"signatures"`
	Memo       string         `json:"memo"`
}

func NewStdTx(msgs []msgs.Msg, fee StdFee, sigs []StdSignature, memo string) StdTx {
	return StdTx{
		Msgs:       msgs,
		Fee:        fee,
		Signatures: sigs,
		Memo:       memo,
	}
}
func (tx StdTx) GetMemo() string               { return tx.Memo }
func (tx StdTx) GetMsgs() []msgs.Msg           { return tx.Msgs }
func (tx StdTx) GetSignatures() []StdSignature { return tx.Signatures }
func (tx StdTx) GetSigners() []msgs.AccAddress {
	seen := map[string]bool{}
	var signers []msgs.AccAddress
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
