package tx

type Option func(*StdSignMsg) *StdSignMsg

func WithSource(source int64) Option {
	return func(txMsg *StdSignMsg) *StdSignMsg {
		txMsg.Source = source
		return txMsg
	}
}

func WithMemo(memo string) Option {
	return func(txMsg *StdSignMsg) *StdSignMsg {
		txMsg.Memo = memo
		return txMsg
	}
}

func WithAcNumAndSequence(accountNum, seq int64) Option {
	return func(txMsg *StdSignMsg) *StdSignMsg {
		txMsg.Sequence = seq
		txMsg.AccountNumber = accountNum
		return txMsg
	}
}

func WithChainID(id string) Option {
	return func(txMsg *StdSignMsg) *StdSignMsg {
		txMsg.ChainID = id
		return txMsg
	}
}
