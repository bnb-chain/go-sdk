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

type HTLTMsg struct {
	From                types.AccAddress `json:"from"`
	To                  types.AccAddress `json:"to"`
	RecipientOtherChain types.HexData    `json:"recipient_other_chain"`
	RandomNumberHash    types.HexData    `json:"random_number_hash"`
	Timestamp           int64            `json:"timestamp"`
	OutAmount           types.Coins      `json:"out_amount"`
	ExpectedIncome      string           `json:"expected_income"`
	HeightSpan          int64            `json:"height_span"`
	CrossChain          bool             `json:"cross_chain"`
}

func NewHTLTMsg(from, to types.AccAddress, recipientOtherChain []byte, randomNumberHash []byte, timestamp int64,
	outAmount types.Coins, expectedIncome string, heightSpan int64, crossChain bool) HTLTMsg {
	return HTLTMsg{
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

func (msg HTLTMsg) Route() string { return AtomicSwapRoute }
func (msg HTLTMsg) Type() string  { return HTLT }
func (msg HTLTMsg) String() string {
	return fmt.Sprintf("HTLT{%v#%v#%v#%v#%v#%v#%v#%v#%v}", msg.From, msg.To, msg.RecipientOtherChain, msg.RandomNumberHash,
		msg.Timestamp, msg.OutAmount, msg.ExpectedIncome, msg.HeightSpan, msg.CrossChain)
}
func (msg HTLTMsg) GetInvolvedAddresses() []types.AccAddress {
	return append(msg.GetSigners(), AtomicSwapCoinsAccAddr)
}
func (msg HTLTMsg) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.From}
}

func (msg HTLTMsg) ValidateBasic() error {
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

func (msg HTLTMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

type DepositHTLTMsg struct {
	From             types.AccAddress `json:"from"`
	To               types.AccAddress `json:"to"`
	OutAmount        types.Coins      `json:"out_amount"`
	RandomNumberHash types.HexData    `json:"random_number_hash"`
}

func NewDepositHTLTMsg(from, to types.AccAddress, outAmount types.Coins, randomNumberHash []byte) DepositHTLTMsg {
	return DepositHTLTMsg{
		From:             from,
		To:               to,
		OutAmount:        outAmount,
		RandomNumberHash: randomNumberHash,
	}
}

func (msg DepositHTLTMsg) Route() string { return AtomicSwapRoute }
func (msg DepositHTLTMsg) Type() string  { return DepositHTLT }
func (msg DepositHTLTMsg) String() string {
	return fmt.Sprintf("depositHTLT{%v#%v#%v#%v}", msg.From, msg.To, msg.OutAmount, msg.RandomNumberHash)
}
func (msg DepositHTLTMsg) GetInvolvedAddresses() []types.AccAddress {
	return append(msg.GetSigners(), AtomicSwapCoinsAccAddr)
}
func (msg DepositHTLTMsg) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.From}
}

func (msg DepositHTLTMsg) ValidateBasic() error {
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

func (msg DepositHTLTMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

type ClaimHTLTMsg struct {
	From             types.AccAddress `json:"from"`
	RandomNumberHash types.HexData    `json:"random_number_hash"`
	RandomNumber     types.HexData    `json:"random_number"`
}

func NewClaimHTLTMsg(from types.AccAddress, randomNumberHash, randomNumber []byte) ClaimHTLTMsg {
	return ClaimHTLTMsg{
		From:             from,
		RandomNumberHash: randomNumberHash,
		RandomNumber:     randomNumber,
	}
}

func (msg ClaimHTLTMsg) Route() string { return AtomicSwapRoute }
func (msg ClaimHTLTMsg) Type() string  { return ClaimHTLT }
func (msg ClaimHTLTMsg) String() string {
	return fmt.Sprintf("claimHashTimeLock{%v#%v#%v}", msg.From, msg.RandomNumberHash, msg.RandomNumber)
}
func (msg ClaimHTLTMsg) GetInvolvedAddresses() []types.AccAddress {
	return append(msg.GetSigners(), AtomicSwapCoinsAccAddr)
}
func (msg ClaimHTLTMsg) GetSigners() []types.AccAddress { return []types.AccAddress{msg.From} }

func (msg ClaimHTLTMsg) ValidateBasic() error {
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

func (msg ClaimHTLTMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

type RefundHTLTMsg struct {
	From             types.AccAddress `json:"from"`
	RandomNumberHash types.HexData    `json:"random_number_hash"`
}

func NewRefundHTLTMsg(from types.AccAddress, randomNumberHash []byte) RefundHTLTMsg {
	return RefundHTLTMsg{
		From:             from,
		RandomNumberHash: randomNumberHash,
	}
}

func (msg RefundHTLTMsg) Route() string { return AtomicSwapRoute }
func (msg RefundHTLTMsg) Type() string  { return RefundHTLT }
func (msg RefundHTLTMsg) String() string {
	return fmt.Sprintf("refundLockedAsset{%v#%v}", msg.From, msg.RandomNumberHash)
}
func (msg RefundHTLTMsg) GetInvolvedAddresses() []types.AccAddress {
	return append(msg.GetSigners(), AtomicSwapCoinsAccAddr)
}
func (msg RefundHTLTMsg) GetSigners() []types.AccAddress { return []types.AccAddress{msg.From} }

func (msg RefundHTLTMsg) ValidateBasic() error {
	if len(msg.From) != types.AddrLen {
		return fmt.Errorf("expected address length is %d, actual length is %d", types.AddrLen, len(msg.From))
	}
	if len(msg.RandomNumberHash) != RandomNumberHashLength {
		return fmt.Errorf("the length of random number hash should be %d", RandomNumberHashLength)
	}
	return nil
}

func (msg RefundHTLTMsg) GetSignBytes() []byte {
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
