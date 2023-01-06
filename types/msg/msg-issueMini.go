package msg

import (
	"github.com/bnb-chain/node/plugins/tokens/issue"
)

const (
	IssueMiniMsgType = issue.IssueMiniMsgType
)

type MiniTokenIssueMsg = issue.IssueMiniMsg

var NewMiniTokenIssueMsg = issue.NewIssueMiniMsg
