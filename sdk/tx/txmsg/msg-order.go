package txmsg

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Order operations
const (
	NewOrder    = "orderNew"
	CancelOrder = "orderCancel"
)

// OrderSide /TimeInForce /OrderType are const, following FIX protocol convention
// Used as Enum
var OrderSide = struct {
	BUY  int8
	SELL int8
}{1, 2}

var sideNames = map[string]int8{
	"BUY":  1,
	"SELL": 2,
}

// IToSide conversion
func IToSide(side int8) string {
	switch side {
	case OrderSide.BUY:
		return "BUY"
	case OrderSide.SELL:
		return "SELL"
	default:
		return "UNKNOWN"
	}
}

// GenerateOrderID generates an order ID
func GenerateOrderID(sequence int64, from AccAddress) string {
	id := fmt.Sprintf("%s-%d", from.String(), sequence)
	return id
}

// IsValidSide validates that a side is valid and supported by the matching engine
func IsValidSide(side int8) bool {
	switch side {
	case OrderSide.BUY, OrderSide.SELL:
		return true
	default:
		return false
	}
}

// SideStringToSideCode converts a string like "BUY" to its internal side code
func SideStringToSideCode(side string) (int8, error) {
	upperSide := strings.ToUpper(side)
	if val, ok := sideNames[upperSide]; ok {
		return val, nil
	}
	return -1, errors.New("side `" + upperSide + "` not found or supported")
}

const (
	_           int8 = iota
	orderMarket int8 = iota
	orderLimit  int8 = iota
)

// OrderType is an enum of order type options supported by the matching engine
var OrderType = struct {
	LIMIT  int8
	MARKET int8
}{orderLimit, orderMarket}

// IToOrderType conversion
func IToOrderType(tpe int8) string {
	switch tpe {
	case OrderType.LIMIT:
		return "LIMIT"
	case OrderType.MARKET:
		return "MARKET"
	default:
		return "UNKNOWN"
	}
}

// IsValidOrderType validates that an order type is valid and supported by the matching engine
func IsValidOrderType(ot int8) bool {
	switch ot {
	case OrderType.LIMIT: // only allow LIMIT for now.
		return true
	default:
		return false
	}
}

const (
	_      int8 = iota
	tifGTC int8 = iota
	_      int8 = iota
	tifIOC int8 = iota
)

// TimeInForce is an enum of TIF (Time in Force) options supported by the matching engine
var TimeInForce = struct {
	GTC int8
	IOC int8
}{tifGTC, tifIOC}

var timeInForceNames = map[string]int8{
	"GTC": tifGTC,
	"IOC": tifIOC,
}

// IsValidTimeInForce validates that a tif code is correct
func IsValidTimeInForce(tif int8) bool {
	switch tif {
	case TimeInForce.GTC, TimeInForce.IOC:
		return true
	default:
		return false
	}
}

// IToTimeInForce conversion
func IToTimeInForce(tif int8) string {
	switch tif {
	case TimeInForce.GTC:
		return "GTC"
	case TimeInForce.IOC:
		return "IOC"
	default:
		return "UNKNOWN"
	}
}

// TifStringToTifCode converts a string like "GTC" to its internal tif code
func TifStringToTifCode(tif string) (int8, error) {
	upperTif := strings.ToUpper(tif)
	if val, ok := timeInForceNames[upperTif]; ok {
		return val, nil
	}
	return -1, errors.New("tif `" + upperTif + "` not found or supported")
}

// CreateOrderMsg def
type CreateOrderMsg struct {
	Sender      AccAddress `json:"sender"`
	ID          string     `json:"id"`
	Symbol      string     `json:"symbol"`
	OrderType   int8       `json:"ordertype"`
	OrderSide   int8       `json:"side"`
	Price       int64      `json:"price"`
	Quantity    int64      `json:"quantity"`
	TimeInForce int8       `json:"timeinforce"`
}

// NewCreateOrderMsg constructs a new CreateOrderMsg
func NewCreateOrderMsg(sender AccAddress, id string, side int8, symbol string, price int64, qty int64) CreateOrderMsg {
	return CreateOrderMsg{
		Sender:      sender,
		ID:          id,
		Symbol:      symbol,
		OrderType:   OrderType.LIMIT, // default
		OrderSide:   side,
		Price:       price,
		Quantity:    qty,
		TimeInForce: TimeInForce.GTC, // default
	}
}

// Type is part of Msg interface
func (msg CreateOrderMsg) Type() string { return NewOrder }

// Get is part of Msg interface
func (msg CreateOrderMsg) Get(key interface{}) (value interface{}) { return nil }

// GetSigners is part of Msg interface
func (msg CreateOrderMsg) GetSigners() []AccAddress { return []AccAddress{msg.Sender} }

// String is part of Msg interface
func (msg CreateOrderMsg) String() string {
	return fmt.Sprintf("CreateOrderMsg{Sender: %v, Id: %v, Symbol: %v, OrderSide: %v, Price: %v, Qty: %v}", msg.Sender, msg.ID, msg.Symbol, msg.OrderSide, msg.Price, msg.Quantity)
}

// GetSignBytes - Get the bytes for the message signer to sign on
func (msg CreateOrderMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// ValidateBasic is used to quickly disqualify obviously invalid messages quickly
func (msg CreateOrderMsg) ValidateBasic() error {
	if len(msg.Sender) == 0 {
		return fmt.Errorf("ErrUnknownAddress %s", msg.Sender.String())
	}

	// `-` is required in the compound order id: <address>-<sequence>
	if len(msg.ID) == 0 || !strings.Contains(msg.ID, "-") {
		return fmt.Errorf("Invalid order ID:%s", msg.ID)
	}

	if msg.Quantity <= 0 {
		return fmt.Errorf("Invalid order Quantity, Zero/Negative Number:%d", msg.Quantity)
	}

	if msg.Price <= 0 {
		return fmt.Errorf("Invalid order Price, Zero/Negative Number:%d", msg.Price)
	}

	if !IsValidOrderType(msg.OrderType) {
		return fmt.Errorf("Invalid order type:%d", msg.OrderType)
	}

	if !IsValidSide(msg.OrderSide) {
		return fmt.Errorf("Invalid side:%d", msg.OrderSide)
	}

	if !IsValidTimeInForce(msg.TimeInForce) {
		return fmt.Errorf("Invalid TimeInForce:%d", msg.TimeInForce)
	}

	return nil
}

// NewCancelOrderMsg constructs a new CancelOrderMsg
func NewCancelOrderMsg(sender AccAddress, symbol, id, refID string) CancelOrderMsg {
	return CancelOrderMsg{
		Sender: sender,
		Symbol: symbol,
		ID:     id,
		RefID:  refID,
	}
}

// CancelOrderMsg represents a message to cancel an open order
type CancelOrderMsg struct {
	Sender AccAddress
	Symbol string `json:"symbol"`
	ID     string `json:"id"`
	RefID  string `json:"refid"`
}

// Type is part of Msg interface
func (msg CancelOrderMsg) Type() string { return CancelOrder }

// Get is part of Msg interface
func (msg CancelOrderMsg) Get(key interface{}) (value interface{}) { return nil }

// GetSigners is part of Msg interface
func (msg CancelOrderMsg) GetSigners() []AccAddress { return []AccAddress{msg.Sender} }

// String is part of Msg interface
func (msg CancelOrderMsg) String() string {
	return fmt.Sprintf("CancelOrderMsg{Sender: %v}", msg.Sender)
}

// GetSignBytes - Get the bytes for the message signer to sign on
func (msg CancelOrderMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// ValidateBasic is used to quickly disqualify obviously invalid messages quickly
func (msg CancelOrderMsg) ValidateBasic() error {
	if len(msg.Sender) == 0 {
		return fmt.Errorf("ErrUnknownAddress %s", msg.Sender.String())
	}

	if len(msg.ID) == 0 || !strings.Contains(msg.ID, "-") {
		return fmt.Errorf("Invalid order ID:%s", msg.ID)
	}

	return nil
}
