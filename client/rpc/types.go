package rpc

import (
	"github.com/binance-chain/go-sdk/types/tx"
	"github.com/tendermint/tendermint/abci/types"
	libbytes "github.com/tendermint/tendermint/libs/bytes"
	libkv "github.com/tendermint/tendermint/libs/kv"
	abci "github.com/tendermint/tendermint/types"
)

type ResultBroadcastTxCommit struct {
	CheckTx   ResponseCheckTx   `json:"check_tx"`
	DeliverTx ResponseDeliverTx `json:"deliver_tx"`
	Hash      libbytes.HexBytes `json:"hash"`
	Height    int64             `json:"height"`
}

func (r ResponseCheckTx) IsErr() bool {
	return r.Code != types.CodeTypeOK
}

func (r *ResultBroadcastTxCommit) complement() {
	r.CheckTx.complement()
	r.DeliverTx.complement()
}

type ResponseCheckTx struct {
	Code      uint32         `json:"code,omitempty"`
	Data      []byte         `json:"data,omitempty"`
	Log       string         `json:"log,omitempty"`
	Info      string         `json:"info,omitempty"`
	GasWanted int64          `json:"gas_wanted,omitempty"`
	GasUsed   int64          `json:"gas_used,omitempty"`
	Events    []types.Event  `json:"events,omitempty"`
	Tags      []libkv.KVPair `json:"tags,omitempty"`
	Codespace string         `json:"codespace,omitempty"`
}

func (r *ResponseCheckTx) complement() {
	if len(r.Tags) > 0 {
		r.Events = []types.Event{{Attributes: r.Tags}}
	} else if len(r.Events) > 0 {
		r.Tags = r.Events[0].Attributes
	}
}

type ResponseDeliverTx struct {
	Code      uint32         `json:"code,omitempty"`
	Data      []byte         `json:"data,omitempty"`
	Log       string         `json:"log,omitempty"`
	Info      string         `json:"info,omitempty"`
	GasWanted int64          `json:"gas_wanted,omitempty"`
	GasUsed   int64          `json:"gas_used,omitempty"`
	Events    []types.Event  `json:"events,omitempty"`
	Tags      []libkv.KVPair `json:"tags,omitempty"`
	Codespace string         `json:"codespace,omitempty"`
}

func (r *ResponseDeliverTx) complement() {
	if len(r.Tags) > 0 {
		r.Events = []types.Event{{Attributes: r.Tags}}
	} else if len(r.Events) > 0 {
		r.Tags = r.Events[0].Attributes
	}
}

type ResultBlockResults struct {
	Height  int64          `json:"height"`
	Results *ABCIResponses `json:"results"`
}

func (r *ResultBlockResults) complement() {
	if r.Results != nil {
		r.Results.complement()
	}
}

type ABCIResponses struct {
	DeliverTx  []*ResponseDeliverTx `json:"DeliverTx"`
	EndBlock   *ResponseEndBlock    `json:"EndBlock"`
	BeginBlock *ResponseBeginBlock  `json:"BeginBlock"`
}

func (r *ABCIResponses) complement() {
	for _, d := range r.DeliverTx {
		d.complement()
	}
	if r.EndBlock != nil {
		r.EndBlock.complement()
	}
	if r.BeginBlock != nil {
		r.BeginBlock.complement()
	}
}

type ResponseEndBlock struct {
	ValidatorUpdates      []types.ValidatorUpdate `json:"validator_updates"`
	ConsensusParamUpdates *types.ConsensusParams  `json:"consensus_param_updates,omitempty"`
	Events                []types.Event           `json:"events,omitempty"`
	Tags                  []libkv.KVPair          `json:"tags,omitempty"`
}

func (r *ResponseEndBlock) complement() {
	if len(r.Tags) > 0 {
		r.Events = []types.Event{{Attributes: r.Tags}}
	} else if len(r.Events) > 0 {
		r.Tags = r.Events[0].Attributes
	}
}

type ResponseBeginBlock struct {
	Events []types.Event  `json:"events,omitempty"`
	Tags   []libkv.KVPair `json:"tags,omitempty"`
}

func (r *ResponseBeginBlock) complement() {
	if len(r.Tags) > 0 {
		r.Events = []types.Event{{Attributes: r.Tags}}
	} else if len(r.Events) > 0 {
		r.Tags = r.Events[0].Attributes
	}
}

type ResultTx struct {
	Hash     libbytes.HexBytes `json:"hash"`
	Height   int64             `json:"height"`
	Index    uint32            `json:"index"`
	TxResult ResponseDeliverTx `json:"tx_result"`
	Tx       abci.Tx           `json:"tx"`
	Proof    abci.TxProof      `json:"proof,omitempty"`
}

func (r *ResultTx) complement() {
	r.TxResult.complement()
}

// Result of searching for txs
type ResultTxSearch struct {
	Txs        []*ResultTx `json:"txs"`
	TotalCount int         `json:"total_count"`
}

func (r *ResultTxSearch) complement() {
	for _, t := range r.Txs {
		t.TxResult.complement()
	}
}

type Info struct {
	Hash   libbytes.HexBytes `json:"hash"`
	Height int64             `json:"height"`
	Tx     tx.Tx             `json:"tx"`
	Result ResponseDeliverTx `json:"result"`
}

func (r *Info) complement() {
	r.Result.complement()
}
