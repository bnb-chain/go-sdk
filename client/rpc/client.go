package rpc

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"sync"
	"time"

	ntypes "github.com/binance-chain/go-sdk/common/types"
	"github.com/tendermint/go-amino"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/rpc/lib/client"
	"github.com/tendermint/tendermint/rpc/lib/types"
	"github.com/tendermint/tendermint/types"
)

const defaultTimeout = 5 * time.Second

func NewRPCClient(nodeURI string) *HTTP {
	return NewHTTP(nodeURI, "/websocket")
}

type HTTP struct {
	remote string
	*WSEvents
}

// NewHTTP takes a remote endpoint in the form tcp://<host>:<port>
// and the websocket path (which always seems to be "/websocket")
func NewHTTP(remote, wsEndpoint string) *HTTP {
	rc := rpcclient.NewJSONRPCClient(remote)
	cdc := rc.Codec()
	ctypes.RegisterAmino(cdc)
	ntypes.RegisterWire(cdc)


	rc.SetCodec(cdc)
	wsEvent := newWSEvents(cdc, remote, wsEndpoint)
	client := &HTTP{
		remote:   remote,
		WSEvents: wsEvent,
	}
	client.Start()
	return client
}

func (c *HTTP) Status() (*ctypes.ResultStatus, error) {
	return c.WSEvents.Status()
}

func (c *HTTP) ABCIInfo() (*ctypes.ResultABCIInfo, error) {
	return c.WSEvents.ABCIInfo()
}

func (c *HTTP) ABCIQuery(path string, data cmn.HexBytes) (*ctypes.ResultABCIQuery, error) {
	return c.ABCIQueryWithOptions(path, data, client.DefaultABCIQueryOptions)
}

func (c *HTTP) ABCIQueryWithOptions(path string, data cmn.HexBytes, opts client.ABCIQueryOptions) (*ctypes.ResultABCIQuery, error) {
	return c.WSEvents.ABCIQueryWithOptions(path, data, opts)
}

func (c *HTTP) BroadcastTxCommit(tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error) {
	return c.WSEvents.BroadcastTxCommit(tx)
}

func (c *HTTP) BroadcastTxAsync(tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	return c.WSEvents.BroadcastTx("broadcast_tx_async", tx)
}

func (c *HTTP) BroadcastTxSync(tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	return c.WSEvents.BroadcastTx("broadcast_tx_sync", tx)
}

func (c *HTTP) UnconfirmedTxs(limit int) (*ctypes.ResultUnconfirmedTxs, error) {
	return c.WSEvents.UnconfirmedTxs(limit)
}

func (c *HTTP) NumUnconfirmedTxs() (*ctypes.ResultUnconfirmedTxs, error) {
	return c.WSEvents.NumUnconfirmedTxs()
}

func (c *HTTP) NetInfo() (*ctypes.ResultNetInfo, error) {
	return c.WSEvents.NetInfo()
}

func (c *HTTP) DumpConsensusState() (*ctypes.ResultDumpConsensusState, error) {
	return c.WSEvents.DumpConsensusState()
}

func (c *HTTP) ConsensusState() (*ctypes.ResultConsensusState, error) {
	return c.WSEvents.ConsensusState()
}

func (c *HTTP) Health() (*ctypes.ResultHealth, error) {
	return c.WSEvents.Health()
}

func (c *HTTP) BlockchainInfo(minHeight, maxHeight int64) (*ctypes.ResultBlockchainInfo, error) {
	return c.WSEvents.BlockchainInfo(minHeight, maxHeight)
}

func (c *HTTP) Genesis() (*ctypes.ResultGenesis, error) {
	return c.WSEvents.Genesis()
}

func (c *HTTP) Block(height *int64) (*ctypes.ResultBlock, error) {
	return c.WSEvents.Block(height)
}

func (c *HTTP) BlockResults(height *int64) (*ctypes.ResultBlockResults, error) {
	return c.WSEvents.BlockResults(height)
}

func (c *HTTP) Commit(height *int64) (*ctypes.ResultCommit, error) {
	return c.WSEvents.Commit(height)
}

func (c *HTTP) Tx(hash []byte, prove bool) (*ctypes.ResultTx, error) {
	return c.WSEvents.Tx(hash, prove)
}

func (c *HTTP) TxSearch(query string, prove bool, page, perPage int) (*ctypes.ResultTxSearch, error) {
	return c.WSEvents.TxSearch(query, prove, page, perPage)
}

func (c *HTTP) Validators(height *int64) (*ctypes.ResultValidators, error) {
	return c.WSEvents.Validators(height)
}

func (c *HTTP) QueryStore(key cmn.HexBytes, storeName string) ([]byte, error) {
	path := fmt.Sprintf("/store/%s/%s", storeName, "key")
	result, err := c.ABCIQuery(path, key)
	if err != nil {
		return nil, err
	}
	resp := result.Response
	if !resp.IsOK() {
		return nil, errors.Errorf(resp.Log)
	}
	return resp.Value, nil
}

/** websocket event stuff here... **/

type WSEvents struct {
	cmn.BaseService
	cdc      *amino.Codec
	remote   string
	endpoint string
	ws       *WSClient

	mtx sync.RWMutex
	// query -> chan

	subscriptionsQuitMap map[string]chan struct{}
	subscriptionsIdMap   map[string]rpctypes.JSONRPCStringID
	subscriptionSet      map[rpctypes.JSONRPCStringID]bool

	responseChanMap sync.Map

	timeout time.Duration
}

func newWSEvents(cdc *amino.Codec, remote, endpoint string) *WSEvents {
	wsEvents := &WSEvents{
		cdc:                  cdc,
		endpoint:             endpoint,
		remote:               remote,
		subscriptionsQuitMap: make(map[string]chan struct{}),
		subscriptionsIdMap:   make(map[string]rpctypes.JSONRPCStringID),
		subscriptionSet:      make(map[rpctypes.JSONRPCStringID]bool),
		timeout:              defaultTimeout,
	}

	wsEvents.BaseService = *cmn.NewBaseService(nil, "WSEvents", wsEvents)
	return wsEvents
}

// OnStart implements cmn.Service by starting WSClient and event loop.
func (w *WSEvents) OnStart() error {
	w.ws = NewWSClient(w.remote, w.endpoint, OnReconnect(func() {
		w.redoSubscriptionsAfter(0 * time.Second)
	}))
	w.ws.SetCodec(w.cdc)

	err := w.ws.Start()
	if err != nil {
		return err
	}

	go w.eventListener()
	return nil
}

// OnStop implements cmn.Service by stopping WSClient.
func (w *WSEvents) OnStop() {
	_ = w.ws.Stop()
}

// Subscribe implements EventsClient by using WSClient to subscribe given
// subscriber to query. By default, returns a channel with cap=1. Error is
// returned if it fails to subscribe.
// Channel is never closed to prevent clients from seeing an erroneus event.
func (w *WSEvents) Subscribe(query string,
	outCapacity ...int) (out chan ctypes.ResultEvent, err error) {
	if _, ok := w.subscriptionsIdMap[query]; ok {
		return nil, fmt.Errorf("already subscribe")
	}

	id, err := w.ws.GenRequestId()
	if err != nil {
		return nil, err
	}
	outCap := 1
	if len(outCapacity) > 0 {
		outCap = outCapacity[0]
	}
	outEvent := make(chan ctypes.ResultEvent, outCap)
	outResp := make(chan rpctypes.RPCResponse, cap(outEvent))
	w.responseChanMap.Store(id, outResp)
	ctx, cancel := w.NewContext()
	defer cancel()
	err = w.ws.Subscribe(ctx, id, query)
	if err != nil {
		w.responseChanMap.Delete(id)
		return nil, err
	}

	quit := make(chan struct{})
	w.mtx.Lock()
	w.subscriptionsQuitMap[query] = quit
	w.subscriptionsIdMap[query] = id
	w.subscriptionSet[id] = true
	w.mtx.Unlock()
	go w.WaitForEventResponse(id, outResp, outEvent, quit)

	return outEvent, nil
}

// Unsubscribe implements EventsClient by using WSClient to unsubscribe given
// subscriber from query.
func (w *WSEvents) Unsubscribe(query string) error {
	ctx, cancel := w.NewContext()
	defer cancel()
	if err := w.ws.Unsubscribe(ctx, w.ws.EmptyRequest(), query); err != nil {
		return err
	}

	w.mtx.Lock()
	if id, ok := w.subscriptionsIdMap[query]; ok {
		delete(w.subscriptionSet, id)
		w.responseChanMap.Delete(id)
	}
	if quit, ok := w.subscriptionsQuitMap[query]; ok {
		close(quit)
	}
	delete(w.subscriptionsIdMap, query)
	delete(w.subscriptionsQuitMap, query)
	w.mtx.Unlock()

	return nil
}

// UnsubscribeAll implements EventsClient by using WSClient to unsubscribe
// given subscriber from all the queries.
func (w *WSEvents) UnsubscribeAll() error {
	ctx, cancel := w.NewContext()
	defer cancel()
	if err := w.ws.UnsubscribeAll(ctx, w.ws.EmptyRequest()); err != nil {
		return err
	}

	w.mtx.Lock()
	for _, id := range w.subscriptionsIdMap {
		w.responseChanMap.Delete(id)
	}
	for _, quit := range w.subscriptionsQuitMap {
		close(quit)
	}
	w.subscriptionSet = make(map[rpctypes.JSONRPCStringID]bool)
	w.subscriptionsQuitMap = make(map[string]chan struct{})
	w.subscriptionsIdMap = make(map[string]rpctypes.JSONRPCStringID)
	w.mtx.Unlock()

	return nil
}

func (w *WSEvents) WaitForEventResponse(requestId interface{}, in chan rpctypes.RPCResponse, eventOut chan ctypes.ResultEvent, quit chan struct{}) {

	for {
		select {
		case <-quit:
			return
		case resp, ok := <-in:
			if !ok {
				w.Logger.Info("channel of event stream is closed", "request id", requestId)
				return
			}
			if resp.Error != nil {
				w.Logger.Error("receive error from event stream", "error", resp.Error)
				continue
			}
			res := new(ctypes.ResultEvent)
			err := w.cdc.UnmarshalJSON(resp.Result, res)
			if err != nil {
				w.Logger.Debug("receive unexpected data from event stream", "result", resp.Result)
				continue
			}
			eventOut <- *res
		}
	}
}

func (w *WSEvents) WaitForResponse(ctx context.Context, outChan chan rpctypes.RPCResponse, result interface{}) error {
	select {
	case resp, ok := <-outChan:
		if !ok {
			return fmt.Errorf("response channel is closed")
		}
		if resp.Error != nil {
			return resp.Error
		}
		return w.cdc.UnmarshalJSON(resp.Result, result)
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (w *WSEvents) SimpleCall(doRpc func(ctx context.Context, id rpctypes.JSONRPCStringID) error, proto interface{}) error {
	id, err := w.ws.GenRequestId()
	if err != nil {
		return err
	}
	outChan := make(chan rpctypes.RPCResponse, 1)
	w.responseChanMap.Store(id, outChan)
	defer close(outChan)
	defer w.responseChanMap.Delete(id)
	ctx, cancel := w.NewContext()
	defer cancel()
	if err = doRpc(ctx, id); err != nil {
		return err
	}
	return w.WaitForResponse(ctx, outChan, proto)
}

func (w *WSEvents) Status() (*ctypes.ResultStatus, error) {
	status := new(ctypes.ResultStatus)
	err := w.SimpleCall(w.ws.Status, status)
	return status, err
}

func (w *WSEvents) ABCIInfo() (*ctypes.ResultABCIInfo, error) {
	info := new(ctypes.ResultABCIInfo)
	err := w.SimpleCall(w.ws.ABCIInfo, info)
	return info, err
}

func (w *WSEvents) ABCIQueryWithOptions(path string, data cmn.HexBytes, opts client.ABCIQueryOptions) (*ctypes.ResultABCIQuery, error) {
	abciQuery := new(ctypes.ResultABCIQuery)
	err := w.SimpleCall(func(ctx context.Context, id rpctypes.JSONRPCStringID) error {
		return w.ws.ABCIQueryWithOptions(ctx, id, path, data, opts)
	}, abciQuery)
	return abciQuery, err
}

func (w *WSEvents) BroadcastTxCommit(tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error) {
	txCommit := new(ctypes.ResultBroadcastTxCommit)
	err := w.SimpleCall(func(ctx context.Context, id rpctypes.JSONRPCStringID) error {
		return w.ws.BroadcastTxCommit(ctx, id, tx)
	}, txCommit)
	return txCommit, err
}

func (w *WSEvents) BroadcastTx(route string, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	txRes := new(ctypes.ResultBroadcastTx)
	err := w.SimpleCall(func(ctx context.Context, id rpctypes.JSONRPCStringID) error {
		return w.ws.BroadcastTx(ctx, id, route, tx)
	}, txRes)
	return txRes, err
}

func (w *WSEvents) UnconfirmedTxs(limit int) (*ctypes.ResultUnconfirmedTxs, error) {
	unConfirmTxs := new(ctypes.ResultUnconfirmedTxs)
	err := w.SimpleCall(func(ctx context.Context, id rpctypes.JSONRPCStringID) error {
		return w.ws.UnconfirmedTxs(ctx, id, limit)
	}, unConfirmTxs)
	return unConfirmTxs, err
}

func (w *WSEvents) NumUnconfirmedTxs() (*ctypes.ResultUnconfirmedTxs, error) {
	numUnConfirmTxs := new(ctypes.ResultUnconfirmedTxs)
	err := w.SimpleCall(w.ws.NumUnconfirmedTxs, numUnConfirmTxs)
	return numUnConfirmTxs, err
}

func (w *WSEvents) NetInfo() (*ctypes.ResultNetInfo, error) {
	netInfo := new(ctypes.ResultNetInfo)
	err := w.SimpleCall(w.ws.NetInfo, netInfo)
	return netInfo, err
}

func (w *WSEvents) DumpConsensusState() (*ctypes.ResultDumpConsensusState, error) {
	consensusState := new(ctypes.ResultDumpConsensusState)
	err := w.SimpleCall(w.ws.DumpConsensusState, consensusState)
	return consensusState, err
}

func (w *WSEvents) ConsensusState() (*ctypes.ResultConsensusState, error) {
	consensusState := new(ctypes.ResultConsensusState)
	err := w.SimpleCall(w.ws.ConsensusState, consensusState)
	return consensusState, err
}

func (w *WSEvents) Health() (*ctypes.ResultHealth, error) {
	health := new(ctypes.ResultHealth)
	err := w.SimpleCall(w.ws.Health, health)
	return health, err
}

func (w *WSEvents) BlockchainInfo(minHeight, maxHeight int64) (*ctypes.ResultBlockchainInfo, error) {

	blocksInfo := new(ctypes.ResultBlockchainInfo)
	err := w.SimpleCall(func(ctx context.Context, id rpctypes.JSONRPCStringID) error {
		return w.ws.BlockchainInfo(ctx, id, minHeight, maxHeight)
	}, blocksInfo)
	return blocksInfo, err
}

func (w *WSEvents) Genesis() (*ctypes.ResultGenesis, error) {

	genesis := new(ctypes.ResultGenesis)
	err := w.SimpleCall(w.ws.Genesis, genesis)
	return genesis, err
}

func (w *WSEvents) Block(height *int64) (*ctypes.ResultBlock, error) {
	block := new(ctypes.ResultBlock)
	err := w.SimpleCall(func(ctx context.Context, id rpctypes.JSONRPCStringID) error {
		return w.ws.Block(ctx, id, height)
	}, block)
	return block, err
}

func (w *WSEvents) BlockResults(height *int64) (*ctypes.ResultBlockResults, error) {

	block := new(ctypes.ResultBlockResults)
	err := w.SimpleCall(func(ctx context.Context, id rpctypes.JSONRPCStringID) error {
		return w.ws.BlockResults(ctx, id, height)
	}, block)
	return block, err
}

func (w *WSEvents) Commit(height *int64) (*ctypes.ResultCommit, error) {
	commit := new(ctypes.ResultCommit)
	err := w.SimpleCall(func(ctx context.Context, id rpctypes.JSONRPCStringID) error {
		return w.ws.Commit(ctx, id, height)
	}, commit)
	return commit, err
}

func (w *WSEvents) Tx(hash []byte, prove bool) (*ctypes.ResultTx, error) {

	tx := new(ctypes.ResultTx)
	err := w.SimpleCall(func(ctx context.Context, id rpctypes.JSONRPCStringID) error {
		return w.ws.Tx(ctx, id, hash, prove)
	}, tx)
	return tx, err
}

func (w *WSEvents) TxSearch(query string, prove bool, page, perPage int) (*ctypes.ResultTxSearch, error) {

	txs := new(ctypes.ResultTxSearch)
	err := w.SimpleCall(func(ctx context.Context, id rpctypes.JSONRPCStringID) error {
		return w.ws.TxSearch(ctx, id, query, prove, page, perPage)
	}, txs)
	return txs, err
}

func (w *WSEvents) Validators(height *int64) (*ctypes.ResultValidators, error) {
	validators := new(ctypes.ResultValidators)
	err := w.SimpleCall(func(ctx context.Context, id rpctypes.JSONRPCStringID) error {
		return w.ws.Validators(ctx, id, height)
	}, validators)
	return validators, err
}

func (w *WSEvents) SetTimeOut(timeout time.Duration) {
	w.timeout = timeout
}

func (w *WSEvents) NewContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), w.timeout)
}

// After being reconnected, it is necessary to redo subscription to server
// otherwise no data will be automatically received.
func (w *WSEvents) redoSubscriptionsAfter(d time.Duration) {
	time.Sleep(d)

	for q, id := range w.subscriptionsIdMap {
		ctx, _ := context.WithTimeout(context.Background(), w.timeout)
		err := w.ws.Subscribe(ctx, id, q)
		if err != nil {
			w.Logger.Error("Failed to resubscribe", "err", err)
		}
	}
}

func (w *WSEvents) eventListener() {
	for {
		select {
		case resp, ok := <-w.ws.ResponsesCh:
			if !ok {
				return
			}
			id, ok := resp.ID.(rpctypes.JSONRPCStringID)
			if !ok {
				w.Logger.Error("unexpected request id type")
				continue
			}
			if exist := w.subscriptionSet[id]; exist {
				// receive ack event, need ignore it
				continue
			}
			idParts := strings.Split(string(id), "#")
			realId := rpctypes.JSONRPCStringID(idParts[0])
			if out, ok := w.responseChanMap.Load(realId); ok {
				outChan, ok := out.(chan rpctypes.RPCResponse)
				if !ok {
					w.Logger.Error("unexpected data type in responseChanMap")
					continue
				}
				select {
				case outChan <- resp:
				default:
					w.Logger.Error("wanted to publish response, but out channel is full", "result", resp.Result)
				}
			}
		case <-w.Quit():
			return
		}
	}
}
