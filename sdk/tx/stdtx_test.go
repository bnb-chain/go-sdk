package tx

import (
	"testing"

	"github.com/BiJie/bnc-go-sdk/sdk/tx/txmsg"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func TestStdTx(t *testing.T) {
	priv := ed25519.GenPrivKey()
	addr := txmsg.AccAddress(priv.PubKey().Address())
	msgs := []txmsg.Msg{txmsg.NewTestMsg(addr)}
	fee := NewStdFee(100, Coin{"BNB", 500})
	sigs := []StdSignature{}

	tx := NewStdTx(msgs, fee, sigs, "")
	require.Equal(t, msgs, tx.GetMsgs())
	require.Equal(t, sigs, tx.GetSignatures())

	feePayer := msgs[0].GetSigners()[0]
	require.Equal(t, addr, feePayer)
}
