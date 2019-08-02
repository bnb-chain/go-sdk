package types

import "time"

type TimeLockRecord struct {
	Id          int64     `json:"id"`
	Description string    `json:"description"`
	Amount      Coins     `json:"amount"`
	LockTime    time.Time `json:"lock_time"`
}

// Params for query 'custom/timelock/timelocks'
type QueryTimeLocksParams struct {
	Account AccAddress
}

// Params for query 'custom/timelock/timelock'
type QueryTimeLockParams struct {
	Account AccAddress
	Id      int64
}
