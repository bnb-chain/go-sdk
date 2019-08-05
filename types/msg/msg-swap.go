package msg

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/tmhash"

	"github.com/binance-chain/go-sdk/common/types"
)

const (
	AtomicSwapRoute = "atomicSwap"
	HTLT            = "HTLT"
	ClaimHTLT       = "claimHTLT"
	RefundHTLT      = "refundHTLT"
)

var (
	// bnb prefix address:  bnb1wxeplyw7x8aahy93w96yhwm7xcq3ke4f8ge93u
	// tbnb prefix address: tbnb1wxeplyw7x8aahy93w96yhwm7xcq3ke4ffasp3d
	AtomicSwapCoinsAccAddr = types.AccAddress(crypto.AddressHash([]byte("BinanceChainAtomicSwapCoins")))
)

type SwapStatus byte
type HexData []byte

const (
	NULL      SwapStatus = 0x00
	Open      SwapStatus = 0x01
	Completed SwapStatus = 0x02
	Expired   SwapStatus = 0x03

	RandomNumberHashLength = 32
	RandomNumberLength     = 32
	Int64Size              = 8
)

func (hexData HexData) String() string {
	str := hex.EncodeToString(hexData)
	if len(str) == 0 {
		return ""
	}
	return "0x" + hex.EncodeToString(hexData)
}

func (hexData HexData) MarshalJSON() ([]byte, error) {
	return json.Marshal(hexData.String())
}

func (hexData *HexData) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(s, "0x") {
		return fmt.Errorf("hex string must prefix with 0x")
	}
	bytesArray, err := hex.DecodeString(s[2:])
	if err != nil {
		return err
	}
	*hexData = bytesArray
	return nil
}

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

type HashTimerLockTransferMsg struct {
	From             types.AccAddress `json:"from"`
	To               types.AccAddress `json:"to"`
	ToOnOtherChain   HexData          `json:"to_on_other_chain"`
	RandomNumberHash HexData          `json:"random_number_hash"`
	Timestamp        int64            `json:"timestamp"`
	OutAmount        types.Coin       `json:"out_amount"`
	InAmount         int64            `json:"in_amount"`
	HeightSpan       int64            `json:"height_span"`
}

func NewHashTimerLockTransferMsg(from, to types.AccAddress, toOnOtherChain []byte, randomNumberHash []byte, timestamp int64,
	outAmount types.Coin, inAmount int64, heightSpan int64) HashTimerLockTransferMsg {
	return HashTimerLockTransferMsg{
		From:             from,
		To:               to,
		ToOnOtherChain:   toOnOtherChain,
		RandomNumberHash: randomNumberHash,
		Timestamp:        timestamp,
		OutAmount:        outAmount,
		InAmount:         inAmount,
		HeightSpan:       heightSpan,
	}
}

func (msg HashTimerLockTransferMsg) Route() string { return AtomicSwapRoute }
func (msg HashTimerLockTransferMsg) Type() string  { return HTLT }
func (msg HashTimerLockTransferMsg) String() string {
	return fmt.Sprintf("hashTimerLockTransfer{%v#%v#%v#%v#%v#%v#%v#%v}", msg.From, msg.To, msg.ToOnOtherChain, msg.RandomNumberHash,
		msg.Timestamp, msg.OutAmount, msg.InAmount, msg.HeightSpan)
}
func (msg HashTimerLockTransferMsg) GetInvolvedAddresses() []types.AccAddress {
	return append(msg.GetSigners(), AtomicSwapCoinsAccAddr)
}
func (msg HashTimerLockTransferMsg) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.From}
}

func (msg HashTimerLockTransferMsg) ValidateBasic() error {
	if len(msg.From) != types.AddrLen {
		return fmt.Errorf("expected address length is %d, actual length is %d", types.AddrLen, len(msg.From))
	}
	if len(msg.To) != types.AddrLen {
		return fmt.Errorf("expected address length is %d, actual length is %d", types.AddrLen, len(msg.To))
	}
	if len(msg.ToOnOtherChain) == 0 || len(msg.ToOnOtherChain) > 32 {
		return fmt.Errorf("the receiver address on other chain shouldn't be nil and its length shouldn't exceed 32")
	}
	if len(msg.RandomNumberHash) != RandomNumberHashLength {
		return fmt.Errorf("the length of random number hash should be %d", RandomNumberHashLength)
	}
	if !msg.OutAmount.IsPositive() {
		return fmt.Errorf("the swapped out coin must be positive")
	}
	if msg.HeightSpan < 360 || msg.HeightSpan > 518400 {
		return fmt.Errorf("the height span should be no less than 360 and no greater than 518400")
	}
	return nil
}

func (msg HashTimerLockTransferMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

type ClaimHashTimerLockMsg struct {
	From             types.AccAddress `json:"from"`
	RandomNumberHash HexData          `json:"random_number_hash"`
	RandomNumber     HexData          `json:"random_number"`
}

func NewClaimHashTimerLockMsg(from types.AccAddress, randomNumberHash, randomNumber []byte) ClaimHashTimerLockMsg {
	return ClaimHashTimerLockMsg{
		From:             from,
		RandomNumberHash: randomNumberHash,
		RandomNumber:     randomNumber,
	}
}

func (msg ClaimHashTimerLockMsg) Route() string { return AtomicSwapRoute }
func (msg ClaimHashTimerLockMsg) Type() string  { return ClaimHTLT }
func (msg ClaimHashTimerLockMsg) String() string {
	return fmt.Sprintf("claimHashTimeLock{%v#%v#%v}", msg.From, msg.RandomNumberHash, msg.RandomNumber)
}
func (msg ClaimHashTimerLockMsg) GetInvolvedAddresses() []types.AccAddress {
	return append(msg.GetSigners(), AtomicSwapCoinsAccAddr)
}
func (msg ClaimHashTimerLockMsg) GetSigners() []types.AccAddress { return []types.AccAddress{msg.From} }

func (msg ClaimHashTimerLockMsg) ValidateBasic() error {
	if len(msg.From) != types.AddrLen {
		return fmt.Errorf("expected address length is %d, actual length is %d", types.AddrLen, len(msg.From))
	}
	if len(msg.RandomNumberHash) != RandomNumberHashLength {
		return fmt.Errorf("the length of random number hash should be %d", RandomNumberHashLength)
	}
	if len(msg.RandomNumber) != RandomNumberLength {
		return fmt.Errorf("the length of random number should be %d", RandomNumberLength)
	}
	return nil
}

func (msg ClaimHashTimerLockMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

type RefundHashTimerLockMsg struct {
	From             types.AccAddress `json:"from"`
	RandomNumberHash HexData          `json:"random_number_hash"`
}

func NewRefundLockedAssetMsg(from types.AccAddress, randomNumberHash []byte) RefundHashTimerLockMsg {
	return RefundHashTimerLockMsg{
		From:             from,
		RandomNumberHash: randomNumberHash,
	}
}

func (msg RefundHashTimerLockMsg) Route() string { return AtomicSwapRoute }
func (msg RefundHashTimerLockMsg) Type() string  { return RefundHTLT }
func (msg RefundHashTimerLockMsg) String() string {
	return fmt.Sprintf("refundLockedAsset{%v#%v}", msg.From, msg.RandomNumberHash)
}
func (msg RefundHashTimerLockMsg) GetInvolvedAddresses() []types.AccAddress {
	return append(msg.GetSigners(), AtomicSwapCoinsAccAddr)
}
func (msg RefundHashTimerLockMsg) GetSigners() []types.AccAddress { return []types.AccAddress{msg.From} }

func (msg RefundHashTimerLockMsg) ValidateBasic() error {
	if len(msg.From) != types.AddrLen {
		return fmt.Errorf("expected address length is %d, actual length is %d", types.AddrLen, len(msg.From))
	}
	if len(msg.RandomNumberHash) != RandomNumberHashLength {
		return fmt.Errorf("the length of random number hash should be %d", RandomNumberHashLength)
	}
	return nil
}

func (msg RefundHashTimerLockMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func CalculateRandomHash(randomNumber []byte, timestamp int64) []byte {
	randomNumberAndTimestamp := make([]byte, RandomNumberLength+Int64Size)
	copy(randomNumberAndTimestamp[:RandomNumberLength], randomNumber)
	binary.BigEndian.PutUint64(randomNumberAndTimestamp[RandomNumberLength:], uint64(timestamp))
	return tmhash.Sum(randomNumberAndTimestamp)
}
