package types

import (
	"encoding/json"

	cmm "github.com/tendermint/tendermint/libs/common"
)

type SwapStatus byte

const (
	NULL      SwapStatus = 0x00
	Open      SwapStatus = 0x01
	Completed SwapStatus = 0x02
	Expired   SwapStatus = 0x03
)

func NewSwapStatusFromString(str string) SwapStatus {
	switch str {
	case "Open", "open":
		return Open
	case "Completed", "completed":
		return Completed
	case "Expired", "expired":
		return Expired
	default:
		return NULL
	}
}

func (status SwapStatus) String() string {
	switch status {
	case Open:
		return "Open"
	case Completed:
		return "Completed"
	case Expired:
		return "Expired"
	default:
		return "NULL"
	}
}

func (status SwapStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(status.String())
}

func (status *SwapStatus) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	*status = NewSwapStatusFromString(s)
	return nil
}

type AtomicSwap struct {
	From      AccAddress `json:"from"`
	To        AccAddress `json:"to"`
	OutAmount Coins      `json:"out_amount"`
	InAmount  Coins      `json:"in_amount"`

	ExpectedIncome      string `json:"expected_income"`
	RecipientOtherChain string `json:"recipient_other_chain"`

	RandomNumberHash cmm.HexBytes `json:"random_number_hash"`
	RandomNumber     cmm.HexBytes `json:"random_number"`
	Timestamp        int64        `json:"timestamp"`

	CrossChain bool `json:"cross_chain"`

	ExpireHeight int64      `json:"expire_height"`
	Index        int64      `json:"index"`
	ClosedTime   int64      `json:"closed_time"`
	Status       SwapStatus `json:"status"`
}

// Params for query 'custom/atomicswap/swapid'
type QuerySwapByID struct {
	SwapID cmm.HexBytes
}

// Params for query 'custom/atomicswap/swapcreator'
type QuerySwapByCreatorParams struct {
	Creator AccAddress
	Limit   int64
	Offset  int64
}

// Params for query 'custom/atomicswap/swaprecipient'
type QuerySwapByRecipientParams struct {
	Recipient AccAddress
	Limit     int64
	Offset    int64
}
