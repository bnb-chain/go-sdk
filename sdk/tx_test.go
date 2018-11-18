package sdk

import (
	"reflect"
	"testing"
)

func TestTxError(t *testing.T) {
	sdk := &SDK{
		dexAPI: &fDexAPI{},
	}

	_, err := sdk.GetTx("")
	if err == nil {
		t.Errorf("GetTx failed, expected `Error` but got %v", err)
	}
}

func TestTx(t *testing.T) {
	sdk := &SDK{
		dexAPI: &fDexAPI{},
	}

	tx, err := sdk.GetTx("52ECED0360605C1F3F336CA20B2C60535B0C72F0")
	if err != nil {
		t.Errorf("GetTx failed, expected no error but got %v", err)
	}

	expected := &Tx{
		Hash: "52ECED0360605C1F3F336CA20B2C60535B0C72F0",
		Log:  "Msg 0: ",
		Data: "eyJ0eXBlIjoiZGV4L05ld09yZGVyUmVzcG9uc2UiLCJ2YWx1ZSI6eyJvcmRlcl9pZCI6ImNvc21vc2FjY2FkZHIxcTY4cGhxN3E2Znl1cDV4MjVtYWdsZjlzeGMydDRoeTQycGE2MjMtMjQwNDAyIn19",
		Code: 0,
	}

	if !reflect.DeepEqual(expected, tx) {
		t.Errorf("GetTx wrong results, expected %v but got %v", expected, tx)
	}

}
