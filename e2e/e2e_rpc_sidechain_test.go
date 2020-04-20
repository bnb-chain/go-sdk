package e2e

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/binance-chain/go-sdk/client/rpc"
	ctypes "github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/keys"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

var (
	sNodeAddr           = "tcp://127.0.0.1:26657"
	sOnceClient         = sync.Once{}
	sTestClientInstance *rpc.HTTP

	jackAddress        = "bnb1pn9l8daaqyhk3jp2lzngl78tvtxuh75xg3hyjf"
	jackMnemonic       = "emotion issue garment picture track thank deny drastic artwork area moral saddle team honey diagram oil unveil train tongue service unlock ivory glove program"

	roseAddress  	   = "bnb108287rc5vq8xkmwvmwg74yp2yrpxrz5tz9fg9e"
	roseMnemonic 	   = "swallow swap dolphin clinic expire swap service cruel armor engage exchange time garage silver cook possible nothing ribbon merry sausage rack term iron solution"

	markAddress        = "bnb1f6um5kaxu598kux66fvq835cxq8l866af0fl3r"
	markMnemonic	   = "crane sport camp tenant broom family load rifle coconut seminar off axis release rival anchor echo clump secret live fat heavy cereal humor mass"

	chainId 		   = "test-chain-72DXjv"
	valAddress 		   = "bva1pn9l8daaqyhk3jp2lzngl78tvtxuh75xgdk5vd"
	valAddress2 	   = "bva1pn9l8daaqyhk3jp2lzngl78tvtxuh75xgdk5vd"

	sideChainId 	   = "bsc"
)

func rpcClient() *rpc.HTTP {
	sOnceClient.Do(func() {
		sTestClientInstance = rpc.NewRPCClient(sNodeAddr, ctypes.ProdNetwork)
	})
	return sTestClientInstance
}

func getRpcClientWithKeyManager() *rpc.HTTP {
	c := rpcClient()
	ctypes.Network = ctypes.ProdNetwork
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

func TestCreateSideChainValidator(t *testing.T)  {
	c := getRpcClientWithKeyManager()

	amount := ctypes.Coin{Denom: "BNB", Amount: 100000000}
	des := msg.Description{Moniker: "mchain"}

	rate, _ := ctypes.NewDecFromStr("1")
	maxRate, _ := ctypes.NewDecFromStr("1")
	maxChangeRate, _ := ctypes.NewDecFromStr("1")

	commissionMsg := ctypes.CommissionMsg{Rate: rate, MaxRate: maxRate, MaxChangeRate: maxChangeRate}

	sideChainId := "bsc"
	sideConsAddr := FromHex("0x9fB29AAc15b9A4B7F17c3385939b007540f4d791")
	sideFeeAddr := FromHex("0x9fB29AAc15b9A4B7F17c3385939b007540f4d791")

	res, err := c.CreateSideChainValidator(amount, des, commissionMsg, sideChainId, sideConsAddr, sideFeeAddr, rpc.Sync, tx.WithChainID(chainId))

	if err != nil {
		fmt.Println(err)
		return
	}

	bz, _ := json.Marshal(res)
	fmt.Println(string(bz))
}

func TestEditSideChainValidator(t *testing.T)  {
	c := getRpcClientWithKeyManager()

	des := msg.Description{Moniker: "mchain"}

	rate, _ := ctypes.NewDecFromStr("2")

	sideConsAddr := FromHex("0xd1B22dCC24C55f4d728E7aaA5c9b5a22e1512C08")
	sideFeeAddr := FromHex("0xd1B22dCC24C55f4d728E7aaA5c9b5a22e1512C08")

	res, err := c.EditSideChainValidatorMsg(sideChainId, des, &rate, sideConsAddr, sideFeeAddr, rpc.Sync, tx.WithChainID(chainId))

	if err != nil {
		fmt.Println(err)
		return
	}

	bz, _ := json.Marshal(res)
	fmt.Println(string(bz))
}

func TestDelegate(t *testing.T)  {
	c := getRpcClientWithKeyManager()

	valAddr, err := ctypes.ValAddressFromBech32(valAddress)

	if err != nil {
		fmt.Println(err)
		return
	}

	amount := ctypes.Coin{Denom: "BNB", Amount: 100000000}

	res, err := c.SideChainDelegate(sideChainId, valAddr, amount, rpc.Sync, tx.WithChainID(chainId))

	if err != nil {
		fmt.Println(err)
		return
	}

	bz, _ := json.Marshal(res)
	fmt.Println(string(bz))
}

func TestRedelegate(t *testing.T)  {
	c := getRpcClientWithKeyManager()

	srcValAddr, err := ctypes.ValAddressFromBech32(valAddress)

	if err != nil {
		fmt.Println(err)
		return
	}

	dstValAddr, err := ctypes.ValAddressFromBech32(valAddress2)

	if err != nil {
		fmt.Println(err)
		return
	}

	amount := ctypes.Coin{Denom: "BNB", Amount: 100000000}

	res, err := c.SideChainRedelegate(sideChainId, srcValAddr, dstValAddr, amount, rpc.Sync, tx.WithChainID(chainId))

	if err != nil {
		fmt.Println(err)
		return
	}

	bz, _ := json.Marshal(res)
	fmt.Println(string(bz))
}

func TestUnbond(t *testing.T)  {
	c := getRpcClientWithKeyManager()

	valAddr, err := ctypes.ValAddressFromBech32(valAddress)

	if err != nil {
		fmt.Println(err)
		return
	}

	amount := ctypes.Coin{Denom: "BNB", Amount: 100000000}

	res, err := c.SideChainUnbond(sideChainId, valAddr, amount, rpc.Sync, tx.WithChainID(chainId))

	if err != nil {
		fmt.Println(err)
		return
	}

	bz, _ := json.Marshal(res)
	fmt.Println(string(bz))
}

func TestQuerySideChainValidator(t *testing.T) {
	c := getRpcClientWithKeyManager()

	valAddr, _ := ctypes.ValAddressFromBech32("bva1pn9l8daaqyhk3jp2lzngl78tvtxuh75xgdk5vd")

	res, err := c.QuerySideChainValidator(sideChainId, valAddr)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(validatorHumanReadableString(*res))
}

func TestQuerySideChainValidators(t *testing.T)  {
	c := getRpcClientWithKeyManager()

	res, err := c.QuerySideChainValidators(sideChainId)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range res {
		fmt.Println(validatorHumanReadableString(v))
	}
}

func TestQuerySideChainDelegation(t *testing.T)  {
	c := getRpcClientWithKeyManager()

	delAddr, _ := ctypes.AccAddressFromBech32(jackAddress)
	valAddr, _ := ctypes.ValAddressFromBech32(valAddress)

	res, err := c.QuerySideChainDelegation(sideChainId, delAddr, valAddr)

	assert.Nil(t, err)
	assert.NotNil(t, res)

	fmt.Println(res)
}

func TestQuerySideChainDelegations(t *testing.T) {
	c := rpcClient()

	delAddr, _ := ctypes.AccAddressFromBech32(jackAddress)

	res, err := c.QuerySideChainDelegations(sideChainId, delAddr)

	assert.Nil(t, err)
	assert.NotNil(t, res)

	fmt.Println(res)
}

func TestQuerySideChainRelegation(t *testing.T)  {
	c := rpcClient()

	delAddr, _ := ctypes.AccAddressFromBech32(jackAddress)
	valSrcAddr, _ := ctypes.ValAddressFromBech32(valAddress)
	valDstAddr, _ := ctypes.ValAddressFromBech32(valAddress2)

	res, err := c.QuerySideChainRedelegation(sideChainId, delAddr, valSrcAddr, valDstAddr)

	assert.Nil(t, err)
	assert.NotNil(t, res)

	fmt.Println(res)
}

func TestQuerySideChainRelegations(t *testing.T)  {
	c := rpcClient()

	delAddr, _ := ctypes.AccAddressFromBech32(jackAddress)

	res, err := c.QuerySideChainRedelegations(sideChainId, delAddr)

	assert.Nil(t, err)
	assert.NotNil(t, res)

	fmt.Println(res)
}

func TestQuerySideChainUnbondingDelegation(t *testing.T) {
	c := rpcClient()

	delAddr, _ := ctypes.AccAddressFromBech32(jackAddress)
	valAddr, _ := ctypes.ValAddressFromBech32(valAddress)

	res, err := c.QuerySideChainUnbondingDelegation(sideChainId, valAddr, delAddr)

	assert.Nil(t, err)
	assert.NotNil(t, res)

	fmt.Println(res)
}

func TestQuerySideChainUnbondingDelegations(t *testing.T)  {
	c := rpcClient()

	delAddr, _ := ctypes.AccAddressFromBech32(jackAddress)

	res, err := c.QuerySideChainUnbondingDelegations(sideChainId, delAddr)

	assert.Nil(t, err)
	assert.NotNil(t, res)

	fmt.Println(res)
}

func validatorHumanReadableString(v ctypes.Validator) (string, error) {
	resp := "Validator \n"
	resp += fmt.Sprintf("Fee Address: %s\n", v.FeeAddr)
	resp += fmt.Sprintf("Operator Address: %s\n", v.OperatorAddr)
	resp += fmt.Sprintf("Validator Consensus Pubkey: %s\n", v.ConsPubKey)
	resp += fmt.Sprintf("Jailed: %v\n", v.Jailed)
	resp += fmt.Sprintf("Status: %s\n", ctypes.BondStatusToString(v.Status))
	resp += fmt.Sprintf("Tokens: %s\n", v.Tokens)
	resp += fmt.Sprintf("Delegator Shares: %s\n", v.DelegatorShares)
	resp += fmt.Sprintf("Description: %s\n", v.Description)
	resp += fmt.Sprintf("Bond Height: %d\n", v.BondHeight)
	resp += fmt.Sprintf("Unbonding Height: %d\n", v.UnbondingHeight)
	resp += fmt.Sprintf("Minimum Unbonding Time: %v\n", v.UnbondingMinTime)
	resp += fmt.Sprintf("Commission: {%s}\n", v.Commission)
	if len(v.SideChainId) != 0 {
		resp += fmt.Sprintf("Distribution Addr: %s\n", v.DistributionAddr)
		resp += fmt.Sprintf("Side Chain Id: %s\n", v.SideChainId)
		resp += fmt.Sprintf("Consensus Addr on Side Chain: %s\n", v.SideConsAddr)
		resp += fmt.Sprintf("Fee Addr on Side Chain: %s\n", v.SideFeeAddr)
	}

	return resp, nil
}