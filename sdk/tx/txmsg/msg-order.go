package txmsg

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	NewOrder    = "orderNew"
	CancelOrder = "orderCancel"
)

// OrderSide/TimeInForce/OrderType are const, following FIX protocol convention
// Used as Enum
var OrderSide = struct {
	BUY  int8
	SELL int8
}{1, 2}

var sideNames = map[string]int8{
	"BUY":  1,
	"SELL": 2,
}

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

type NewOrderMsg struct {
	Sender      AccAddress `json:"sender"`
	Id          string     `json:"id"`
	Symbol      string     `json:"symbol"`
	OrderType   int8       `json:"ordertype"`
	OrderSide   int8       `json:"side"`
	Price       int64      `json:"price"`
	Quantity    int64      `json:"quantity"`
	TimeInForce int8       `json:"timeinforce"`
}

// NewNewOrderMsg constructs a new NewOrderMsg
func NewNewOrderMsg(sender AccAddress, id string, side int8, symbol string, price int64, qty int64) NewOrderMsg {
	return NewOrderMsg{
		Sender:      sender,
		Id:          id,
		Symbol:      symbol,
		OrderType:   OrderType.LIMIT, // default
		OrderSide:   side,
		Price:       price,
		Quantity:    qty,
		TimeInForce: TimeInForce.GTC, // default
	}
}

// var _ sdk.Msg = NewOrderMsg{}

// nolint
func (msg NewOrderMsg) Type() string                            { return NewOrder }
func (msg NewOrderMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg NewOrderMsg) GetSigners() []AccAddress                { return []AccAddress{msg.Sender} }
func (msg NewOrderMsg) String() string {
	return fmt.Sprintf("NewOrderMsg{Sender: %v, Id: %v, Symbol: %v, OrderSide: %v, Price: %v, Qty: %v}", msg.Sender, msg.Id, msg.Symbol, msg.OrderSide, msg.Price, msg.Quantity)
}

// NewCancelOrderMsg constructs a new CancelOrderMsg
func NewCancelOrderMsg(sender AccAddress, symbol, id, refId string) CancelOrderMsg {
	return CancelOrderMsg{
		Sender: sender,
		Symbol: symbol,
		Id:     id,
		RefId:  refId,
	}
}

// CancelOrderMsg represents a message to cancel an open order
type CancelOrderMsg struct {
	Sender AccAddress
	Symbol string `json:"symbol"`
	Id     string `json:"id"`
	RefId  string `json:"refid"`
}

// var _ sdk.Msg = CancelOrderMsg{}

// nolint
func (msg CancelOrderMsg) Type() string                            { return CancelOrder }
func (msg CancelOrderMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg CancelOrderMsg) GetSigners() []AccAddress                { return []AccAddress{msg.Sender} }
func (msg CancelOrderMsg) String() string {
	return fmt.Sprintf("CancelOrderMsg{Sender: %v}", msg.Sender)
}

// GetSignBytes - Get the bytes for the message signer to sign on
func (msg NewOrderMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
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
func (msg NewOrderMsg) ValidateBasic() error {
	// if len(msg.Sender) == 0 {
	// 	return sdk.ErrUnknownAddress(msg.Sender.String()).TraceSDK("")
	// }
	// // `-` is required in the compound order id: <address>-<sequence>
	// if len(msg.Id) == 0 || !strings.Contains(msg.Id, "-") {
	// 	return types.ErrInvalidOrderParam("Id", fmt.Sprintf("Invalid order ID:%s", msg.Id))
	// }
	// if msg.Quantity <= 0 {
	// 	return types.ErrInvalidOrderParam("Quantity", fmt.Sprintf("Zero/Negative Number:%d", msg.Quantity))
	// }
	// if msg.Price <= 0 {
	// 	return types.ErrInvalidOrderParam("Price", fmt.Sprintf("Zero/Negative Number:%d", msg.Quantity))
	// }
	// if !IsValidOrderType(msg.OrderType) {
	// 	return types.ErrInvalidOrderParam("OrderType", fmt.Sprintf("Invalid order type:%d", msg.OrderType))
	// }
	// if !IsValidSide(msg.OrderSide) {
	// 	return types.ErrInvalidOrderParam("OrderSide", fmt.Sprintf("Invalid side:%d", msg.OrderSide))
	// }
	// if !IsValidTimeInForce(msg.TimeInForce) {
	// 	return types.ErrInvalidOrderParam("TimeInForce", fmt.Sprintf("Invalid TimeInForce:%d", msg.TimeInForce))
	// }

	return nil
}

// ValidateBasic is used to quickly disqualify obviously invalid messages quickly
func (msg CancelOrderMsg) ValidateBasic() error {
	// if len(msg.Sender) == 0 {
	// 	return sdk.ErrUnknownAddress(msg.Sender.String()).TraceSDK("")
	// }
	// if len(msg.Id) == 0 || !strings.Contains(msg.Id, "-") {
	// 	return types.ErrInvalidOrderParam("Id", fmt.Sprintf("Invalid order ID:%s", msg.Id))
	// }
	return nil
}
