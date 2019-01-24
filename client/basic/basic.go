package basic

import (
	"encoding/json"
	"fmt"
	"gopkg.in/resty.v1"
	"net/http"

	"github.com/binance-chain/go-sdk/types/tx"
)

type BasicClient interface {
	Get(path string, qp map[string]string) ([]byte, error)
	Post(path string, body interface{}, param map[string]string) ([]byte, error)

	GetTx(txHash string) (*tx.TxResult, error)
	PostTx(hexTx []byte, param map[string]string) ([]tx.TxCommitResult, error)
}

type client struct {
	apiUrl string
}

func NewClient(apiUrl string) BasicClient {
	return &client{apiUrl}
}

func (c *client) Get(path string, qp map[string]string) ([]byte, error) {
	resp, err := resty.R().SetQueryParams(qp).Get(c.apiUrl + path)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() >= http.StatusMultipleChoices || resp.StatusCode() < http.StatusOK {
		err = fmt.Errorf("bad response, status code %d, response: %s", resp.StatusCode(), string(resp.Body()))
	}
	return resp.Body(), err
}

// Post generic method
func (c *client) Post(path string, body interface{}, param map[string]string) ([]byte, error) {
	resp, err := resty.R().
		SetHeader("Content-Type", "text/plain").
		SetBody(body).
		SetQueryParams(param).
		Post(c.apiUrl + path)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode() >= http.StatusMultipleChoices {
		err = fmt.Errorf("bad response, status code %d, response: %s", resp.StatusCode(), string(resp.Body()))
	}
	return resp.Body(), err
}

// GetTx returns transaction details
func (c *client) GetTx(txHash string) (*tx.TxResult, error) {
	if txHash == "" {
		return nil, fmt.Errorf("Invalid tx hash %s ", txHash)
	}

	qp := map[string]string{}
	resp, err := c.Get("/tx/"+txHash, qp)
	if err != nil {
		return nil, err
	}

	var txResult tx.TxResult
	if err := json.Unmarshal(resp, &txResult); err != nil {
		return nil, err
	}

	return &txResult, nil
}

// PostTx returns transaction details
func (c *client) PostTx(hexTx []byte, param map[string]string) ([]tx.TxCommitResult, error) {
	if len(hexTx) == 0 {
		return nil, fmt.Errorf("Invalid tx  %s", hexTx)
	}

	body := hexTx
	resp, err := c.Post("/broadcast", body, param)
	if err != nil {
		return nil, err
	}
	txResult := make([]tx.TxCommitResult, 0)
	if err := json.Unmarshal(resp, &txResult); err != nil {
		return nil, err
	}

	return txResult, nil
}
