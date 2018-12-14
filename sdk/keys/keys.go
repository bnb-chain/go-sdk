package keys

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/BiJie/bnc-go-sdk/sdk/tx"
	"github.com/BiJie/bnc-go-sdk/sdk/tx/txmsg"
	"github.com/cosmos/go-bip39"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"io/ioutil"
	"strings"
)

const (
	defaultBIP39Passphrase = ""
)

type KeyManager interface {
	Sign(tx.StdSignMsg) ([]byte, error)
	GetPrivKey() tmcrypto.PrivKey
	GetAddr() txmsg.AccAddress
}

func NewMnemonicKeyManager(mnemonic string) (KeyManager, error) {
	k := keyManager{}
	err := k.recoveryFromKMnemonic(mnemonic)
	return &k, err
}

func NewKeyStoreKeyManager(file string, auth string) (KeyManager, error) {
	k := keyManager{}
	err := k.recoveryFromKeyStore(file, auth)
	return &k, err
}

func NewPrivateKeyManager(wifKey string) (KeyManager, error) {
	k := keyManager{}
	err := k.recoveryFromPrivateKey(wifKey)
	return &k, err
}

type keyManager struct {
	recoverType string
	privKey     tmcrypto.PrivKey
	addr        txmsg.AccAddress
}

func (m *keyManager) recoveryFromKMnemonic(mnemonic string) error {
	words := strings.Split(mnemonic, " ")
	if len(words) != 12 && len(words) != 24 {
		return fmt.Errorf("mnemonic length should either be 12 or 24")
	}
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, defaultBIP39Passphrase)
	if err != nil {
		return err
	}
	// create master key and derive first key:
	masterPriv, ch := ComputeMastersFromSeed(seed)
	derivedPriv, err := DerivePrivateKeyForPath(masterPriv, ch, FullFundraiserPath)
	if err != nil {
		return err
	}
	priKey := secp256k1.PrivKeySecp256k1(derivedPriv)
	addr := txmsg.AccAddress(priKey.PubKey().Address())
	if err != nil {
		return err
	}
	m.addr = addr
	m.privKey = priKey
	return nil
}

func (m *keyManager) recoveryFromKeyStore(keystoreFile string, auth string) error {
	if auth == "" {
		return fmt.Errorf("Password is missing ")
	}
	keyJson, err := ioutil.ReadFile(keystoreFile)
	if err != nil {
		return err
	}
	var encryptedKey EncryptedKeyJSONV1
	err = json.Unmarshal(keyJson, &encryptedKey)
	if err != nil {
		return err
	}
	keyBytes, err := decryptKeyV1(&encryptedKey, auth)
	if err != nil {
		return err
	}
	if len(keyBytes) != 32 {
		return fmt.Errorf("Len of Keybytes is not equal to 32 ")
	}
	var keyBytesArray [32]byte
	copy(keyBytesArray[:], keyBytes[:32])
	priKey := secp256k1.PrivKeySecp256k1(keyBytesArray)
	addr := txmsg.AccAddress(priKey.PubKey().Address())
	m.addr = addr
	m.privKey = priKey
	return nil
}

func (m *keyManager) recoveryFromPrivateKey(privateKey string) error {
	priBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return err
	}

	if len(priBytes) != 32 {
		return fmt.Errorf("Len of Keybytes is not equal to 32 ")
	}
	var keyBytesArray [32]byte
	copy(keyBytesArray[:], priBytes[:32])
	priKey := secp256k1.PrivKeySecp256k1(keyBytesArray)
	addr := txmsg.AccAddress(priKey.PubKey().Address())
	m.addr = addr
	m.privKey = priKey
	return nil
}

func (m *keyManager) Sign(msg tx.StdSignMsg) ([]byte, error) {
	sig, err := m.makeSignature(msg)
	if err != nil {
		return nil, err
	}
	newTx := tx.NewStdTx(msg.Msgs, []tx.StdSignature{sig}, msg.Memo, msg.Source, msg.Data)
	bz, err := tx.Cdc.MarshalBinary(&newTx)
	if err != nil {
		return nil, err
	}
	//return bz, nil
	return []byte(hex.EncodeToString(bz)), nil
}

func (m *keyManager) GetPrivKey() tmcrypto.PrivKey {
	return m.privKey
}

func (m *keyManager) GetAddr() txmsg.AccAddress {
	return m.addr
}

func (m *keyManager) makeSignature(msg tx.StdSignMsg) (sig tx.StdSignature, err error) {
	if err != nil {
		return
	}
	sigBytes, err := m.privKey.Sign(msg.Bytes())
	if err != nil {
		return
	}
	return tx.StdSignature{
		AccountNumber: msg.AccountNumber,
		Sequence:      msg.Sequence,
		PubKey:        m.privKey.PubKey(),
		Signature:     sigBytes,
	}, nil
}
