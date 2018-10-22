package sdk

import resty "gopkg.in/resty.v1"

// DexAPI wrapper
type DexAPI struct {
	baseURL string
}

// IDexAPI methods
type IDexAPI interface {
	Get(path string, qp map[string]string) ([]byte, error)
	Post(path string, qp map[string]string, body []byte) ([]byte, error)
}

// Get generic method
func (api *DexAPI) Get(path string, qp map[string]string) ([]byte, error) {
	resp, err := resty.R().SetQueryParams(qp).Get(api.baseURL + path)
	if err != nil {
		return nil, err
	}

	return resp.Body(), nil
}

// Post generic method
func (api *DexAPI) Post(path string, qp map[string]string, body []byte) ([]byte, error) {
	resp, err := resty.R().SetQueryParams(qp).Post(api.baseURL + path)
	if err != nil {
		return nil, err
	}

	return resp.Body(), nil
}
