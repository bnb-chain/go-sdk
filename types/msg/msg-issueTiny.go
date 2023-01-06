package msg

import (
	"github.com/bnb-chain/node/plugins/tokens/issue"
)

const (
	IssueTinyMsgType = issue.IssueTinyMsgType
)

type TinyTokenIssueMsg = issue.IssueTinyMsg

var NewTinyTokenIssueMsg = issue.NewIssueTinyMsg
