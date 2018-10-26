package tx

import (
	tmcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
)

// Sign returns signature
func (tx *Tx) Sign(privKeyBytes []byte, signMsg StdSignMsg) (txBytes []byte, err error) {
	priv, err := cryptoAmino.PrivKeyFromBytes(privKeyBytes)
	if err != nil {
		return nil, err
	}

	sig, err := priv.Sign(tmcrypto.Sha256(signMsg.Bytes()))
	if err != nil {
		return nil, err
	}

	// return sig, nil
	sigs := []StdSignature{{
		PubKey:        priv.PubKey(),
		AccountNumber: signMsg.AccountNumber,
		Sequence:      signMsg.Sequence,
		Signature:     sig,
	}}

	stdTx := NewStdTx(signMsg.Msgs, signMsg.Fee, sigs, signMsg.Memo)
	return Cdc.MarshalBinaryLengthPrefixed(stdTx)
}
