package tx

type StdFee struct {
	Amount Coins `json:"amount"`
	Gas    int64 `json:"gas"`
}

func NewStdFee(gas int64, amount ...Coin) StdFee {
	return StdFee{
		Amount: amount,
		Gas:    gas,
	}
}

// fee bytes for signing later
func (fee StdFee) Bytes() []byte {
	if len(fee.Amount) == 0 {
		fee.Amount = Coins{}
	}
	bz, err := Cdc.MarshalJSON(fee)
	if err != nil {
		panic(err)
	}
	return bz
}
