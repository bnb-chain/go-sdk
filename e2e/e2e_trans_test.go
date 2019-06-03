package e2e

import (
	"fmt"
	"github.com/binance-chain/go-sdk/client/transaction"
	"testing"
	time2 "time"

	"github.com/stretchr/testify/assert"

	sdk "github.com/binance-chain/go-sdk/client"
	"github.com/binance-chain/go-sdk/common"
	ctypes "github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/keys"
	"github.com/binance-chain/go-sdk/types/msg"
	txtype "github.com/binance-chain/go-sdk/types/tx"
)

// After bnbchain integration_test.sh has runned
func TestTransProcess(t *testing.T) {
	//----- Recover account ---------
	mnemonic := "test mnemonic"
	baeUrl := "test base url"
	keyManager, err := keys.NewMnemonicKeyManager(mnemonic)
	assert.NoError(t, err)
	validatorMnemonics := []string{
		"test mnemonic",
	}
	testAccount1 := keyManager.GetAddr()
	testKeyManager2, _ := keys.NewKeyManager()
	testAccount2 := testKeyManager2.GetAddr()
	testKeyManager3, _ := keys.NewKeyManager()
	testAccount3 := testKeyManager3.GetAddr()

	//-----   Init sdk  -------------
	client, err := sdk.NewDexClient(baeUrl, ctypes.TestNetwork, keyManager)
	assert.NoError(t, err)
	nativeSymbol := msg.NativeToken

	//-----  Get account  -----------
	account, err := client.GetAccount(testAccount1.String())
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.True(t, len(account.Balances) > 0)

	//----- Get Markets  ------------
	markets, err := client.GetMarkets(ctypes.NewMarketsQuery().WithLimit(1))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(markets))
	tradeSymbol := markets[0].QuoteAssetSymbol
	if tradeSymbol == "BNB" {
		tradeSymbol = markets[0].BaseAssetSymbol
	}

	//-----  Get Depth  ----------
	depth, err := client.GetDepth(ctypes.NewDepthQuery(tradeSymbol, nativeSymbol))
	assert.NoError(t, err)
	assert.True(t, depth.Height > 0)

	//----- Get Kline
	kline, err := client.GetKlines(ctypes.NewKlineQuery(tradeSymbol, nativeSymbol, "1h").WithLimit(1))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(kline))

	//-----  Get Ticker 24h  -----------
	ticker24h, err := client.GetTicker24h(ctypes.NewTicker24hQuery().WithSymbol(tradeSymbol, nativeSymbol))
	assert.NoError(t, err)
	assert.True(t, len(ticker24h) > 0)

	//-----  Get Tokens  -----------
	tokens, err := client.GetTokens()
	assert.NoError(t, err)
	fmt.Printf("GetTokens: %v \n", tokens)

	//-----  Get Trades  -----------
	fmt.Println(testAccount1.String())
	trades, err := client.GetTrades(ctypes.NewTradesQuery(testAccount1.String(), true).WithSymbol(tradeSymbol, nativeSymbol))
	assert.NoError(t, err)
	fmt.Printf("GetTrades: %v \n", trades)

	//-----  Get Time    -----------
	time, err := client.GetTime()
	assert.NoError(t, err)
	fmt.Printf("Get time: %v \n", time)

	//----- Create order -----------
	createOrderResult, err := client.CreateOrder(tradeSymbol, nativeSymbol, msg.OrderSide.BUY, 100000000, 100000000, true, transaction.WithSource(100),transaction.WithMemo("test memo"))
	assert.NoError(t, err)
	assert.True(t, true, createOrderResult.Ok)

	//---- Get Open Order ---------
	openOrders, err := client.GetOpenOrders(ctypes.NewOpenOrdersQuery(testAccount1.String(), true))
	assert.NoError(t, err)
	assert.True(t, len(openOrders.Order) > 0)
	orderId := openOrders.Order[0].ID
	fmt.Printf("GetOpenOrders:  %v \n", openOrders)

	//---- Get Order    ------------
	order, err := client.GetOrder(orderId)
	assert.NoError(t, err)
	assert.Equal(t, common.CombineSymbol(tradeSymbol, nativeSymbol), order.Symbol)

	//---- Cancel Order  ---------
	time2.Sleep(2 * time2.Second)
	cancelOrderResult, err := client.CancelOrder(tradeSymbol, nativeSymbol, orderId, true)
	assert.NoError(t, err)
	assert.True(t, cancelOrderResult.Ok)
	fmt.Printf("cancelOrderResult:  %v \n", cancelOrderResult)

	//---- Get Close Order---------
	closedOrders, err := client.GetClosedOrders(ctypes.NewClosedOrdersQuery(testAccount1.String(), true).WithSymbol(tradeSymbol, nativeSymbol))
	assert.NoError(t, err)
	assert.True(t, len(closedOrders.Order) > 0)
	fmt.Printf("GetClosedOrders: %v \n", closedOrders)

	//----    Get tx      ---------
	tx, err := client.GetTx(openOrders.Order[0].TransactionHash)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	fmt.Printf("GetTx: %v\n", tx)

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
	time2.Sleep(2 * time2.Second)
	issueresult, err := client.GetTx(issue.Hash)
	assert.NoError(t, err)
	assert.True(t, issueresult.Code == txtype.CodeOk)

	//--- mint token -----------
	mint, err := client.MintToken(issue.Symbol, 100000000, true)
	assert.NoError(t, err)
	fmt.Printf("Mint token: %v\n", mint)

	//---- Submit Proposal ------
	time2.Sleep(2 * time2.Second)
	listTradingProposal, err := client.SubmitListPairProposal("New trading pair", msg.ListTradingPairParams{issue.Symbol, nativeSymbol, 1000000000, "my trade", time2.Now().Add(1 * time2.Hour)}, 200000000000, 20*time2.Second, true)
	assert.NoError(t, err)
	fmt.Printf("Submit list trading pair: %v\n", listTradingProposal)

	//---  check submit proposal success ---
	time2.Sleep(2 * time2.Second)
	submitPorposalStatus, err := client.GetTx(listTradingProposal.Hash)
	assert.NoError(t, err)
	assert.True(t, submitPorposalStatus.Code == txtype.CodeOk)

	//---- Vote Proposal  -------
	for _, m := range validatorMnemonics {
		time2.Sleep(1 * time2.Second)
		k, err := keys.NewMnemonicKeyManager(m)
		assert.NoError(t, err)
		client, err := sdk.NewDexClient(baeUrl, ctypes.TestNetwork, k)
		assert.NoError(t, err)
		vote, err := client.VoteProposal(listTradingProposal.ProposalId, msg.OptionYes, true)
		assert.NoError(t, err)
		fmt.Printf("Vote: %v\n", vote)
	}

	//--- List trade pair ------
	time2.Sleep(20 * time2.Second)
	l, err := client.ListPair(listTradingProposal.ProposalId, issue.Symbol, nativeSymbol, 1000000000, true)
	assert.NoError(t, err)
	fmt.Printf("List trading pair: %v\n", l)
}
