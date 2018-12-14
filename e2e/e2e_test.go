package e2e

import (
	"fmt"
	"testing"
	time2 "time"

	"github.com/BiJie/bnc-go-sdk/sdk"
	"github.com/BiJie/bnc-go-sdk/sdk/api"
	"github.com/BiJie/bnc-go-sdk/sdk/keys"
	"github.com/BiJie/bnc-go-sdk/sdk/tx/txmsg"
	"github.com/stretchr/testify/assert"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// After bnbchain integration_test.sh has runned
func TestAllProcess(t *testing.T) {
	//----- Recover account ---------
	mnemonic := "lock globe panda armed mandate fabric couple dove climb step stove price recall decrease fire sail ring media enhance excite deny valid ceiling arm"
	keyManager, err := keys.NewMnemonicKeyManager(mnemonic)
	assert.NoError(t, err)
	testAccount1 := keyManager.GetAddr()
	_, testAccount2 := PrivAndAddr()


	//-----   Init sdk  -------------
	bnc, _ := sdk.NewBncSDK("http://dex-api.fdgahl.cn", "chain-bnb", keyManager)
	nativeSymbol := txmsg.NativeToken
	//-----  Get account  -----------

	account, err := bnc.GetAccount(testAccount1.String())
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.True(t, len(account.Balances) > 0)

	//----- Get Markets  ------------
	markets, err := bnc.GetMarkets(api.NewMarketsQuery().WithLimit(1))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(markets))
	tradeSymbol := markets[0].TradeAsset
	if markets[0].QuoteAsset != "BNB" {
		tradeSymbol = markets[0].QuoteAsset
	}

	//-----  Get Depth  -----------
	depth, err := bnc.GetDepth(api.NewDepthQuery(tradeSymbol, nativeSymbol))
	assert.NoError(t, err)
	assert.True(t, depth.Height > 0)

	//----- Get Kline
	kline, err := bnc.GetKlines(api.NewKlineQuery(tradeSymbol, nativeSymbol, "1h").WithLimit(1))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(kline))

	//-----  Get Ticker 24h  -----------
	ticker24h, err := bnc.GetTicker24h(api.NewTicker24hQuery().WithSymbol(tradeSymbol, nativeSymbol))
	assert.NoError(t, err)
	assert.True(t, len(ticker24h) > 0)

	//-----  Get Tokens  -----------
	tokens, err := bnc.GetTokens()
	assert.NoError(t, err)
	fmt.Printf("GetTokens: %v \n", tokens)

	//-----  Get Trades  -----------
	fmt.Println(testAccount1.String())
	trades, err := bnc.GetTrades(api.NewTradesQuery(testAccount1.String()).WithSymbol(tradeSymbol, nativeSymbol))
	assert.NoError(t, err)
	fmt.Printf("GetTrades: %v \n", trades)

	//-----  Get Time    -----------
	time, err := bnc.GetTime()
	assert.NoError(t, err)
	fmt.Printf("Get time: %v \n", time)

	//----- Create order -----------
	createOrderResult, err := bnc.CreateOrder(tradeSymbol, nativeSymbol, txmsg.OrderSide.BUY, 100000000, 100000000, true)
	assert.NoError(t, err)
	assert.True(t, true, createOrderResult.Ok)

	//---- Get Order    ------------
	order, err := bnc.GetOrder(createOrderResult.OrderId)
	assert.NoError(t, err)
	assert.Equal(t, createOrderResult.OrderId, order.ID)
	assert.Equal(t, api.CombineSymbol(tradeSymbol, nativeSymbol), order.Symbol)

	//---- Get Open Order ---------
	openOrders, err := bnc.GetOpenOrders(api.NewOpenOrdersQuery(testAccount1.String()))
	assert.NoError(t, err)
	assert.True(t, len(openOrders.Order) > 0)
	fmt.Printf("GetOpenOrders:  %v \n", openOrders)

	//---- Cancle Order  ---------
	cancleOrderResult, err := bnc.CancelOrder(tradeSymbol, nativeSymbol, createOrderResult.OrderId, createOrderResult.OrderId, true)
	assert.NoError(t, err)
	assert.True(t, cancleOrderResult.Ok)
	fmt.Printf("cancleOrderResult:  %v \n", cancleOrderResult)

	//---- Get Close Order---------
	closedOrders, err := bnc.GetClosedOrders(api.NewClosedOrdersQuery(testAccount1.String()).WithSymbol(tradeSymbol, nativeSymbol))
	assert.NoError(t, err)
	assert.True(t, len(closedOrders.Order) > 0)
	fmt.Printf("GetClosedOrders: %v \n", closedOrders)

	//----    Get tx      ---------
	tx, err := bnc.GetTx(openOrders.Order[0].TransactionHash)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	fmt.Printf("GetTx: %v\n", tx)

	//----   Send tx  -----------
	send, err := bnc.SendToken(testAccount2, nativeSymbol, 10000000000, true)
	assert.NoError(t, err)
	assert.True(t, send.Ok)
	fmt.Printf("Send token: %v\n", send)

	//---    Get test2 account-----
	newTestAccout2, err := bnc.GetAccount(testAccount2.String())
	assert.NoError(t, err)
	for _, c := range newTestAccout2.Balances {
		if c.Symbol == nativeSymbol {
			fmt.Printf("test account BNB: %s \n", c.Free)
		}
	}

	//----   Freeze Token ---------
	freeze, err := bnc.FreezeToken(nativeSymbol, 100000000, true)
	assert.NoError(t, err)
	assert.True(t, freeze.Ok)
	fmt.Printf("freeze token: %v\n", freeze)

	//----   Unfreeze Token ---------
	unfreeze, err := bnc.UnfreezeToken(nativeSymbol, 100000000, true)
	assert.NoError(t, err)
	assert.True(t, unfreeze.Ok)
	fmt.Printf("Unfreeze token: %v\n", unfreeze)

	//----   issue token ---------
	issue, err := bnc.IssueToken("SDK-Token", "sdk", 10000000000000000, true, false)
	assert.NoError(t, err)
	fmt.Printf("Issue token: %v\n", issue)

	//---  check issue success ---
	time2.Sleep(2 * time2.Second)
	issueresult, err := bnc.GetTx(issue.Hash)
	assert.NoError(t, err)
	assert.True(t, issueresult.Code == api.CodeOk)

	//---- Submit Proposal ------
	listTradingProposal, err := bnc.SubmitListPairProposal("New trading pair", txmsg.ListTradingPairParams{issue.Symbol, nativeSymbol, 1000000000, "my trade", time2.Now().Add(1 * time2.Hour)}, 200000000000, true)
	fmt.Println(err)
	assert.NoError(t, err)
	fmt.Printf("Submit list trading pair: %v\n", listTradingProposal)

	//---- Vote Proposal  -------
	time2.Sleep(10 * time2.Second)
	vote, err := bnc.VoteProposal(listTradingProposal.ProposalId, txmsg.OptionYes, true)
	assert.NoError(t, err)
	fmt.Printf("Vote: %v\n", vote)

	//---- Get new markets
	//time2.Sleep(1 * time2.Minute)
	//markets, err = bnc.GetMarkets(&api.MarketsQuery{Limit: 1, Offset: 0})
	//assert.NoError(t, err)
	//fmt.Printf("New markets: %v \n ", markets)

}

func PrivAndAddr() (tmcrypto.PrivKey, txmsg.AccAddress) {
	priv := secp256k1.GenPrivKey()
	addr := txmsg.AccAddress(priv.PubKey().Address())
	return priv, addr
}
