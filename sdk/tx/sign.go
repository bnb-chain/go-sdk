package tx

import (
	tmcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
)

// Sign message, prepare signatures and return HEX format marshalled stdtx
func (tx *Tx) Sign(privKeyBytes []byte, signMsg StdSignMsg) (hexSignedStdTx []byte, err error) {
	// validate all messages
	for _, msg := range signMsg.Msgs {
		err = msg.ValidateBasic()
		if err != nil {
			return nil, err
		}
	}

	priv, err := cryptoAmino.PrivKeyFromBytes(privKeyBytes)
	if err != nil {
		return nil, err
	}

	sig, err := priv.Sign(tmcrypto.Sha256(signMsg.Bytes()))
	if err != nil {
		return nil, err
	}

	sigs := []StdSignature{{
		PubKey:        priv.PubKey(),
		AccountNumber: signMsg.AccountNumber,
		Sequence:      signMsg.Sequence,
		Signature:     sig,
	}}

	stdTx := NewStdTx(signMsg.Msgs, signMsg.Fee, sigs, signMsg.Memo)

	stdTxBytes, err := Cdc.MarshalBinaryLengthPrefixed(stdTx)
	if err != nil {
		return nil, err
	}

	return EncodeHex(stdTxBytes), nil
}
