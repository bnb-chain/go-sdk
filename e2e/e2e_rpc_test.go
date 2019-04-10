package e2e

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	ctypes "github.com/binance-chain/go-sdk/common/types"
	"github.com/tendermint/tendermint/types"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/binance-chain/go-sdk/client/rpc"
	"github.com/stretchr/testify/assert"
	tmquery "github.com/tendermint/tendermint/libs/pubsub/query"
)

var (
	nodeAddr      = "tcp://seed-pre-s3.binance.org:80"
	badAddr       = "tcp://127.0.0.1:80"
	testTxHash    = "9E9E6EA3FA13684DD260DB627144EABDB50F2C205DE733447C5E8415311670C9"
	testTxHeight  = 960284
	testAddress   = "tbnb1l6vgk5yyxcalm06gdsg55ay4pjkfueazkvwh58"
	testTradePair = "X00-243_BNB"

	onceClient         = sync.Once{}
	testClientInstance *rpc.HTTP
)

func defaultClient() *rpc.HTTP {
	onceClient.Do(func() {
		testClientInstance = rpc.NewRPCClient(nodeAddr)
	})
	return testClientInstance
}

func TestRPCStatus(t *testing.T) {
	c := defaultClient()
	status, err := c.Status()
	assert.NoError(t, err)
	bz, err := json.Marshal(status)
	fmt.Println(string(bz))
}

func TestRPCABCIInfo(t *testing.T) {
	c := defaultClient()
	info, err := c.ABCIInfo()
	assert.NoError(t, err)
	bz, err := json.Marshal(info)
	fmt.Println(string(bz))
}

func TestUnconfirmedTxs(t *testing.T) {
	c := defaultClient()
	txs, err := c.UnconfirmedTxs(10)
	assert.NoError(t, err)
	bz, err := json.Marshal(txs)
	fmt.Println(string(bz))
}

func TestNumUnconfirmedTxs(t *testing.T) {
	c := defaultClient()
	numTxs, err := c.NumUnconfirmedTxs()
	assert.NoError(t, err)
	bz, err := json.Marshal(numTxs)
	fmt.Println(string(bz))
}

func TestNetInfo(t *testing.T) {
	c := defaultClient()
	netInfo, err := c.NetInfo()
	assert.NoError(t, err)
	bz, err := json.Marshal(netInfo)
	fmt.Println(string(bz))
}

func TestDumpConsensusState(t *testing.T) {
	c := defaultClient()
	state, err := c.DumpConsensusState()
	assert.NoError(t, err)
	bz, err := json.Marshal(state)
	fmt.Println(string(bz))
}

func TestConsensusState(t *testing.T) {
	c := defaultClient()
	state, err := c.ConsensusState()
	assert.NoError(t, err)
	bz, err := json.Marshal(state)
	fmt.Println(string(bz))
}

func TestHealth(t *testing.T) {
	c := defaultClient()
	health, err := c.Health()
	assert.NoError(t, err)
	bz, err := json.Marshal(health)
	fmt.Println(string(bz))
}

func TestBlockchainInfo(t *testing.T) {
	c := defaultClient()
	blockInfos, err := c.BlockchainInfo(1, 5)
	assert.NoError(t, err)
	bz, err := json.Marshal(blockInfos)
	fmt.Println(string(bz))
}

func TestGenesis(t *testing.T) {
	c := defaultClient()
	genesis, err := c.Genesis()
	assert.NoError(t, err)
	bz, err := json.Marshal(genesis)
	fmt.Println(string(bz))
}

func TestBlock(t *testing.T) {
	c := defaultClient()
	block, err := c.Block(nil)
	assert.NoError(t, err)
	bz, err := json.Marshal(block)
	fmt.Println(string(bz))
}

func TestBlockResults(t *testing.T) {
	c := defaultClient()
	block, err := c.BlockResults(nil)
	assert.NoError(t, err)
	bz, err := json.Marshal(block)
	fmt.Println(string(bz))
}

func TestCommit(t *testing.T) {
	c := defaultClient()
	commit, err := c.Commit(nil)
	assert.NoError(t, err)
	bz, err := json.Marshal(commit)
	fmt.Println(string(bz))
}

func TestTx(t *testing.T) {
	c := defaultClient()
	bz, err := hex.DecodeString(testTxHash)
	assert.NoError(t, err)

	tx, err := c.Tx(bz, false)
	assert.NoError(t, err)
	bz, err = json.Marshal(tx)
	fmt.Println(string(bz))
}

func TestReconnection(t *testing.T) {
	c := defaultClient()
	status, err := c.Status()
	assert.NoError(t, err)
	bz, err := json.Marshal(status)
	fmt.Println(string(bz))
	time.Sleep(10 * time.Second)
	status, err = c.Status()
	assert.Error(t, err)
	fmt.Println(err)
	time.Sleep(10 * time.Second)
	status, err = c.Status()
	assert.Error(t, err)
	fmt.Println(err)
	bz, err = json.Marshal(status)
	fmt.Println(string(bz))
}

func TestTxSearch(t *testing.T) {
	c := defaultClient()

	tx, err := c.TxSearch(fmt.Sprintf("tx.height=%d", testTxHeight), false, 1, 10)
	assert.NoError(t, err)
	bz, err := json.Marshal(tx)
	fmt.Println(string(bz))
}

func TestValidators(t *testing.T) {
	c := defaultClient()
	validators, err := c.Validators(nil)
	assert.NoError(t, err)
	bz, err := json.Marshal(validators)
	fmt.Println(string(bz))
}

func TestBadNodeAddr(t *testing.T) {
	c := rpc.NewRPCClient(badAddr)
	_, err := c.Validators(nil)
	assert.Error(t, err, "context deadline exceeded")
}

func TestSetTimeOut(t *testing.T) {
	c := rpc.NewRPCClient(badAddr)
	c.SetTimeOut(1 * time.Second)
	before := time.Now()
	_, err := c.Validators(nil)
	duration := time.Now().Sub(before).Seconds()
	assert.True(t, duration > 1)
	assert.True(t, duration < 2)
	assert.Error(t, err, "context deadline exceeded")
}

func TestSubscribeEvent(t *testing.T) {
	c := defaultClient()
	query := "tm.event = 'CompleteProposal'"
	_, err := tmquery.New(query)
	assert.NoError(t, err)
	out, err := c.Subscribe(query, 10)
	assert.NoError(t, err)
	noMoreEvent := make(chan struct{}, 1)
	go func() {
		for {
			select {
			case o := <-out:
				bz, err := json.Marshal(o)
				assert.NoError(t, err)
				fmt.Println(string(bz))
			case <-noMoreEvent:
				fmt.Println("no more event after")
			}
		}
	}()
	time.Sleep(100 * time.Second)
	err = c.Unsubscribe(query)
	noMoreEvent <- struct{}{}
	assert.NoError(t, err)
	time.Sleep(3 * time.Second)
}

func TestSubscribeEventTwice(t *testing.T) {
	c := defaultClient()
	query := "tm.event = 'CompleteProposal'"
	_, err := tmquery.New(query)
	assert.NoError(t, err)
	_, err = c.Subscribe(query, 10)
	assert.NoError(t, err)
	_, err = c.Subscribe(query, 10)
	assert.Error(t, err)
}

func TestReceiveWithRequestId(t *testing.T) {
	c := defaultClient()
	c.SetTimeOut(5 * time.Second)
	w := sync.WaitGroup{}
	w.Add(10)
	testCases := []func(t *testing.T){
		TestRPCStatus,
		TestRPCABCIInfo,
		TestUnconfirmedTxs,
		TestNumUnconfirmedTxs,
		TestNetInfo,
		TestDumpConsensusState,
		TestConsensusState,
		TestHealth,
		TestBlockchainInfo,
		TestGenesis,
		TestBlock,
		TestBlockResults,
		TestCommit,
		TestTx,
		TestTxSearch,
		//TestValidators,
	}
	for i := 0; i < 10; i++ {
		testFuncIndex := rand.Intn(len(testCases))
		go func() {
			testCases[testFuncIndex](t)
			w.Done()
		}()
	}
	w.Wait()
}

func TestListAllTokens(t *testing.T) {
	c := defaultClient()
	tokens, err := c.ListAllTokens(1, 10)
	assert.NoError(t, err)
	bz, err := json.Marshal(tokens)
	fmt.Println(string(bz))
}

func TestGetTokenInfo(t *testing.T) {
	c := defaultClient()
	token, err := c.GetTokenInfo("BNB")
	assert.NoError(t, err)
	bz, err := json.Marshal(token)
	fmt.Println(string(bz))
}

func TestGetAccount(t *testing.T) {
	ctypes.Network = ctypes.TestNetwork
	c := defaultClient()
	acc, err := ctypes.AccAddressFromBech32(testAddress)
	assert.NoError(t, err)
	account, err := c.GetAccount(acc)
	assert.NoError(t, err)
	bz, err := json.Marshal(account)
	fmt.Println(string(bz))
}

func TestGetBalances(t *testing.T) {
	ctypes.Network = ctypes.TestNetwork
	c := defaultClient()
	acc, err := ctypes.AccAddressFromBech32(testAddress)
	assert.NoError(t, err)
	balances, err := c.GetBalances(acc)
	assert.NoError(t, err)
	bz, err := json.Marshal(balances)
	fmt.Println(string(bz))
}

func TestGetBalance(t *testing.T) {
	ctypes.Network = ctypes.TestNetwork
	c := defaultClient()
	acc, err := ctypes.AccAddressFromBech32(testAddress)
	assert.NoError(t, err)
	balance, err := c.GetBalance(acc, "BNB")
	assert.NoError(t, err)
	bz, err := json.Marshal(balance)
	fmt.Println(string(bz))
}

func TestGetFees(t *testing.T) {
	c := defaultClient()
	fees, err := c.GetFee()
	assert.NoError(t, err)
	bz, err := json.Marshal(fees)
	fmt.Println(string(bz))
}

func TestGetOpenOrder(t *testing.T) {
	ctypes.Network = ctypes.TestNetwork
	acc, err := ctypes.AccAddressFromBech32(testAddress)
	assert.NoError(t, err)
	c := defaultClient()
	openorders, err := c.GetOpenOrders(acc, testTradePair)
	assert.NoError(t, err)
	bz, err := json.Marshal(openorders)
	assert.NoError(t, err)
	fmt.Println(string(bz))
}

func TestGetTradePair(t *testing.T){
	c := defaultClient()
	trades, err := c.GetTradingPairs(0,10)
	assert.NoError(t, err)
	bz, err := json.Marshal(trades)
	fmt.Println(string(bz))
}

func TestGetDepth(t *testing.T){
	c := defaultClient()
	depth, err := c.GetDepth(testTradePair)
	assert.NoError(t, err)
	bz, err := json.Marshal(depth)
	fmt.Println(string(bz))
}

func TestBroadcastTxCommit(t *testing.T){
	c := defaultClient()
	txstring:="cc01f0625dee0a4c2a2c87fa0a220a14443c2367e8e2edfc93aac1700bf843ef8be69c56120a0a03424e421080a3c34712220a1487dbcff17c64291c2b3538806c72a4d3a0ef6128120a0a03424e421080a3c34712700a26eb5ae9872102942fb6ffe96f001a15931e0702dd1c10370ffb568fd962039f0c4d2d45b53e9712408454253a4cf0e8f868276dfe2caa96b4ed7f94e8abace386b3fd69c454f7aa7d3a088e482328b94d991b6e6f1449cdb34e2a90bb81d102d0dac55488b35650ec18bcd82820021a04746573742001"
	txbyte,err:=hex.DecodeString(txstring)
	assert.NoError(t,err)
	res,err:=c.BroadcastTxCommit(types.Tx(txbyte))
	assert.NoError(t,err)
	fmt.Println(res)
}

func TestGetStakeValidators(t *testing.T){
	c := defaultClient()
	ctypes.Network = ctypes.TestNetwork
	vals,err:=c.GetStakeValidators()
	assert.NoError(t,err)
	bz, err := json.Marshal(vals)
	fmt.Println(string(bz))
}




