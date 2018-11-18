package sdk

import (
	"encoding/json"
	"fmt"

	tmcrypto "github.com/tendermint/tendermint/crypto"
)

// Tx def
type Tx struct {
	Hash string `json:"hash"`
	Log  string `json:"log"`
	Data string `json:"data"`
	Code int32  `json:"code"`
}

// GetTx returns transaction details
func (sdk *SDK) GetTx(txHash string) (*Tx, error) {
	if txHash == "" {
		return nil, fmt.Errorf("Invalid tx hash %s", txHash)
	}

	qp := map[string]string{}
	resp, err := sdk.dexAPI.Get("/tx/"+txHash, qp)
	if err != nil {
		return nil, err
	}

	var tx Tx
	if err := json.Unmarshal(resp, &tx); err != nil {
		return nil, err
	}

	return &tx, nil
}

// SignTx returns signature
func (sdk *SDK) SignTx(priv tmcrypto.PrivKey, msg []byte) (sig []byte, pub tmcrypto.PubKey, err error) {
	sig, err = priv.Sign(msg)
	if err != nil {
		return nil, nil, err
	}
	pub = priv.PubKey()
	return sig, pub, nil
}
