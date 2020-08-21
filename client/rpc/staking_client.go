package rpc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
	"github.com/tendermint/go-amino"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
)

var (
	StakeStoreKey               = "stake"
	StakeScStoreKey             = "sc"
	SideChainStorePrefixByIdKey = []byte{0x01}
	ValidatorsKey               = []byte{0x21}
	DelegationKey               = []byte{0x31}
	RedelegationKey             = []byte{0x34}
	UnbondingDelegationKey      = []byte{0x32}
	PoolKey                     = []byte{0x01}
)

type StakingClient interface {
	CreateSideChainValidator(delegation types.Coin, description msg.Description, commission types.CommissionMsg, sideChainId string, sideConsAddr []byte, sideFeeAddr []byte, syncType SyncType, options ...tx.Option) (*coretypes.ResultBroadcastTx, error)
	EditSideChainValidator(sideChainId string, description msg.Description, commissionRate *types.Dec, sideFeeAddr []byte, syncType SyncType, options ...tx.Option) (*coretypes.ResultBroadcastTx, error)
	SideChainDelegate(sideChainId string, valAddr types.ValAddress, delegation types.Coin, syncType SyncType, options ...tx.Option) (*coretypes.ResultBroadcastTx, error)
	SideChainRedelegate(sideChainId string, valSrcAddr types.ValAddress, valDstAddr types.ValAddress, amount types.Coin, syncType SyncType, options ...tx.Option) (*coretypes.ResultBroadcastTx, error)
	SideChainUnbond(sideChainId string, valAddr types.ValAddress, amount types.Coin, syncType SyncType, options ...tx.Option) (*coretypes.ResultBroadcastTx, error)
	SideChainUnjail(sideChainId string, valAddr types.ValAddress, syncType SyncType, options ...tx.Option) (*coretypes.ResultBroadcastTx, error)

	QuerySideChainValidator(sideChainId string, valAddr types.ValAddress) (*types.Validator, error)
	QuerySideChainTopValidators(sideChainId string, top int) ([]types.Validator, error)
	QuerySideChainDelegation(sideChainId string, delAddr types.AccAddress, valAddr types.ValAddress) (*types.DelegationResponse, error)
	QuerySideChainDelegations(sideChainId string, delAddr types.AccAddress) ([]types.DelegationResponse, error)
	QuerySideChainRedelegation(sideChainId string, delAddr types.AccAddress, valSrcAddr types.ValAddress, valDstAddr types.ValAddress) (*types.Redelegation, error)
	QuerySideChainRedelegations(sideChainId string, delAddr types.AccAddress) ([]types.Redelegation, error)
	QuerySideChainUnbondingDelegation(sideChainId string, valAddr types.ValAddress, delAddr types.AccAddress) (*types.UnbondingDelegation, error)
	QuerySideChainUnbondingDelegations(sideChainId string, delAddr types.AccAddress) ([]types.UnbondingDelegation, error)
	GetSideChainUnBondingDelegationsByValidator(sideChainId string, valAddr types.ValAddress) ([]types.UnbondingDelegation, error)
	GetSideChainRedelegationsByValidator(sideChainId string, valAddr types.ValAddress) ([]types.Redelegation, error)
	GetSideChainPool(sideChainId string) (*types.Pool, error)
	GetSideChainAllValidatorsCount(sideChainId string, jailInvolved bool) (int, error)
}

type bechValidator struct {
	FeeAddr      types.AccAddress `json:"fee_addr"`                   // the bech32 address for fee collection
	OperatorAddr types.ValAddress `json:"operator_address"`           // the bech32 address of the validator's operator
	ConsPubKey   string           `json:"consensus_pubkey,omitempty"` // the bech32 consensus public key of the validator
	Jailed       bool             `json:"jailed"`                     // has the validator been jailed from bonded status?

	Status          types.BondStatus `json:"status"`           // validator status (bonded/unbonding/unbonded)
	Tokens          types.Dec        `json:"tokens"`           // delegated tokens (incl. self-delegation)
	DelegatorShares types.Dec        `json:"delegator_shares"` // total shares issued to a validator's delegators

	Description        types.Description `json:"description"`           // description terms for the validator
	BondHeight         int64             `json:"bond_height"`           // earliest height as a bonded validator
	BondIntraTxCounter int16             `json:"bond_intra_tx_counter"` // block-local tx index of validator change

	UnbondingHeight  int64     `json:"unbonding_height"` // if unbonding, height at which this validator has begun unbonding
	UnbondingMinTime time.Time `json:"unbonding_time"`   // if unbonding, min time for the validator to complete unbonding

	Commission types.Commission `json:"commission"` // commission parameters

	DistributionAddr types.AccAddress `json:"distribution_addr,omitempty"` // the address receives rewards from the side address, and distribute rewards to delegators. It's auto generated
	SideChainId      string           `json:"side_chain_id,omitempty"`     // side chain id to distinguish different side chains
	SideConsAddr     string           `json:"side_cons_addr,omitempty"`    // consensus address of the side chain validator, this replaces the `ConsPubKey`
	SideFeeAddr      string           `json:"side_fee_addr,omitempty"`     // fee address on the side chain
}

func (c *HTTP) CreateSideChainValidator(delegation types.Coin, description msg.Description, commission types.CommissionMsg,
	sideChainId string, sideConsAddr []byte, sideFeeAddr []byte, syncType SyncType, options ...tx.Option) (*coretypes.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyMissingError
	}

	valOpAddr := types.ValAddress(c.key.GetAddr())

	m := msg.NewCreateSideChainValidatorMsg(valOpAddr, delegation, description, commission, sideChainId, sideConsAddr, sideFeeAddr)

	return c.Broadcast(m, syncType, options...)
}

func (c *HTTP) EditSideChainValidator(sideChainId string, description msg.Description, commissionRate *types.Dec,
	sideFeeAddr []byte, syncType SyncType, options ...tx.Option) (*coretypes.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyMissingError
	}

	valOpAddr := types.ValAddress(c.key.GetAddr())

	m := msg.NewEditSideChainValidatorMsg(sideChainId, valOpAddr, description, commissionRate, sideFeeAddr)

	return c.Broadcast(m, syncType, options...)
}

func (c *HTTP) SideChainDelegate(sideChainId string, valAddr types.ValAddress, delegation types.Coin, syncType SyncType,
	options ...tx.Option) (*coretypes.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyMissingError
	}

	delAddr := c.key.GetAddr()

	m := msg.NewSideChainDelegateMsg(sideChainId, delAddr, valAddr, delegation)

	return c.Broadcast(m, syncType, options...)
}

func (c *HTTP) SideChainRedelegate(sideChainId string, valSrcAddr types.ValAddress, valDstAddr types.ValAddress, amount types.Coin,
	syncType SyncType, options ...tx.Option) (*coretypes.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyMissingError
	}

	if bytes.Equal(valSrcAddr, valDstAddr) {
		return nil, fmt.Errorf("cannot redelegate to the same validator")
	}

	delAddr := c.key.GetAddr()

	m := msg.NewSideChainRedelegateMsg(sideChainId, delAddr, valSrcAddr, valDstAddr, amount)

	return c.Broadcast(m, syncType, options...)
}

func (c *HTTP) SideChainUnbond(sideChainId string, valAddr types.ValAddress, amount types.Coin, syncType SyncType,
	options ...tx.Option) (*coretypes.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyMissingError
	}

	delAddr := c.key.GetAddr()

	m := msg.NewSideChainUndelegateMsg(sideChainId, delAddr, valAddr, amount)

	return c.Broadcast(m, syncType, options...)
}

func (c *HTTP) SideChainUnjail(sideChainId string, valAddr types.ValAddress, syncType SyncType, options ...tx.Option) (*coretypes.ResultBroadcastTx, error) {
	if c.key == nil {
		return nil, KeyMissingError
	}

	m := msg.NewMsgSideChainUnjail(valAddr, sideChainId)

	return c.Broadcast(m, syncType, options...)
}

//Query a validator
func (c *HTTP) QuerySideChainValidator(sideChainId string, valAddr types.ValAddress) (*types.Validator, error) {
	storePrefix, err := c.getSideChainStorePrefixKey(sideChainId)

	if err != nil {
		return nil, err
	}

	key := append(storePrefix, getValidatorKey(valAddr)...)

	bz, err := c.QueryStore(key, StakeStoreKey)

	if err != nil {
		return nil, err
	}

	if len(bz) == 0 {
		return nil, EmptyResultError
	}

	var validator types.Validator

	err = c.cdc.UnmarshalBinaryLengthPrefixed(bz, &validator)

	if err != nil {
		return nil, err
	}

	return &validator, nil
}

func (c *HTTP) QuerySideChainTopValidators(sideChainId string, top int) ([]types.Validator, error) {
	if top > 50 || top < 1 {
		return nil, fmt.Errorf("top must be between 1 and 50")
	}

	params := types.QueryTopValidatorsParams{
		BaseParams: types.NewBaseParams(sideChainId),
		Top:        top,
	}

	bz, err := json.Marshal(params)

	if err != nil {
		return nil, err
	}

	res, err := c.QueryWithData("custom/stake/topValidators", bz)
	if err != nil {
		return nil, err
	}

	var validators = make([]types.Validator, 0)

	if len(res) == 0 {
		return validators, nil
	}

	var bvs []bechValidator
	if err = c.cdc.UnmarshalJSON(res, &bvs); err != nil {
		return nil, err
	}

	for _, v := range bvs {
		validator := types.Validator{
			FeeAddr:            v.FeeAddr,
			OperatorAddr:       v.OperatorAddr,
			ConsPubKey:         v.ConsPubKey,
			Jailed:             v.Jailed,
			Status:             v.Status,
			Tokens:             v.Tokens,
			DelegatorShares:    v.DelegatorShares,
			Description:        v.Description,
			BondHeight:         v.BondHeight,
			BondIntraTxCounter: v.BondIntraTxCounter,
			UnbondingHeight:    v.UnbondingHeight,
			UnbondingMinTime:   v.UnbondingMinTime,
			Commission:         v.Commission,
		}

		if len(v.SideChainId) != 0 {
			validator.DistributionAddr = v.DistributionAddr
			validator.SideChainId = v.SideChainId
			if sideConsAddr, err := decodeSideChainAddress(v.SideConsAddr); err != nil {
				return nil, err
			} else {
				validator.SideConsAddr = sideConsAddr
			}
			if sideFeeAddr, err := decodeSideChainAddress(v.SideFeeAddr); err != nil {
				return nil, err
			} else {
				validator.SideFeeAddr = sideFeeAddr
			}
		}

		validators = append(validators, validator)
	}

	return validators, nil
}

//Query a delegation based on address and validator address
func (c *HTTP) QuerySideChainDelegation(sideChainId string, delAddr types.AccAddress, valAddr types.ValAddress) (*types.DelegationResponse, error) {
	params := types.QueryBondsParams{
		BaseParams:    types.NewBaseParams(sideChainId),
		DelegatorAddr: delAddr,
		ValidatorAddr: valAddr,
	}

	bz, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	response, err := c.QueryWithData("custom/stake/delegation", bz)
	if err != nil {
		return nil, err
	} else if len(response) == 0 {
		return nil, fmt.Errorf("No delegation found ")
	}

	var delResponse types.DelegationResponse
	if err := c.cdc.UnmarshalJSON(response, &delResponse); err != nil {
		return nil, err
	}

	return &delResponse, nil
}

//Query all delegations made from one delegator
func (c *HTTP) QuerySideChainDelegations(sideChainId string, delAddr types.AccAddress) ([]types.DelegationResponse, error) {
	params := types.QueryDelegatorParams{
		BaseParams:    types.NewBaseParams(sideChainId),
		DelegatorAddr: delAddr,
	}

	var delegationResponses []types.DelegationResponse
	delegationResponses = make([]types.DelegationResponse, 0)

	bz, err := json.Marshal(params)
	if err != nil {
		return delegationResponses, err
	}

	response, err := c.QueryWithData("custom/stake/delegatorDelegations", bz)
	if err != nil {
		return delegationResponses, err
	} else if len(response) == 0 {
		return delegationResponses, fmt.Errorf("No delegation found with delegator-addr %s ", delAddr)
	}

	if err := c.cdc.UnmarshalJSON(response, &delegationResponses); err != nil {
		return delegationResponses, err
	}

	return delegationResponses, nil
}

//Query a redelegation record based on delegator and a source and destination validator address
func (c *HTTP) QuerySideChainRedelegation(sideChainId string, delAddr types.AccAddress, valSrcAddr types.ValAddress,
	valDstAddr types.ValAddress) (*types.Redelegation, error) {
	storePrefix, err := c.getSideChainStorePrefixKey(sideChainId)
	if err != nil {
		return nil, err
	}

	redKey := getREDKey(delAddr, valSrcAddr, valDstAddr)
	key := append(storePrefix, redKey...)
	res, err := c.QueryStore(key, StakeStoreKey)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, EmptyResultError
	}

	result, err := types.UnmarshalRED(c.cdc, redKey, res)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

//Query all redelegations records for one delegator
func (c *HTTP) QuerySideChainRedelegations(sideChainId string, delAddr types.AccAddress) ([]types.Redelegation, error) {
	storePrefix, err := c.getSideChainStorePrefixKey(sideChainId)
	if err != nil {
		return nil, err
	}

	key := append(storePrefix, getREDsKey(delAddr)...)
	resKVs, err := c.QueryStoreSubspace(key, StakeStoreKey)
	if err != nil {
		return nil, err
	}

	var redels = make([]types.Redelegation, 0)

	if len(resKVs) == 0 {
		return redels, nil
	}

	for _, kv := range resKVs {
		k := kv.Key[len(storePrefix):]
		red, err := types.UnmarshalRED(c.cdc, k, kv.Value)
		if err != nil {
			return redels, err
		}
		redels = append(redels, red)
	}

	return redels, nil
}

//Query an unbonding-delegation record based on delegator and validator address
func (c *HTTP) QuerySideChainUnbondingDelegation(sideChainId string, valAddr types.ValAddress, delAddr types.AccAddress) (*types.UnbondingDelegation, error) {
	storePrefix, err := c.getSideChainStorePrefixKey(sideChainId)
	if err != nil {
		return nil, err
	}

	ubdKey := getUBDKey(delAddr, valAddr)
	key := append(storePrefix, ubdKey...)
	res, err := c.QueryStore(key, StakeStoreKey)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, EmptyResultError
	}

	ubd, err := unmarshalUBD(c.cdc, ubdKey, res)

	if err != nil {
		return nil, err
	}

	return &ubd, nil
}

//Query all unbonding-delegations records for one delegator
func (c *HTTP) QuerySideChainUnbondingDelegations(sideChainId string, delAddr types.AccAddress) ([]types.UnbondingDelegation, error) {
	storePrefix, err := c.getSideChainStorePrefixKey(sideChainId)
	if err != nil {
		return nil, err
	}

	key := append(storePrefix, getUBDsKey(delAddr)...)

	resKVs, err := c.QueryStoreSubspace(key, StakeStoreKey)
	if err != nil {
		return nil, err
	}

	var ubds = make([]types.UnbondingDelegation, 0)

	if len(resKVs) == 0 {
		return ubds, nil
	}

	for _, kv := range resKVs {
		k := kv.Key[len(storePrefix):]
		ubd, err := unmarshalUBD(c.cdc, k, kv.Value)
		if err != nil {
			return nil, err
		}
		ubds = append(ubds, ubd)
	}

	return ubds, nil
}

func (c *HTTP) GetSideChainUnBondingDelegationsByValidator(sideChainId string, valAddr types.ValAddress) ([]types.UnbondingDelegation, error) {
	params := types.QueryValidatorParams{
		BaseParams:    types.NewBaseParams(sideChainId),
		ValidatorAddr: valAddr,
	}

	bz, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	response, err := c.QueryWithData("custom/stake/validatorUnbondingDelegations", bz)
	if err != nil {
		return nil, err
	}

	var ubds = make([]types.UnbondingDelegation, 0)

	if len(response) == 0 {
		return ubds, nil
	}

	if err = c.cdc.UnmarshalJSON(response, &ubds); err != nil {
		return nil, err
	}

	return ubds, nil
}

func (c *HTTP) GetSideChainRedelegationsByValidator(sideChainId string, valAddr types.ValAddress) ([]types.Redelegation, error) {
	params := types.QueryValidatorParams{
		BaseParams:    types.NewBaseParams(sideChainId),
		ValidatorAddr: valAddr,
	}

	bz, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	response, err := c.QueryWithData("custom/stake/validatorRedelegations", bz)
	if err != nil {
		return nil, err
	}

	var reds = make([]types.Redelegation, 0)

	if len(response) == 0 {
		return reds, nil
	}

	if err = c.cdc.UnmarshalJSON(response, &reds); err != nil {
		return nil, err
	}

	return reds, nil
}

func (c *HTTP) GetSideChainPool(sideChainId string) (*types.Pool, error) {
	storePrefix, err := c.getSideChainStorePrefixKey(sideChainId)

	if err != nil {
		return nil, err
	}

	key := append(storePrefix, PoolKey...)
	res, err := c.QueryStore(key, StakeStoreKey)

	if len(res) == 0 {
		zeroDec, err := types.NewDecFromStr("0")
		if err != nil {
			return nil, err
		}
		return &types.Pool{
			LooseTokens:  zeroDec,
			BondedTokens: zeroDec,
		}, nil
	}

	var pool types.Pool
	err = c.cdc.UnmarshalBinaryLengthPrefixed(res, &pool)
	if err != nil {
		return nil, err
	}

	return &pool, nil
}

func (c *HTTP) GetSideChainAllValidatorsCount(sideChainId string, jailInvolved bool) (int, error) {
	params := types.NewBaseParams(sideChainId)

	bz, err := json.Marshal(params)
	if err != nil {
		return 0, err
	}

	path := "custom/stake/allUnJailValidatorsCount"
	if jailInvolved {
		path = "custom/stake/allValidatorsCount"
	}
	response, err := c.QueryWithData(path, bz)

	if err != nil {
		return 0, err
	}

	count := strings.ReplaceAll(string(response), "\"", "")

	return strconv.Atoi(count)
}

func (c *HTTP) getSideChainStorePrefixKey(sideChainId string) ([]byte, error) {
	key := append(SideChainStorePrefixByIdKey, []byte(sideChainId)...)
	result, err := c.QueryStore(key, StakeScStoreKey)

	if err != nil {
		return nil, err
	} else if len(result) == 0 {
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

func decodeSideChainAddress(addr string) ([]byte, error) {
	if strings.HasPrefix(addr, "0x") {
		return hex.DecodeString(addr[2:])
	} else {
		return hex.DecodeString(addr)
	}
}
