package notification

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type NotifyRequest struct {
	TransactionID string `json:"transaction_id"`
	Phone         string `json:"phone"`
	DeviceToken   string `json:"device_token"`
	SMSMessage    string `json:"sms_message"`
	PushTitle     string `json:"push_title"`
	PushBody      string `json:"push_body"`
}



func (c *Client) Notify(req NotifyRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Post(c.baseURL+"/v1/notify", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}