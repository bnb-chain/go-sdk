package msg

import (
	"github.com/bnb-chain/node/plugins/account"
)

const (
	AccountFlagsRoute      = account.AccountFlagsRoute
	SetAccountFlagsMsgType = account.SetAccountFlagsMsgType
)

type SetAccountFlagsMsg = account.SetAccountFlagsMsg

var NewSetAccountFlagsMsg = account.NewSetAccountFlagsMsg
