package rpc

import (
	"bytes"
	"fmt"
	types "github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
	"github.com/tendermint/go-amino"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"time"
)

var (
	SideChainStoreName          = "stake"
	ValidatorsKey               = []byte{0x21}
	DelegationKey               = []byte{0x31}
	SideChainStorePrefixByIdKey = []byte{0x51}
	RedelegationKey             = []byte{0x34}
	UnbondingDelegationKey      = []byte{0x32}
	DelegationTokenDemon		= "BNB"
)

func (c *HTTP) CreateSideChainValidator(delegation types.Coin, description msg.Description, commission types.CommissionMsg, sideChainId string, sideConsAddr []byte, sideFeeAddr []byte, syncType SyncType, options ...tx.Option) (*coretypes.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyManagerMissingError
	}

	if len(description.Moniker) == 0 {
		return nil, fmt.Errorf("Moniker in description is missing ")
	}

	if err := checkDelegationCoin(delegation); err != nil {
		return nil, err
	}

	valOpAddr := types.ValAddress(c.key.GetAddr())

	m := msg.NewCreateSideChainValidatorMsg(valOpAddr, delegation, description, commission, sideChainId, sideConsAddr, sideFeeAddr, )

	return c.broadcast(m, syncType, options...)
}

func (c *HTTP) EditSideChainValidatorMsg(sideChainId string, description msg.Description, commissionRate *types.Dec, sideConsAddr, sideFeeAddr []byte, syncType SyncType, options ...tx.Option) (*coretypes.ResultBroadcastTx, error)  {
	if c.key == nil {
		return nil, KeyManagerMissingError
	}

	valOpAddr := types.ValAddress(c.key.GetAddr())

	m := msg.NewEditSideChainValidatorMsg(sideChainId, valOpAddr, description, commissionRate, sideConsAddr, sideFeeAddr)

	return c.broadcast(m, syncType, options...)
}

func (c *HTTP) SideChainDelegate(sideChainId string, valAddr types.ValAddress, delegation types.Coin, syncType SyncType, options ...tx.Option) (*coretypes.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyManagerMissingError
	}

	if err := checkDelegationCoin(delegation); err != nil {
		return nil, err
	}

	delAddr := c.key.GetAddr()

	m := msg.NewSideChainDelegateMsg(sideChainId, delAddr, valAddr, delegation)

	return c.broadcast(m, syncType, options...)
}

func (c *HTTP) SideChainRedelegate(sideChainId string, valSrcAddr types.ValAddress, valDstAddr types.ValAddress, amount types.Coin, syncType SyncType, options ...tx.Option) (*coretypes.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyManagerMissingError
	}

	if bytes.Equal(valSrcAddr, valDstAddr) {
		return nil, fmt.Errorf("cannot redelegate to the same validator")
	}

	if err := checkDelegationCoin(amount); err != nil {
		return nil, err
	}

	delAddr := c.key.GetAddr()

	m := msg.NewSideChainRedelegateMsg(sideChainId, delAddr, valSrcAddr, valDstAddr, amount)

	return c.broadcast(m, syncType, options...)
}

func (c *HTTP) SideChainUnbond(sideChainId string, valAddr types.ValAddress, amount types.Coin, syncType SyncType, options ...tx.Option) (*coretypes.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyManagerMissingError
	}

	if err := checkDelegationCoin(amount); err != nil {
		return nil, fmt.Errorf("Unbond token must be %s ", DelegationTokenDemon)
	}

	delAddr := c.key.GetAddr()

	m := msg.NewSideChainUndelegateMsg(sideChainId, delAddr, valAddr, amount)

	return c.broadcast(m, syncType, options...)
}

//Query a validator
func (c *HTTP) QuerySideChainValidator(sideChainId string, valAddr types.ValAddress) (*types.Validator, error)  {
	storePrefix, err := c.getSideChainStorePrefixKey(sideChainId, SideChainStoreName)

	if err != nil {
		return nil, err
	}

	keyPrefix, err := c.QueryStore(storePrefix, SideChainStoreName)
	if err != nil {
		return nil, err
	}

	key := append(keyPrefix, getValidatorKey(valAddr)...)

	bz, err := c.QueryStore(key, SideChainStoreName)

	if err != nil {
		return nil, err
	}

	var validator types.Validator

	err = c.cdc.UnmarshalBinaryLengthPrefixed(bz, &validator)

	if err != nil {
		return nil, err
	}

	return &validator, nil
}

//Query for all validators
func (c *HTTP) QuerySideChainValidators(sideChainId string) ([]types.Validator, error) {
	storePrefix, err := c.getSideChainStorePrefixKey(sideChainId, SideChainStoreName)

	if err != nil {
		return nil, err
	}

	keyPrefix, err := c.QueryStore(storePrefix, SideChainStoreName)
	if err != nil {
		return nil, err
	}

	key := append(keyPrefix, ValidatorsKey...)

	resKVs, err := c.QueryStoreSubspace(key, SideChainStoreName)

	if err != nil {
		return nil, err
	}

	var validators []types.Validator
	for _, kv := range resKVs {
		var validator types.Validator
		c.cdc.MustUnmarshalBinaryLengthPrefixed(kv.Value, &validator)
		validators = append(validators, validator)
	}

	return validators, err
}

//Query a delegation based on address and validator address
func (c *HTTP) QuerySideChainDelegation(sideChainId string, delAddr types.AccAddress, valAddr types.ValAddress) (*types.Delegation, error) {
	storePrefix, err := c.getSideChainStorePrefixKey(sideChainId, SideChainStoreName)

	if err != nil {
		return nil, err
	}

	delegateKey := getDelegationKey(delAddr, valAddr)

	key := append(storePrefix, delegateKey...)
	res, err := c.QueryStore(key, SideChainStoreName)
	if err != nil {
		return nil, err
	}

	delegation, err := types.UnmarshalDelegation(c.cdc, delegateKey, res)

	return &delegation, nil
}

//Query all delegations made from one delegator
func (c *HTTP) QuerySideChainDelegations(sideChainId string, delAddr types.AccAddress) ([]types.Delegation, error) {
	storePrefix, err := c.getSideChainStorePrefixKey(sideChainId, SideChainStoreName)

	if err != nil {
		return nil, err
	}

	key := append(storePrefix, getDelegationsKey(delAddr)...)

	resKVS, err := c.QueryStoreSubspace(key, SideChainStoreName)
	if err != nil {
		return nil, err
	}

	var delegations []types.Delegation
	for _, kv := range resKVS {
		k := kv.Key[len(storePrefix):]
		delegation, err := types.UnmarshalDelegation(c.cdc, k, kv.Value)
		if err != nil {
			return nil, err
		}
		delegations = append(delegations, delegation)
	}

	return delegations, nil
}

//Query a redelegation record based on delegator and a source and destination validator address
func (c *HTTP) QuerySideChainRedelegation(sideChainId string, delAddr types.AccAddress, valSrcAddr types.ValAddress, valDstAddr types.ValAddress) (*types.Redelegation, error) {
	storePrefix, err := c.getSideChainStorePrefixKey(sideChainId, SideChainStoreName)
	if err != nil {
		return nil, err
	}

	redKey := getREDKey(delAddr, valSrcAddr, valDstAddr)
	key := append(storePrefix, redKey...)
	res, err := c.QueryStore(key, SideChainStoreName)
	if err != nil {
		return nil, err
	}

	if len(res) > 0 {
		result, err := types.UnmarshalRED(c.cdc, redKey, res)
		if err != nil {
			return nil, err
		}

		return &result, nil
	}

	return &types.Redelegation{}, fmt.Errorf("Query result is empty ")
}

//Query all redelegations records for one delegator
func (c *HTTP) QuerySideChainRedelegations(sideChainId string, delAddr types.AccAddress) ([]types.Redelegation, error) {
	storePrefix, err := c.getSideChainStorePrefixKey(sideChainId, SideChainStoreName)
	if err != nil {
		return nil, err
	}

	key := append(storePrefix, getREDsKey(delAddr)...)
	resKVs, err := c.QueryStoreSubspace(key, SideChainStoreName)
	if err != nil {
		return nil, err
	}

	var redels []types.Redelegation
	for _, kv := range resKVs {
		k := kv.Key[len(storePrefix):]
		red, err := types.UnmarshalRED(c.cdc, k, kv.Value)
		if err != nil {
			panic(err)
		}
		redels = append(redels, red)
	}

	if redels != nil && len(redels) > 0 {
		return redels, nil
	}else{
		return nil, fmt.Errorf("Query result is empty ")
	}
}

//Query an unbonding-delegation record based on delegator and validator address
func (c *HTTP) QuerySideChainUnbondingDelegation(sideChainId string, valAddr types.ValAddress, delAddr types.AccAddress) (*types.UnbondingDelegation, error) {
	storePrefix, err := c.getSideChainStorePrefixKey(sideChainId, SideChainStoreName)
	if err != nil {
		return nil, err
	}

	ubdKey := getUBDKey(delAddr, valAddr)
	key := append(storePrefix, ubdKey...)
	res, err := c.QueryStore(key, SideChainStoreName)
	if err != nil {
		return nil, err
	}

	ubd, err := unmarshalUBD(c.cdc, ubdKey, res)

	if err != nil {
		return nil, err
	}

	return &ubd, nil
}

//Query all unbonding-delegations records for one delegator
func (c *HTTP) QuerySideChainUnbondingDelegations(sideChainId string, delAddr types.AccAddress) ([]types.UnbondingDelegation, error) {
	storePrefix, err := c.getSideChainStorePrefixKey(sideChainId, SideChainStoreName)
	if err != nil {
		return nil, err
	}

	key := append(storePrefix, getUBDsKey(delAddr)...)

	resKVs, err := c.QueryStoreSubspace(key, SideChainStoreName)
	if err != nil {
		return nil, err
	}

	var ubds []types.UnbondingDelegation
	for _, kv := range resKVs {
		k := kv.Key[len(storePrefix):]
		ubd, err := unmarshalUBD(c.cdc, k, kv.Value)
		if err != nil{
			return nil, err
		}
		ubds = append(ubds, ubd)
	}

	return ubds, nil
}

func (c *HTTP) getSideChainStorePrefixKey(sideChainId string, storeName string) ([]byte, error) {
	key := append(SideChainStorePrefixByIdKey, []byte(sideChainId)...)
	result, err := c.QueryStore(key, storeName)

	if err != nil {
		return nil, err
	}else if len(result) == 0 {
		return nil, fmt.Errorf("Invalid side-chain-id %s ", sideChainId)
	}

	return result, nil
}

func getValidatorKey(operatorAddr types.ValAddress) []byte {
	return append(ValidatorsKey, operatorAddr.Bytes()...)
}

func getDelegationKey(delAddr types.AccAddress, valAddr types.ValAddress) []byte {
	return append(getDelegationsKey(delAddr), valAddr.Bytes()...)
}

func getDelegationsKey(delAddr types.AccAddress) []byte {
	return append(DelegationKey, delAddr.Bytes()...)
}

func getREDKey(delAddr types.AccAddress, valSrcAddr, valDstAddr types.ValAddress) []byte {
	key := make([]byte, 1+types.AddrLen*3)

	copy(key[0:types.AddrLen+1], getREDsKey(delAddr.Bytes()))
	copy(key[types.AddrLen+1:2*types.AddrLen+1], valSrcAddr.Bytes())
	copy(key[2*types.AddrLen+1:3*types.AddrLen+1], valDstAddr.Bytes())

	return key
}

// gets the prefix keyspace for redelegations from a delegator
func getREDsKey(delAddr types.AccAddress) []byte {
	return append(RedelegationKey, delAddr.Bytes()...)
}

// gets the key for an unbonding delegation by delegator and validator addr
// VALUE: stake/types.UnbondingDelegation
func getUBDKey(delAddr types.AccAddress, valAddr types.ValAddress) []byte {
	return append(
		getUBDsKey(delAddr.Bytes()),
		valAddr.Bytes()...)
}

// gets the prefix for all unbonding delegations from a delegator
func getUBDsKey(delAddr types.AccAddress) []byte {
	return append(UnbondingDelegationKey, delAddr.Bytes()...)
}

type ubdValue struct {
	CreationHeight int64
	MinTime        time.Time
	InitialBalance types.Coin
	Balance        types.Coin
}

// unmarshal a unbonding delegation from a store key and value
func unmarshalUBD(cdc *amino.Codec, key, value []byte) (ubd types.UnbondingDelegation, err error) {
	var storeValue ubdValue
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &storeValue)
	if err != nil {
		return
	}

	addrs := key[1:] // remove prefix bytes
	if len(addrs) != 2*types.AddrLen {
		err = fmt.Errorf("unexpected address length for this (address, validator) pair")
		return
	}
	delAddr := types.AccAddress(addrs[:types.AddrLen])
	valAddr := types.ValAddress(addrs[types.AddrLen:])

	return types.UnbondingDelegation{
		DelegatorAddr:  delAddr,
		ValidatorAddr:  valAddr,
		CreationHeight: storeValue.CreationHeight,
		MinTime:        storeValue.MinTime,
		InitialBalance: storeValue.InitialBalance,
		Balance:        storeValue.Balance,
	}, nil
}

func checkDelegationCoin(coin types.Coin) error {
	if coin.Denom != DelegationTokenDemon {
		return fmt.Errorf("Delegation token must be %s ", DelegationTokenDemon)
	}

	return nil
}