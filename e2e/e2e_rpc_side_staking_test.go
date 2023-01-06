package e2e

import (
	"encoding/hex"
	"sync"
	"testing"

	"github.com/bnb-chain/go-sdk/client/rpc"
	ctypes "github.com/bnb-chain/go-sdk/common/types"
	"github.com/bnb-chain/go-sdk/keys"
	"github.com/bnb-chain/go-sdk/types"
	"github.com/bnb-chain/go-sdk/types/msg"
	"github.com/bnb-chain/go-sdk/types/tx"
	"github.com/stretchr/testify/assert"
)

var (
	sNodeAddr           = "tcp://127.0.0.1:26657"
	sOnceClient         = sync.Once{}
	sTestClientInstance *rpc.HTTP

	jackAddress  = "bnb1lrzg56jhtkqu7fmca3394vdx00r7apx4gwvj6w"
	jackMnemonic = "orphan thing pelican flee spray sense sketch dutch opinion vessel fringe surround hurt theory hospital provide foil love stock vast shrug detail harbor pattern"

	roseAddress  = "bnb1rxnydtfjccaz2tck7wrentntdylrnnqzmvqvwn"
	roseMnemonic = "earth hamster near become enlist degree foil crucial weapon poverty mad purity chest lucky equal jazz pony either knee cloud drive badge jacket caught"

	markAddress  = "bnb1sh4cfzvcut9nywffs6gs5zkyt4pzeej6k84klt"
	markMnemonic = "depend water drink monitor earn praise permit autumn board cable impact wink wolf sting middle misery bridge stamp close very robust slam annual verify"

	chainId     = "test-chain-qUlw6e"
	valAddress  = "bva1lrzg56jhtkqu7fmca3394vdx00r7apx4gjdzy2"
	valAddress2 = "bva1rxnydtfjccaz2tck7wrentntdylrnnqzmspush"
)

func rpcClient() *rpc.HTTP {
	sOnceClient.Do(func() {
		sTestClientInstance = rpc.NewRPCClient(sNodeAddr, ctypes.ProdNetwork)
	})
	return sTestClientInstance
}

func getRpcClientWithKeyManager() *rpc.HTTP {
	c := rpcClient()
	ctypes.SetNetwork(ctypes.ProdNetwork)
	keyManager, _ := keys.NewMnemonicKeyManager(jackMnemonic)
	c.SetKeyManager(keyManager)
	return c
}

// FromHex returns the bytes represented by the hexadecimal string s.
// s may be prefixed with "0x".
func FromHex(s string) []byte {
	if has0xPrefix(s) {
		s = s[2:]
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return Hex2Bytes(s)
}

// Hex2Bytes returns the bytes represented by the hexadecimal string str.
func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}

// has0xPrefix validates str begins with '0x' or '0X'.
func has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

func TestCreateSideChainValidator(t *testing.T) {
	c := getRpcClientWithKeyManager()

	amount := ctypes.Coin{Denom: "BNB", Amount: 100000000}
	des := msg.Description{Moniker: "mchain"}

	rate, _ := ctypes.NewDecFromStr("1")
	maxRate, _ := ctypes.NewDecFromStr("1")
	maxChangeRate, _ := ctypes.NewDecFromStr("1")

	commissionMsg := ctypes.CommissionMsg{Rate: rate, MaxRate: maxRate, MaxChangeRate: maxChangeRate}

	sideChainId := types.RialtoNet
	sideConsAddr := FromHex("0x9fB29AAc15b9A4B7F17c3385939b007540f4d791")
	sideFeeAddr := FromHex("0xd1B22dCC24C55f4d728E7aaA5c9b5a22e1512C08")

	res, err := c.CreateSideChainValidator(amount, des, commissionMsg, sideChainId, sideConsAddr, sideFeeAddr, rpc.Sync, tx.WithChainID(chainId))

	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestEditSideChainValidator(t *testing.T) {
	c := getRpcClientWithKeyManager()

	des := msg.Description{Moniker: "mchain"}

	rate, _ := ctypes.NewDecFromStr("2")

	sideFeeAddr := FromHex("0xd1B22dCC24C55f4d728E7aaA5c9b5a22e1512C08")
	consAddr := FromHex("0xd1B22dCC24C55f4d728E7aaA5c9b5a22e1512C08")
	res, err := c.EditSideChainValidator(types.RialtoNet, des, &rate, sideFeeAddr, consAddr, rpc.Sync, tx.WithChainID(chainId))

	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestDelegate(t *testing.T) {
	c := getRpcClientWithKeyManager()

	valAddr, err := ctypes.ValAddressFromBech32(valAddress)

	assert.NoError(t, err)

	amount := ctypes.Coin{Denom: "BNB", Amount: 100000000}

	res, err := c.SideChainDelegate(types.RialtoNet, valAddr, amount, rpc.Sync, tx.WithChainID(chainId))

	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestRedelegate(t *testing.T) {
	c := getRpcClientWithKeyManager()

	srcValAddr, err := ctypes.ValAddressFromBech32(valAddress)

	assert.NoError(t, err)

	dstValAddr, err := ctypes.ValAddressFromBech32(valAddress2)

	assert.NoError(t, err)

	amount := ctypes.Coin{Denom: "BNB", Amount: 100000000}

	res, err := c.SideChainRedelegate(types.RialtoNet, srcValAddr, dstValAddr, amount, rpc.Sync, tx.WithChainID(chainId))

	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestUnbond(t *testing.T) {
	c := getRpcClientWithKeyManager()

	valAddr, err := ctypes.ValAddressFromBech32(valAddress)

	assert.NoError(t, err)

	amount := ctypes.Coin{Denom: "BNB", Amount: 100000000}

	res, err := c.SideChainUnbond(types.RialtoNet, valAddr, amount, rpc.Sync, tx.WithChainID(chainId))

	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestQuerySideChainValidator(t *testing.T) {
	c := getRpcClientWithKeyManager()

	valAddr, err := ctypes.ValAddressFromBech32(valAddress)

	assert.Nil(t, err)

	res, err := c.QuerySideChainValidator(types.RialtoNet, valAddr)

	if res == nil {
		assert.Equal(t, rpc.EmptyResultError, err)
	} else {
		assert.NotNil(t, res.OperatorAddr)
	}
}

func TestQuerySideChainTopValidators(t *testing.T) {
	c := getRpcClientWithKeyManager()
	_, err := c.QuerySideChainTopValidators(types.RialtoNet, 5)
	assert.NoError(t, err)
}

func TestQuerySideChainDelegation(t *testing.T) {
	c := getRpcClientWithKeyManager()

	delAddr, _ := ctypes.AccAddressFromBech32(jackAddress)
	valAddr, _ := ctypes.ValAddressFromBech32(valAddress)

	res, err := c.QuerySideChainDelegation(types.RialtoNet, delAddr, valAddr)

	if res == nil {
		assert.Equal(t, rpc.EmptyResultError, err)
	} else {
		assert.NotNil(t, res.ValidatorAddr)
	}
}

func TestQuerySideChainDelegations(t *testing.T) {
	c := rpcClient()

	delAddr, _ := ctypes.AccAddressFromBech32(jackAddress)

	_, err := c.QuerySideChainDelegations(types.RialtoNet, delAddr)
	assert.Nil(t, err)
}

func TestQuerySideChainRelegation(t *testing.T) {
	c := rpcClient()

	delAddr, _ := ctypes.AccAddressFromBech32(jackAddress)
	valSrcAddr, _ := ctypes.ValAddressFromBech32(valAddress)
	valDstAddr, _ := ctypes.ValAddressFromBech32(valAddress2)

	res, err := c.QuerySideChainRedelegation(types.RialtoNet, delAddr, valSrcAddr, valDstAddr)

	if res == nil {
		assert.Equal(t, rpc.EmptyResultError, err)
	} else {
		assert.NotNil(t, res.DelegatorAddr)
	}
}

func TestQuerySideChainRelegations(t *testing.T) {
	c := rpcClient()
	delAddr, err := ctypes.AccAddressFromBech32(jackAddress)
	assert.Nil(t, err)
	_, err = c.QuerySideChainRedelegations(types.RialtoNet, delAddr)
	assert.Nil(t, err)
}

func TestQuerySideChainUnbondingDelegation(t *testing.T) {
	c := rpcClient()

	delAddr, _ := ctypes.AccAddressFromBech32(jackAddress)
	valAddr, _ := ctypes.ValAddressFromBech32(valAddress)

	res, err := c.QuerySideChainUnbondingDelegation(types.RialtoNet, valAddr, delAddr)

	if res == nil {
		assert.Equal(t, rpc.EmptyResultError, err)
	} else {
		assert.NotNil(t, res.DelegatorAddr)
	}
}

func TestQuerySideChainUnbondingDelegations(t *testing.T) {
	c := rpcClient()
	delAddr, _ := ctypes.AccAddressFromBech32(jackAddress)
	_, err := c.QuerySideChainUnbondingDelegations(types.RialtoNet, delAddr)
	assert.Nil(t, err)
}

func TestGetSideChainUnBondingDelegationsByValidator(t *testing.T) {
	c := getRpcClientWithKeyManager()
	valAddr, _ := ctypes.ValAddressFromBech32(jackAddress)
	_, err := c.GetSideChainUnBondingDelegationsByValidator(types.RialtoNet, valAddr)
	assert.Nil(t, err)
}

func TestGetSideChainRedelegationsByValidator(t *testing.T) {
	c := getRpcClientWithKeyManager()
	valAddr, _ := ctypes.ValAddressFromBech32(jackAddress)
	_, err := c.GetSideChainRedelegationsByValidator(types.RialtoNet, valAddr)
	assert.Nil(t, err)
}

func TestGetSideChainId(t *testing.T) {
	c := getRpcClientWithKeyManager()
	_, err := c.GetSideChainPool(types.RialtoNet)
	assert.Nil(t, err)
}

func TestGetSideChainAllValidatorsCount(t *testing.T) {
	c := getRpcClientWithKeyManager()
	_, err := c.GetSideChainAllValidatorsCount(types.RialtoNet, false)
	assert.Nil(t, err)
}
