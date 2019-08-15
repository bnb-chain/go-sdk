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
	DepositHTLT     = "depositHTLT"
	ClaimHTLT       = "claimHTLT"
	RefundHTLT      = "refundHTLT"

	Int64Size                    = 8
	RandomNumberHashLength       = 32
	RandomNumberLength           = 32
	MaxRecipientOtherChainLength = 32
	MaxExpectedIncomeLength      = 64
	MinimumHeightSpan            = 360
	MaximumHeightSpan            = 518400
)

var (
	// bnb prefix address:  bnb1wxeplyw7x8aahy93w96yhwm7xcq3ke4f8ge93u
	// tbnb prefix address: tbnb1wxeplyw7x8aahy93w96yhwm7xcq3ke4ffasp3d
	AtomicSwapCoinsAccAddr = types.AccAddress(crypto.AddressHash([]byte("BinanceChainAtomicSwapCoins")))
)

type HashTimerLockedTransferMsg struct {
	From                types.AccAddress `json:"from"`
	To                  types.AccAddress `json:"to"`
	RecipientOtherChain types.HexData    `json:"recipient_other_chain"`
	RandomNumberHash    types.HexData    `json:"random_number_hash"`
	Timestamp           int64            `json:"timestamp"`
	OutAmount           types.Coin       `json:"out_amount"`
	ExpectedIncome      string           `json:"expected_income"`
	HeightSpan          int64            `json:"height_span"`
	CrossChain          bool             `json:"cross_chain"`
}

func NewHashTimerLockedTransferMsg(from, to types.AccAddress, recipientOtherChain []byte, randomNumberHash []byte, timestamp int64,
	outAmount types.Coin, expectedIncome string, heightSpan int64, crossChain bool) HashTimerLockedTransferMsg {
	return HashTimerLockedTransferMsg{
		From:                from,
		To:                  to,
		RecipientOtherChain: recipientOtherChain,
		RandomNumberHash:    randomNumberHash,
		Timestamp:           timestamp,
		OutAmount:           outAmount,
		ExpectedIncome:      expectedIncome,
		HeightSpan:          heightSpan,
		CrossChain:          crossChain,
	}
}

func (msg HashTimerLockedTransferMsg) Route() string { return AtomicSwapRoute }
func (msg HashTimerLockedTransferMsg) Type() string  { return HTLT }
func (msg HashTimerLockedTransferMsg) String() string {
	return fmt.Sprintf("HTLT{%v#%v#%v#%v#%v#%v#%v#%v#%v}", msg.From, msg.To, msg.RecipientOtherChain, msg.RandomNumberHash,
		msg.Timestamp, msg.OutAmount, msg.ExpectedIncome, msg.HeightSpan, msg.CrossChain)
}
func (msg HashTimerLockedTransferMsg) GetInvolvedAddresses() []types.AccAddress {
	return append(msg.GetSigners(), AtomicSwapCoinsAccAddr)
}
func (msg HashTimerLockedTransferMsg) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.From}
}

func (msg HashTimerLockedTransferMsg) ValidateBasic() error {
	if len(msg.From) != types.AddrLen {
		return fmt.Errorf("Expected address length is %d, actual length is %d", types.AddrLen, len(msg.From))
	}
	if len(msg.To) != types.AddrLen {
		return fmt.Errorf("Expected address length is %d, actual length is %d", types.AddrLen, len(msg.To))
	}
	if !msg.CrossChain && len(msg.RecipientOtherChain) != 0 {
		return fmt.Errorf("Must leave recipient address on other chain to empty for single chain swap")
	}
	if msg.CrossChain && len(msg.RecipientOtherChain) == 0 {
		return fmt.Errorf("Missing recipient address on other chain for cross chain swap")
	}
	if len(msg.RecipientOtherChain) > MaxRecipientOtherChainLength {
		return fmt.Errorf("The length of recipient address on other chain should be less than %d", MaxRecipientOtherChainLength)
	}
	if len(msg.ExpectedIncome) > MaxExpectedIncomeLength {
		return fmt.Errorf("The length of expected income should be less than %d", MaxExpectedIncomeLength)
	}
	if len(msg.RandomNumberHash) != RandomNumberHashLength {
		return fmt.Errorf("The length of random number hash should be %d", RandomNumberHashLength)
	}
	if !msg.OutAmount.IsPositive() {
		return fmt.Errorf("The swapped out coin must be positive")
	}
	if msg.HeightSpan < MinimumHeightSpan || msg.HeightSpan > MaximumHeightSpan {
		return fmt.Errorf("The height span should be no less than 360 and no greater than 518400")
	}
	return nil
}

func (msg HashTimerLockedTransferMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

type DepositHashTimerLockedTransferMsg struct {
	From             types.AccAddress `json:"from"`
	To               types.AccAddress `json:"to"`
	OutAmount        types.Coin       `json:"out_amount"`
	RandomNumberHash types.HexData    `json:"random_number_hash"`
}

func NewDepositHashTimerLockedTransferMsg(from, to types.AccAddress, outAmount types.Coin, randomNumberHash []byte) DepositHashTimerLockedTransferMsg {
	return DepositHashTimerLockedTransferMsg{
		From:             from,
		To:               to,
		OutAmount:        outAmount,
		RandomNumberHash: randomNumberHash,
	}
}

func (msg DepositHashTimerLockedTransferMsg) Route() string { return AtomicSwapRoute }
func (msg DepositHashTimerLockedTransferMsg) Type() string  { return DepositHTLT }
func (msg DepositHashTimerLockedTransferMsg) String() string {
	return fmt.Sprintf("depositHTLT{%v#%v#%v#%v}", msg.From, msg.To, msg.OutAmount, msg.RandomNumberHash)
}
func (msg DepositHashTimerLockedTransferMsg) GetInvolvedAddresses() []types.AccAddress {
	return append(msg.GetSigners(), AtomicSwapCoinsAccAddr)
}
func (msg DepositHashTimerLockedTransferMsg) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.From}
}

func (msg DepositHashTimerLockedTransferMsg) ValidateBasic() error {
	if len(msg.From) != types.AddrLen {
		return fmt.Errorf("Expected address length is %d, actual length is %d", types.AddrLen, len(msg.From))
	}
	if len(msg.To) != types.AddrLen {
		return fmt.Errorf("Expected address length is %d, actual length is %d", types.AddrLen, len(msg.To))
	}
	if len(msg.RandomNumberHash) != RandomNumberHashLength {
		return fmt.Errorf("The length of random number hash should be %d", RandomNumberHashLength)
	}
	if !msg.OutAmount.IsPositive() {
		return fmt.Errorf("The swapped out coin must be positive")
	}
	return nil
}

func (msg DepositHashTimerLockedTransferMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

type ClaimHashTimerLockedTransferMsg struct {
	From             types.AccAddress `json:"from"`
	RandomNumberHash types.HexData    `json:"random_number_hash"`
	RandomNumber     types.HexData    `json:"random_number"`
}

func NewClaimHashTimerLockedTransferMsg(from types.AccAddress, randomNumberHash, randomNumber []byte) ClaimHashTimerLockedTransferMsg {
	return ClaimHashTimerLockedTransferMsg{
		From:             from,
		RandomNumberHash: randomNumberHash,
		RandomNumber:     randomNumber,
	}
}

func (msg ClaimHashTimerLockedTransferMsg) Route() string { return AtomicSwapRoute }
func (msg ClaimHashTimerLockedTransferMsg) Type() string  { return ClaimHTLT }
func (msg ClaimHashTimerLockedTransferMsg) String() string {
	return fmt.Sprintf("claimHashTimeLock{%v#%v#%v}", msg.From, msg.RandomNumberHash, msg.RandomNumber)
}
func (msg ClaimHashTimerLockedTransferMsg) GetInvolvedAddresses() []types.AccAddress {
	return append(msg.GetSigners(), AtomicSwapCoinsAccAddr)
}
func (msg ClaimHashTimerLockedTransferMsg) GetSigners() []types.AccAddress { return []types.AccAddress{msg.From} }

func (msg ClaimHashTimerLockedTransferMsg) ValidateBasic() error {
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

func (msg ClaimHashTimerLockedTransferMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

type RefundHashTimerLockedTransferMsg struct {
	From             types.AccAddress `json:"from"`
	RandomNumberHash types.HexData    `json:"random_number_hash"`
}

func NewRefundHashTimerLockedTransferMsg(from types.AccAddress, randomNumberHash []byte) RefundHashTimerLockedTransferMsg {
	return RefundHashTimerLockedTransferMsg{
		From:             from,
		RandomNumberHash: randomNumberHash,
	}
}

func (msg RefundHashTimerLockedTransferMsg) Route() string { return AtomicSwapRoute }
func (msg RefundHashTimerLockedTransferMsg) Type() string  { return RefundHTLT }
func (msg RefundHashTimerLockedTransferMsg) String() string {
	return fmt.Sprintf("refundLockedAsset{%v#%v}", msg.From, msg.RandomNumberHash)
}
func (msg RefundHashTimerLockedTransferMsg) GetInvolvedAddresses() []types.AccAddress {
	return append(msg.GetSigners(), AtomicSwapCoinsAccAddr)
}
func (msg RefundHashTimerLockedTransferMsg) GetSigners() []types.AccAddress { return []types.AccAddress{msg.From} }

func (msg RefundHashTimerLockedTransferMsg) ValidateBasic() error {
	if len(msg.From) != types.AddrLen {
		return fmt.Errorf("expected address length is %d, actual length is %d", types.AddrLen, len(msg.From))
	}
	if len(msg.RandomNumberHash) != RandomNumberHashLength {
		return fmt.Errorf("the length of random number hash should be %d", RandomNumberHashLength)
	}
	return nil
}

func (msg RefundHashTimerLockedTransferMsg) GetSignBytes() []byte {
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
