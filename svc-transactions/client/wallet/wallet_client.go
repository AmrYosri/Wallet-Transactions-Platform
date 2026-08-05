package wallet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

type BalanceChangeRequest struct {
	Type   string `json:"type"`
	Amount int64  `json:"amount"`
}

type BalanceChangeResponse struct {
	WalletID      string `json:"wallet_id"`
	Type          string `json:"type"`
	Amount        int64  `json:"amount"`
	BalanceBefore int64  `json:"balance_before"`
	BalanceAfter  int64  `json:"balance_after"`
}

func (c *Client) ApplyBalanceChange(ctx context.Context, walletID string, changeType string, amount int64) (*BalanceChangeResponse, error) {
	reqBody := BalanceChangeRequest{Type: changeType, Amount: amount}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/v1/wallets/%s/balance", c.baseURL, walletID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wallet service returned status %d", resp.StatusCode)
	}

	var result BalanceChangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
