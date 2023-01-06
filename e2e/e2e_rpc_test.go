package e2e

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/bnb-chain/go-sdk/client/rpc"
	"github.com/bnb-chain/go-sdk/client/transaction"
	ctypes "github.com/bnb-chain/go-sdk/common/types"
	"github.com/bnb-chain/go-sdk/keys"
	stypes "github.com/bnb-chain/go-sdk/types"
	"github.com/bnb-chain/go-sdk/types/msg"
	"github.com/bnb-chain/go-sdk/types/tx"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/common"
	tmquery "github.com/tendermint/tendermint/libs/pubsub/query"
	"github.com/tendermint/tendermint/types"
)

var (
	nodeAddr           = "tcp://data-seed-pre-0-s3.binance.org:80"
	badAddr            = "tcp://127.0.0.1:80"
	testTxHash         = "6165507B990A1CCAD1758512382C6B76F952CC945ABB84D9BF18160C11DE902A"
	testTxHeight       = int64(34762951)
	testAddress        = "tbnb1e803p76n4rtyeclef7pg3295nurwfuwsw8l36m"
	testDelAddr        = "tbnb12hlquylu78cjylk5zshxpdj6hf3t0tahwjt3ex"
	testTradePair      = "PPC-00A_BNB"
	testTradeSymbol    = "000-0E1"
	testTxStr          = "xxx"
	testNewOwner       = "tbnb1rtzy6szuyzcj4amfn6uarvne8a5epxrdc28nhr"
	mnemonic           = "test mnemonic"
	onceClient         = sync.Once{}
	testClientInstance *rpc.HTTP

	scParams = `[{"type": "params/StakeParamSet","value": {"unbonding_time": "604800000000000","max_validators": 11,"bond_denom": "BNB","min_self_delegation": "5000000000000","min_delegation_change": "100000000"}},{"type": "params/SlashParamSet","value": {"max_evidence_age": "259200000000000","signed_blocks_window": "0","min_signed_per_window": "0","double_sign_unbond_duration": "9223372036854775807","downtime_unbond_duration": "172800000000000","too_low_del_unbond_duration": "86400000000000","slash_fraction_double_sign": "0","slash_fraction_downtime": "0","double_sign_slash_amount": "1000000000000","downtime_slash_amount": "5000000000","submitter_reward": "100000000000","downtime_slash_fee": "1000000000"}},{"type": "params/OracleParamSet","value": {"ConsensusNeeded": "70000000"}},{"type": "params/IbcParamSet","value": {"relayer_fee": "1000000"}}]`
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
	statuses := []ctypes.ProposalStatus{
		ctypes.StatusDepositPeriod,
		ctypes.StatusVotingPeriod,
		ctypes.StatusPassed,
		ctypes.StatusRejected,
	}
	for _, s := range statuses {
		proposals, err := c.GetProposals(s, 100)
		assert.NoError(t, err)
		for _, p := range proposals {
			assert.Equal(t, p.GetStatus(), s)
		}
		bz, err := json.Marshal(proposals)
		fmt.Println(string(bz))
	}
}
func TestRPCGetTimelocks(t *testing.T) {
	c := defaultClient()
	acc, err := ctypes.AccAddressFromBech32(testAddress)
	assert.NoError(t, err)
	records, err := c.GetTimelocks(acc)
	assert.NoError(t, err)
	fmt.Println(len(records))
	for _, record := range records {
		fmt.Println(record)
	}
}

func TestRPCGetTimelock(t *testing.T) {
	c := defaultClient()
	acc, err := ctypes.AccAddressFromBech32(testAddress)
	assert.NoError(t, err)
	record, err := c.GetTimelock(acc, 1)
	assert.NoError(t, err)
	fmt.Println(record)

}

func TestRPCGetProposal(t *testing.T) {
	c := defaultClient()
	proposal, err := c.GetProposal(int64(100))
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
	block, err := c.BlockResults(&testTxHeight)
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

func TestClaimTx(t *testing.T) {
	c := defaultClient()
	bz, err := hex.DecodeString(testTxHash)
	assert.NoError(t, err)

	rawTx, err := c.Tx(bz, false)
	assert.NoError(t, err)
	claimTx, err := rpc.ParseTx(tx.Cdc, rawTx.Tx)
	claimMsg := claimTx.GetMsgs()[0].(msg.ClaimMsg)
	packages, err := msg.ParseClaimPayload(claimMsg.Payload)
	assert.NoError(t, err)
	newBz, err := json.Marshal(packages)
	assert.NoError(t, err)
	fmt.Println(string(newBz))
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
	c.SetTimeOut(2 * time.Second)
	w := sync.WaitGroup{}
	w.Add(100)
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
	for i := 0; i < 100; i++ {
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
	tokens, err := c.ListAllTokens(0, 10)
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
	ctypes.SetNetwork(ctypes.TestNetwork)
	c := defaultClient()
	acc, err := ctypes.AccAddressFromBech32(testAddress)
	assert.NoError(t, err)
	account, err := c.GetAccount(acc)
	assert.NoError(t, err)
	bz, err := json.Marshal(account)
	fmt.Println(string(bz))
	fmt.Println(hex.EncodeToString(account.GetAddress().Bytes()))
}

func TestNoneExistGetAccount(t *testing.T) {
	ctypes.SetNetwork(ctypes.TestNetwork)
	c := defaultClient()
	acc, err := keys.NewKeyManager()
	account, err := c.GetAccount(acc.GetAddr())
	assert.NoError(t, err)
	assert.Nil(t, account)
}

func TestGetBalances(t *testing.T) {
	ctypes.SetNetwork(ctypes.TestNetwork)
	c := defaultClient()
	acc, err := ctypes.AccAddressFromBech32(testAddress)
	assert.NoError(t, err)
	balances, err := c.GetBalances(acc)
	assert.Equal(t, 1, len(balances))
	assert.NoError(t, err)
	bz, err := json.Marshal(balances)
	fmt.Println(string(bz))
}

func TestNoneExistGetBalances(t *testing.T) {
	ctypes.SetNetwork(ctypes.TestNetwork)
	c := defaultClient()
	acc, _ := keys.NewKeyManager()
	balances, err := c.GetBalances(acc.GetAddr())
	assert.NoError(t, err)
	bz, err := json.Marshal(balances)
	fmt.Println(string(bz))
}

func TestGetBalance(t *testing.T) {
	ctypes.SetNetwork(ctypes.TestNetwork)
	c := defaultClient()
	acc, err := ctypes.AccAddressFromBech32(testAddress)
	assert.NoError(t, err)
	balance, err := c.GetBalance(acc, "BNB")
	assert.NoError(t, err)
	bz, err := json.Marshal(balance)
	fmt.Println(string(bz))
}

func TestNoneExistGetBalance(t *testing.T) {
	ctypes.SetNetwork(ctypes.TestNetwork)
	c := defaultClient()
	acc, _ := keys.NewKeyManager()
	balance, err := c.GetBalance(acc.GetAddr(), "BNB")
	assert.NoError(t, err)
	assert.Equal(t, ctypes.Fixed8Zero, balance.Free)
	assert.Equal(t, ctypes.Fixed8Zero, balance.Locked)
	assert.Equal(t, ctypes.Fixed8Zero, balance.Frozen)
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

func TestSendToken(t *testing.T) {
	c := defaultClient()
	ctypes.SetNetwork(ctypes.TestNetwork)
	keyManager, err := keys.NewMnemonicKeyManager(mnemonic)
	assert.NoError(t, err)
	c.SetKeyManager(keyManager)
	testacc, err := ctypes.AccAddressFromBech32(testAddress)
	assert.NoError(t, err)
	res, err := c.SendToken([]msg.Transfer{{testacc, []ctypes.Coin{{"BNB", 100000}}}}, rpc.Sync, transaction.WithMemo("123"))
	assert.NoError(t, err)
	bz, err := json.Marshal(res)
	fmt.Println(string(bz))
}

func TestQuerySideChainParam(t *testing.T) {
	c := defaultClient()
	ctypes.SetNetwork(ctypes.TestNetwork)
	params, err := c.GetSideChainParams("chapel")
	assert.NoError(t, err)
	bz, _ := json.Marshal(params)
	fmt.Println(string(bz))
}

func TestSubmitSideProposal(t *testing.T) {
	c := defaultClient()
	ctypes.SetNetwork(ctypes.TestNetwork)

	keyManager, err := keys.NewMnemonicKeyManager(mnemonic)
	assert.NoError(t, err)
	c.SetKeyManager(keyManager)

	iScPrams := make([]msg.SCParam, 0)

	err = tx.Cdc.UnmarshalJSON([]byte(scParams), &iScPrams)
	assert.NoError(t, err)

	res, err := c.SideChainSubmitSCParamsProposal("title", msg.SCChangeParams{SCParams: iScPrams, Description: "des"}, ctypes.Coins{{stypes.NativeSymbol, 5e11}}, 5*time.Second, "rialto", rpc.Sync)
	assert.NoError(t, err)
	assert.True(t, res.Code == 0)
	proposalIdStr := string(res.Data)
	id, err := strconv.ParseInt(proposalIdStr, 10, 64)
	assert.NoError(t, err)
	res, err = c.SideChainVote(int64(id), msg.OptionYes, "rialto", rpc.Sync)
	assert.NoError(t, err)
	assert.True(t, res.Code == 0)
}

func TestSubmitCSCProposal(t *testing.T) {
	c := defaultClient()
	ctypes.SetNetwork(ctypes.TestNetwork)

	keyManager, err := keys.NewMnemonicKeyManager(mnemonic)
	assert.NoError(t, err)
	c.SetKeyManager(keyManager)

	cscPrams := msg.CSCParamChange{
		Key:    common.RandStr(common.RandIntn(255) + 1),
		Value:  hex.EncodeToString(common.RandBytes(common.RandIntn(255) + 1)),
		Target: hex.EncodeToString(common.RandBytes(20)),
	}

	res, err := c.SideChainSubmitCSCParamsProposal("title", cscPrams, ctypes.Coins{{stypes.NativeSymbol, 5e8}}, 5*time.Second, "chapel", rpc.Sync)
	assert.NoError(t, err)
	assert.True(t, res.Code == 0)

	proposalIdStr := string(res.Data)
	id, err := strconv.ParseInt(proposalIdStr, 10, 64)
	assert.NoError(t, err)
	res, err = c.SideChainDeposit(int64(id), ctypes.Coins{{"BNB", 1e8}}, "chapel", rpc.Sync)
	assert.NoError(t, err)
	assert.True(t, res.Code == 0)

}

func TestTransferTokenOwnership(t *testing.T) {
	c := defaultClient()
	ctypes.SetNetwork(ctypes.TestNetwork)
	keyManager, err := keys.NewMnemonicKeyManager(mnemonic)
	assert.NoError(t, err)
	c.SetKeyManager(keyManager)
	fmt.Println(keyManager.GetAddr().String())
	newOwner, err := ctypes.AccAddressFromBech32(testNewOwner)
	assert.NoError(t, err)
	result, err := c.TransferTokenOwnership(testTradeSymbol, newOwner, rpc.Commit)
	assert.NoError(t, err)
	bz, _ := json.Marshal(result)
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
	ctypes.SetNetwork(ctypes.TestNetwork)
	vals, err := c.GetStakeValidators()
	assert.NoError(t, err)
	bz, err := json.Marshal(vals)
	fmt.Println(string(bz))
}

func TestGetDelegatorUnbondingDelegations(t *testing.T) {
	c := defaultClient()
	ctypes.SetNetwork(ctypes.TestNetwork)
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
	c.SetTimeOut(2 * time.Second)
	w := sync.WaitGroup{}
	w.Add(100)
	for i := 0; i < 100; i++ {
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

func TestListAllMiniTokens(t *testing.T) {
	c := defaultClient()
	tokens, err := c.ListAllMiniTokens(0, 10)
	assert.NoError(t, err)
	bz, err := json.Marshal(tokens)
	fmt.Println(string(bz))
}

func TestGetMiniTokenInfo(t *testing.T) {
	c := defaultClient()
	tokens, err := c.ListAllMiniTokens(0, 10)
	assert.NoError(t, err)
	if len(tokens) > 0 {
		token, err := c.GetMiniTokenInfo(tokens[0].Symbol)
		assert.NoError(t, err)
		bz, err := json.Marshal(token)
		fmt.Println(string(bz))
	}
}
