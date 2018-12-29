package keys

import (
	"bytes"
	"testing"

	"github.com/binance-chain/go-sdk/tx"
	"github.com/binance-chain/go-sdk/tx/txmsg"
	"github.com/stretchr/testify/assert"
)

func TestRecoveryFromKeyWordsNoError(t *testing.T) {
	mnemonic := "bottom quick strong ranch section decide pepper broken oven demand coin run jacket curious business achieve mule bamboo remain vote kid rigid bench rubber"
	keyManger, err := NewMnemonicKeyManager(mnemonic)
	assert.NoError(t, err)
	acc := keyManger.GetAddr()
	key := keyManger.GetPrivKey()
	if acc.String() != "bnc1ddt3ls9fjcd8mh69ujdg3fxc89qle2a7a9x6sk" {
		t.Fatalf("RecoveryFromKeyWords get unstable account")
	}
	if key == nil {
		t.Fatalf("Failed to recover private key")
	}
}

func TestRecoveryFromKeyBaseNoError(t *testing.T) {
	file := "testkeystore.json"
	planText := []byte("Test msg")
	keyManager, err := NewKeyStoreKeyManager(file, "Zjubfd@123")
	assert.NoError(t, err)
	sigs, err := keyManager.GetPrivKey().Sign(planText)
	assert.NoError(t, err)
	valid := keyManager.GetPrivKey().PubKey().VerifyBytes(planText, sigs)
	assert.True(t, valid)
}

func TestRecoveryPrivateKeyNoError(t *testing.T) {
	planText := []byte("Test msg")
	priv := "9579fff0cab07a4379e845a890105004ba4c8276f8ad9d22082b2acbf02d884b"
	keyManager, err := NewPrivateKeyManager(priv)
	assert.NoError(t, err)
	sigs, err := keyManager.GetPrivKey().Sign(planText)
	assert.NoError(t, err)
	valid := keyManager.GetPrivKey().PubKey().VerifyBytes(planText, sigs)
	assert.True(t, valid)
}

func TestSignTxNoError(t *testing.T) {
	test1Mnemonic := "swift slam quote sail high remain mandate sample now stamp title among fiscal captain joy puppy ghost arrow attract ozone situate install gain mean"
	test2Mnemonic := "bottom quick strong ranch section decide pepper broken oven demand coin run jacket curious business achieve mule bamboo remain vote kid rigid bench rubber"

	test1KeyManger, err := NewMnemonicKeyManager(test1Mnemonic)
	assert.NoError(t, err)
	test2KeyManagr, err := NewMnemonicKeyManager(test2Mnemonic)
	assert.NoError(t, err)

	test1Addr := test1KeyManger.GetAddr()
	test2Addr := test2KeyManagr.GetAddr()
	testCases := []struct {
		msg         txmsg.Msg
		keyManager  KeyManager
		accountNUm  int64
		sequence    int64
		expectHexTx string
		errMsg      string
	}{
		{txmsg.CreateSendMsg(test1Addr, test2KeyManagr.GetAddr(),
			txmsg.Coins{txmsg.Coin{Denom: "BNB", Amount: 100000000000000}}),
			test1KeyManger,
			0,
			1,
			"c601f0625dee0a522a2c87fa0a250a141d0e3086e8e4e0a53c38a90d55bd58b34d57d2fa120d0a03424e42108080d287e2bc2d12250a146b571fc0a9961a7ddf45e49a88a4d83941fcabbe120d0a03424e42108080d287e2bc2d126c0a26eb5ae98721027e69d96640300433654e016d218a8d7ffed751023d8efe81e55dedbd6754c9711240a4d54e307ae18bd8783ab4dca7b220f9a57fb944e1bf77d5f50bea6dec00699760ccbeb3fb67df4f58b3a2edbf6b884476ac0beba1662bfcc51989789be4351a2002",
			"send message sign error",
		},
		{
			txmsg.NewTokenIssueMsg(test2Addr, "Bitcoin", "BTC", 1000000000000000, true),
			test2KeyManagr,
			1,
			0,
			"a701f0625dee0a3317efab800a146b571fc0a9961a7ddf45e49a88a4d83941fcabbe1207426974636f696e1a03425443208080b4ccd4dfc6032801126c0a26eb5ae9872103d8f33449356d58b699f6b16a498bd391aa5e051085415d0fe1873939bc1d2e3a124005083827b50bbb538ca1c1d43cc269bac51db692687cfd143d5c9ab5e361a64614f877f7790aa34ca45df006a794d94d1a5fe18057ed6b5be0ac8a187301a0bb1802",
			"issue message sign error",
		},
		{txmsg.NewMsgSubmitProposal("list BTC/BNB", "{\"base_asset_symbol\":\"BTC-86A\",\"quote_asset_symbol\":\"BNB\",\"init_price\":100000000,\"description\":\"list BTC/BNB\",\"expire_time\":\"2018-12-24T00:46:05+08:00\"}", txmsg.ProposalTypeListTradingPair, test1Addr, txmsg.Coins{txmsg.Coin{Denom: "BNB", Amount: 200000000000}}),
			test1KeyManger,
			0,
			2,
			"c802f0625dee0ad301b42d614e0a0c6c697374204254432f424e421298017b22626173655f61737365745f73796d626f6c223a224254432d383641222c2271756f74655f61737365745f73796d626f6c223a22424e42222c22696e69745f7072696365223a3130303030303030302c226465736372697074696f6e223a226c697374204254432f424e42222c226578706972655f74696d65223a22323031382d31322d32345430303a34363a30352b30383a3030227d180422141d0e3086e8e4e0a53c38a90d55bd58b34d57d2fa2a0c0a03424e421080c0ee8ed20b126c0a26eb5ae98721027e69d96640300433654e016d218a8d7ffed751023d8efe81e55dedbd6754c9711240594b7e2e44cc210815d68cbea769666de3d2053b35a3e5d120016350f1ae2b8d4c8ea1ce124dd1223a8c646703b35d21776880f7812741821930e16cbbc227d32004",
			"submit proposal sign error",
		},
		{
			txmsg.NewMsgVote(test1Addr, 1, txmsg.OptionYes),
			test1KeyManger,
			0,
			3,
			"9201f0625dee0a1ea1cadd36080212141d0e3086e8e4e0a53c38a90d55bd58b34d57d2fa1801126c0a26eb5ae98721027e69d96640300433654e016d218a8d7ffed751023d8efe81e55dedbd6754c9711240db2d60b2b8fd86eb8183a6ad146ed476abe814bd256516c4b22df027fbf9717d01a269ea5d04878a70a20db1573de5853580a52025049103772d1baf5d52a3dd2006",
			"vote proposal sign error",
		},
		{
			txmsg.NewDexListMsg(test2Addr, 1, "BTC-86A", "BNB", 100000000),
			test2KeyManagr,
			1,
			2,
			"a501f0625dee0a2fb41de13f0a146b571fc0a9961a7ddf45e49a88a4d83941fcabbe10021a074254432d3836412203424e42288084af5f126e0a26eb5ae9872103d8f33449356d58b699f6b16a498bd391aa5e051085415d0fe1873939bc1d2e3a12409933cf7c386c24bb178b071e0a59fb1d2e1357acfa444722ca0689551c9ba3816494e171d285dec9322e2fdc89c57aebe59af74f35db5c5da1c9738502b8f97c18022004",
			"List tradimg sign error",
		},
		{txmsg.NewCreateOrderMsg(test1Addr, "1D0E3086E8E4E0A53C38A90D55BD58B34D57D2FA-5", 1, "BTC-86A_BNB", 100000000, 1000000000),
			test1KeyManger,
			0,
			4,
			"d801f0625dee0a64ce6dc0430a141d0e3086e8e4e0a53c38a90d55bd58b34d57d2fa122a314430453330383645384534453041353343333841393044353542443538423334443537443246412d351a0b4254432d3836415f424e4220042802308084af5f3880a8d6b9074002126c0a26eb5ae98721027e69d96640300433654e016d218a8d7ffed751023d8efe81e55dedbd6754c971124014ef51e96c0471057f5ad619deaf55acb35364a1df3254284d3c05913ce8cb0870315b2562a5fe0a1d5ea7cdb26fcd6f1d6ba38188e5d13b1abddf5e91e538412008",
			"Create order sign error",
		},
		{
			txmsg.NewCancelOrderMsg(test1Addr, "BTC-86A_BNB", "1D0E3086E8E4E0A53C38A90D55BD58B34D57D2FA-5", "1D0E3086E8E4E0A53C38A90D55BD58B34D57D2FA-5"),
			test1KeyManger,
			0,
			5,
			"f301f0625dee0a7f166e681b0a141d0e3086e8e4e0a53c38a90d55bd58b34d57d2fa120b4254432d3836415f424e421a2a314430453330383645384534453041353343333841393044353542443538423334443537443246412d35222a314430453330383645384534453041353343333841393044353542443538423334443537443246412d35126c0a26eb5ae98721027e69d96640300433654e016d218a8d7ffed751023d8efe81e55dedbd6754c97112408ab86cb8a06db80bb7ed59e8371ff633919e0b642ae0609c44759b7cd1f4d0bf2d9f1230234dd6c02bb75d787afd37db9a6594feb93c7797569960b393d71b6b200a",
			"Cancel order sign error",
		},
		{
			txmsg.NewFreezeMsg(test1Addr, "BNB", 100000000),
			test1KeyManger,
			0,
			10,
			"9a01f0625dee0a26e774b32d0a200a141d0e3086e8e4e0a53c38a90d55bd58b34d57d2fa1203424e42188084af5f126c0a26eb5ae98721027e69d96640300433654e016d218a8d7ffed751023d8efe81e55dedbd6754c9711240b8190928bbfb3bed2bf968e37c022bd254f65de5554bec250349303b79a4c5d0367a76fd7a544f7d2c47af014aadd010ce249236a22f229ee9862539c5d37a7c2014",
			"Freeze token sign error",
		},
		{
			txmsg.NewUnfreezeMsg(test1Addr, "BNB", 100000000),
			test1KeyManger,
			0,
			11,
			"9a01f0625dee0a266515ff0d0a200a141d0e3086e8e4e0a53c38a90d55bd58b34d57d2fa1203424e42188084af5f126c0a26eb5ae98721027e69d96640300433654e016d218a8d7ffed751023d8efe81e55dedbd6754c9711240f219b84d2b2e5c8bd361330f9d9d43d926d8dad11684116eeff11d4a4099139279e8957f00d09f8f7567350dbbe3029ae52ba0fd65818ca89840bb0aad517e762016",
			"Unfreeze token sign error",
		},
		{
			txmsg.NewTokenBurnMsg(test1Addr, "BNB", 100000000),
			test1KeyManger,
			0,
			12,
			"9a01f0625dee0a267ed2d2a00a200a141d0e3086e8e4e0a53c38a90d55bd58b34d57d2fa1203424e42188084af5f126c0a26eb5ae98721027e69d96640300433654e016d218a8d7ffed751023d8efe81e55dedbd6754c97112406eb54c161f39a613e533d4a94970d4d9c5894882543d692bb038cbefcb14ae747d15908770b298189ce92a427426f31def19c3d700efcc764c455ee1af0617c22018",
			"Burn token sign error",
		},
		{
			txmsg.NewMintMsg(test2Addr, "BTC-86A", 100000000),
			test2KeyManagr,
			1,
			5,
			"9e01f0625dee0a28467e08290a146b571fc0a9961a7ddf45e49a88a4d83941fcabbe12074254432d383641188084af5f126e0a26eb5ae9872103d8f33449356d58b699f6b16a498bd391aa5e051085415d0fe1873939bc1d2e3a124075028fb03d98557483853afd12c4d06a642d5ac0f7534bfb8f1bfff367950bdf7365c49e7cfc7e7393c5b4d66154adc9a9b3f35adcf455e114b219273c69d0a51802200a",
			"Mint token sign error",
		},
	}
	for _, c := range testCases {
		signMsg := tx.StdSignMsg{
			ChainID:       "bnbchain-1000",
			AccountNumber: c.accountNUm,
			Sequence:      c.sequence,
			Memo:          "",
			Msgs:          []txmsg.Msg{c.msg},
			Source:        0,
		}

		signResult, err := c.keyManager.Sign(signMsg)
		assert.NoError(t, err)
		expectHexTx := c.expectHexTx
		assert.True(t, bytes.Equal(signResult, []byte(expectHexTx)), c.errMsg)
	}
}
