package types

type MarketDepth struct {
	Bids   [][]string `json:"bids"` // "bids": [ [ "0.0024", "10" ] ]
	Asks   [][]string `json:"asks"` // "asks": [ [ "0.0024", "10" ] ]
	Height int64      `json:"height"`
}
