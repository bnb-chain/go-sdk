package keys

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	ctypes "github.com/bnb-chain/go-sdk/common/types"
	"github.com/bnb-chain/go-sdk/types/msg"
	"github.com/bnb-chain/go-sdk/types/tx"
)

func TestRecoveryFromKeyWordsNoError(t *testing.T) {
	mnemonic := "bottom quick strong ranch section decide pepper broken oven demand coin run jacket curious business achieve mule bamboo remain vote kid rigid bench rubber"
	keyManager, err := NewMnemonicKeyManager(mnemonic)
	assert.NoError(t, err)
	acc := keyManager.GetAddr()
	key := keyManager.GetPrivKey()
	assert.Equal(t, "bnb1ddt3ls9fjcd8mh69ujdg3fxc89qle2a7km33aa", acc.String())
	assert.NotNil(t, key)
	customPathKey, err := NewMnemonicPathKeyManager(mnemonic, "1'/1/1")
	assert.NoError(t, err)
	assert.Equal(t, "bnb1c67nwp7u5adl7gw0ffn3d47kttcm4crjy9mrye", customPathKey.GetAddr().String())
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

	test1KeyManager, err := NewMnemonicKeyManager(test1Mnemonic)
	assert.NoError(t, err)
	test2KeyManager, err := NewMnemonicKeyManager(test2Mnemonic)
	assert.NoError(t, err)

	test1Addr := test1KeyManager.GetAddr()
	test2Addr := test2KeyManager.GetAddr()
	testCases := []struct {
		msg         msg.Msg
		keyManager  KeyManager
		accountNUm  int64
		sequence    int64
		expectHexTx string
		errMsg      string
	}{
		{msg.CreateSendMsg(test1Addr, ctypes.Coins{ctypes.Coin{Denom: "BNB", Amount: 100000000000000}}, []msg.Transfer{{test2KeyManager.GetAddr(), ctypes.Coins{ctypes.Coin{Denom: "BNB", Amount: 100000000000000}}}}),
			test1KeyManager,
			0,
			1,
			"c601f0625dee0a522a2c87fa0a250a141d0e3086e8e4e0a53c38a90d55bd58b34d57d2fa120d0a03424e42108080e983b1de1612250a146b571fc0a9961a7ddf45e49a88a4d83941fcabbe120d0a03424e42108080e983b1de16126c0a26eb5ae98721027e69d96640300433654e016d218a8d7ffed751023d8efe81e55dedbd6754c97112408b23eecfa8237a27676725173e58154e6c204bb291b31c3b7b507c8f04e2773909ba70e01b54f4bd0bc76669f5712a5a66b9508acdf3aa5e4fde75fbe57622a12001",
			"send message sign error",
		},
		{
			msg.NewTokenIssueMsg(test2Addr, "Bitcoin", "BTC", 1000000000000000, true),
			test2KeyManager,
			1,
			0,
			"a701f0625dee0a3317efab800a146b571fc0a9961a7ddf45e49a88a4d83941fcabbe1207426974636f696e1a034254432080809aa6eaafe3012801126c0a26eb5ae9872103d8f33449356d58b699f6b16a498bd391aa5e051085415d0fe1873939bc1d2e3a12403686586a55f8c50a11ae6f09c35734f09830a566823846f9333c3e53f6d83d4a2cd3a0542c37a8d28b474f563d44223ba6c2b7cf260539b7b85020999ebe2c001801",
			"issue message sign error",
		},
		{msg.NewMsgSubmitProposal("list BTC/BNB", "{\"base_asset_symbol\":\"BTC-86A\",\"quote_asset_symbol\":\"BNB\",\"init_price\":100000000,\"description\":\"list BTC/BNB\",\"expire_time\":\"2018-12-24T00:46:05+08:00\"}", msg.ProposalTypeListTradingPair, test1Addr, ctypes.Coins{ctypes.Coin{Denom: "BNB", Amount: 200000000000}}, time.Second),
			test1KeyManager,
			0,
			2,
			"ce02f0625dee0ad901b42d614e0a0c6c697374204254432f424e421298017b22626173655f61737365745f73796d626f6c223a224254432d383641222c2271756f74655f61737365745f73796d626f6c223a22424e42222c22696e69745f7072696365223a3130303030303030302c226465736372697074696f6e223a226c697374204254432f424e42222c226578706972655f74696d65223a22323031382d31322d32345430303a34363a30352b30383a3030227d180422141d0e3086e8e4e0a53c38a90d55bd58b34d57d2fa2a0c0a03424e421080a0b787e905308094ebdc03126c0a26eb5ae98721027e69d96640300433654e016d218a8d7ffed751023d8efe81e55dedbd6754c9711240ebac8c34f27e9dc0719167c4ad87bc2e3e1437022c3287030425db8f4233c3b80938dfa16f555738ba97e92aa7a15ebb6ac8baa5d799118cccf503302d166df92002",
			"submit proposal sign error",
		},
		{
			msg.NewMsgVote(test1Addr, 1, msg.OptionYes),
			test1KeyManager,
			0,
			3,
			"9201f0625dee0a1ea1cadd36080112141d0e3086e8e4e0a53c38a90d55bd58b34d57d2fa1801126c0a26eb5ae98721027e69d96640300433654e016d218a8d7ffed751023d8efe81e55dedbd6754c97112405a2394615da1d744e991b0cc0f188ac2d7108259e15ba0b992729a8ee54e77e5641dd2ed846d59be468be2daeb628c3a6633187225a0ce85db884c965467baf52003",
			"vote proposal sign error",
		},
		{
			msg.NewDexListMsg(test2Addr, 1, "BTC-86A", "BNB", 100000000),
			test2KeyManager,
			1,
			2,
			"a501f0625dee0a2fb41de13f0a146b571fc0a9961a7ddf45e49a88a4d83941fcabbe10011a074254432d3836412203424e422880c2d72f126e0a26eb5ae9872103d8f33449356d58b699f6b16a498bd391aa5e051085415d0fe1873939bc1d2e3a1240ce897838dd4d70d3c337b62ac1f60ec022bb2bf281fb77eca95adcd64de0a6aa574b128c5f732545c8b0e62dcd1cd90f9898c5f7781ae9d64042859c4e40558b18012002",
			"List tradimg sign error",
		},
		{msg.NewCreateOrderMsg(test1Addr, "1D0E3086E8E4E0A53C38A90D55BD58B34D57D2FA-5", 1, "BTC-86A_BNB", 100000000, 1000000000),
			test1KeyManager,
			0,
			4,
			"d801f0625dee0a64ce6dc0430a141d0e3086e8e4e0a53c38a90d55bd58b34d57d2fa122a314430453330383645384534453041353343333841393044353542443538423334443537443246412d351a0b4254432d3836415f424e42200228013080c2d72f388094ebdc034001126c0a26eb5ae98721027e69d96640300433654e016d218a8d7ffed751023d8efe81e55dedbd6754c97112409fe317e036f2bdc8c87a0138dc52367faef80ea1d6e21a35634b17a82ed7be632c9cb03f865f6f8a6872736ccab716a157f3cb99339afa55686aa455dc134f6a2004",
			"Create order sign error",
		},
		{
			msg.NewCancelOrderMsg(test1Addr, "BTC-86A_BNB", "1D0E3086E8E4E0A53C38A90D55BD58B34D57D2FA-5"),
			test1KeyManager,
			0,
			5,
			"c701f0625dee0a53166e681b0a141d0e3086e8e4e0a53c38a90d55bd58b34d57d2fa120b4254432d3836415f424e421a2a314430453330383645384534453041353343333841393044353542443538423334443537443246412d35126c0a26eb5ae98721027e69d96640300433654e016d218a8d7ffed751023d8efe81e55dedbd6754c9711240fe2fd18630317849bd1d4ae064f8c4fd95f6186bdb61e2b73a5fb5e93ac7794d4a990ba943694659df9d3f49d5312fec020b80148677f3e95fd6d88486bba19d2005",
			"Cancel order sign error",
		},
		{
			msg.NewFreezeMsg(test1Addr, "BNB", 100000000),
			test1KeyManager,
			0,
			10,
			"9801f0625dee0a24e774b32d0a141d0e3086e8e4e0a53c38a90d55bd58b34d57d2fa1203424e421880c2d72f126c0a26eb5ae98721027e69d96640300433654e016d218a8d7ffed751023d8efe81e55dedbd6754c971124013142dc274677af4f09d4be295f1855709b7608e1a9d4cc76aa23103c092ce1915c0ed51fc3c8a8510b57a7a8e8d532c6f5d1159cdb0e7333dda0d0a9e55cac4200a",
			"Freeze token sign error",
		},
		{
			msg.NewUnfreezeMsg(test1Addr, "BNB", 100000000),
			test1KeyManager,
			0,
			11,
			"9801f0625dee0a246515ff0d0a141d0e3086e8e4e0a53c38a90d55bd58b34d57d2fa1203424e421880c2d72f126c0a26eb5ae98721027e69d96640300433654e016d218a8d7ffed751023d8efe81e55dedbd6754c9711240af320bdecb27fe5f7f89abcb9ffc11df2479859fa66586654f36425d908ff7a32921c9af658f86c63db797981a3110db33c6033db017bdca5ca87b1c440c8fc6200b",
			"Unfreeze token sign error",
		},
		{
			msg.NewTokenBurnMsg(test1Addr, "BNB", 100000000),
			test1KeyManager,
			0,
			12,
			"9801f0625dee0a247ed2d2a00a141d0e3086e8e4e0a53c38a90d55bd58b34d57d2fa1203424e421880c2d72f126c0a26eb5ae98721027e69d96640300433654e016d218a8d7ffed751023d8efe81e55dedbd6754c971124066f3e784fa602d7697dd46bf89a17e82db7ab89da1e72b3c253cd14ee073628c7cfbe1b05dab541bd162687e5bd390bfed029c99792e69015c1f86deb399fde6200c",
			"Burn token sign error",
		},
		{
			msg.NewMintMsg(test2Addr, "BTC-86A", 100000000),
			test2KeyManager,
			1,
			5,
			"9e01f0625dee0a28467e08290a146b571fc0a9961a7ddf45e49a88a4d83941fcabbe12074254432d3836411880c2d72f126e0a26eb5ae9872103d8f33449356d58b699f6b16a498bd391aa5e051085415d0fe1873939bc1d2e3a124073b5f00488861a7abdf2274fb719add8b8c9a0bbe16c46bc3ac844671df7c705080d8b061d37bbfcec93d9ed05b3601fde94adadc5086f828402c9f91ce55d1518012005",
			"Mint token sign error",
		},
	}
	for _, c := range testCases {
		signMsg := tx.StdSignMsg{
			ChainID:       "bnbchain-1000",
			AccountNumber: c.accountNUm,
			Sequence:      c.sequence,
			Memo:          "",
			Msgs:          []msg.Msg{c.msg},
			Source:        0,
		}

		rawSignResult, err := c.keyManager.Sign(signMsg)
		signResult := []byte(hex.EncodeToString(rawSignResult))
		assert.NoError(t, err)
		expectHexTx := c.expectHexTx
		assert.True(t, bytes.Equal(signResult, []byte(expectHexTx)), c.errMsg)
	}
}

func TestExportAsKeyStoreNoError(t *testing.T) {
	defer os.Remove("TestGenerateKeyStoreNoError.json")
	km, err := NewKeyManager()
	assert.NoError(t, err)
	encryPlain1, err := km.GetPrivKey().Sign([]byte("test plain"))
	assert.NoError(t, err)
	keyJSONV1, err := km.ExportAsKeyStore("testpassword")
	assert.NoError(t, err)
	bz, err := json.Marshal(keyJSONV1)
	assert.NoError(t, err)
	err = ioutil.WriteFile("TestGenerateKeyStoreNoError.json", bz, 0660)
	assert.NoError(t, err)
	newkm, err := NewKeyStoreKeyManager("TestGenerateKeyStoreNoError.json", "testpassword")
	assert.NoError(t, err)
	encryPlain2, err := newkm.GetPrivKey().Sign([]byte("test plain"))
	assert.NoError(t, err)
	assert.True(t, bytes.Equal(encryPlain1, encryPlain2))
}

func TestExportAsMnemonicNoError(t *testing.T) {
	km, err := NewKeyManager()
	assert.NoError(t, err)
	encryPlain1, err := km.GetPrivKey().Sign([]byte("test plain"))
	assert.NoError(t, err)
	mnemonic, err := km.ExportAsMnemonic()
	assert.NoError(t, err)
	newkm, err := NewMnemonicKeyManager(mnemonic)
	assert.NoError(t, err)
	encryPlain2, err := newkm.GetPrivKey().Sign([]byte("test plain"))
	assert.NoError(t, err)
	assert.True(t, bytes.Equal(encryPlain1, encryPlain2))
	_, err = newkm.ExportAsMnemonic()
	assert.NoError(t, err)
}

func TestExportAsPrivateKeyNoError(t *testing.T) {
	km, err := NewKeyManager()
	assert.NoError(t, err)
	encryPlain1, err := km.GetPrivKey().Sign([]byte("test plain"))
	assert.NoError(t, err)
	pk, err := km.ExportAsPrivateKey()
	assert.NoError(t, err)
	newkm, err := NewPrivateKeyManager(pk)
	assert.NoError(t, err)
	encryPlain2, err := newkm.GetPrivKey().Sign([]byte("test plain"))
	assert.NoError(t, err)
	assert.True(t, bytes.Equal(encryPlain1, encryPlain2))
}

func TestExportAsMnemonicyError(t *testing.T) {
	km, err := NewPrivateKeyManager("9579fff0cab07a4379e845a890105004ba4c8276f8ad9d22082b2acbf02d884b")
	assert.NoError(t, err)
	_, err = km.ExportAsMnemonic()
	assert.Error(t, err)
	file := "testkeystore.json"
	km, err = NewKeyStoreKeyManager(file, "Zjubfd@123")
	assert.NoError(t, err)
	_, err = km.ExportAsMnemonic()
	assert.Error(t, err)
}
