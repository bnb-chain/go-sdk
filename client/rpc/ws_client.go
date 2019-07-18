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
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/rpc/lib/types"
	"github.com/tendermint/tendermint/types"

	"github.com/binance-chain/go-sdk/common/uuid"
	"github.com/binance-chain/go-sdk/types/tx"
)

const (
	defaultMaxReconnectAttempts    = -1
	defaultMaxReconnectBackOffTime = 10 * time.Second
	defaultWriteWait               = 100 * time.Millisecond
	defaultReadWait                = 0
	defaultPingPeriod              = 0
	defaultDialPeriod              = 1 * time.Second

	protoHTTP  = "http"
	protoHTTPS = "https"
	protoWSS   = "wss"
	protoWS    = "ws"
	protoTCP   = "tcp"
)

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

// OnStop implements cmn.Service by stopping WSClient.
func (w *WSEvents) PendingRequest() int {
	size := 0
	w.responseChanMap.Range(func(key, value interface{}) bool {
		size++
		return true
	})
	return size
}

func (c *WSEvents) IsActive() bool {
	return c.ws.IsActive()
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
			return errors.New("response channel is closed")
		}
		if resp.Error != nil {
			return resp.Error
		}
		return w.cdc.UnmarshalJSON(resp.Result, result)
	case <-ctx.Done():
		w.ws.reconnectAfter <- ctx.Err()
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

func (w *WSEvents) TxInfoSearch(query string, prove bool, page, perPage int) ([]tx.Info, error) {

	txs := new(ctypes.ResultTxSearch)
	err := w.SimpleCall(func(ctx context.Context, id rpctypes.JSONRPCStringID) error {
		return w.ws.TxSearch(ctx, id, query, prove, page, perPage)
	}, txs)
	if err != nil {
		return nil, err
	}
	return FormatTxResults(w.cdc, txs.Txs)
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

// WSClient is a WebSocket client. The methods of WSClient are safe for use by
// multiple goroutines.
type WSClient struct {
	cmn.BaseService

	conn *websocket.Conn
	cdc  *amino.Codec

	Address  string // IP:PORT or /path/to/socket
	Endpoint string // /websocket/url/endpoint
	Dialer   func(string, string) (net.Conn, error)

	// Single user facing channel to read RPCResponses from, closed only when the client is being stopped.
	ResponsesCh chan rpctypes.RPCResponse

	// Callback, which will be called each time after successful reconnect.
	onReconnect func()

	// internal channels
	send           chan rpctypes.RPCRequest // user requests
	reconnectAfter chan error               // reconnect requests

	wg sync.WaitGroup

	mtx     sync.RWMutex
	dialing atomic.Value

	// Maximum reconnect attempts (0 or greater; default: 25).
	// Less than 0 means always try to reconnect.
	maxReconnectAttempts int

	// Time allowed to write a message to the server. 0 means block until operation succeeds.
	writeWait time.Duration

	// Time allowed to read the next message from the server. 0 means block until operation succeeds.
	readWait time.Duration

	// Send pings to server with this period. Must be less than readWait. If 0, no pings will be sent.
	pingPeriod time.Duration

	// Support both ws and wss protocols
	protocol string
}

// NewWSClient returns a new client. See the commentary on the func(*WSClient)
// functions for a detailed description of how to configure ping period and
// pong wait time. The endpoint argument must begin with a `/`.
func NewWSClient(remoteAddr, endpoint string, options ...func(*WSClient)) *WSClient {
	protocol, addr, dialer := makeHTTPDialer(remoteAddr)
	// default to ws protocol, unless wss is explicitly specified
	if protocol != "wss" {
		protocol = "ws"
	}

	c := &WSClient{
		cdc:      amino.NewCodec(),
		Address:  addr,
		Dialer:   dialer,
		Endpoint: endpoint,

		maxReconnectAttempts: defaultMaxReconnectAttempts,
		readWait:             defaultReadWait,
		writeWait:            defaultWriteWait,
		pingPeriod:           defaultPingPeriod,
		protocol:             protocol,
	}
	c.dialing.Store(true)
	c.BaseService = *cmn.NewBaseService(nil, "WSClient", c)
	for _, option := range options {
		option(c)
	}
	return c
}

// OnReconnect sets the callback, which will be called every time after
// successful reconnect.
func OnReconnect(cb func()) func(*WSClient) {
	return func(c *WSClient) {
		c.onReconnect = cb
	}
}

// String returns WS client full address.
func (c *WSClient) String() string {
	return fmt.Sprintf("%s (%s)", c.Address, c.Endpoint)
}

// OnStart implements cmn.Service by dialing a server and creating read and
// write routines.
func (c *WSClient) OnStart() error {

	c.ResponsesCh = make(chan rpctypes.RPCResponse)

	c.send = make(chan rpctypes.RPCRequest)
	// 1 additional error may come from the read/write
	// goroutine depending on which failed first.
	c.reconnectAfter = make(chan error, 1)
	// capacity for 1 request. a user won't be able to send more because the send
	// channel is unbuffered.

	err := c.dial()
	if err != nil {
		c.wg.Add(1)
		go c.dialRoutine()
	} else {
		c.dialing.Store(false)
	}

	c.startReadWriteRoutines()
	go c.reconnectRoutine()

	return nil
}

// Stop overrides cmn.Service#Stop. There is no other way to wait until Quit
// channel is closed.
func (c *WSClient) Stop() error {
	if err := c.BaseService.Stop(); err != nil {
		return err
	}
	// only close user-facing channels when we can't write to them
	c.wg.Wait()
	close(c.ResponsesCh)

	return nil
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
// ResponsesCh, errors, if any, on ErrorsCh. Will block until send succeeds or
// ctx.Done is closed.
func (c *WSClient) Send(ctx context.Context, request rpctypes.RPCRequest) error {
	select {
	case c.send <- request:
		c.Logger.Info("sent a request", "req", request)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Call the given method. See Send description.
func (c *WSClient) Call(ctx context.Context, method string, id rpctypes.JSONRPCStringID, params map[string]interface{}) error {
	if c.IsDialing() {
		return errors.New("websocket is dialing, can't send any request")
	}
	request, err := rpctypes.MapToRequest(c.cdc, id, method, params)
	if err != nil {
		return err
	}
	return c.Send(ctx, request)
}

// CallWithArrayParams the given method with params in a form of array. See
// Send description.
func (c *WSClient) CallWithArrayParams(ctx context.Context, method string, params []interface{}) error {
	request, err := rpctypes.ArrayToRequest(c.cdc, rpctypes.JSONRPCStringID("ws-client"), method, params)
	if err != nil {
		return err
	}
	return c.Send(ctx, request)
}

func (c *WSClient) GetConnection() *websocket.Conn {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.conn
}

func (c *WSClient) SetConnection(conn *websocket.Conn) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if c.conn != nil {
		c.conn.Close()
	}
	c.conn = conn
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

func (c *WSClient) EmptyRequest() rpctypes.JSONRPCStringID {
	return rpctypes.JSONRPCStringID("")
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
	c.SetConnection(conn)
	return nil
}

// reconnect tries to redial up to maxReconnectAttempts with exponential
// backoff.
func (c *WSClient) reconnect() error {
	attempt := 0

	c.dialing.Store(true)
	defer func() {
		c.dialing.Store(false)
	}()
	backOffDuration := 1 * time.Second
	for {
		// will never overflow until doomsday
		backOffDuration := time.Duration(attempt)*time.Second + backOffDuration
		if backOffDuration > defaultMaxReconnectBackOffTime {
			backOffDuration = defaultMaxReconnectBackOffTime
		}
		c.Logger.Info("reconnecting", "attempt", attempt+1, "backoff_duration", backOffDuration)
		time.Sleep(backOffDuration)

		err := c.dial()
		if err != nil {
			c.Logger.Error("failed to redial", "err", err)
		} else {
			c.Logger.Info("reconnected")
			if c.onReconnect != nil {
				go c.onReconnect()
			}
			return nil
		}

		attempt++

		if c.maxReconnectAttempts >= 0 && attempt > c.maxReconnectAttempts {
			return errors.Wrap(err, "reached maximum reconnect attempts")
		}
	}
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

func (c *WSClient) reconnectRoutine() {
	for {
		select {
		case originalError := <-c.reconnectAfter:
			if err := c.reconnect(); err != nil {
				c.Logger.Error("failed to reconnect", "err", err, "original_err", originalError)
				c.Stop()
				return
			}
			// drain reconnectAfter
		LOOP:
			for {
				select {
				case <-c.reconnectAfter:
				default:
					break LOOP
				}
			}

		case <-c.Quit():
			return
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

	defer func() {
		ticker.Stop()
		if err := c.GetConnection().Close(); err != nil {
			// ignore error; it will trigger in tests
			// likely because it's closing an already closed connection
		}
	}()

	for {
		select {
		case request := <-c.send:
			if c.writeWait > 0 {
				if err := c.GetConnection().SetWriteDeadline(time.Now().Add(c.writeWait)); err != nil {
					c.Logger.Error("failed to set write deadline", "err", err)
				}
			}
			if err := c.GetConnection().WriteJSON(request); err != nil {
				c.Logger.Error("failed to send request", "err", err)
				c.reconnectAfter <- err
			}
		case <-ticker.C:
			if c.writeWait > 0 {
				if err := c.GetConnection().SetWriteDeadline(time.Now().Add(c.writeWait)); err != nil {
					c.Logger.Error("failed to set write deadline", "err", err)
				}
			}
			if err := c.GetConnection().WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				c.Logger.Error("failed to write ping", "err", err)
				c.reconnectAfter <- err
				continue
			}
			c.Logger.Debug("sent ping")
		case <-c.Quit():
			if err := c.GetConnection().WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
				c.Logger.Error("failed to write message", "err", err)
			}
			return
		}
	}
}

// The client ensures that there is at most one reader to a connection by
// executing all reads from this goroutine.
func (c *WSClient) readRoutine() {
	defer func() {
		if err := c.GetConnection().Close(); err != nil {
			// ignore error; it will trigger in tests
			// likely because it's closing an already closed connection
		}
	}()
	c.wg.Wait()
	c.GetConnection().SetPongHandler(func(string) error {
		c.Logger.Debug("got pong")
		return nil
	})

	for {
		// reset deadline for every message type (control or data)
		if c.readWait > 0 {
			if err := c.GetConnection().SetReadDeadline(time.Now().Add(c.readWait)); err != nil {
				c.Logger.Error("failed to set read deadline", "err", err)
			}
		}
		_, data, err := c.GetConnection().ReadMessage()
		if err != nil {
			c.Logger.Error("failed to read response", "err", err)
			c.reconnectAfter <- err
			continue
		}

		var response rpctypes.RPCResponse
		err = json.Unmarshal(data, &response)
		if err != nil {
			c.Logger.Error("failed to parse response", "err", err, "data", string(data))
			continue
		}
		c.Logger.Info("got response", "resp", response.Result)
		// Combine a non-blocking read on BaseService.Quit with a non-blocking write on ResponsesCh to avoid blocking
		// c.wg.Wait() in c.Stop(). Note we rely on Quit being closed so that it sends unlimited Quit signals to stop
		// both readRoutine and writeRoutine
		select {
		case <-c.Quit():
			return
		case c.ResponsesCh <- response:
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

func (c *WSClient) ABCIQueryWithOptions(ctx context.Context, id rpctypes.JSONRPCStringID, path string, data cmn.HexBytes, opts client.ABCIQueryOptions) error {
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
		return net.DialTimeout(protocol, address, defaultTimeout)
	}
}

// parse the indexed txs into an array of Info
func FormatTxResults(cdc *amino.Codec, res []*ctypes.ResultTx) ([]tx.Info, error) {
	var err error
	out := make([]tx.Info, len(res))
	for i := range res {
		out[i], err = formatTxResult(cdc, res[i])
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func formatTxResult(cdc *amino.Codec, res *ctypes.ResultTx) (tx.Info, error) {
	parsedTx, err := ParseTx(cdc, res.Tx)
	if err != nil {
		return tx.Info{}, err
	}

	return tx.Info{
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
