package msg

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	sdk "github.com/binance-chain/go-sdk/common/types"
)

const (
	RouteBridge = "bridge"

	BindMsgType        = "crossBind"
	TransferOutMsgType = "crossTransferOut"
)

// EthereumAddress defines a standard ethereum address
type EthereumAddress gethCommon.Address

// NewEthereumAddress is a constructor function for EthereumAddress
func NewEthereumAddress(address string) EthereumAddress {
	return EthereumAddress(gethCommon.HexToAddress(address))
}

func (ethAddr EthereumAddress) IsEmpty() bool {
	addrValue := big.NewInt(0)
	addrValue.SetBytes(ethAddr[:])

	return addrValue.Cmp(big.NewInt(0)) == 0
}

// Route should return the name of the module
func (ethAddr EthereumAddress) String() string {
	return gethCommon.Address(ethAddr).String()
}

// MarshalJSON marshals the ethereum address to JSON
func (ethAddr EthereumAddress) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%v\"", ethAddr.String())), nil
}

// UnmarshalJSON unmarshals an ethereum address
func (ethAddr *EthereumAddress) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(reflect.TypeOf(gethCommon.Address{}), input, ethAddr[:])
}

type TransferInClaim struct {
	ContractAddress   EthereumAddress   `json:"contract_address"`
	RefundAddresses   []EthereumAddress `json:"refund_addresses"`
	ReceiverAddresses []sdk.AccAddress  `json:"receiver_addresses"`
	Amounts           []int64           `json:"amounts"`
	Symbol            string            `json:"symbol"`
	RelayFee          sdk.Coin          `json:"relay_fee"`
	ExpireTime        int64             `json:"expire_time"`
}

type TransferOutRefundClaim struct {
	RefundAddress sdk.AccAddress `json:"refund_address"`
	Amount        sdk.Coin       `json:"amount"`
	RefundReason  RefundReason   `json:"refund_reason"`
}

type UpdateBindClaim struct {
	Status          BindStatus      `json:"status"`
	Symbol          string          `json:"symbol"`
	ContractAddress EthereumAddress `json:"contract_address"`
}

type SkipSequenceClaim struct {
	ClaimType ClaimType `json:"claim_type"`
	Sequence  int64     `json:"sequence"`
}

type RefundReason uint16

const (
	UnboundToken        RefundReason = 1
	Timeout             RefundReason = 2
	InsufficientBalance RefundReason = 3
	Unknown             RefundReason = 4
)

type BindMsg struct {
	From             sdk.AccAddress  `json:"from"`
	Symbol           string          `json:"symbol"`
	Amount           int64           `json:"amount"`
	ContractAddress  EthereumAddress `json:"contract_address"`
	ContractDecimals int8            `json:"contract_decimals"`
	ExpireTime       int64           `json:"expire_time"`
}

func NewBindMsg(from sdk.AccAddress, symbol string, amount int64, contractAddress EthereumAddress, contractDecimals int8, expireTime int64) BindMsg {
	return BindMsg{
		From:             from,
		Amount:           amount,
		Symbol:           symbol,
		ContractAddress:  contractAddress,
		ContractDecimals: contractDecimals,
		ExpireTime:       expireTime,
	}
}

func (msg BindMsg) Route() string { return RouteBridge }
func (msg BindMsg) Type() string  { return BindMsgType }
func (msg BindMsg) String() string {
	return fmt.Sprintf("Bind{%v#%s#%d$%s#%d#%d}", msg.From, msg.Symbol, msg.Amount, msg.ContractAddress.String(), msg.ContractDecimals, msg.ExpireTime)
}
func (msg BindMsg) GetInvolvedAddresses() []sdk.AccAddress { return msg.GetSigners() }
func (msg BindMsg) GetSigners() []sdk.AccAddress           { return []sdk.AccAddress{msg.From} }

func (msg BindMsg) ValidateBasic() error {
	if len(msg.From) != sdk.AddrLen {
		return fmt.Errorf("address length should be %d", sdk.AddrLen)
	}

	if len(msg.Symbol) == 0 {
		return fmt.Errorf("symbol should not be empty")
	}

	if msg.Amount <= 0 {
		return fmt.Errorf("amount should be larger than 0")
	}

	if msg.ContractAddress.IsEmpty() {
		return fmt.Errorf("contract address should not be empty")
	}

	if msg.ContractDecimals < 0 {
		return fmt.Errorf("decimal should be no less than 0")
	}

	if msg.ExpireTime <= 0 {
		return fmt.Errorf("expire time should be larger than 0")
	}

	return nil
}

func (msg BindMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

type TransferOutMsg struct {
	From       sdk.AccAddress  `json:"from"`
	To         EthereumAddress `json:"to"`
	Amount     sdk.Coin        `json:"amount"`
	ExpireTime int64           `json:"expire_time"`
}

func NewTransferOutMsg(from sdk.AccAddress, to EthereumAddress, amount sdk.Coin, expireTime int64) TransferOutMsg {
	return TransferOutMsg{
		From:       from,
		To:         to,
		Amount:     amount,
		ExpireTime: expireTime,
	}
}

func (msg TransferOutMsg) Route() string { return RouteBridge }
func (msg TransferOutMsg) Type() string  { return TransferOutMsgType }
func (msg TransferOutMsg) String() string {
	return fmt.Sprintf("TransferOut{%v#%s#%s#%d}", msg.From, msg.To.String(), msg.Amount.String(), msg.ExpireTime)
}
func (msg TransferOutMsg) GetInvolvedAddresses() []sdk.AccAddress { return msg.GetSigners() }
func (msg TransferOutMsg) GetSigners() []sdk.AccAddress           { return []sdk.AccAddress{msg.From} }
func (msg TransferOutMsg) ValidateBasic() error {
	if len(msg.From) != sdk.AddrLen {
		return fmt.Errorf("address length should be %d", sdk.AddrLen)
	}

	if msg.To.IsEmpty() {
		return fmt.Errorf("to address should not be empty")
	}

	if !msg.Amount.IsPositive() {
		return fmt.Errorf("amount should be positive")
	}

	if msg.ExpireTime <= 0 {
		return fmt.Errorf("expire time should be larger than 0")
	}

	return nil
}
func (msg TransferOutMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

type BindStatus int8

const (
	BindStatusSuccess          BindStatus = 0
	BindStatusRejected         BindStatus = 1
	BindStatusTimeout          BindStatus = 2
	BindStatusInvalidParameter BindStatus = 3
)
