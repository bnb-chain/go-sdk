package tx

const (
	CodeOk int32 = 0
)

// TxResult def
type TxResult struct {
	Hash string `json:"hash"`
	Log  string `json:"log"`
	Data string `json:"data"`
	Code int32  `json:"code"`
}

// TxCommitResult for POST tx results
type TxCommitResult struct {
	Ok   bool   `json:"ok"`
	Log  string `json:"log"`
	Hash string `json:"hash"`
	Code int32  `json:"code"`
	Data string `json:"data"`
}
