package e2e

import (
	"fmt"
	"github.com/bnb-chain/go-sdk/types"
	"strings"
	"testing"
	time2 "time"

	"github.com/bnb-chain/go-sdk/client/rpc"
	"github.com/tendermint/tendermint/crypto"

	"github.com/stretchr/testify/assert"

	"github.com/tendermint/tendermint/types/time"

	sdk "github.com/bnb-chain/go-sdk/client"
	ctypes "github.com/bnb-chain/go-sdk/common/types"
	"github.com/bnb-chain/go-sdk/keys"
	"github.com/bnb-chain/go-sdk/types/msg"
	txtype "github.com/bnb-chain/go-sdk/types/tx"
)

// After bnbchain integration_test.sh has runned
func TestTransProcess(t *testing.T) {
	//----- Recover account ---------
	mnemonic := "test mnemonic"
	baeUrl := "testnet-dex.binance.org"
	keyManager, err := keys.NewMnemonicKeyManager(mnemonic)
	assert.NoError(t, err)
	testAccount1 := keyManager.GetAddr()
	testKeyManager2, _ := keys.NewKeyManager()
	testAccount2 := testKeyManager2.GetAddr()
	testKeyManager3, _ := keys.NewKeyManager()
	testAccount3 := testKeyManager3.GetAddr()

	//-----   Init sdk  -------------
	client, err := sdk.NewDexClient(baeUrl, ctypes.TestNetwork, keyManager)

	assert.NoError(t, err)
	nativeSymbol := types.NativeSymbol

	//---- set Account flags
	addFlags, err := client.AddAccountFlags([]ctypes.FlagOption{ctypes.TransferMemoCheckerFlag}, true)
	assert.NoError(t, err)
	fmt.Printf("Set account flags: %v \n", addFlags)
	accn, _ := client.GetAccount(client.GetKeyManager().GetAddr().String())
	fmt.Println(accn)
	setFlags, err := client.SetAccountFlags(0, true)
	assert.NoError(t, err)
	fmt.Printf("Set account flags: %v \n", setFlags)

	//-----  Get account  -----------
	account, err := client.GetAccount(testAccount1.String())
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.True(t, len(account.Balances) > 0)

	//-----  Get Tokens  -----------
	tokens, err := client.GetTokens(ctypes.NewTokensQuery().WithLimit(101))
	assert.NoError(t, err)
	fmt.Printf("GetTokens: %v \n", tokens)

	//-----  Get Time    -----------
	time, err := client.GetTime()
	assert.NoError(t, err)
	fmt.Printf("Get time: %v \n", time)

	//-----   time lock  -----------
	lockResult, err := client.TimeLock("test lock", ctypes.Coins{{"BNB", 100000000}}, int64(time2.Now().Add(65*time2.Second).Unix()), true)
	assert.NoError(t, err)
	fmt.Printf("timelock %d \n", lockResult.LockId)

	//----- time relock ---------
	relockResult, err := client.TimeReLock(lockResult.LockId, "test lock", ctypes.Coins{{"BNB", 200000000}}, int64(time2.Now().Add(70*time2.Second).Unix()), true)
	assert.NoError(t, err)
	fmt.Printf("timelock %d \n", relockResult.LockId)

	//------ time unlock --------
	time2.Sleep(75 * time2.Second)
	unlockResult, err := client.TimeUnLock(relockResult.LockId, true)
	assert.NoError(t, err)
	fmt.Printf("timelock %d \n", unlockResult.LockId)

	//----   Send tx  -----------
	send, err := client.SendToken([]msg.Transfer{{testAccount2, []ctypes.Coin{{nativeSymbol, 100000000}}}, {testAccount3, []ctypes.Coin{{nativeSymbol, 100000000}}}}, true)
	assert.NoError(t, err)
	assert.True(t, send.Ok)
	fmt.Printf("Send token: %v\n", send)

	//---    Get test2 account-----
	newTestAccout2, err := client.GetAccount(testAccount2.String())
	assert.NoError(t, err)
	for _, c := range newTestAccout2.Balances {
		if c.Symbol == nativeSymbol {
			fmt.Printf("test account BNB: %s \n", c.Free)
		}
	}

	//----   Freeze Token ---------
	freeze, err := client.FreezeToken(nativeSymbol, 100, true)
	assert.NoError(t, err)
	assert.True(t, freeze.Ok)
	fmt.Printf("freeze token: %v\n", freeze)

	//----   Unfreeze Token ---------
	unfreeze, err := client.UnfreezeToken(nativeSymbol, 100, true)
	assert.NoError(t, err)
	assert.True(t, unfreeze.Ok)
	fmt.Printf("Unfreeze token: %v\n", unfreeze)

	//----   issue token ---------
	issue, err := client.IssueToken("Client-Token", "sdk", 10000000000, true, true)
	assert.NoError(t, err)
	fmt.Printf("Issue token: %v\n", issue)

	//---  check issue success ---
	time2.Sleep(4 * time2.Second)
	issueresult, err := client.GetTx(issue.Hash)
	assert.NoError(t, err)
	assert.True(t, issueresult.Code == txtype.CodeOk)

	//--- mint token -----------
	mint, err := client.MintToken(issue.Symbol, 100000000, true)
	assert.NoError(t, err)
	fmt.Printf("Mint token: %v\n", mint)

	//---- Submit Proposal ------
	time2.Sleep(4 * time2.Second)
	listTradingProposal, err := client.SubmitListPairProposal("New trading pair", msg.ListTradingPairParams{issue.Symbol, nativeSymbol, 1000000000, "my trade", time2.Now().Add(1 * time2.Hour)}, 200000000000, 20*time2.Second, true)
	assert.NoError(t, err)
	fmt.Printf("Submit list trading pair: %v\n", listTradingProposal)

	//---  check submit proposal success ---
	time2.Sleep(4 * time2.Second)
	submitPorposalStatus, err := client.GetTx(listTradingProposal.Hash)
	assert.NoError(t, err)
	assert.True(t, submitPorposalStatus.Code == txtype.CodeOk)

	//----   issue mini token ---------
	time2.Sleep(4 * time2.Second)
	miniIssue, err := client.IssueMiniToken("Mini-Client-Token", "msdk", 10000000000, true, true, "http://test.sdk")
	assert.NoError(t, err)
	fmt.Printf("Issue mini token: %v\n", miniIssue)

	//----   issue tiny token ---------
	time2.Sleep(4 * time2.Second)
	tinyIssue, err := client.IssueMiniToken("Tiny-Client-Token", "tsdk", 10000000000, true, true, "http://test.sdk")
	assert.NoError(t, err)
	fmt.Printf("Issue tiny token: %v\n", tinyIssue)

	//----   set mini token uri ---------
	time2.Sleep(4 * time2.Second)
	setUri, err := client.SetURI(miniIssue.Symbol, "http://test-uri.sdk", true)
	assert.NoError(t, err)
	fmt.Printf("Set mini token uri: %v\n", setUri)

	//-----  Get Mini Tokens  -----------
	time2.Sleep(4 * time2.Second)
	miniTokens, err := client.GetMiniTokens(ctypes.NewTokensQuery().WithLimit(101))
	assert.NoError(t, err)
	fmt.Printf("Get Mini Tokens: %v \n", miniTokens)
}

func TestAtomicSwap(t *testing.T) {
	mnemonic1 := "test mnemonic1"
	mnemonic2 := "test mnemonic2"
	baeUrl := "testnet-dex.binance.org"
	keyManager, err := keys.NewMnemonicKeyManager(mnemonic1)
	assert.NoError(t, err)
	testAccount1 := keyManager.GetAddr()
	testKeyManager2, err := keys.NewMnemonicKeyManager(mnemonic2)
	assert.NoError(t, err)
	testAccount2 := testKeyManager2.GetAddr()

	client, err := sdk.NewDexClient(baeUrl, ctypes.TestNetwork, keyManager)
	assert.NoError(t, err)

	randomNumber := crypto.CRandBytes(32)
	timestamp := int64(time.Now().Unix())
	randomNumberHash := msg.CalculateRandomHash(randomNumber, timestamp)
	recipientOtherChain := "0x491e71b619878c083eaf2894718383c7eb15eb17"
	senderOtherChain := "0x833914c3A745d924bf71d98F9F9Ae126993E3C88"
	amount := ctypes.Coins{ctypes.Coin{"BNB", 10000}}
	expetedIncome := "10000:BNB"
	heightSpan := int64(1000)
	_, err = client.HTLT(testAccount2, recipientOtherChain, senderOtherChain, randomNumberHash, timestamp, amount, expetedIncome, heightSpan, true, true)
	assert.NoError(t, err)
	time2.Sleep(4 * time2.Second)
	swapID := msg.CalculateSwapID(randomNumberHash, testAccount1, senderOtherChain)
	_, err = client.ClaimHTLT(swapID, randomNumber, true)
	assert.NoError(t, err)
	time2.Sleep(4 * time2.Second)

	randomNumber = crypto.CRandBytes(32)
	timestamp = int64(time.Now().Unix())
	randomNumberHash = msg.CalculateRandomHash(randomNumber, timestamp)
	heightSpan = int64(360)
	_, err = client.HTLT(testAccount2, recipientOtherChain, senderOtherChain, randomNumberHash, timestamp, amount, expetedIncome, heightSpan, true, true)
	assert.NoError(t, err)
	time2.Sleep(4 * time2.Second)
	swapID1 := msg.CalculateSwapID(randomNumberHash, testAccount1, senderOtherChain)
	_, err = client.RefundHTLT(swapID1, true)
	assert.Error(t, err)
	time2.Sleep(4 * time2.Second)
	assert.True(t, strings.Contains(err.Error(), "is still not reached"))

	randomNumber = crypto.CRandBytes(32)
	timestamp = int64(time.Now().Unix())
	randomNumberHash = msg.CalculateRandomHash(randomNumber, timestamp)
	amount = ctypes.Coins{ctypes.Coin{"BNB", 10000}}
	expetedIncome = "1000:BTC-271"
	heightSpan = int64(1000)
	_, err = client.HTLT(testAccount2, "", "", randomNumberHash, timestamp, amount, expetedIncome, heightSpan, false, true)
	assert.NoError(t, err)
	time2.Sleep(4 * time2.Second)
	swapID2 := msg.CalculateSwapID(randomNumberHash, testAccount1, "")
	depositAmount := ctypes.Coins{ctypes.Coin{"BTC-271", 1000}}
	client1, err := sdk.NewDexClient(baeUrl, ctypes.TestNetwork, testKeyManager2)
	assert.NoError(t, err)
	_, err = client1.DepositHTLT(swapID2, depositAmount, true)
	assert.NoError(t, err)
	time2.Sleep(4 * time2.Second)
	_, err = client.ClaimHTLT(swapID2, randomNumber, true)
	assert.NoError(t, err)
	time2.Sleep(4 * time2.Second)

	c := rpc.NewRPCClient("tcp://seed-pre-s3.binance.org:80", ctypes.TestNetwork)
	swap, err := c.GetSwapByID(swapID)
	assert.NoError(t, err)

	randomNumberHashList, err := c.GetSwapByCreator(swap.From.String(), 0, 100)
	assert.NoError(t, err)
	assert.True(t, len(randomNumberHashList) > 0)

	randomNumberHashList, err = c.GetSwapByRecipient(swap.To.String(), 0, 100)
	assert.NoError(t, err)
	assert.True(t, len(randomNumberHashList) > 0)
}
