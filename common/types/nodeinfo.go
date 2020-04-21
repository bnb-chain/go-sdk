package types

import (
	"time"

	libbytes "github.com/tendermint/tendermint/libs/bytes"
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
	Network  string            `json:"network"`  // network/chain ID
	Version  string            `json:"version"`  // major.minor.revision
	Channels libbytes.HexBytes `json:"channels"` // channels this node knows about

	// ASCIIText fields
	Moniker string        `json:"moniker"` // arbitrary moniker
	Other   NodeInfoOther `json:"other"`   // other application specific data
}

type ValidatorInfo struct {
	Address     libbytes.HexBytes `json:"address"`
	PubKey      []uint8           `json:"pub_key"`
	VotingPower int64             `json:"voting_power"`
}
type SyncInfo struct {
	LatestBlockHash   libbytes.HexBytes `json:"latest_block_hash"`
	LatestAppHash     libbytes.HexBytes `json:"latest_app_hash"`
	LatestBlockHeight int64             `json:"latest_block_height"`
	LatestBlockTime   time.Time         `json:"latest_block_time"`
	CatchingUp        bool              `json:"catching_up"`
}

type NodeInfoOther struct {
	AminoVersion     string `json:"amino_version"`
	P2PVersion       string `json:"p2p_version"`
	ConsensusVersion string `json:"consensus_version"`
	RPCVersion       string `json:"rpc_version"`
	TxIndex          string `json:"tx_index"`
	RPCAddress       string `json:"rpc_address"`
}
