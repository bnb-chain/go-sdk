package types

import "time"

type TimeLockRecord struct {
	Id          int64     `json:"id"`
	Description string    `json:"description"`
	Amount      Coins `json:"amount"`
	LockTime    time.Time `json:"lock_time"`
}

type TimeLockRecords []TimeLockRecord

func (a TimeLockRecords) Len() int {
	return len(a)
}
func (a TimeLockRecords) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a TimeLockRecords) Less(i, j int) bool {
	return a[i].Id < a[j].Id
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