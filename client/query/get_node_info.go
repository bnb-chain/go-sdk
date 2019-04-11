package query

import (
	"encoding/json"
	"time"

	cmn "github.com/tendermint/tendermint/libs/common"
)

// Account definition
type ResultStatus struct {
	NodeInfo      NodeInfo      `json:"node_info"`
	SyncInfo      SyncInfo      `json:"sync_info"`
	ValidatorInfo ValidatorInfo `json:"validator_info"`
}

type NodeInfo struct {
	// Authenticate
	// TODO: replace with NetAddress
	ID         string `json:"id"`          // authenticated identifier
	ListenAddr string `json:"listen_addr"` // accepting incoming

	// Check compatibility.
	// Channels are HexBytes so easier to read as JSON
	Network  string       `json:"network"`  // network/chain ID
	Version  string       `json:"version"`  // major.minor.revision
	Channels cmn.HexBytes `json:"channels"` // channels this node knows about

	// ASCIIText fields
	Moniker string        `json:"moniker"` // arbitrary moniker
	Other   NodeInfoOther `json:"other"`   // other application specific data
}

type ValidatorInfo struct {
	Address     cmn.HexBytes `json:"address"`
	PubKey      []uint8      `json:"pub_key"`
	VotingPower int64        `json:"voting_power"`
}
type SyncInfo struct {
	LatestBlockHash   cmn.HexBytes `json:"latest_block_hash"`
	LatestAppHash     cmn.HexBytes `json:"latest_app_hash"`
	LatestBlockHeight int64        `json:"latest_block_height"`
	LatestBlockTime   time.Time    `json:"latest_block_time"`
	CatchingUp        bool         `json:"catching_up"`
}

type NodeInfoOther struct {
	AminoVersion     string `json:"amino_version"`
	P2PVersion       string `json:"p2p_version"`
	ConsensusVersion string `json:"consensus_version"`
	RPCVersion       string `json:"rpc_version"`
	TxIndex          string `json:"tx_index"`
	RPCAddress       string `json:"rpc_address"`
}

func (c *client) GetNodeInfo() (*ResultStatus, error) {
	qp := map[string]string{}
	resp, err := c.baseClient.Get("/node-info", qp)
	if err != nil {
		return nil, err
	}
	var resultStatus ResultStatus
	if err := json.Unmarshal(resp, &resultStatus); err != nil {
		return nil, err
	}

	return &resultStatus, nil
}
