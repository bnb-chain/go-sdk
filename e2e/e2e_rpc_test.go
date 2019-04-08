package e2e

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	ctypes "github.com/binance-chain/go-sdk/common/types"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/binance-chain/go-sdk/client/rpc"
	"github.com/stretchr/testify/assert"
	tmquery "github.com/tendermint/tendermint/libs/pubsub/query"
)

var (
	nodeAddr     = "tcp://data-seed-pre-1-s3.binance.org:80"
	badAddr      = "tcp://127.0.0.1:80"
	testTxHash   = "A9DBDB2052FEEA13B953B40F8E6D3D0B0D0C592A9A0736A99BA4A4C31A3E33C8"
	testTxHeight = 6064550

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
	time.Sleep(10*time.Second)
	status, err = c.Status()
	assert.Error(t, err)
	fmt.Println(err)
	time.Sleep(10*time.Second)
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
	tokens, err := c.ListAllTokens(1,10)
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
	acc,err:=ctypes.AccAddressFromBech32("tbnb1z7sr92ar6njy9f80r4zl5rjtgm0hsej4686asa")
	assert.NoError(t,err)
	account,err:=c.GetAccount(acc)
	assert.NoError(t,err)
	bz, err := json.Marshal(account)
	fmt.Println(string(bz))
}

