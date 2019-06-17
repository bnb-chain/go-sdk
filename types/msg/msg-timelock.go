package msg

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/binance-chain/go-sdk/common/types"
	"github.com/tendermint/tendermint/crypto"
)

const (
	MaxTimeLockDescriptionLength = 128
	MinLockTime                  = 60 * time.Second

	InitialRecordId = 1
)

var (
	TimeLockCoinsAccAddr = types.AccAddress(crypto.AddressHash([]byte("BinanceChainTimeLockCoins")))
)

type TimeLockMsg struct {
	From        types.AccAddress `json:"from"`
	Description string           `json:"description"`
	Amount      types.Coins      `json:"amount"`
	LockTime    int64            `json:"lock_time"`
}

func NewTimeLockMsg(from types.AccAddress, description string, amount types.Coins, lockTime int64) TimeLockMsg {
	return TimeLockMsg{
		From:        from,
		Description: description,
		Amount:      amount,
		LockTime:    lockTime,
	}
}

func (msg TimeLockMsg) Route() string { return MsgRoute }
func (msg TimeLockMsg) Type() string  { return "timeLock" }
func (msg TimeLockMsg) String() string {
	return fmt.Sprintf("TimeLock{%s#%v#%v#%v}", msg.From, msg.Description, msg.Amount, msg.LockTime)
}
func (msg TimeLockMsg) GetInvolvedAddresses() []types.AccAddress {
	return []types.AccAddress{msg.From, TimeLockCoinsAccAddr}
}
func (msg TimeLockMsg) GetSigners() []types.AccAddress { return []types.AccAddress{msg.From} }

func (msg TimeLockMsg) ValidateBasic() error {
	if len(msg.Description) == 0 || len(msg.Description) > MaxTimeLockDescriptionLength {
		return fmt.Errorf("length of description(%d) should be larger than 0 and be less than or equal to %d",
			len(msg.Description), MaxTimeLockDescriptionLength)
	}

	if msg.LockTime <= 0 {
		return fmt.Errorf("lock time(%d) should be larger than 0", msg.LockTime)
	}

	if !msg.Amount.IsValid() {
		return fmt.Errorf("amount %v is invalid", msg.Amount)
	}

	if !msg.Amount.IsPositive() {
		return fmt.Errorf("amount %v can't be negative", msg.Amount)
	}

	return nil
}

func (msg TimeLockMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

type TimeRelockMsg struct {
	From        types.AccAddress `json:"from"`
	Id          int64            `json:"time_lock_id"`
	Description string           `json:"description"`
	Amount      types.Coins      `json:"amount"`
	LockTime    int64            `json:"lock_time"`
}

func NewTimeRelockMsg(from types.AccAddress, id int64, description string, amount types.Coins, lockTime int64) TimeRelockMsg {
	return TimeRelockMsg{
		From:        from,
		Id:          id,
		Description: description,
		Amount:      amount,
		LockTime:    lockTime,
	}
}

func (msg TimeRelockMsg) Route() string { return MsgRoute }
func (msg TimeRelockMsg) Type() string  { return "timeRelock" }
func (msg TimeRelockMsg) String() string {
	return fmt.Sprintf("TimeRelock{%v#%s#%v#%v#%v}", msg.Id, msg.From, msg.Description, msg.Amount, msg.LockTime)
}
func (msg TimeRelockMsg) GetInvolvedAddresses() []types.AccAddress {
	return []types.AccAddress{msg.From, TimeLockCoinsAccAddr}
}
func (msg TimeRelockMsg) GetSigners() []types.AccAddress { return []types.AccAddress{msg.From} }

func (msg TimeRelockMsg) ValidateBasic() error {
	if msg.Id < InitialRecordId {
		return fmt.Errorf("time lock id should not be less than %d", InitialRecordId)
	}

	if len(msg.Description) > MaxTimeLockDescriptionLength {
		return fmt.Errorf("length of description(%d) should be less than or equal to %d",
			len(msg.Description), MaxTimeLockDescriptionLength)
	}

	if msg.LockTime < 0 {
		return fmt.Errorf("lock time(%d) should not be less than 0", msg.LockTime)
	}

	if !msg.Amount.IsValid() {
		return fmt.Errorf("amount %v is invalid", msg.Amount)
	}

	if !msg.Amount.IsNotNegative() {
		return fmt.Errorf("amount %v can't be negative", msg.Amount)
	}

	if len(msg.Description) == 0 &&
		msg.Amount.IsZero() &&
		msg.LockTime == 0 {
		return fmt.Errorf("nothing to update for time lock")
	}

	return nil
}

func (msg TimeRelockMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

type TimeUnlockMsg struct {
	From types.AccAddress `json:"from"`
	Id   int64            `json:"time_lock_id"`
}

func NewTimeUnlockMsg(from types.AccAddress, id int64) TimeUnlockMsg {
	return TimeUnlockMsg{
		From: from,
		Id:   id,
	}
}

func (msg TimeUnlockMsg) Route() string { return MsgRoute }
func (msg TimeUnlockMsg) Type() string  { return "timeUnlock" }
func (msg TimeUnlockMsg) String() string {
	return fmt.Sprintf("TimeUnlock{%s#%v}", msg.From, msg.Id)
}
func (msg TimeUnlockMsg) GetInvolvedAddresses() []types.AccAddress {
	return []types.AccAddress{msg.From, TimeLockCoinsAccAddr}
}
func (msg TimeUnlockMsg) GetSigners() []types.AccAddress { return []types.AccAddress{msg.From} }

func (msg TimeUnlockMsg) ValidateBasic() error {
	if msg.Id < InitialRecordId {
		return fmt.Errorf("time lock id should not be less than %d", InitialRecordId)
	}
	return nil
}

func (msg TimeUnlockMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}
