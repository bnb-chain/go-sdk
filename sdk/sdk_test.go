package sdk

import "testing"

func TestNewSDK(t *testing.T) {
	_, err := NewSDK("")

	if err == nil {
		t.Errorf("NewSDK failed, expected `Error` but got %v", err)
	}

	sdk, err2 := NewSDK("http://localhost")
	if err2 != nil || sdk == nil {
		t.Errorf("NewSDK failed, expected `sdk` instance but got %v", err2)
	}
}
