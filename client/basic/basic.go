package basic

import (
	"encoding/json"
	"fmt"
	"gopkg.in/resty.v1"
	"net/http"
	"net/url"
	"time"

	"github.com/binance-chain/go-sdk/types"
	"github.com/binance-chain/go-sdk/types/tx"
	"github.com/gorilla/websocket"
)

const (
	MaxReadWaitTime = 30 * time.Second
)

type BasicClient interface {
	Get(path string, qp map[string]string) ([]byte, error)
	Post(path string, body interface{}, param map[string]string) ([]byte, error)

	GetTx(txHash string) (*tx.TxResult, error)
	PostTx(hexTx []byte, param map[string]string) ([]tx.TxCommitResult, error)
	WsGet(path string, constructMsg func([]byte) (interface{}, error), closeCh <-chan struct{}) (<-chan interface{}, error)
}

type client struct {
	baseUrl string
	apiUrl  string
}

func NewClient(baseUrl string) BasicClient {
	return &client{baseUrl: baseUrl, apiUrl: fmt.Sprintf("%s://%s", types.DefaultApiSchema, baseUrl+types.DefaultAPIVersionPrefix)}
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

func (c *client) WsGet(path string, constructMsg func([]byte) (interface{}, error), closeCh <-chan struct{}) (<-chan interface{}, error) {
	u := url.URL{Scheme: types.DefaultWSSchema, Host: c.baseUrl, Path: fmt.Sprintf("%s/%s", types.DefaultWSPrefix, path)}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	conn.SetPingHandler(nil)
	conn.SetPongHandler(
		func(string) error {
			conn.SetReadDeadline(time.Now().Add(MaxReadWaitTime))
			return nil
		})
	messages := make(chan interface{}, 0)
	finish := make(chan struct{}, 0)
	keepAliveCh := time.NewTicker(30 * time.Minute)
	pingTicker := time.NewTicker(10 * time.Second)
	go func() {
		defer conn.Close()
		defer close(messages)
		defer keepAliveCh.Stop()
		defer pingTicker.Stop()
		select {
		case <-closeCh:
			return
		case <-finish:
			return
		}
	}()
	go func() {
		writeMsg := func(m interface{}) bool {
			select {
			case <-closeCh:
				// already closed by user
				return true
			default:
			}
			messages <- m
			return false
		}
		for {
			select {
			case <-closeCh:
				conn.WriteControl(websocket.CloseMessage, nil, time.Now().Add(time.Second))
				return
			case <-keepAliveCh.C:
				conn.WriteJSON(&struct {
					Method string
				}{"keepAlive"})
			case <-pingTicker.C:
				conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second))
			default:
				response := WSResponse{}
				err := conn.ReadJSON(&response)
				if err != nil {
					if closed := writeMsg(err); !closed {
						close(finish)
					}
					return
				}
				bz, err := json.Marshal(response.Data)
				if err != nil {
					if closed := writeMsg(err); !closed {
						close(finish)
					}
					return
				}
				msg, err := constructMsg(bz)
				if err != nil {
					if closed := writeMsg(err); !closed {
						close(finish)
					}
					return
				}
				//Todo delete condition when ws do not return account and order in the same time.
				if msg != nil {
					if closed := writeMsg(msg); closed {
						return
					}
				}
			}
		}
	}()
	return messages, nil
}

type WSResponse struct {
	Stream string
	Data   interface{}
}
