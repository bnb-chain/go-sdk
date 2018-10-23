package sdk

import "fmt"

// SDK wrapper
type SDK struct {
	dexAPI IDexAPI
}

// NewSDK init
func NewSDK(baseURL string) (*SDK, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("Invalid baseURL %s", baseURL)
	}

	return &SDK{
		dexAPI: &DexAPI{baseURL},
	}, nil
}

// ToMapStrStr conversion
func ToMapStrStr(m map[string]interface{}) map[string]string {
	mStrStr := make(map[string]string)
	for k, v := range m {
		strValue := fmt.Sprintf("%v", v)

		// skip empty values
		if strValue != "" && strValue != "0" {
			mStrStr[k] = strValue
		}
	}
	return mStrStr
}
