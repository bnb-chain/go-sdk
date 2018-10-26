package tx

import (
	"fmt"
	"testing"

	"./txmsg"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func TestStdSignBytes(t *testing.T) {
	priv := ed25519.GenPrivKey()
	addr := txmsg.AccAddress(priv.PubKey().Address())
	msgs := []txmsg.Msg{txmsg.NewTestMsg(addr)}
	fee := NewStdFee(5000, Coin{"atom", 150})
	signMsg := StdSignMsg{
		3,
		"1234",
		fee,
		"memo",
		msgs,
		6,
	}
	require.Equal(t, fmt.Sprintf("{\"account_number\":\"3\",\"chain_id\":\"1234\",\"fee\":{\"amount\":[{\"amount\":\"150\",\"denom\":\"atom\"}],\"gas\":\"5000\"},\"memo\":\"memo\",\"msgs\":[[\"%s\"]],\"sequence\":\"6\"}", addr), string(signMsg.Bytes()))
}
