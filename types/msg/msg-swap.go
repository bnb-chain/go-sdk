package msg

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

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

type HashTimerLockTransferMsg struct {
	From                types.AccAddress `json:"from"`
	To                  types.AccAddress `json:"to"`
	RecipientOtherChain types.HexData    `json:"recipient_other_chain"`
	RandomNumberHash    types.HexData    `json:"random_number_hash"`
	Timestamp           int64            `json:"timestamp"`
	OutAmount           types.Coin       `json:"out_amount"`
	InAmountOtherChain  int64            `json:"in_amount_other_chain"`
	HeightSpan          int64            `json:"height_span"`
}

func NewHashTimerLockTransferMsg(from, to types.AccAddress, recipientOtherChain []byte, randomNumberHash []byte, timestamp int64,
	outAmount types.Coin, inAmountOtherChain int64, heightSpan int64) HashTimerLockTransferMsg {
	return HashTimerLockTransferMsg{
		From:                from,
		To:                  to,
		RecipientOtherChain: recipientOtherChain,
		RandomNumberHash:    randomNumberHash,
		Timestamp:           timestamp,
		OutAmount:           outAmount,
		InAmountOtherChain:  inAmountOtherChain,
		HeightSpan:          heightSpan,
	}
}

func (msg HashTimerLockTransferMsg) Route() string { return AtomicSwapRoute }
func (msg HashTimerLockTransferMsg) Type() string  { return HTLT }
func (msg HashTimerLockTransferMsg) String() string {
	return fmt.Sprintf("hashTimerLockTransfer{%v#%v#%v#%v#%v#%v#%v#%v}", msg.From, msg.To, msg.RecipientOtherChain, msg.RandomNumberHash,
		msg.Timestamp, msg.OutAmount, msg.InAmountOtherChain, msg.HeightSpan)
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
	if len(msg.RecipientOtherChain) > 32 {
		return fmt.Errorf("the length of recipient address on other chain shouldn't exceed 32")
	}
	if len(msg.RandomNumberHash) != types.RandomNumberHashLength {
		return fmt.Errorf("the length of random number hash should be %d", types.RandomNumberHashLength)
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
	RandomNumberHash types.HexData    `json:"random_number_hash"`
	RandomNumber     types.HexData    `json:"random_number"`
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
	if len(msg.RandomNumberHash) != types.RandomNumberHashLength {
		return fmt.Errorf("the length of random number hash should be %d", types.RandomNumberHashLength)
	}
	if len(msg.RandomNumber) != types.RandomNumberLength {
		return fmt.Errorf("the length of random number should be %d", types.RandomNumberLength)
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
	RandomNumberHash types.HexData    `json:"random_number_hash"`
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
	if len(msg.RandomNumberHash) != types.RandomNumberHashLength {
		return fmt.Errorf("the length of random number hash should be %d", types.RandomNumberHashLength)
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
	randomNumberAndTimestamp := make([]byte, types.RandomNumberLength+types.Int64Size)
	copy(randomNumberAndTimestamp[:types.RandomNumberLength], randomNumber)
	binary.BigEndian.PutUint64(randomNumberAndTimestamp[types.RandomNumberLength:], uint64(timestamp))
	return tmhash.Sum(randomNumberAndTimestamp)
}
