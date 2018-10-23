package sdk

import (
	"reflect"
	"testing"
)

func TestGetTime(t *testing.T) {
	sdk := &SDK{
		dexAPI: &fDexAPI{},
	}

	depth, err := sdk.GetTime()
	if err != nil {
		t.Errorf("GetTime failed, expected no error but got %v", err)
	}

	expected := "14:20:00T"

	if !reflect.DeepEqual(expected, depth) {
		t.Errorf("GetTime wrong results, expected %v but got %v", expected, depth)
	}

}
