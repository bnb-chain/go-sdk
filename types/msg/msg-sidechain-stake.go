package msg

import (
	"bytes"
	"fmt"

	"github.com/binance-chain/go-sdk/common/types"
)

const (
	TypeCreateSideChainValidator = "side_create_validator"
	TypeEditSideChainValidator   = "side_edit_validator"
	TypeSideChainDelegate        = "side_delegate"
	TypeSideChainRedelegate      = "side_redelegate"
	TypeSideChainUndelegate      = "side_undelegate"

	SideChainStakeMsgRoute = "stake"
	SideChainAddrLen       = 20

	MinDelegationAmount = 1e8
)

type CreateSideChainValidatorMsg struct {
	Description   Description         `json:"description"`
	Commission    types.CommissionMsg `json:"commission"`
	DelegatorAddr types.AccAddress    `json:"delegator_address"`
	ValidatorAddr types.ValAddress    `json:"validator_address"`
	Delegation    types.Coin          `json:"delegation"`
	SideChainId   string              `json:"side_chain_id"`
	SideConsAddr  []byte              `json:"side_cons_addr"`
	SideFeeAddr   []byte              `json:"side_fee_addr"`
}

func NewCreateSideChainValidatorMsg(valAddr types.ValAddress, delegation types.Coin,
	description Description, commission types.CommissionMsg, sideChainId string, sideConsAddr []byte, sideFeeAddr []byte) CreateSideChainValidatorMsg {
	return NewMsgCreateSideChainValidatorOnBehalfOf(types.AccAddress(valAddr), valAddr, delegation, description, commission, sideChainId, sideConsAddr, sideFeeAddr)
}

func NewMsgCreateSideChainValidatorOnBehalfOf(delegatorAddr types.AccAddress, valAddr types.ValAddress, delegation types.Coin,
	description Description, commission types.CommissionMsg, sideChainId string, sideConsAddr []byte, sideFeeAddr []byte) CreateSideChainValidatorMsg {
	return CreateSideChainValidatorMsg{
		Description:   description,
		Commission:    commission,
		DelegatorAddr: delegatorAddr,
		ValidatorAddr: valAddr,
		Delegation:    delegation,
		SideChainId:   sideChainId,
		SideConsAddr:  sideConsAddr,
		SideFeeAddr:   sideFeeAddr,
	}
}

func (msg CreateSideChainValidatorMsg) Route() string { return SideChainStakeMsgRoute }

func (msg CreateSideChainValidatorMsg) Type() string { return TypeCreateSideChainValidator }

func (msg CreateSideChainValidatorMsg) GetSigners() []types.AccAddress {
	addrs := []types.AccAddress{msg.DelegatorAddr}

	if !bytes.Equal(msg.DelegatorAddr.Bytes(), msg.ValidatorAddr.Bytes()) {
		addrs = append(addrs, types.AccAddress(msg.ValidatorAddr))
	}
	return addrs
}

func (msg CreateSideChainValidatorMsg) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return MustSortJSON(bz)
}

func (msg CreateSideChainValidatorMsg) GetInvolvedAddresses() []types.AccAddress {
	return msg.GetSigners()
}

func (msg CreateSideChainValidatorMsg) ValidateBasic() error {
	//self-delegator address length is 20
	if len(msg.DelegatorAddr) != types.AddrLen {
		return fmt.Errorf("Expected delegator address length is %d, actual length is %d ", types.AddrLen, len(msg.DelegatorAddr))
	}

	//validator operator address length is 20
	if len(msg.ValidatorAddr) != types.AddrLen {
		return fmt.Errorf("Expected validator address length is %d, actual length is %d ", types.AddrLen, len(msg.ValidatorAddr))
	}

	//description check
	if msg.Description == (Description{}) {
		return fmt.Errorf("description must be included")
	}
	if _, err := msg.Description.EnsureLength(); err != nil {
		return err
	}

	//commission check
	commission := types.NewCommission(msg.Commission.Rate, msg.Commission.MaxRate, msg.Commission.MaxChangeRate)
	if err := commission.Validate(); err != nil {
		return err
	}

	//side chain id length 1-20
	if len(msg.SideChainId) == 0 || len(msg.SideChainId) > MaxSideChainIdLength {
		return fmt.Errorf("side chain id must be included and max length is %d bytes", MaxSideChainIdLength)
	}

	//sideConsAddr length should between 16 - 64
	if err := checkSideChainAddr("SideConsAddr", msg.SideConsAddr); err != nil {
		return err
	}

	//sideFeeAddr length should between 16 - 64
	if err := checkSideChainAddr("SideFeeAddr", msg.SideFeeAddr); err != nil {
		return err
	}

	return nil
}

func checkSideChainAddr(addrName string, addr []byte) error {
	if len(addr) != SideChainAddrLen {
		return fmt.Errorf("Expected %s length is %d, got %d ", addrName, SideChainAddrLen, len(addr))
	}

	return nil
}

//----------------------------------------------------------------------------

type EditSideChainValidatorMsg struct {
	Description   Description      `json:"description"`
	ValidatorAddr types.ValAddress `json:"address"`

	CommissionRate *types.Dec `json:"commission_rate"`

	SideChainId string `json:"side_chain_id"`
	SideFeeAddr []byte `json:"side_fee_addr"`
}

func NewEditSideChainValidatorMsg(sideChainId string, validatorAddr types.ValAddress, description Description, commissionRate *types.Dec, sideFeeAddr []byte) EditSideChainValidatorMsg {
	return EditSideChainValidatorMsg{
		Description:    description,
		ValidatorAddr:  validatorAddr,
		CommissionRate: commissionRate,
		SideChainId:    sideChainId,
		SideFeeAddr:    sideFeeAddr,
	}
}

func (msg EditSideChainValidatorMsg) Route() string { return SideChainStakeMsgRoute }

func (msg EditSideChainValidatorMsg) Type() string { return TypeEditSideChainValidator }

func (msg EditSideChainValidatorMsg) GetSigners() []types.AccAddress {
	return []types.AccAddress{types.AccAddress(msg.ValidatorAddr)}
}

func (msg EditSideChainValidatorMsg) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return MustSortJSON(bz)
}

func (msg EditSideChainValidatorMsg) ValidateBasic() error {
	//validator operator address length is 20
	if len(msg.ValidatorAddr) != types.AddrLen {
		return fmt.Errorf("Expected validator address length is %d, actual length is %d ", types.AddrLen, len(msg.ValidatorAddr))
	}

	//description check
	if msg.Description == (Description{}) {
		return fmt.Errorf("description must be included. if you do not want to edit the description, assign each field to [do-not-modify]")
	}
	if _, err := msg.Description.EnsureLength(); err != nil {
		return err
	}

	//commission rate is between 0 and 1
	if msg.CommissionRate != nil {
		if msg.CommissionRate.GT(types.OneDec()) || msg.CommissionRate.LT(types.ZeroDec()) {
			return fmt.Errorf("commission rate must be between 0 and 1 (inclusive)")
		}
	}

	//side chain id length 1 - 20
	if len(msg.SideChainId) == 0 || len(msg.SideChainId) > MaxSideChainIdLength {
		return fmt.Errorf("side chain id must be included and max length is %d bytes", MaxSideChainIdLength)
	}

	//sideFeeAddr length should between 16 - 64
	if err := checkSideChainAddr("SideFeeAddr", msg.SideFeeAddr); err != nil {
		return err
	}

	return nil
}

func (msg EditSideChainValidatorMsg) GetInvolvedAddresses() []types.AccAddress {
	return msg.GetSigners()
}

//----------------------------------------------------------------------------

type SideChainDelegateMsg struct {
	DelegatorAddr types.AccAddress `json:"delegator_addr"`
	ValidatorAddr types.ValAddress `json:"validator_addr"`
	Delegation    types.Coin       `json:"delegation"`

	SideChainId string `json:"side_chain_id"`
}

func NewSideChainDelegateMsg(sideChainId string, delAddr types.AccAddress, valAddr types.ValAddress, delegation types.Coin) SideChainDelegateMsg {
	return SideChainDelegateMsg{
		DelegatorAddr: delAddr,
		ValidatorAddr: valAddr,
		Delegation:    delegation,
		SideChainId:   sideChainId,
	}
}

func (msg SideChainDelegateMsg) Route() string { return SideChainStakeMsgRoute }

func (msg SideChainDelegateMsg) Type() string { return TypeSideChainDelegate }

func (msg SideChainDelegateMsg) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.DelegatorAddr}
}

func (msg SideChainDelegateMsg) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return MustSortJSON(bz)
}

func (msg SideChainDelegateMsg) ValidateBasic() error {
	//delegator address is 20
	if len(msg.DelegatorAddr) != types.AddrLen {
		return fmt.Errorf("Expected delegator address length is %d, actual length is %d ", types.AddrLen, len(msg.DelegatorAddr))
	}

	//validator operator address length is 20
	if len(msg.ValidatorAddr) != types.AddrLen {
		return fmt.Errorf("Expected validator address length is %d, actual length is %d ", types.AddrLen, len(msg.ValidatorAddr))
	}

	//delegation amount should be greater than or equal to 1e8
	if msg.Delegation.Amount < MinDelegationAmount {
		return fmt.Errorf("delegation must not be less than %d ", msg.Delegation.Amount)
	}

	//side chain id length 1 - 20
	if len(msg.SideChainId) == 0 || len(msg.SideChainId) > MaxSideChainIdLength {
		return fmt.Errorf("side chain id must be included and max length is %d bytes", MaxSideChainIdLength)
	}
	return nil
}

func (msg SideChainDelegateMsg) GetInvolvedAddresses() []types.AccAddress {
	return []types.AccAddress{msg.DelegatorAddr, types.AccAddress(msg.ValidatorAddr)}
}

//----------------------------------------------------------------------------

type SideChainRedelegateMsg struct {
	DelegatorAddr    types.AccAddress `json:"delegator_addr"`
	ValidatorSrcAddr types.ValAddress `json:"validator_src_addr"`
	ValidatorDstAddr types.ValAddress `json:"validator_dst_addr"`
	Amount           types.Coin       `json:"amount"`
	SideChainId      string           `json:"side_chain_id"`
}

func NewSideChainRedelegateMsg(sideChainId string, delegatorAddr types.AccAddress, valSrcAddr types.ValAddress, valDstAddr types.ValAddress, amount types.Coin) SideChainRedelegateMsg {
	return SideChainRedelegateMsg{
		DelegatorAddr:    delegatorAddr,
		ValidatorSrcAddr: valSrcAddr,
		ValidatorDstAddr: valDstAddr,
		Amount:           amount,
		SideChainId:      sideChainId,
	}
}

func (msg SideChainRedelegateMsg) Route() string { return SideChainStakeMsgRoute }

func (msg SideChainRedelegateMsg) Type() string { return TypeSideChainRedelegate }

func (msg SideChainRedelegateMsg) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.DelegatorAddr}
}

func (msg SideChainRedelegateMsg) GetSignBytes() []byte {
	b := MsgCdc.MustMarshalJSON(msg)
	return MustSortJSON(b)
}

func (msg SideChainRedelegateMsg) GetInvolvedAddresses() []types.AccAddress {
	return []types.AccAddress{msg.DelegatorAddr, types.AccAddress(msg.ValidatorSrcAddr), types.AccAddress(msg.DelegatorAddr)}
}

func (msg SideChainRedelegateMsg) ValidateBasic() error {
	//delegator address length is 20
	if len(msg.DelegatorAddr) != types.AddrLen {
		return fmt.Errorf("Expected delegator address length is %d, actual length is %d ", types.AddrLen, len(msg.DelegatorAddr))
	}

	//source validator address length is 20
	if len(msg.ValidatorSrcAddr) != types.AddrLen {
		return fmt.Errorf("Expected ValidatorSrcAddr length is %d, actual length is %d ", types.AddrLen, len(msg.ValidatorSrcAddr))
	}

	//destination validator address length is 20
	if len(msg.ValidatorDstAddr) != types.AddrLen {
		return fmt.Errorf("Expected ValidatorDstAddr length is %d, actual length is %d ", types.AddrLen, len(msg.ValidatorDstAddr))
	}

	//source address and destination address are not the same one
	if bytes.Equal(msg.ValidatorSrcAddr, msg.ValidatorDstAddr) {
		return fmt.Errorf("cannot redelegate to the same validator")
	}

	//redelegation amount is greater than or equal to 1e8
	if msg.Amount.Amount < MinDelegationAmount {
		return fmt.Errorf("redelegation amount must not be less than %f ", MinDelegationAmount)
	}

	//side chain id length 1 - 20
	if len(msg.SideChainId) == 0 || len(msg.SideChainId) > MaxSideChainIdLength {
		return fmt.Errorf("side chain id must be included and max length is %d bytes", MaxSideChainIdLength)
	}
	return nil
}

//----------------------------------------------------------------------------

type SideChainUndelegateMsg struct {
	DelegatorAddr types.AccAddress `json:"delegator_addr"`
	ValidatorAddr types.ValAddress `json:"validator_addr"`
	Amount        types.Coin       `json:"amount"`
	SideChainId   string           `json:"side_chain_id"`
}

func NewSideChainUndelegateMsg(sideChainId string, delegatorAddr types.AccAddress, valAddr types.ValAddress, amount types.Coin) SideChainUndelegateMsg {
	return SideChainUndelegateMsg{
		DelegatorAddr: delegatorAddr,
		ValidatorAddr: valAddr,
		Amount:        amount,
		SideChainId:   sideChainId,
	}
}

func (msg SideChainUndelegateMsg) Route() string { return SideChainStakeMsgRoute }

func (msg SideChainUndelegateMsg) Type() string { return TypeSideChainUndelegate }

func (msg SideChainUndelegateMsg) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.DelegatorAddr}
}

func (msg SideChainUndelegateMsg) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return MustSortJSON(bz)
}

func (msg SideChainUndelegateMsg) ValidateBasic() error {
	//delegator address length is 20
	if len(msg.DelegatorAddr) != types.AddrLen {
		return fmt.Errorf("Expected delegator address length is %d, actual length is %d ", types.AddrLen, len(msg.DelegatorAddr))
	}

	//validator operator address length is 20
	if len(msg.ValidatorAddr) != types.AddrLen {
		return fmt.Errorf("Expected validator address length is %d, actual length is %d ", types.AddrLen, len(msg.ValidatorAddr))
	}

	//undelegate amount must be positive
	if msg.Amount.Amount <= 0 {
		return fmt.Errorf("undelegation amount must be positive")
	}

	//side chain id length 1 - 20
	if len(msg.SideChainId) == 0 || len(msg.SideChainId) > MaxSideChainIdLength {
		return fmt.Errorf("side chain id must be included and max length is %d bytes", MaxSideChainIdLength)
	}
	return nil
}

func (msg SideChainUndelegateMsg) GetInvolvedAddresses() []types.AccAddress {
	return []types.AccAddress{msg.DelegatorAddr, types.AccAddress(msg.ValidatorAddr)}
}
