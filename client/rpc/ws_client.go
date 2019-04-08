package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/rcrowley/go-metrics"

	"github.com/tendermint/go-amino"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/rpc/client"
	types "github.com/tendermint/tendermint/rpc/lib/types"
	ctypes "github.com/tendermint/tendermint/types"

	"github.com/binance-chain/go-sdk/common/uuid"
)

const (
	defaultMaxReconnectAttempts    = 25
	defaultMaxReconnectBackOffTime = 600 * time.Second
	defaultWriteWait               = 100 * time.Millisecond
	defaultReadWait                = 0
	defaultPingPeriod              = 0

	protoHTTP  = "http"
	protoHTTPS = "https"
	protoWSS   = "wss"
	protoWS    = "ws"
	protoTCP   = "tcp"
)

// WSClient is a WebSocket client. The methods of WSClient are safe for use by
// multiple goroutines.
type WSClient struct {
	cmn.BaseService

	conn *websocket.Conn
	cdc  *amino.Codec

	Address  string // IP:PORT or /path/to/socket
	Endpoint string // /websocket/url/endpoint
	Dialer   func(string, string) (net.Conn, error)

	// Time between sending a ping and receiving a pong. See
	// https://godoc.org/github.com/rcrowley/go-metrics#Timer.
	PingPongLatencyTimer metrics.Timer

	// Single user facing channel to read RPCResponses from, closed only when the client is being stopped.
	ResponsesCh chan types.RPCResponse

	// Callback, which will be called each time after successful reconnect.
	onReconnect func()

	// internal channels
	send            chan types.RPCRequest // user requests
	reconnectAfter  chan error            // reconnect requests
	readRoutineQuit chan struct{}         // a way for readRoutine to close writeRoutine

	wg sync.WaitGroup

	mtx            sync.RWMutex
	sentLastPingAt time.Time
	reconnecting   bool

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
		cdc:                  amino.NewCodec(),
		Address:              addr,
		Dialer:               dialer,
		Endpoint:             endpoint,
		PingPongLatencyTimer: metrics.NewTimer(),

		maxReconnectAttempts: defaultMaxReconnectAttempts,
		readWait:             defaultReadWait,
		writeWait:            defaultWriteWait,
		pingPeriod:           defaultPingPeriod,
		protocol:             protocol,
	}
	c.BaseService = *cmn.NewBaseService(nil, "WSClient", c)
	for _, option := range options {
		option(c)
	}
	return c
}

// MaxReconnectAttempts sets the maximum number of reconnect attempts before returning an error.
// It should only be used in the constructor and is not Goroutine-safe.
func MaxReconnectAttempts(max int) func(*WSClient) {
	return func(c *WSClient) {
		c.maxReconnectAttempts = max
	}
}

// ReadWait sets the amount of time to wait before a websocket read times out.
// It should only be used in the constructor and is not Goroutine-safe.
func ReadWait(readWait time.Duration) func(*WSClient) {
	return func(c *WSClient) {
		c.readWait = readWait
	}
}

// WriteWait sets the amount of time to wait before a websocket write times out.
// It should only be used in the constructor and is not Goroutine-safe.
func WriteWait(writeWait time.Duration) func(*WSClient) {
	return func(c *WSClient) {
		c.writeWait = writeWait
	}
}

// PingPeriod sets the duration for sending websocket pings.
// It should only be used in the constructor - not Goroutine-safe.
func PingPeriod(pingPeriod time.Duration) func(*WSClient) {
	return func(c *WSClient) {
		c.pingPeriod = pingPeriod
	}
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
	err := c.dial()
	if err != nil {
		return err
	}

	c.ResponsesCh = make(chan types.RPCResponse)

	c.send = make(chan types.RPCRequest)
	// 1 additional error may come from the read/write
	// goroutine depending on which failed first.
	c.reconnectAfter = make(chan error, 1)
	// capacity for 1 request. a user won't be able to send more because the send
	// channel is unbuffered.

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

// IsReconnecting returns true if the client is reconnecting right now.
func (c *WSClient) IsReconnecting() bool {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.reconnecting
}

// IsActive returns true if the client is running and not reconnecting.
func (c *WSClient) IsActive() bool {
	return c.IsRunning() && !c.IsReconnecting()
}

// Send the given RPC request to the server. Results will be available on
// ResponsesCh, errors, if any, on ErrorsCh. Will block until send succeeds or
// ctx.Done is closed.
func (c *WSClient) Send(ctx context.Context, request types.RPCRequest) error {
	select {
	case c.send <- request:
		c.Logger.Info("sent a request", "req", request)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Call the given method. See Send description.
func (c *WSClient) Call(ctx context.Context, method string, id types.JSONRPCStringID, params map[string]interface{}) error {
	if c.IsReconnecting(){
		return fmt.Errorf("websocket is reconnecting, can't send any request")
	}
	request, err := types.MapToRequest(c.cdc, id, method, params)
	if err != nil {
		return err
	}
	return c.Send(ctx, request)
}

// CallWithArrayParams the given method with params in a form of array. See
// Send description.
func (c *WSClient) CallWithArrayParams(ctx context.Context, method string, params []interface{}) error {
	request, err := types.ArrayToRequest(c.cdc, types.JSONRPCStringID("ws-client"), method, params)
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

func (c *WSClient) GenRequestId() (types.JSONRPCStringID, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return types.JSONRPCStringID(id.String()), nil
}

func (c *WSClient) EmptyRequest() types.JSONRPCStringID {
	return types.JSONRPCStringID("")
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
	c.conn = conn
	return nil
}

// reconnect tries to redial up to maxReconnectAttempts with exponential
// backoff.
func (c *WSClient) reconnect() error {
	attempt := 0

	c.mtx.Lock()
	c.reconnecting = true
	c.mtx.Unlock()
	defer func() {
		c.mtx.Lock()
		c.reconnecting = false
		c.mtx.Unlock()
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

		if c.maxReconnectAttempts >=0 && attempt > c.maxReconnectAttempts {
			return errors.Wrap(err, "reached maximum reconnect attempts")
		}
	}
}

func (c *WSClient) startReadWriteRoutines() {
	c.wg.Add(2)
	c.readRoutineQuit = make(chan struct{})
	go c.readRoutine()
	go c.writeRoutine()
}

func (c *WSClient) reconnectRoutine() {
	for {
		select {
		case originalError := <-c.reconnectAfter:
			// wait until writeRoutine and readRoutine finish
			c.wg.Wait()
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
		if err := c.conn.Close(); err != nil {
			// ignore error; it will trigger in tests
			// likely because it's closing an already closed connection
		}
		c.wg.Done()
	}()

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
				c.reconnectAfter <- err
				// add request to the backlog, so we don't lose it
				//c.backlog <- request
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
				c.reconnectAfter <- err
				return
			}
			c.mtx.Lock()
			c.sentLastPingAt = time.Now()
			c.mtx.Unlock()
			c.Logger.Debug("sent ping")
		case <-c.readRoutineQuit:
			return
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
	defer func() {
		if err := c.conn.Close(); err != nil {
			// ignore error; it will trigger in tests
			// likely because it's closing an already closed connection
		}
		c.wg.Done()
	}()

	c.conn.SetPongHandler(func(string) error {
		// gather latency stats
		c.mtx.RLock()
		t := c.sentLastPingAt
		c.mtx.RUnlock()
		c.PingPongLatencyTimer.UpdateSince(t)

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
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				return
			}

			c.Logger.Error("failed to read response", "err", err)
			close(c.readRoutineQuit)
			c.reconnectAfter <- err
			return
		}

		var response types.RPCResponse
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
		case c.ResponsesCh <- response:
		}
	}
}

///////////////////////////////////////////////////////////////////////////////
// Predefined methods

// Subscribe to a query. Note the server must have a "subscribe" route
// defined.
func (c *WSClient) Subscribe(ctx context.Context, id types.JSONRPCStringID, query string) error {
	params := map[string]interface{}{"query": query}
	return c.Call(ctx, "subscribe", id, params)
}

// Unsubscribe from a query. Note the server must have a "unsubscribe" route
// defined.
func (c *WSClient) Unsubscribe(ctx context.Context, id types.JSONRPCStringID, query string) error {
	params := map[string]interface{}{"query": query}

	return c.Call(ctx, "unsubscribe", id, params)
}

// UnsubscribeAll from all. Note the server must have a "unsubscribe_all" route
// defined.
func (c *WSClient) UnsubscribeAll(ctx context.Context, id types.JSONRPCStringID) error {
	params := map[string]interface{}{}
	return c.Call(ctx, "unsubscribe_all", id, params)
}

func (c *WSClient) Status(ctx context.Context, id types.JSONRPCStringID) error {
	return c.Call(ctx, "status", id, map[string]interface{}{})
}

func (c *WSClient) ABCIInfo(ctx context.Context, id types.JSONRPCStringID) error {
	return c.Call(ctx, "abci_info", id, map[string]interface{}{})
}

func (c *WSClient) ABCIQueryWithOptions(ctx context.Context, id types.JSONRPCStringID, path string, data cmn.HexBytes, opts client.ABCIQueryOptions) error {
	return c.Call(ctx, "abci_query", id, map[string]interface{}{"path": path, "data": data, "height": opts.Height, "prove": opts.Prove})
}

func (c *WSClient) BroadcastTxCommit(ctx context.Context, id types.JSONRPCStringID, tx ctypes.Tx) error {
	return c.Call(ctx, "broadcast_tx_commit", id, map[string]interface{}{"tx": tx})
}

func (c *WSClient) BroadcastTx(ctx context.Context, id types.JSONRPCStringID, route string, tx ctypes.Tx) error {
	return c.Call(ctx, route, id, map[string]interface{}{"tx": tx})
}

func (c *WSClient) UnconfirmedTxs(ctx context.Context, id types.JSONRPCStringID, limit int) error {
	return c.Call(ctx, "unconfirmed_txs", id, map[string]interface{}{"limit": limit})
}

func (c *WSClient) NumUnconfirmedTxs(ctx context.Context, id types.JSONRPCStringID) error {
	return c.Call(ctx, "num_unconfirmed_txs", id, map[string]interface{}{})
}

func (c *WSClient) NetInfo(ctx context.Context, id types.JSONRPCStringID) error {
	return c.Call(ctx, "net_info", id, map[string]interface{}{})
}

func (c *WSClient) DumpConsensusState(ctx context.Context, id types.JSONRPCStringID) error {
	return c.Call(ctx, "dump_consensus_state", id, map[string]interface{}{})
}

func (c *WSClient) ConsensusState(ctx context.Context, id types.JSONRPCStringID) error {
	return c.Call(ctx, "consensus_state", id, map[string]interface{}{})
}

func (c *WSClient) Health(ctx context.Context, id types.JSONRPCStringID) error {
	return c.Call(ctx, "health", id, map[string]interface{}{})
}

func (c *WSClient) BlockchainInfo(ctx context.Context, id types.JSONRPCStringID, minHeight, maxHeight int64) error {
	return c.Call(ctx, "blockchain", id,
		map[string]interface{}{"minHeight": minHeight, "maxHeight": maxHeight})
}

func (c *WSClient) Genesis(ctx context.Context, id types.JSONRPCStringID) error {
	return c.Call(ctx, "genesis", id, map[string]interface{}{})
}

func (c *WSClient) Block(ctx context.Context, id types.JSONRPCStringID, height *int64) error {
	return c.Call(ctx, "block", id, map[string]interface{}{"height": height})
}

func (c *WSClient) BlockResults(ctx context.Context, id types.JSONRPCStringID, height *int64) error {
	return c.Call(ctx, "block_results", id, map[string]interface{}{"height": height})
}

func (c *WSClient) Commit(ctx context.Context, id types.JSONRPCStringID, height *int64) error {
	return c.Call(ctx, "commit", id, map[string]interface{}{"height": height})
}
func (c *WSClient) Tx(ctx context.Context, id types.JSONRPCStringID, hash []byte, prove bool) error {
	params := map[string]interface{}{
		"hash":  hash,
		"prove": prove,
	}
	return c.Call(ctx, "tx", id, params)
}

func (c *WSClient) TxSearch(ctx context.Context, id types.JSONRPCStringID, query string, prove bool, page, perPage int) error {
	params := map[string]interface{}{
		"query":    query,
		"prove":    prove,
		"page":     page,
		"per_page": perPage,
	}
	return c.Call(ctx, "tx_search", id, params)
}

func (c *WSClient) Validators(ctx context.Context, id types.JSONRPCStringID, height *int64) error {
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
		return net.Dial(protocol, address)
	}
}
