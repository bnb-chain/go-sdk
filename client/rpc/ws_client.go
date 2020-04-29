package rpc

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"sync/atomic"

	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/tendermint/go-amino"
	libbytes "github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/libs/log"
	libservice "github.com/tendermint/tendermint/libs/service"
	"github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpctypes "github.com/tendermint/tendermint/rpc/lib/types"
	"github.com/tendermint/tendermint/types"

	"github.com/binance-chain/go-sdk/common/uuid"
	"github.com/binance-chain/go-sdk/types/tx"
)

const (
	defaultWsAliveCheckPeriod = 1 * time.Second
	defaultWriteWait          = 500 * time.Millisecond
	defaultReadWait           = 0
	defaultPingPeriod         = 0
	defaultDialPeriod         = 1 * time.Second

	protoHTTP  = "http"
	protoHTTPS = "https"
	protoWSS   = "wss"
	protoWS    = "ws"
	protoTCP   = "tcp"

	EmptyRequest = rpctypes.JSONRPCStringID("")
)

/** websocket event stuff here... **/
type WSEvents struct {
	libservice.BaseService
	cdc      *amino.Codec
	remote   string
	endpoint string
	ws       *WSClient

	mtx       sync.RWMutex
	reconnect chan *WSClient

	subscriptionsQuitMap map[string]chan struct{}
	subscriptionsIdMap   map[string]rpctypes.JSONRPCStringID
	subscriptionSet      map[rpctypes.JSONRPCStringID]bool

	responsesCh chan rpctypes.RPCResponse

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
		timeout:              DefaultTimeout,
		responsesCh:          make(chan rpctypes.RPCResponse),
		reconnect:            make(chan *WSClient),
	}

	wsEvents.BaseService = *libservice.NewBaseService(nil, "WSEvents", wsEvents)
	return wsEvents
}

// OnStart implements libservice.Service by starting WSClient and event loop.
func (w *WSEvents) OnStart() error {
	wsClient := NewWSClient(w.remote, w.endpoint, w.responsesCh, setOnDialSuccess(w.redoSubscriptionsAfter))
	wsClient.SetCodec(w.cdc)
	err := wsClient.Start()
	if err != nil {
		return err
	}
	w.setWsClient(wsClient)

	go w.eventListener()
	go w.reconnectRoutine()
	return nil
}

func (w *WSEvents) SetLogger(logger log.Logger) {
	w.BaseService.SetLogger(logger)
	wsClient := w.getWsClient()
	if wsClient != nil {
		w.getWsClient().SetLogger(logger)
	}
}

// OnStop implements libservice.Service by stopping WSClient.
func (w *WSEvents) OnStop() {
	w.getWsClient().Stop()
}

// OnStop implements libservice.Service by stopping WSClient.
func (w *WSEvents) PendingRequest() int {
	size := 0
	w.responseChanMap.Range(func(key, value interface{}) bool {
		size++
		return true
	})
	return size
}

func (c *WSEvents) IsActive() bool {
	return c.getWsClient().IsActive()
}

// Subscribe implements EventsClient by using WSClient to subscribe given
// subscriber to query. By default, returns a channel with cap=1. Error is
// returned if it fails to subscribe.
// Channel is never closed to prevent clients from seeing an erroneus event.
func (w *WSEvents) Subscribe(query string,
	outCapacity ...int) (out chan ctypes.ResultEvent, err error) {
	if _, ok := w.subscriptionsIdMap[query]; ok {
		return nil, errors.New("already subscribe")
	}

	id, err := w.getWsClient().GenRequestId()
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
	err = w.getWsClient().Subscribe(ctx, id, query)
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
	if err := w.getWsClient().Unsubscribe(ctx, EmptyRequest, query); err != nil {
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
	if err := w.getWsClient().UnsubscribeAll(ctx, EmptyRequest); err != nil {
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

func (w *WSEvents) setWsClient(wsClient *WSClient) {
	w.mtx.Lock()
	defer w.mtx.Unlock()
	w.ws = wsClient
}

func (w *WSEvents) getWsClient() *WSClient {
	w.mtx.Lock()
	defer w.mtx.Unlock()
	return w.ws
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

func (w *WSEvents) WaitForResponse(ctx context.Context, outChan chan rpctypes.RPCResponse, result interface{}, ws *WSClient) error {
	select {
	case resp, ok := <-outChan:
		if !ok {
			return errors.New("response channel is closed")
		}
		if resp.Error != nil {
			return resp.Error
		}
		return w.cdc.UnmarshalJSON(resp.Result, result)
	case <-ctx.Done():
		w.reconnect <- ws
		return ctx.Err()
	}
}

func (w *WSEvents) SimpleCall(doRpc func(ctx context.Context, id rpctypes.JSONRPCStringID) error, ws *WSClient, proto interface{}) error {
	id, err := ws.GenRequestId()
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
	return w.WaitForResponse(ctx, outChan, proto, ws)
}

func (w *WSEvents) Status() (*ctypes.ResultStatus, error) {
	status := new(ctypes.ResultStatus)
	wsClient := w.getWsClient()
	err := w.SimpleCall(wsClient.Status, wsClient, status)
	return status, err
}

func (w *WSEvents) ABCIInfo() (*ctypes.ResultABCIInfo, error) {
	info := new(ctypes.ResultABCIInfo)
	wsClient := w.getWsClient()
	err := w.SimpleCall(wsClient.ABCIInfo, wsClient, info)
	return info, err
}

func (w *WSEvents) ABCIQueryWithOptions(path string, data libbytes.HexBytes, opts client.ABCIQueryOptions) (*ctypes.ResultABCIQuery, error) {
	abciQuery := new(ctypes.ResultABCIQuery)
	wsClient := w.getWsClient()
	err := w.SimpleCall(func(ctx context.Context, id rpctypes.JSONRPCStringID) error {
		return wsClient.ABCIQueryWithOptions(ctx, id, path, data, opts)
	}, wsClient, abciQuery)
	return abciQuery, err
}

func (w *WSEvents) BroadcastTxCommit(tx types.Tx) (*ResultBroadcastTxCommit, error) {
	txCommit := new(ResultBroadcastTxCommit)
	wsClient := w.getWsClient()
	err := w.SimpleCall(func(ctx context.Context, id rpctypes.JSONRPCStringID) error {
		return wsClient.BroadcastTxCommit(ctx, id, tx)
	}, wsClient, txCommit)
	if err == nil {
		txCommit.complement()
	}
	return txCommit, err
}

func (w *WSEvents) BroadcastTx(route string, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	txRes := new(ctypes.ResultBroadcastTx)
	wsClient := w.getWsClient()
	err := w.SimpleCall(func(ctx context.Context, id rpctypes.JSONRPCStringID) error {
		return wsClient.BroadcastTx(ctx, id, route, tx)
	}, wsClient, txRes)
	return txRes, err
}

func (w *WSEvents) UnconfirmedTxs(limit int) (*ctypes.ResultUnconfirmedTxs, error) {
	unConfirmTxs := new(ctypes.ResultUnconfirmedTxs)
	wsClient := w.getWsClient()
	err := w.SimpleCall(func(ctx context.Context, id rpctypes.JSONRPCStringID) error {
		return wsClient.UnconfirmedTxs(ctx, id, limit)
	}, wsClient, unConfirmTxs)
	return unConfirmTxs, err
}

func (w *WSEvents) NumUnconfirmedTxs() (*ctypes.ResultUnconfirmedTxs, error) {
	numUnConfirmTxs := new(ctypes.ResultUnconfirmedTxs)
	wsClient := w.getWsClient()
	err := w.SimpleCall(wsClient.NumUnconfirmedTxs, wsClient, numUnConfirmTxs)
	return numUnConfirmTxs, err
}

func (w *WSEvents) NetInfo() (*ctypes.ResultNetInfo, error) {
	netInfo := new(ctypes.ResultNetInfo)
	wsClient := w.getWsClient()
	err := w.SimpleCall(wsClient.NetInfo, wsClient, netInfo)
	return netInfo, err
}

func (w *WSEvents) DumpConsensusState() (*ctypes.ResultDumpConsensusState, error) {
	consensusState := new(ctypes.ResultDumpConsensusState)
	wsClient := w.getWsClient()
	err := w.SimpleCall(wsClient.DumpConsensusState, wsClient, consensusState)
	return consensusState, err
}

func (w *WSEvents) ConsensusState() (*ctypes.ResultConsensusState, error) {
	consensusState := new(ctypes.ResultConsensusState)
	wsClient := w.getWsClient()
	err := w.SimpleCall(wsClient.ConsensusState, wsClient, consensusState)
	return consensusState, err
}

func (w *WSEvents) Health() (*ctypes.ResultHealth, error) {
	health := new(ctypes.ResultHealth)
	wsClient := w.getWsClient()
	err := w.SimpleCall(wsClient.Health, wsClient, health)
	return health, err
}

func (w *WSEvents) BlockchainInfo(minHeight, maxHeight int64) (*ctypes.ResultBlockchainInfo, error) {

	blocksInfo := new(ctypes.ResultBlockchainInfo)
	wsClient := w.getWsClient()
	err := w.SimpleCall(func(ctx context.Context, id rpctypes.JSONRPCStringID) error {
		return wsClient.BlockchainInfo(ctx, id, minHeight, maxHeight)
	}, wsClient, blocksInfo)
	return blocksInfo, err
}

func (w *WSEvents) Genesis() (*ctypes.ResultGenesis, error) {

	genesis := new(ctypes.ResultGenesis)
	wsClient := w.getWsClient()
	err := w.SimpleCall(wsClient.Genesis, wsClient, genesis)
	return genesis, err
}

func (w *WSEvents) Block(height *int64) (*ctypes.ResultBlock, error) {
	block := new(ctypes.ResultBlock)
	wsClient := w.getWsClient()
	err := w.SimpleCall(func(ctx context.Context, id rpctypes.JSONRPCStringID) error {
		return wsClient.Block(ctx, id, height)
	}, wsClient, block)
	return block, err
}

func (w *WSEvents) BlockResults(height *int64) (*ResultBlockResults, error) {

	block := new(ResultBlockResults)
	wsClient := w.getWsClient()
	err := w.SimpleCall(func(ctx context.Context, id rpctypes.JSONRPCStringID) error {
		return wsClient.BlockResults(ctx, id, height)
	}, wsClient, block)
	if err == nil {
		block.complement()
	}
	return block, err
}

func (w *WSEvents) Commit(height *int64) (*ctypes.ResultCommit, error) {
	commit := new(ctypes.ResultCommit)
	wsClient := w.getWsClient()
	err := w.SimpleCall(func(ctx context.Context, id rpctypes.JSONRPCStringID) error {
		return wsClient.Commit(ctx, id, height)
	}, wsClient, commit)
	return commit, err
}

func (w *WSEvents) Tx(hash []byte, prove bool) (*ResultTx, error) {

	tx := new(ResultTx)
	wsClient := w.getWsClient()
	err := w.SimpleCall(func(ctx context.Context, id rpctypes.JSONRPCStringID) error {
		return wsClient.Tx(ctx, id, hash, prove)
	}, wsClient, tx)
	if err == nil {
		tx.complement()
	}
	return tx, err
}

func (w *WSEvents) TxSearch(query string, prove bool, page, perPage int) (*ResultTxSearch, error) {

	txs := new(ResultTxSearch)
	wsClient := w.getWsClient()
	err := w.SimpleCall(func(ctx context.Context, id rpctypes.JSONRPCStringID) error {
		return wsClient.TxSearch(ctx, id, query, prove, page, perPage)
	}, wsClient, txs)
	if err == nil {
		txs.complement()
	}
	return txs, err
}

func (w *WSEvents) TxInfoSearch(query string, prove bool, page, perPage int) ([]Info, error) {

	txs := new(ResultTxSearch)
	wsClient := w.getWsClient()
	err := w.SimpleCall(func(ctx context.Context, id rpctypes.JSONRPCStringID) error {
		return wsClient.TxSearch(ctx, id, query, prove, page, perPage)
	}, wsClient, txs)
	if err != nil {
		return nil, err
	}
	return FormatTxResults(w.cdc, txs.Txs)
}

func (w *WSEvents) Validators(height *int64) (*ctypes.ResultValidators, error) {
	validators := new(ctypes.ResultValidators)
	wsClient := w.getWsClient()
	err := w.SimpleCall(func(ctx context.Context, id rpctypes.JSONRPCStringID) error {
		return wsClient.Validators(ctx, id, height)
	}, wsClient, validators)
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
func (w *WSEvents) redoSubscriptionsAfter() {

	for q, id := range w.subscriptionsIdMap {
		ctx, _ := context.WithTimeout(context.Background(), w.timeout)
		err := w.getWsClient().Subscribe(ctx, id, q)
		if err != nil {
			w.Logger.Error("Failed to resubscribe", "err", err)
		}
	}
}

func (w *WSEvents) reconnectRoutine() {
	checkTicker := time.NewTicker(defaultWsAliveCheckPeriod)
	defer checkTicker.Stop()
	for {
		select {
		case <-w.Quit():
			return
		case wsClient := <-w.reconnect:
			if wsClient.IsRunning() {
				w.Logger.Error("try to stop wsClient since context deadline exceed", "server", wsClient)
				// try to stop, may stop twice, but make no harm.
				wsClient.Stop()
			}
		case <-checkTicker.C:
			if !w.getWsClient().IsRunning() {
				w.Logger.Info("ws client have been stopped, try start new one", "server", w.getWsClient())
				wsClient := NewWSClient(w.remote, w.endpoint, w.responsesCh, setOnDialSuccess(w.redoSubscriptionsAfter))
				wsClient.SetCodec(w.cdc)
				err := wsClient.Start()
				// should not happen
				if err != nil {
					w.Logger.Error("wsClient start failed", "err", err)
					continue
				}
				w.Logger.Info("ws client reconnect success", "server", w.getWsClient())
				w.setWsClient(wsClient)
			}
		}
	}
}

func (w *WSEvents) eventListener() {
	for {
		select {
		case resp, ok := <-w.responsesCh:
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

// WSClient is a WebSocket client. The methods of WSClient are safe for use by
// multiple goroutines.
type WSClient struct {
	libservice.BaseService

	conn *websocket.Conn
	cdc  *amino.Codec

	Address  string // IP:PORT or /path/to/socket
	Endpoint string // /websocket/url/endpoint
	Dialer   func(string, string) (net.Conn, error)

	// Single user facing channel to read RPCResponses from
	responsesCh chan<- rpctypes.RPCResponse

	// internal channels
	send chan rpctypes.RPCRequest // user requests

	wg sync.WaitGroup

	mtx     sync.RWMutex
	dialing atomic.Value

	// Time allowed to write a message to the server. 0 means block until operation succeeds.
	writeWait time.Duration

	// Time allowed to read the next message from the server. 0 means block until operation succeeds.
	readWait time.Duration

	// Send pings to server with this period. Must be less than readWait. If 0, no pings will be sent.
	pingPeriod time.Duration

	// Support both ws and wss protocols
	protocol string

	onDialSuccess func()
}

// NewWSClient returns a new client. See the commentary on the func(*WSClient)
// functions for a detailed description of how to configure ping period and
// pong wait time. The endpoint argument must begin with a `/`.
func NewWSClient(remoteAddr, endpoint string, responsesCh chan<- rpctypes.RPCResponse, options ...func(*WSClient)) *WSClient {
	protocol, addr, dialer := makeHTTPDialer(remoteAddr)
	// default to ws protocol, unless wss is explicitly specified
	if protocol != "wss" {
		protocol = "ws"
	}

	c := &WSClient{
		cdc:         amino.NewCodec(),
		Address:     addr,
		Dialer:      dialer,
		Endpoint:    endpoint,
		readWait:    defaultReadWait,
		writeWait:   defaultWriteWait,
		pingPeriod:  defaultPingPeriod,
		protocol:    protocol,
		responsesCh: responsesCh,
		send:        make(chan rpctypes.RPCRequest),
	}
	c.dialing.Store(true)
	c.BaseService = *libservice.NewBaseService(nil, "WSClient", c)
	for _, option := range options {
		option(c)
	}
	return c
}

// String returns WS client full address.
func (c *WSClient) String() string {
	if c.conn != nil {
		return fmt.Sprintf("%s, local: %s, remote: %s", c.Address, c.conn.LocalAddr().String(), c.conn.RemoteAddr().String())
	}
	return fmt.Sprintf("%s (%s)", c.Address, c.Endpoint)
}

// OnStart implements libservice.Service by dialing a server and creating read and
// write routines.
func (c *WSClient) OnStart() error {
	err := c.dial()
	if err != nil {
		c.wg.Add(1)
		go c.dialRoutine()
	} else {
		c.dialing.Store(false)
	}

	c.startReadWriteRoutines()
	return nil
}

func (c *WSClient) OnStop() {
	if c.conn != nil {
		c.conn.Close()
	}
	return
}

// IsDialing returns true if the client is dialing right now.
func (c *WSClient) IsDialing() bool {
	return c.dialing.Load().(bool)
}

// IsActive returns true if the client is running and not dialing.
func (c *WSClient) IsActive() bool {
	return c.IsRunning() && !c.IsDialing()
}

// Send the given RPC request to the server. Results will be available on
// responsesCh, errors, if any, on ErrorsCh. Will block until send succeeds or
// ctx.Done is closed.
func (c *WSClient) Send(ctx context.Context, request rpctypes.RPCRequest) error {
	select {
	case c.send <- request:
		c.Logger.Debug("sent a request", "req", request)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Call the given method. See Send description.
func (c *WSClient) Call(ctx context.Context, method string, id rpctypes.JSONRPCStringID, params map[string]interface{}) error {
	if !c.IsActive() {
		return errors.New("websocket client is dialing or stopped, can't send any request")
	}
	request, err := rpctypes.MapToRequest(c.cdc, id, method, params)
	if err != nil {
		return err
	}
	return c.Send(ctx, request)
}

func (c *WSClient) Codec() *amino.Codec {
	return c.cdc
}

func (c *WSClient) SetCodec(cdc *amino.Codec) {
	c.cdc = cdc
}

func (c *WSClient) GenRequestId() (rpctypes.JSONRPCStringID, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return rpctypes.JSONRPCStringID(id.String()), nil
}

///////////////////////////////////////////////////////////////////////////////
// Private methods

func (c *WSClient) dial() error {
	dialer := &websocket.Dialer{
		NetDial: c.Dialer,
		Proxy:   http.ProxyFromEnvironment,
	}
	rHeader := http.Header{}
	conn, _, err := dialer.Dial(c.protocol+"://"+c.Address+c.Endpoint, rHeader)
	if err != nil {
		return err
	}
	// only do once during the lifecycle of WSClient
	c.conn = conn
	if c.onDialSuccess != nil {
		go c.onDialSuccess()
	}
	return nil
}

func (c *WSClient) startReadWriteRoutines() {
	go c.readRoutine()
	go c.writeRoutine()
}

func (c *WSClient) dialRoutine() {
	dialTicker := time.NewTicker(defaultDialPeriod)
	defer dialTicker.Stop()
	c.dialing.Store(true)
	defer func() {
		c.dialing.Store(false)
		c.wg.Done()
	}()
	for {
		select {
		case <-c.Quit():
			return
		case <-dialTicker.C:
			err := c.dial()
			if err == nil {
				return
			}
		}
	}
}

// The client ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *WSClient) writeRoutine() {
	c.wg.Wait()
	var ticker *time.Ticker
	if c.pingPeriod > 0 {
		// ticker with a predefined period
		ticker = time.NewTicker(c.pingPeriod)
	} else {
		// ticker that never fires
		ticker = &time.Ticker{C: make(<-chan time.Time)}
	}
	defer ticker.Stop()

	for {
		select {
		case request := <-c.send:
			if c.writeWait > 0 {
				if err := c.conn.SetWriteDeadline(time.Now().Add(c.writeWait)); err != nil {
					c.Logger.Error("failed to set write deadline", "err", err)
				}
			}
			if err := c.conn.WriteJSON(request); err != nil {
				c.Logger.Error("failed to send request", "err", err)
				c.Stop()
				return
			}
		case <-ticker.C:
			if c.writeWait > 0 {
				if err := c.conn.SetWriteDeadline(time.Now().Add(c.writeWait)); err != nil {
					c.Logger.Error("failed to set write deadline", "err", err)
				}
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				c.Logger.Error("failed to write ping", "err", err)
				c.Stop()
				return
			}
			c.Logger.Debug("sent ping")
		case <-c.Quit():
			if err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
				c.Logger.Error("failed to write message", "err", err)
			}
			return
		}
	}
}

// The client ensures that there is at most one reader to a connection by
// executing all reads from this goroutine.
func (c *WSClient) readRoutine() {
	c.wg.Wait()
	c.conn.SetPongHandler(func(string) error {
		c.Logger.Debug("got pong")
		return nil
	})

	for {
		// reset deadline for every message type (control or data)
		if c.readWait > 0 {
			if err := c.conn.SetReadDeadline(time.Now().Add(c.readWait)); err != nil {
				c.Logger.Error("failed to set read deadline", "err", err)
			}
		}
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			c.Logger.Error("failed to read response", "err", err)
			c.Stop()
			return
		}

		var response rpctypes.RPCResponse
		err = json.Unmarshal(data, &response)
		if err != nil {
			c.Logger.Error("failed to parse response", "err", err, "data", string(data))
			continue
		}
		// Combine a non-blocking read on BaseService.Quit with a non-blocking write on responsesCh to avoid blocking
		// c.wg.Wait() in c.Stop(). Note we rely on Quit being closed so that it sends unlimited Quit signals to stop
		// both readRoutine and writeRoutine
		select {
		case <-c.Quit():
			return
		case c.responsesCh <- response:
		}
	}
}

///////////////////////////////////////////////////////////////////////////////
// Predefined methods

// Subscribe to a query. Note the server must have a "subscribe" route
// defined.
func (c *WSClient) Subscribe(ctx context.Context, id rpctypes.JSONRPCStringID, query string) error {
	params := map[string]interface{}{"query": query}
	return c.Call(ctx, "subscribe", id, params)
}

// Unsubscribe from a query. Note the server must have a "unsubscribe" route
// defined.
func (c *WSClient) Unsubscribe(ctx context.Context, id rpctypes.JSONRPCStringID, query string) error {
	params := map[string]interface{}{"query": query}

	return c.Call(ctx, "unsubscribe", id, params)
}

// UnsubscribeAll from all. Note the server must have a "unsubscribe_all" route
// defined.
func (c *WSClient) UnsubscribeAll(ctx context.Context, id rpctypes.JSONRPCStringID) error {
	params := map[string]interface{}{}
	return c.Call(ctx, "unsubscribe_all", id, params)
}

func (c *WSClient) Status(ctx context.Context, id rpctypes.JSONRPCStringID) error {
	return c.Call(ctx, "status", id, map[string]interface{}{})
}

func (c *WSClient) ABCIInfo(ctx context.Context, id rpctypes.JSONRPCStringID) error {
	return c.Call(ctx, "abci_info", id, map[string]interface{}{})
}

func (c *WSClient) ABCIQueryWithOptions(ctx context.Context, id rpctypes.JSONRPCStringID, path string, data libbytes.HexBytes, opts client.ABCIQueryOptions) error {
	return c.Call(ctx, "abci_query", id, map[string]interface{}{"path": path, "data": data, "height": opts.Height, "prove": opts.Prove})
}

func (c *WSClient) BroadcastTxCommit(ctx context.Context, id rpctypes.JSONRPCStringID, tx types.Tx) error {
	return c.Call(ctx, "broadcast_tx_commit", id, map[string]interface{}{"tx": tx})
}

func (c *WSClient) BroadcastTx(ctx context.Context, id rpctypes.JSONRPCStringID, route string, tx types.Tx) error {
	return c.Call(ctx, route, id, map[string]interface{}{"tx": tx})
}

func (c *WSClient) UnconfirmedTxs(ctx context.Context, id rpctypes.JSONRPCStringID, limit int) error {
	return c.Call(ctx, "unconfirmed_txs", id, map[string]interface{}{"limit": limit})
}

func (c *WSClient) NumUnconfirmedTxs(ctx context.Context, id rpctypes.JSONRPCStringID) error {
	return c.Call(ctx, "num_unconfirmed_txs", id, map[string]interface{}{})
}

func (c *WSClient) NetInfo(ctx context.Context, id rpctypes.JSONRPCStringID) error {
	return c.Call(ctx, "net_info", id, map[string]interface{}{})
}

func (c *WSClient) DumpConsensusState(ctx context.Context, id rpctypes.JSONRPCStringID) error {
	return c.Call(ctx, "dump_consensus_state", id, map[string]interface{}{})
}

func (c *WSClient) ConsensusState(ctx context.Context, id rpctypes.JSONRPCStringID) error {
	return c.Call(ctx, "consensus_state", id, map[string]interface{}{})
}

func (c *WSClient) Health(ctx context.Context, id rpctypes.JSONRPCStringID) error {
	return c.Call(ctx, "health", id, map[string]interface{}{})
}

func (c *WSClient) BlockchainInfo(ctx context.Context, id rpctypes.JSONRPCStringID, minHeight, maxHeight int64) error {
	return c.Call(ctx, "blockchain", id,
		map[string]interface{}{"minHeight": minHeight, "maxHeight": maxHeight})
}

func (c *WSClient) Genesis(ctx context.Context, id rpctypes.JSONRPCStringID) error {
	return c.Call(ctx, "genesis", id, map[string]interface{}{})
}

func (c *WSClient) Block(ctx context.Context, id rpctypes.JSONRPCStringID, height *int64) error {
	return c.Call(ctx, "block", id, map[string]interface{}{"height": height})
}

func (c *WSClient) BlockResults(ctx context.Context, id rpctypes.JSONRPCStringID, height *int64) error {
	return c.Call(ctx, "block_results", id, map[string]interface{}{"height": height})
}

func (c *WSClient) Commit(ctx context.Context, id rpctypes.JSONRPCStringID, height *int64) error {
	return c.Call(ctx, "commit", id, map[string]interface{}{"height": height})
}
func (c *WSClient) Tx(ctx context.Context, id rpctypes.JSONRPCStringID, hash []byte, prove bool) error {
	params := map[string]interface{}{
		"hash":  hash,
		"prove": prove,
	}
	return c.Call(ctx, "tx", id, params)
}

func (c *WSClient) TxSearch(ctx context.Context, id rpctypes.JSONRPCStringID, query string, prove bool, page, perPage int) error {
	params := map[string]interface{}{
		"query":    query,
		"prove":    prove,
		"page":     page,
		"per_page": perPage,
	}
	return c.Call(ctx, "tx_search", id, params)
}

func (c *WSClient) Validators(ctx context.Context, id rpctypes.JSONRPCStringID, height *int64) error {
	return c.Call(ctx, "validators", id, map[string]interface{}{"height": height})
}

func setOnDialSuccess(onDial func()) func(c *WSClient) {
	return func(c *WSClient) {
		c.onDialSuccess = onDial
	}
}

func makeHTTPDialer(remoteAddr string) (string, string, func(string, string) (net.Conn, error)) {
	// protocol to use for http operations, to support both http and https
	clientProtocol := protoHTTP

	parts := strings.SplitN(remoteAddr, "://", 2)
	var protocol, address string
	if len(parts) == 1 {
		// default to tcp if nothing specified
		protocol, address = protoTCP, remoteAddr
	} else if len(parts) == 2 {
		protocol, address = parts[0], parts[1]
	} else {
		// return a invalid message
		msg := fmt.Sprintf("Invalid addr: %s", remoteAddr)
		return clientProtocol, msg, func(_ string, _ string) (net.Conn, error) {
			return nil, errors.New(msg)
		}
	}

	// accept http as an alias for tcp and set the client protocol
	switch protocol {
	case protoHTTP, protoHTTPS:
		clientProtocol = protocol
		protocol = protoTCP
	case protoWS, protoWSS:
		clientProtocol = protocol
	}

	// replace / with . for http requests (kvstore domain)
	trimmedAddress := strings.Replace(address, "/", ".", -1)
	return clientProtocol, trimmedAddress, func(proto, addr string) (net.Conn, error) {
		return net.DialTimeout(protocol, address, DefaultTimeout)
	}
}

// parse the indexed txs into an array of Info
func FormatTxResults(cdc *amino.Codec, res []*ResultTx) ([]Info, error) {
	var err error
	out := make([]Info, len(res))
	for i := range res {
		out[i], err = formatTxResult(cdc, res[i])
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func formatTxResult(cdc *amino.Codec, res *ResultTx) (Info, error) {
	parsedTx, err := ParseTx(cdc, res.Tx)
	if err != nil {
		return Info{}, err
	}
	return Info{
		Hash:   res.Hash,
		Height: res.Height,
		Tx:     parsedTx,
		Result: res.TxResult,
	}, nil
}

func ParseTx(cdc *amino.Codec, txBytes []byte) (tx.Tx, error) {
	var parsedTx tx.StdTx
	err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &parsedTx)

	if err != nil {
		return nil, err
	}

	return parsedTx, nil
}
