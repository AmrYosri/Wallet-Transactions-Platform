package transactions

import "time"

type TransactionRequest struct {
	WalletID string `json:"wallet_id"`
	Amount   int64  `json:"amount"`
	RequestID string `json:"request_id"`
}

type TransactionResponse struct {
	ID           string    `json:"id"`
	WalletID     string    `json:"wallet_id"`
	Type         string    `json:"type"`
	Amount       int64     `json:"amount"`
	Currency     string    `json:"currency"`
	BalanceBefore int64     `json:"balance_before"`
	BalanceAfter int64     `json:"balance_after"`
	CreatedAt    time.Time `json:"created_at"`
}

type TransactionListResponse struct {
	WalletID     string                `json:"wallet_id"`
	Transactions []TransactionResponse `json:"transactions"`
}
