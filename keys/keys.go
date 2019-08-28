package keys

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/cosmos/go-bip39"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/sha3"

	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/binance-chain/go-sdk/common"
	"github.com/binance-chain/go-sdk/common/ledger"
	"github.com/binance-chain/go-sdk/common/types"
	ctypes "github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/common/uuid"
	"github.com/binance-chain/go-sdk/types/tx"
	"github.com/tendermint/tendermint/crypto"
)

const (
	defaultBIP39Passphrase = ""
)

type KeyManager interface {
	Sign(tx.StdSignMsg) ([]byte, error)
	GetPrivKey() crypto.PrivKey
	GetAddr() ctypes.AccAddress

	ExportAsMnemonic() (string, error)
	ExportAsPrivateKey() (string, error)
	ExportAsKeyStore(password string) (*EncryptedKeyJSON, error)
}

func NewMnemonicKeyManager(mnemonic string) (KeyManager, error) {
	k := keyManager{}
	err := k.recoveryFromMnemonic(mnemonic, FullPath)
	return &k, err
}

// The full path is  "purpose' / coin_type' / account' / change / address_index".
// "purpose' / coin_type'" is fixed as "44'/714'/", user can customize the rest part.
func NewMnemonicPathKeyManager(mnemonic, keyPath string) (KeyManager, error) {
	k := keyManager{}
	err := k.recoveryFromMnemonic(mnemonic, BIP44Prefix+keyPath)
	return &k, err
}

func NewKeyStoreKeyManager(file string, auth string) (KeyManager, error) {
	k := keyManager{}
	err := k.recoveryFromKeyStore(file, auth)
	return &k, err
}

func NewPrivateKeyManager(priKey string) (KeyManager, error) {
	k := keyManager{}
	err := k.recoveryFromPrivateKey(priKey)
	return &k, err
}

func NewLedgerKeyManager(path ledger.DerivationPath) (KeyManager, error) {
	k := keyManager{}
	err := k.recoveryFromLedgerKey(path)
	return &k, err
}

type keyManager struct {
	privKey  crypto.PrivKey
	addr     ctypes.AccAddress
	mnemonic string
}

func (m *keyManager) ExportAsMnemonic() (string, error) {
	if m.mnemonic == "" {
		return "", fmt.Errorf("This key manager is not recover from mnemonic or anto generated ")
	}
	return m.mnemonic, nil
}

func (m *keyManager) ExportAsPrivateKey() (string, error) {
	secpPrivateKey, ok := m.privKey.(secp256k1.PrivKeySecp256k1)
	if !ok {
		return "", fmt.Errorf(" Only PrivKeySecp256k1 key is supported ")
	}
	return hex.EncodeToString(secpPrivateKey[:]), nil
}

func (m *keyManager) ExportAsKeyStore(password string) (*EncryptedKeyJSON, error) {
	return generateKeyStore(m.GetPrivKey(), password)
}

func NewKeyManager() (KeyManager, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return nil, err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}
	return NewMnemonicKeyManager(mnemonic)
}

func (m *keyManager) recoveryFromMnemonic(mnemonic, keyPath string) error {
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
	derivedPriv, err := DerivePrivateKeyForPath(masterPriv, ch, keyPath)
	if err != nil {
		return err
	}
	priKey := secp256k1.PrivKeySecp256k1(derivedPriv)
	addr := ctypes.AccAddress(priKey.PubKey().Address())
	if err != nil {
		return err
	}
	m.addr = addr
	m.privKey = priKey
	m.mnemonic = mnemonic
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
	var encryptedKey EncryptedKeyJSON
	err = json.Unmarshal(keyJson, &encryptedKey)
	if err != nil {
		return err
	}
	keyBytes, err := decryptKey(&encryptedKey, auth)
	if err != nil {
		return err
	}
	if len(keyBytes) != 32 {
		return fmt.Errorf("Len of Keybytes is not equal to 32 ")
	}
	var keyBytesArray [32]byte
	copy(keyBytesArray[:], keyBytes[:32])
	priKey := secp256k1.PrivKeySecp256k1(keyBytesArray)
	addr := ctypes.AccAddress(priKey.PubKey().Address())
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
	addr := ctypes.AccAddress(priKey.PubKey().Address())
	m.addr = addr
	m.privKey = priKey
	return nil
}

func (m *keyManager) recoveryFromLedgerKey(path ledger.DerivationPath) error {
	if ledger.DiscoverLedger == nil {
		return fmt.Errorf("no Ledger discovery function defined, please make sure you have added ledger to build tags and cgo is enabled")
	}

	device, err := ledger.DiscoverLedger()
	if err != nil {
		return fmt.Errorf("failed to find ledger device: %s", err.Error())
	}

	pkl, err := ledger.GenLedgerSecp256k1Key(path, device)
	if err != nil {
		return fmt.Errorf("failed to create PrivKeyLedgerSecp256k1: %s", err.Error())
	}

	addr := types.AccAddress(pkl.PubKey().Address())
	m.addr = addr
	m.privKey = pkl
	return nil
}

func (m *keyManager) Sign(msg tx.StdSignMsg) ([]byte, error) {
	sig, err := m.makeSignature(msg)
	if err != nil {
		return nil, err
	}
	newTx := tx.NewStdTx(msg.Msgs, []tx.StdSignature{sig}, msg.Memo, msg.Source, msg.Data)
	bz, err := tx.Cdc.MarshalBinaryLengthPrefixed(&newTx)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func (m *keyManager) GetPrivKey() crypto.PrivKey {
	return m.privKey
}

func (m *keyManager) GetAddr() ctypes.AccAddress {
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

func generateKeyStore(privateKey crypto.PrivKey, password string) (*EncryptedKeyJSON, error) {
	addr := ctypes.AccAddress(privateKey.PubKey().Address())
	salt, err := common.GenerateRandomBytes(32)
	if err != nil {
		return nil, err
	}
	iv, err := common.GenerateRandomBytes(16)
	if err != nil {
		return nil, err
	}
	scryptParamsJSON := make(map[string]interface{}, 4)
	scryptParamsJSON["prf"] = "hmac-sha256"
	scryptParamsJSON["dklen"] = 32
	scryptParamsJSON["salt"] = hex.EncodeToString(salt)
	scryptParamsJSON["c"] = 262144

	cipherParamsJSON := cipherparamsJSON{IV: hex.EncodeToString(iv)}
	derivedKey := pbkdf2.Key([]byte(password), salt, 262144, 32, sha256.New)
	encryptKey := derivedKey[:32]
	secpPrivateKey, ok := privateKey.(secp256k1.PrivKeySecp256k1)
	if !ok {
		return nil, fmt.Errorf(" Only PrivKeySecp256k1 key is supported ")
	}
	cipherText, err := aesCTRXOR(encryptKey, secpPrivateKey[:], iv)
	if err != nil {
		return nil, err
	}

	hasher := sha3.NewLegacyKeccak512()
	hasher.Write(derivedKey[16:32])
	hasher.Write(cipherText)
	mac := hasher.Sum(nil)

	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	cryptoStruct := CryptoJSON{
		Cipher:       "aes-256-ctr",
		CipherText:   hex.EncodeToString(cipherText),
		CipherParams: cipherParamsJSON,
		KDF:          "pbkdf2",
		KDFParams:    scryptParamsJSON,
		MAC:          hex.EncodeToString(mac),
	}
	return &EncryptedKeyJSON{
		Address: addr.String(),
		Crypto:  cryptoStruct,
		Id:      id.String(),
		Version: 1,
	}, nil
}
