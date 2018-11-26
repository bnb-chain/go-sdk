package txmsg

import (
	"encoding/json"
	"fmt"
)

// TestMsg type for testing
type TestMsg struct {
	signers []AccAddress
}

// NewTestMsg to creates a TestMsg
func NewTestMsg(addrs ...AccAddress) *TestMsg {
	return &TestMsg{
		signers: addrs,
	}
}

// Route is part of Msg interface
func (msg *TestMsg) Route() string { return "TestMsg" }

// Type is part of Msg interface
func (msg *TestMsg) Type() string { return "TestMsg" }

// GetSignBytes is part of Msg interface
func (msg *TestMsg) GetSignBytes() []byte {
	bz, err := json.Marshal(msg.signers)
	if err != nil {
		panic(err)
	}
	return bz
}

// GetSigners is part of Msg interface
func (msg *TestMsg) GetSigners() []AccAddress {
	return msg.signers
}

// GetInvolvedAddresses as part of the Msg interface
func (msg *TestMsg) GetInvolvedAddresses() []AccAddress {
	return msg.GetSigners()
}

// String is part of Msg interface
func (msg *TestMsg) String() string {
	return fmt.Sprintf("TestMsg")
}

// ValidateBasic is part of Msg interface
func (msg *TestMsg) ValidateBasic() error { return nil }
