package msg

import (
	"github.com/bnb-chain/node/plugins/tokens/timelock"
)

const (
	MaxTimeLockDescriptionLength = timelock.MaxDescriptionLength
	MinLockTime                  = timelock.MinLockTime
)

var (
	TimeLockCoinsAccAddr = timelock.TimeLockCoinsAccAddr

	NewTimeLockMsg   = timelock.NewTimeLockMsg
	NewTimeRelockMsg = timelock.NewTimeRelockMsg
	NewTimeUnlockMsg = timelock.NewTimeUnlockMsg
)

type (
	TimeLockMsg   = timelock.TimeLockMsg
	TimeRelockMsg = timelock.TimeRelockMsg
	TimeUnlockMsg = timelock.TimeUnlockMsg
)
