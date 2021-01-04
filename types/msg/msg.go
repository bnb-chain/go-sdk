package msg

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/binance-chain/go-sdk/common/types"

	"github.com/pkg/errors"

	"github.com/binance-chain/go-sdk/common"
)

// constants
const (
	DotBSuffix                    = ".B"
	NativeToken                   = "BNB"
	NativeTokenDotBSuffixed       = "BNB" + DotBSuffix
	Decimals                int8  = 8
	MaxTotalSupply          int64 = 9000000000000000000 // 90 billions with 8 decimal digits

	TokenSymbolMaxLen          = 8
	TokenSymbolMinLen          = 2
	TokenSymbolTxHashSuffixLen = 3

	MiniTokenSymbolMaxLen          = 8
	MiniTokenSymbolMinLen          = 2
	MiniTokenSymbolSuffixLen       = 4
	MiniTokenSymbolMSuffix         = "M"
	MiniTokenSymbolTxHashSuffixLen = 3
	MaxMiniTokenNameLength         = 32
	MaxTokenURILength              = 2048
)

// Msg - Transactions messages must fulfill the Msg
type Msg interface {
	// Return the message type.
	// Must be alphanumeric or empty.
	Route() string

	// Returns a human-readable string for the message, intended for utilization
	// within tags
	Type() string

	// ValidateBasic does a simple validation check that
	// doesn't require access to any other information.
	ValidateBasic() error

	// Get the canonical byte representation of the Msg.
	GetSignBytes() []byte

	// Signers returns the addrs of signers that must sign.
	// CONTRACT: All signatures must be present to be valid.
	// CONTRACT: Returns addrs in some deterministic order.
	GetSigners() []types.AccAddress

	// Get involved addresses of this msg so that we can publish account balance change
	GetInvolvedAddresses() []types.AccAddress
}

// ValidateSymbol utility
func ValidateSymbol(symbol string) error {
	if len(symbol) == 0 {
		return errors.New("suffixed token symbol cannot be empty")
	}

	// suffix exception for native token (less drama in existing tests)
	if symbol == NativeToken ||
		symbol == NativeTokenDotBSuffixed {
		return nil
	}

	parts, err := splitSuffixedTokenSymbol(symbol)
	if err != nil {
		return err
	}

	symbolPart := parts[0]

	// since the native token was given a suffix exception above, do not allow it to have a suffix
	if symbolPart == NativeToken ||
		symbolPart == NativeTokenDotBSuffixed {
		return errors.New("native token symbol should not be suffixed with tx hash")
	}

	if strings.HasSuffix(symbolPart, DotBSuffix) {
		symbolPart = strings.TrimSuffix(symbolPart, DotBSuffix)
	}

	// check len without .B suffix
	if len(symbolPart) < TokenSymbolMinLen {
		return fmt.Errorf("token symbol part is too short, got %d chars", len(symbolPart))
	}
	if len(symbolPart) > TokenSymbolMaxLen {
		return fmt.Errorf("token symbol part is too long, got %d chars", len(symbolPart))
	}

	if !common.IsAlphaNum(symbolPart) {
		return errors.New("token symbol part should be alphanumeric")
	}

	txHashPart := parts[1]

	if len(txHashPart) != TokenSymbolTxHashSuffixLen {
		return fmt.Errorf("token symbol tx hash suffix must be %d chars in length, got %d", TokenSymbolTxHashSuffixLen, len(txHashPart))
	}

	// prohibit non-hexadecimal chars in the suffix part
	isHex, err := regexp.MatchString(fmt.Sprintf("[0-9A-F]{%d}", TokenSymbolTxHashSuffixLen), txHashPart)
	if err != nil {
		return err
	}
	if !isHex {
		return fmt.Errorf("token symbol tx hash suffix must be hex with a length of %d", TokenSymbolTxHashSuffixLen)
	}

	return nil
}

func splitSuffixedTokenSymbol(suffixed string) ([]string, error) {
	// as above, the native token symbol is given an exception - it is not required to be suffixed
	if suffixed == NativeToken ||
		suffixed == NativeTokenDotBSuffixed {
		return []string{suffixed, ""}, nil
	}

	split := strings.SplitN(suffixed, "-", 2)

	if len(split) != 2 {
		return nil, errors.New("suffixed token symbol must contain a hyphen ('-')")
	}

	if strings.Contains(split[1], "-") {
		return nil, errors.New("suffixed token symbol must contain just one hyphen ('-')")
	}

	return split, nil
}

type StatusText int

const (
	PendingStatusText StatusText = iota
	SuccessStatusText
	FailedStatusText
)

var StatusTextToString = [...]string{"pending", "success", "failed"}
var StringToStatusText = map[string]StatusText{
	"pending": PendingStatusText,
	"success": SuccessStatusText,
	"failed":  FailedStatusText,
}

func (text StatusText) String() string {
	return StatusTextToString[text]
}

func (text StatusText) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%v\"", text.String())), nil
}

func (text *StatusText) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	stringKey, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'pending' in this case.
	*text = StringToStatusText[stringKey]
	return nil
}

type Status struct {
	Text       StatusText `json:"text"`
	FinalClaim string     `json:"final_claim"`
}

type Prophecy struct {
	ID     string `json:"id"`
	Status Status `json:"status"`

	//WARNING: Mappings are nondeterministic in Amino,
	// an so iterating over them could result in consensus failure. New code should not iterate over the below 2 mappings.

	//This is a mapping from a claim to the list of validators that made that claim.
	ClaimValidators map[string][]types.ValAddress `json:"claim_validators"`
	//This is a mapping from a validator bech32 address to their claim
	ValidatorClaims map[string]string `json:"validator_claims"`
}

// DBProphecy is what the prophecy becomes when being saved to the database.
//  Tendermint/Amino does not support maps so we must serialize those variables into bytes.
type DBProphecy struct {
	ID              string `json:"id"`
	Status          Status `json:"status"`
	ValidatorClaims []byte `json:"validator_claims"`
}

// SerializeForDB serializes a prophecy into a DBProphecy
func (prophecy Prophecy) SerializeForDB() (DBProphecy, error) {
	validatorClaims, err := json.Marshal(prophecy.ValidatorClaims)
	if err != nil {
		return DBProphecy{}, err
	}

	return DBProphecy{
		ID:              prophecy.ID,
		Status:          prophecy.Status,
		ValidatorClaims: validatorClaims,
	}, nil
}

// DeserializeFromDB deserializes a DBProphecy into a prophecy
func (dbProphecy DBProphecy) DeserializeFromDB() (Prophecy, error) {
	var validatorClaims map[string]string
	if err := json.Unmarshal(dbProphecy.ValidatorClaims, &validatorClaims); err != nil {
		return Prophecy{}, err
	}

	var claimValidators = map[string][]types.ValAddress{}
	for addr, claim := range validatorClaims {
		valAddr, err := types.ValAddressFromBech32(addr)
		if err != nil {
			panic(fmt.Errorf("unmarshal validator address err, address=%s", addr))
		}
		claimValidators[claim] = append(claimValidators[claim], valAddr)
	}

	return Prophecy{
		ID:              dbProphecy.ID,
		Status:          dbProphecy.Status,
		ClaimValidators: claimValidators,
		ValidatorClaims: validatorClaims,
	}, nil
}

//Validate and check if it's mini token
func IsValidMiniTokenSymbol(symbol string) bool {
	return ValidateMiniTokenSymbol(symbol) == nil
}

func ValidateMiniTokenSymbol(symbol string) error {
	if len(symbol) == 0 {
		return errors.New("suffixed token symbol cannot be empty")
	}

	parts, err := splitSuffixedTokenSymbol(symbol)
	if err != nil {
		return err
	}

	symbolPart := parts[0]

	// check len without suffix
	if len(symbolPart) < MiniTokenSymbolMinLen {
		return fmt.Errorf("mini-token symbol part is too short, got %d chars", len(symbolPart))
	}
	if len(symbolPart) > MiniTokenSymbolMaxLen {
		return fmt.Errorf("mini-token symbol part is too long, got %d chars", len(symbolPart))
	}

	if !common.IsAlphaNum(symbolPart) {
		return errors.New("mini-token symbol part should be alphanumeric")
	}

	suffixPart := parts[1]

	if len(suffixPart) != MiniTokenSymbolSuffixLen {
		return fmt.Errorf("mini-token symbol suffix must be %d chars in length, got %d", MiniTokenSymbolSuffixLen, len(suffixPart))
	}

	if suffixPart[len(suffixPart)-1:] != MiniTokenSymbolMSuffix {
		return fmt.Errorf("mini-token symbol suffix must end with M")
	}

	// prohibit non-hexadecimal chars in the suffix part
	isHex, err := regexp.MatchString(fmt.Sprintf("[0-9A-F]{%d}M", MiniTokenSymbolTxHashSuffixLen), suffixPart)
	if err != nil {
		return err
	}
	if !isHex {
		return fmt.Errorf("mini-token symbol tx hash suffix must be hex with a length of %d", MiniTokenSymbolTxHashSuffixLen)
	}

	return nil
}
