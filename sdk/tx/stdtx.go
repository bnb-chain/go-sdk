package tx

import (
	"./txmsg"
)

type StdTx struct {
	Msgs       []txmsg.Msg    `json:"msg"`
	Fee        StdFee         `json:"fee"`
	Signatures []StdSignature `json:"signatures"`
	Memo       string         `json:"memo"`
}

func NewStdTx(msgs []txmsg.Msg, fee StdFee, sigs []StdSignature, memo string) StdTx {
	return StdTx{
		Msgs:       msgs,
		Fee:        fee,
		Signatures: sigs,
		Memo:       memo,
	}
}
func (tx StdTx) GetMemo() string               { return tx.Memo }
func (tx StdTx) GetMsgs() []txmsg.Msg          { return tx.Msgs }
func (tx StdTx) GetSignatures() []StdSignature { return tx.Signatures }
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
