package e2e

import (
	"fmt"
	"testing"
	time2 "time"

	"github.com/stretchr/testify/assert"

	sdk "github.com/binance-chain/go-sdk"
	"github.com/binance-chain/go-sdk/api"
	"github.com/binance-chain/go-sdk/common/crypto"
	"github.com/binance-chain/go-sdk/common/crypto/secp256k1"
	"github.com/binance-chain/go-sdk/keys"
	"github.com/binance-chain/go-sdk/tx/txmsg"
)

// After bnbchain integration_test.sh has runned
func TestAllProcess(t *testing.T) {
	//----- Recover account ---------
	mnemonic := "mutual peasant quote avoid mandate current timber pipe prize wheel snow tattoo surround method weasel annual laundry evidence ecology cage medal perfect trial misery"
	keyManager, err := keys.NewMnemonicKeyManager(mnemonic)
	assert.NoError(t, err)
	testAccount1 := keyManager.GetAddr()
	_, testAccount2 := PrivAndAddr()

	//-----   Init sdk  -------------
	client, _ := sdk.NewBncClient("http://dex-api.fdgahl.cn", "chain-bnb", keyManager)
	nativeSymbol := txmsg.NativeToken
	//-----  Get account  -----------

	account, err := client.GetAccount(testAccount1.String())
	fmt.Println(err)
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.True(t, len(account.Balances) > 0)

	//----- Get Markets  ------------
	markets, err := client.GetMarkets(api.NewMarketsQuery().WithLimit(1))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(markets))
	tradeSymbol := markets[0].TradeAsset
	if markets[0].QuoteAsset != "BNB" {
		tradeSymbol = markets[0].QuoteAsset
	}

	//-----  Get Depth  -----------
	depth, err := client.GetDepth(api.NewDepthQuery(tradeSymbol, nativeSymbol))
	assert.NoError(t, err)
	assert.True(t, depth.Height > 0)

	//----- Get Kline
	kline, err := client.GetKlines(api.NewKlineQuery(tradeSymbol, nativeSymbol, "1h").WithLimit(1))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(kline))

	//-----  Get Ticker 24h  -----------
	ticker24h, err := client.GetTicker24h(api.NewTicker24hQuery().WithSymbol(tradeSymbol, nativeSymbol))
	assert.NoError(t, err)
	assert.True(t, len(ticker24h) > 0)

	//-----  Get Tokens  -----------
	tokens, err := client.GetTokens()
	assert.NoError(t, err)
	fmt.Printf("GetTokens: %v \n", tokens)

	//-----  Get Trades  -----------
	fmt.Println(testAccount1.String())
	trades, err := client.GetTrades(api.NewTradesQuery(testAccount1.String()).WithSymbol(tradeSymbol, nativeSymbol))
	assert.NoError(t, err)
	fmt.Printf("GetTrades: %v \n", trades)

	//-----  Get Time    -----------
	time, err := client.GetTime()
	assert.NoError(t, err)
	fmt.Printf("Get time: %v \n", time)

	//----- Create order -----------
	createOrderResult, err := client.CreateOrder(tradeSymbol, nativeSymbol, txmsg.OrderSide.BUY, 10000, 100000000, true)
	fmt.Println(err)
	assert.NoError(t, err)
	assert.True(t, true, createOrderResult.Ok)

	//---- Get Open Order ---------
	openOrders, err := client.GetOpenOrders(api.NewOpenOrdersQuery(testAccount1.String()))
	assert.NoError(t, err)
	assert.True(t, len(openOrders.Order) > 0)
	orderId := openOrders.Order[0].ID
	fmt.Printf("GetOpenOrders:  %v \n", openOrders)

	//---- Get Order    ------------
	order, err := client.GetOrder(orderId)
	assert.NoError(t, err)
	assert.Equal(t, api.CombineSymbol(tradeSymbol, nativeSymbol), order.Symbol)

	//---- Cancle Order  ---------
	cancleOrderResult, err := client.CancelOrder(tradeSymbol, nativeSymbol, orderId, orderId, true)
	assert.NoError(t, err)
	assert.True(t, cancleOrderResult.Ok)
	fmt.Printf("cancleOrderResult:  %v \n", cancleOrderResult)

	//---- Get Close Order---------
	closedOrders, err := client.GetClosedOrders(api.NewClosedOrdersQuery(testAccount1.String()).WithSymbol(tradeSymbol, nativeSymbol))
	assert.NoError(t, err)
	assert.True(t, len(closedOrders.Order) > 0)
	fmt.Printf("GetClosedOrders: %v \n", closedOrders)

	//----    Get tx      ---------
	tx, err := client.GetTx(openOrders.Order[0].TransactionHash)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	fmt.Printf("GetTx: %v\n", tx)

	//----   Send tx  -----------
	send, err := client.SendToken(testAccount2, nativeSymbol, 10000000000, true)
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
	freeze, err := client.FreezeToken(nativeSymbol, 100000000, true)
	assert.NoError(t, err)
	assert.True(t, freeze.Ok)
	fmt.Printf("freeze token: %v\n", freeze)

	//----   Unfreeze Token ---------
	unfreeze, err := client.UnfreezeToken(nativeSymbol, 100000000, true)
	assert.NoError(t, err)
	assert.True(t, unfreeze.Ok)
	fmt.Printf("Unfreeze token: %v\n", unfreeze)

	//----   issue token ---------
	issue, err := client.IssueToken("Client-Token", "sdk", 10000000000000000, true, true)
	assert.NoError(t, err)
	fmt.Printf("Issue token: %v\n", issue)

	//---  check issue success ---
	time2.Sleep(2 * time2.Second)
	issueresult, err := client.GetTx(issue.Hash)
	assert.NoError(t, err)
	assert.True(t, issueresult.Code == api.CodeOk)

	//--- mint token -----------
	mint,err:=client.MintToken(issue.Symbol,100000000, true)
	assert.NoError(t,err)
	fmt.Printf("Mint token: %v\n", mint)

	//---- Submit Proposal ------
	listTradingProposal, err := client.SubmitListPairProposal("New trading pair", txmsg.ListTradingPairParams{issue.Symbol, nativeSymbol, 1000000000, "my trade", time2.Now().Add(1 * time2.Hour)}, 200000000000, true)
	assert.NoError(t, err)
	fmt.Printf("Submit list trading pair: %v\n", listTradingProposal)

	//---  check submit proposal success ---
	//time2.Sleep(2 * time2.Second)
	//submitPorposalStatus, err := client.GetTx(listTradingProposal.Hash)
	//assert.NoError(t, err)
	//assert.True(t, submitPorposalStatus.Code == api.CodeOk)

	//---- Vote Proposal  -------
	//time2.Sleep(10 * time2.Second)
	//vote, err := client.VoteProposal(listTradingProposal.ProposalId, txmsg.OptionYes, true)
	//assert.NoError(t, err)
	//fmt.Printf("Vote: %v\n", vote)

	//---- Get new markets
	//time2.Sleep(1 * time2.Minute)
	//markets, err = client.GetMarkets(&api.MarketsQuery{Limit: 1, Offset: 0})
	//assert.NoError(t, err)
	//fmt.Printf("New markets: %v \n ", markets)

}

func PrivAndAddr() (crypto.PrivKey, txmsg.AccAddress) {
	priv := secp256k1.GenPrivKey()
	addr := txmsg.AccAddress(priv.PubKey().Address())
	return priv, addr
}
