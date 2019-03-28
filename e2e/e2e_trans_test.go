package e2e

import (
	"flag"
	"fmt"
	"os"
	"testing"
	time2 "time"

	"github.com/stretchr/testify/assert"

	sdk "github.com/binance-chain/go-sdk/client"
	"github.com/binance-chain/go-sdk/client/query"
	"github.com/binance-chain/go-sdk/common"
	"github.com/binance-chain/go-sdk/keys"
	"github.com/binance-chain/go-sdk/types"
	"github.com/binance-chain/go-sdk/types/msg"
	tx2 "github.com/binance-chain/go-sdk/types/tx"
)

var (
	mnemonic = flag.String("mnemonic", "test mnemonic", "mnemonic of a test account")

	client      sdk.DexClient
	tradeSymbol string
)

func TestMain(t *testing.M) {
	flag.Parse()
	keyManager, err := keys.NewMnemonicKeyManager(*mnemonic)
	if err != nil {
		panic(fmt.Sprintf("failed to construct keymanager %v", err))
	}
	client, err := sdk.NewDexClient("testnet-dex.binance.org", types.TestNetwork, keyManager)
	if err != nil {
		panic(fmt.Sprintf("failed to construct client %v", err))
	}
	markets, err := client.GetMarkets(query.NewMarketsQuery().WithLimit(1))
	if err != nil {
		panic(fmt.Sprintf("failed to get markets %v", err))
	}
	if len(markets) == 0 {
		panic("the chain do not have any markets")
	}
	tradeSymbol = markets[0].TradeAsset
	if markets[0].QuoteAsset != "BNB" {
		tradeSymbol = markets[0].QuoteAsset
	}
	os.Exit(t.Run())
}

func TestGetAccount(t *testing.T) {
	testAccount := client.GetKeyManager().GetAddr()
	account, err := client.GetAccount(testAccount.String())
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.True(t, len(account.Balances) > 0)
}

func TestGetMarkets(t *testing.T) {
	markets, err := client.GetMarkets(query.NewMarketsQuery().WithLimit(1))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(markets))
}

func TestGetDepth(t *testing.T) {
	depth, err := client.GetDepth(query.NewDepthQuery(tradeSymbol, msg.NativeToken))
	assert.NoError(t, err)
	assert.True(t, depth.Height > 0)
}

func TestGetKline(t *testing.T) {
	kline, err := client.GetKlines(query.NewKlineQuery(tradeSymbol, msg.NativeToken, "1h").WithLimit(1))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(kline))
}

func TestGetTicker(t *testing.T) {
	ticker24h, err := client.GetTicker24h(query.NewTicker24hQuery().WithSymbol(tradeSymbol, msg.NativeToken))
	assert.NoError(t, err)
	assert.True(t, len(ticker24h) > 0)
}

func TestGetTokens(t *testing.T) {
	tokens, err := client.GetTokens()
	assert.NoError(t, err)
	fmt.Printf("GetTokens: %v \n", tokens)
}

func TestGetTrades(t *testing.T) {
	account := client.GetKeyManager().GetAddr().String()
	trades, err := client.GetTrades(query.NewTradesQuery(account, true).WithSymbol(tradeSymbol, msg.NativeToken))
	assert.NoError(t, err)
	fmt.Printf("GetTrades: %v \n", trades)
}

func TestGetTime(t *testing.T) {
	time, err := client.GetTime()
	assert.NoError(t, err)
	fmt.Printf("Get time: %v \n", time)
}

func TestTransProcess(t *testing.T) {
	//----- Recover account ---------
	testAccount1 := client.GetKeyManager().GetAddr()
	testKeyManager2, _ := keys.NewKeyManager()
	testAccount2 := testKeyManager2.GetAddr()
	testKeyManager3, _ := keys.NewKeyManager()
	testAccount3 := testKeyManager3.GetAddr()

	//----- Create order -----------
	createOrderResult, err := client.CreateOrder(tradeSymbol, msg.NativeToken, msg.OrderSide.SELL, 30000000000, 100000000, true)
	assert.NoError(t, err)
	assert.True(t, true, createOrderResult.Ok)

	//---- Get Open Order ---------
	openOrders, err := client.GetOpenOrders(query.NewOpenOrdersQuery(testAccount1.String(), true))
	assert.NoError(t, err)
	assert.True(t, len(openOrders.Order) > 0)
	orderId := openOrders.Order[0].ID
	fmt.Printf("GetOpenOrders:  %v \n", openOrders)

	//---- Get Order    ------------
	order, err := client.GetOrder(orderId)
	assert.NoError(t, err)
	assert.Equal(t, common.CombineSymbol(tradeSymbol, msg.NativeToken), order.Symbol)

	//---- Cancle Order  ---------
	time2.Sleep(2 * time2.Second)
	cancleOrderResult, err := client.CancelOrder(tradeSymbol, msg.NativeToken, orderId, true)
	assert.NoError(t, err)
	assert.True(t, cancleOrderResult.Ok)
	fmt.Printf("cancleOrderResult:  %v \n", cancleOrderResult)

	//---- Get Close Order---------
	closedOrders, err := client.GetClosedOrders(query.NewClosedOrdersQuery(testAccount1.String(), true).WithSymbol(tradeSymbol, msg.NativeToken))
	assert.NoError(t, err)
	assert.True(t, len(closedOrders.Order) > 0)
	fmt.Printf("GetClosedOrders: %v \n", closedOrders)

	//----    Get tx      ---------
	tx, err := client.GetTx(openOrders.Order[0].TransactionHash)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	fmt.Printf("GetTx: %v\n", tx)

	//----   Send tx  -----------
	send, err := client.SendToken([]msg.Transfer{{testAccount2, []types.Coin{{msg.NativeToken, 100000000}}}, {testAccount3, []types.Coin{{msg.NativeToken, 100000000}}}}, true)
	assert.NoError(t, err)
	assert.True(t, send.Ok)
	fmt.Printf("Send token: %v\n", send)

	//---    Get test2 account-----
	newTestAccout2, err := client.GetAccount(testAccount2.String())
	assert.NoError(t, err)
	for _, c := range newTestAccout2.Balances {
		if c.Symbol == msg.NativeToken {
			fmt.Printf("test account BNB: %s \n", c.Free)
		}
	}

	//----   Freeze Token ---------
	freeze, err := client.FreezeToken(msg.NativeToken, 100, true)
	assert.NoError(t, err)
	assert.True(t, freeze.Ok)
	fmt.Printf("freeze token: %v\n", freeze)

	//----   Unfreeze Token ---------
	unfreeze, err := client.UnfreezeToken(msg.NativeToken, 100, true)
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
	assert.True(t, issueresult.Code == tx2.CodeOk)

	//--- mint token -----------
	mint, err := client.MintToken(issue.Symbol, 100000000, true)
	assert.NoError(t, err)
	fmt.Printf("Mint token: %v\n", mint)

	//---- Submit Proposal ------
	time2.Sleep(2 * time2.Second)
	listTradingProposal, err := client.SubmitListPairProposal("New trading pair", msg.ListTradingPairParams{issue.Symbol, msg.NativeToken, 1000000000, "my trade", time2.Now().Add(1 * time2.Hour)}, 200000000000, true)
	assert.NoError(t, err)
	fmt.Printf("Submit list trading pair: %v\n", listTradingProposal)

	//---  check submit proposal success ---
	time2.Sleep(2 * time2.Second)
	submitPorposalStatus, err := client.GetTx(listTradingProposal.Hash)
	assert.NoError(t, err)
	assert.True(t, submitPorposalStatus.Code == tx2.CodeOk)

	//---- Vote Proposal  -------
	//time2.Sleep(2 * time2.Second)
	//vote, err := client.VoteProposal(listTradingProposal.ProposalId, msg.OptionYes, true)
	//assert.NoError(t, err)
	//fmt.Printf("Vote: %v\n", vote)

	//--- List trade pair ------
	//l, err := client.ListPair(listTradingProposal.ProposalId, issue.Symbol, msg.NativeToken, 1000000000, true)
	//assert.NoError(t, err)
	//fmt.Printf("List trading pair: %v\n", l)

}
