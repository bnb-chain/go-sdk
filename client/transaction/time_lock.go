package transaction

import (
	"strconv"

	"github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
)

type TimeLockResult struct {
	tx.TxCommitResult
	LockId int64 `json:"lock_id"`
}

func (c *client) TimeLock(description string, amount types.Coins, lockTime int64, sync bool, options ...Option) (*TimeLockResult, error) {
	fromAddr := c.keyManager.GetAddr()

	lockMsg := msg.NewTimeLockMsg(fromAddr, description, amount, lockTime)
	commit, err := c.broadcastMsg(lockMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	var lockId int64
	if commit.Ok && sync {
		lockId, err = strconv.ParseInt(string(commit.Data), 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return &TimeLockResult{*commit, lockId}, err
}

type TimeUnLockResult struct {
	tx.TxCommitResult
	LockId int64 `json:"lock_id"`
}

func (c *client) TimeUnLock(id int64, sync bool, options ...Option) (*TimeUnLockResult, error) {
	fromAddr := c.keyManager.GetAddr()

	unlockMsg := msg.NewTimeUnlockMsg(fromAddr, id)
	err := unlockMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := c.broadcastMsg(unlockMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	var lockId int64
	if commit.Ok && sync {
		lockId, err = strconv.ParseInt(string(commit.Data), 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return &TimeUnLockResult{*commit, lockId}, err
}

type TimeReLockResult struct {
	tx.TxCommitResult
	LockId int64 `json:"lock_id"`
}

func (c *client) TimeReLock(id int64, description string, amount types.Coins, lockTime int64, sync bool, options ...Option) (*TimeReLockResult, error) {
	fromAddr := c.keyManager.GetAddr()

	relockMsg := msg.NewTimeRelockMsg(fromAddr, id, description, amount, lockTime)
	err := relockMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := c.broadcastMsg(relockMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	var lockId int64
	if commit.Ok && sync {
		lockId, err = strconv.ParseInt(string(commit.Data), 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return &TimeReLockResult{*commit, lockId}, err
}
