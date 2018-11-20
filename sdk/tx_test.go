package sdk

import (
	"reflect"
	"testing"
)

func TestGetTxError(t *testing.T) {
	sdk := &SDK{
		dexAPI: &fDexAPI{},
	}

	_, err := sdk.GetTx("")
	if err == nil {
		t.Errorf("GetTx failed, expected `Error` but got %v", err)
	}
}

func TestGetTx(t *testing.T) {
	sdk := &SDK{
		dexAPI: &fDexAPI{},
	}

	tx, err := sdk.GetTx("52ECED0360605C1F3F336CA20B2C60535B0C72F0")
	if err != nil {
		t.Errorf("GetTx failed, expected no error but got %v", err)
	}

	expected := &TxResult{
		Hash: "52ECED0360605C1F3F336CA20B2C60535B0C72F0",
		Log:  "Msg 0: ",
		Data: "eyJ0eXBlIjoiZGV4L05ld09yZGVyUmVzcG9uc2UiLCJ2YWx1ZSI6eyJvcmRlcl9pZCI6ImNvc21vc2FjY2FkZHIxcTY4cGhxN3E2Znl1cDV4MjVtYWdsZjlzeGMydDRoeTQycGE2MjMtMjQwNDAyIn19",
		Code: 0,
	}

	if !reflect.DeepEqual(expected, tx) {
		t.Errorf("GetTx wrong results, expected %v but got %v", expected, tx)
	}

}

func TestPostTx(t *testing.T) {
	sdk := &SDK{
		dexAPI: &fDexAPI{},
	}

	txResult, err := sdk.PostTx([]byte("eyJ0eXBlIjoiZGV4L05ld09yZGVyUmVzcG9uc2UiLCJ2YWx1ZSI6eyJvcmRlcl9pZCI6ImNvc21vc2FjY2FkZHIxcTY4cGhxN3E2Znl1cDV4MjVtYWdsZjlzeGMydDRoeTQycGE2MjMtMjQwNDAyIn19"))
	if err != nil {
		t.Errorf("PostTx failed, expected no error but got %v", err)
	}

	expected := []*TxCommitResult{
		&TxCommitResult{
			Ok:   true,
			Code: 0,
			Data: "eyJ0eXBlIjoiZGV4L05ld09yZGVyUmVzcG9uc2UiLCJ2YWx1ZSI6eyJvcmRlcl9pZCI6ImNvc21vc2FjY2FkZHIxcTY4cGhxN3E2Znl1cDV4MjVtYWdsZjlzeGMydDRoeTQycGE2MjMtMjQwNDAyIn19",
			Log:  "logABC",
			Hash: "52ECED0360605C1F3F336CA20B2C60535B0C72F0",
		},
	}

	if !reflect.DeepEqual(expected, txResult) {
		t.Errorf("PostTx wrong results, expected %v but got %v", expected, txResult)
	}

}
