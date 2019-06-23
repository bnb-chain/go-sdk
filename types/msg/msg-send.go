package msg

import (
	"encoding/json"
	"fmt"

	"github.com/binance-chain/go-sdk/common/types"
)

// SendMsg - high level transaction of the coin module
type SendMsg struct {
	Inputs  []Input  `json:"inputs"`
	Outputs []Output `json:"outputs"`
}

// NewMsgSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgSend(in []Input, out []Output) SendMsg {
	return SendMsg{Inputs: in, Outputs: out}
}

func (msg SendMsg) Route() string { return "bank" } // TODO: "bank/send"
func (msg SendMsg) Type() string  { return "send" }

// Implements Msg.
func (msg SendMsg) ValidateBasic() error {
	if len(msg.Inputs) == 0 {
		return fmt.Errorf("Len of inputs is less than 1 ")
	}
	if len(msg.Outputs) == 0 {
		return fmt.Errorf("Len of outputs is less than 1 ")
	}
	// make sure all inputs and outputs are individually valid
	var totalIn, totalOut types.Coins
	for _, in := range msg.Inputs {
		if err := in.ValidateBasic(); err != nil {
			return err
		}
		totalIn = totalIn.Plus(in.Coins)
	}
	for _, out := range msg.Outputs {
		if err := out.ValidateBasic(); err != nil {
			return err
		}
		totalOut = totalOut.Plus(out.Coins)
	}
	// make sure inputs and outputs match
	if !totalIn.IsEqual(totalOut) {
		return fmt.Errorf("inputs %v and outputs %v don't match", totalIn, totalOut)
	}
	return nil
}

// Implements Msg.
func (msg SendMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// Implements Msg.
func (msg SendMsg) GetSigners() []types.AccAddress {
	addrs := make([]types.AccAddress, len(msg.Inputs))
	for i, in := range msg.Inputs {
		addrs[i] = in.Address
	}
	return addrs
}

func (msg SendMsg) GetInvolvedAddresses() []types.AccAddress {
	numOfInputs := len(msg.Inputs)
	numOfOutputs := len(msg.Outputs)
	addrs := make([]types.AccAddress, numOfInputs+numOfOutputs, numOfInputs+numOfOutputs)
	for i, in := range msg.Inputs {
		addrs[i] = in.Address
	}
	for i, out := range msg.Outputs {
		addrs[i+numOfInputs] = out.Address
	}
	return addrs
}

//----------------------------------------
// Input

// Transaction Input
type Input struct {
	Address types.AccAddress `json:"address"`
	Coins   types.Coins      `json:"coins"`
}

// Return bytes to sign for Input
func (in Input) GetSignBytes() []byte {
	bin, err := MsgCdc.MarshalJSON(in)
	if err != nil {
		panic(err)
	}
	return MustSortJSON(bin)
}

// ValidateBasic - validate transaction input
func (in Input) ValidateBasic() error {
	if len(in.Address) == 0 {
		return fmt.Errorf("Len of input address is less than 1 ")
	}
	if !in.Coins.IsValid() {
		return fmt.Errorf("Inputs coins %v is invalid ", in.Coins)
	}
	if !in.Coins.IsPositive() {
		return fmt.Errorf("Inputs coins %v is negative ", in.Coins)
	}
	return nil
}

// NewInput - create a transaction input, used with SendMsg
func NewInput(addr types.AccAddress, coins types.Coins) Input {
	input := Input{
		Address: addr,
		Coins:   coins,
	}
	return input
}

//----------------------------------------
// Output

// Transaction Output
type Output struct {
	Address types.AccAddress `json:"address"`
	Coins   types.Coins      `json:"coins"`
}

// Return bytes to sign for Output
func (out Output) GetSignBytes() []byte {
	bin, err := MsgCdc.MarshalJSON(out)
	if err != nil {
		panic(err)
	}
	return MustSortJSON(bin)
}

// ValidateBasic - validate transaction output
func (out Output) ValidateBasic() error {
	if len(out.Address) == 0 {
		return fmt.Errorf("Len output %d should is less than 1 ", 0)
	}
	if !out.Coins.IsValid() {
		return fmt.Errorf("Coins is invalid ")
	}
	if !out.Coins.IsPositive() {
		return fmt.Errorf(" Coins is negative ")
	}
	return nil
}

// NewOutput - create a transaction output, used with SendMsg
func NewOutput(addr types.AccAddress, coins types.Coins) Output {
	output := Output{
		Address: addr,
		Coins:   coins,
	}
	return output
}

type Transfer struct {
	ToAddr types.AccAddress
	Coins  types.Coins
}

func CreateSendMsg(from types.AccAddress, fromCoins types.Coins, transfers []Transfer) SendMsg {
	input := NewInput(from, fromCoins)

	output := make([]Output, 0, len(transfers))
	for _, t := range transfers {
		t.Coins = t.Coins.Sort()
		output = append(output, NewOutput(t.ToAddr, t.Coins))
	}
	msg := NewMsgSend([]Input{input}, output)
	return msg
}
