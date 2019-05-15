package e2e

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"os/exec"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	tmquery "github.com/tendermint/tendermint/libs/pubsub/query"
	"github.com/tendermint/tendermint/types"

	"github.com/binance-chain/go-sdk/client/rpc"
	ctypes "github.com/binance-chain/go-sdk/common/types"
)

var (
	nodeAddr           = "tcp://127.0.0.1:80"
	badAddr            = "tcp://127.0.0.1:80"
	testTxHash         = "A27C20143E6B7D8160B50883F81132C1DFD0072FF2C1FE71E0158FBD001E23E4"
	testTxHeight       = 8669273
	testAddress        = "tbnb1l6vgk5yyxcalm06gdsg55ay4pjkfueazkvwh58"
	testDelAddr        = "tbnb12hlquylu78cjylk5zshxpdj6hf3t0tahwjt3ex"
	testTradePair      = "X00-243_BNB"
	testTxStr          = "xxx"
	onceClient         = sync.Once{}
	testClientInstance *rpc.HTTP
)

func startBnbchaind(t *testing.T) *exec.Cmd {
	cmd := exec.Command("bnbchaind", "start", "--home", "testnoded")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Start()
	assert.NoError(t, err)
	// wait for completely start
	time.Sleep(15 * time.Second)
	return cmd
}

func defaultClient() *rpc.HTTP {
	onceClient.Do(func() {
		testClientInstance = rpc.NewRPCClient(nodeAddr, ctypes.TestNetwork)
	})
	return testClientInstance
}

func TestRPCGetProposals(t *testing.T) {
	c := defaultClient()
	statuses:= []ctypes.ProposalStatus{
		ctypes.StatusDepositPeriod,
		ctypes.StatusVotingPeriod,
		ctypes.StatusPassed,
		ctypes.StatusRejected,
	}
	for _,s:=range statuses{
		proposals, err := c.GetProposals(s, 100)
		assert.NoError(t, err)
		for _,p:=range proposals{
			assert.Equal(t,p.GetStatus(),s)
		}
		bz, err := json.Marshal(proposals)
		fmt.Println(string(bz))
	}
}

func TestRPCGetProposal(t *testing.T) {
	c := defaultClient()
	proposal, err := c.GetProposal(int64(1))
	assert.NoError(t, err)
	bz, err := json.Marshal(proposal)
	fmt.Println(string(bz))
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
	repeatNum := 10
	c := defaultClient()

	// Find error
	time.Sleep(1 * time.Second)
	for i := 0; i < repeatNum; i++ {
		_, err := c.Status()
		assert.Error(t, err)
	}

	// Reconnect and find no error
	cmd := startBnbchaind(t)

	for i := 0; i < repeatNum; i++ {
		status, err := c.Status()
		assert.NoError(t, err)
		bz, err := json.Marshal(status)
		fmt.Println(string(bz))
	}

	// kill process
	err := cmd.Process.Kill()
	assert.NoError(t, err)
	err = cmd.Process.Release()
	assert.NoError(t, err)
	time.Sleep(1 * time.Second)

	// Find error
	for i := 0; i < repeatNum; i++ {
		_, err := c.Status()
		assert.Error(t, err)
	}

	// Restart bnbchain
	cmd = startBnbchaind(t)

	// Find no error
	for i := 0; i < repeatNum; i++ {
		status, err := c.Status()
		assert.NoError(t, err)
		bz, _ := json.Marshal(status)
		fmt.Println(string(bz))
	}

	// Stop bnbchain
	cmd.Process.Kill()
	cmd.Process.Release()
}

func TestTxSearch(t *testing.T) {
	c := defaultClient()
	tx, err := c.TxInfoSearch(fmt.Sprintf("tx.height=%d", testTxHeight), false, 1, 10)
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
	c := rpc.NewRPCClient(badAddr, ctypes.TestNetwork)
	_, err := c.Validators(nil)
	assert.Error(t, err, "context deadline exceeded")
}

func TestSetTimeOut(t *testing.T) {
	c := rpc.NewRPCClient(badAddr, ctypes.TestNetwork)
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
	time.Sleep(10 * time.Second)
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
	c.SetTimeOut(1 * time.Second)
	w := sync.WaitGroup{}
	w.Add(2000)
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
		//TestBlockResults,
		TestCommit,
		//TestTx,
		//TestTxSearch,
		TestValidators,
	}
	for i := 0; i < 2000; i++ {
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
	fmt.Println(hex.EncodeToString(account.GetAddress().Bytes()))

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

func TestGetTradePair(t *testing.T) {
	c := defaultClient()
	trades, err := c.GetTradingPairs(0, 10)
	assert.NoError(t, err)
	bz, err := json.Marshal(trades)
	fmt.Println(string(bz))
}

func TestGetDepth(t *testing.T) {
	c := defaultClient()
	depth, err := c.GetDepth(testTradePair)
	assert.NoError(t, err)
	bz, err := json.Marshal(depth)
	fmt.Println(string(bz))
}

func TestBroadcastTxCommit(t *testing.T) {
	c := defaultClient()
	txbyte, err := hex.DecodeString(testTxStr)
	assert.NoError(t, err)
	res, err := c.BroadcastTxCommit(types.Tx(txbyte))
	assert.NoError(t, err)
	fmt.Println(res)
}

func TestGetStakeValidators(t *testing.T) {
	c := defaultClient()
	ctypes.Network = ctypes.TestNetwork
	vals, err := c.GetStakeValidators()
	assert.NoError(t, err)
	bz, err := json.Marshal(vals)
	fmt.Println(string(bz))
}

func TestGetDelegatorUnbondingDelegations(t *testing.T) {
	c := defaultClient()
	ctypes.Network = ctypes.TestNetwork
	acc, err := ctypes.AccAddressFromBech32(testDelAddr)
	assert.NoError(t, err)
	vals, err := c.GetDelegatorUnbondingDelegations(acc)
	assert.NoError(t, err)
	bz, err := json.Marshal(vals)
	fmt.Println(string(bz))
}

func TestNoRequestLeakInBadNetwork(t *testing.T) {
	c := rpc.NewRPCClient(badAddr, ctypes.TestNetwork)
	c.SetTimeOut(1 * time.Second)
	w := sync.WaitGroup{}
	w.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			c.GetFee()
			w.Done()
		}()
	}
	w.Wait()
	assert.Equal(t, c.PendingRequest(), 0)
}

func TestNoRequestLeakInGoodNetwork(t *testing.T) {
	c := defaultClient()
	c.SetTimeOut(1 * time.Second)
	w := sync.WaitGroup{}
	w.Add(3000)
	for i := 0; i < 3000; i++ {
		go func() {
			_, err := c.Block(nil)
			assert.NoError(t, err)
			//bz, err := json.Marshal(fees)
			//fmt.Println(string(bz))
			w.Done()
		}()
	}
	w.Wait()
	assert.Equal(t, c.PendingRequest(), 0)
}
