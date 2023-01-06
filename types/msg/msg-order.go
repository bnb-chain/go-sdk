package msg

import (
	"github.com/bnb-chain/node/plugins/dex/order"
)

const (
	RouteNewOrder    = order.RouteNewOrder
	RouteCancelOrder = order.RouteCancelOrder
)

var (
	OrderSide            = order.Side
	SideStringToSideCode = order.SideStringToSideCode
	GenerateOrderID      = order.GenerateOrderID
	IsValidSide          = order.IsValidSide
)

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

var (
	OrderType          = order.OrderType
	IsValidOrderType   = order.IsValidOrderType
	TimeInForce        = order.TimeInForce
	IsValidTimeInForce = order.IsValidTimeInForce
	TifStringToTifCode = order.TifStringToTifCode
)

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

// IToTimeInForce conversion
func IToTimeInForce(tif int8) string {
	switch tif {
	case TimeInForce.GTE:
		return "GTE"
	case TimeInForce.IOC:
		return "IOC"
	default:
		return "UNKNOWN"
	}
}

type (
	CreateOrderMsg = order.NewOrderMsg
	CancelOrderMsg = order.CancelOrderMsg
)

var (
	NewCreateOrderMsg = order.NewNewOrderMsg
	NewCancelOrderMsg = order.NewCancelOrderMsg
)
