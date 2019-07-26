package msg

import (
	"encoding/json"
	"fmt"

	"github.com/binance-chain/go-sdk/common/types"
)

const (
	AccountFlagsRoute      = "accountFlags"
	SetAccountFlagsMsgType = "setAccountFlags"
)

type SetAccountFlagsMsg struct {
	From  types.AccAddress `json:"from"`
	Flags uint64           `json:"flags"`
}

func NewSetAccountFlagsMsg(from types.AccAddress, flags uint64) SetAccountFlagsMsg {
	return SetAccountFlagsMsg{
		From:  from,
		Flags: flags,
	}
}

func (msg SetAccountFlagsMsg) Route() string { return AccountFlagsRoute }
func (msg SetAccountFlagsMsg) Type() string  { return SetAccountFlagsMsgType }
func (msg SetAccountFlagsMsg) String() string {
	return fmt.Sprintf("setAccountFlags{%v#%v}", msg.From, msg.Flags)
}
func (msg SetAccountFlagsMsg) GetInvolvedAddresses() []types.AccAddress { return msg.GetSigners() }
func (msg SetAccountFlagsMsg) GetSigners() []types.AccAddress           { return []types.AccAddress{msg.From} }

func (msg SetAccountFlagsMsg) ValidateBasic() error {
	if len(msg.From) != types.AddrLen {
		return fmt.Errorf("Expected address length is %d, actual length is %d", types.AddrLen, len(msg.From))
	}
	return nil
}

func (msg SetAccountFlagsMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}
