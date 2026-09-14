// Package apperror — the error type and catalog for svc-transactions.
package apperror

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AppError is what every layer returns and the REST layer writes out.
// Status picks the HTTP code but never appears in the JSON.
type AppError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
	Status  int                    `json:"-"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// WithDetails returns a copy with details attached.
// Copy, not mutation — the vars below are shared across every goroutine.
func (e *AppError) WithDetails(d map[string]interface{}) *AppError {
	c := *e
	c.Details = d
	return &c
}

func New(code, message string, status int) *AppError {
	return &AppError{Code: code, Message: message, Status: status}
}

// Write renders any error as JSON. Anything that isn't an AppError becomes a
// generic 500, so no internal message reaches the client.
func Write(w http.ResponseWriter, err error) {
	appErr, ok := err.(*AppError)
	if !ok {
		appErr = Internal
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Status)
	json.NewEncoder(w).Encode(appErr)
}

var (
	// own errors
	InvalidAmount       = New("INVALID_AMOUNT", "amount must be positive", http.StatusBadRequest)
	InvalidRequestBody  = New("INVALID_REQUEST_BODY", "invalid request body", http.StatusBadRequest)
	TransactionNotFound = New("TRANSACTION_NOT_FOUND", "transaction not found", http.StatusNotFound)

	// relayed from svc-wallet — same codes, so the client sees the real reason
	// rather than a generic "wallet call failed"
	WalletNotFound     = New("WALLET_NOT_FOUND", "wallet not found", http.StatusNotFound)
	InsufficientFunds  = New("INSUFFICIENT_FUNDS", "insufficient funds", http.StatusUnprocessableEntity)
	LimitExceeded      = New("LIMIT_EXCEEDED", "amount exceeds per-transaction limit", http.StatusUnprocessableEntity)
	MaxBalanceExceeded = New("MAX_BALANCE_EXCEEDED", "resulting balance exceeds maximum", http.StatusUnprocessableEntity)
	InvalidType        = New("INVALID_TYPE", "invalid transaction type", http.StatusBadRequest)
	DuplicateRequest   = New("DUPLICATE_REQUEST", "request already processed", http.StatusConflict)

	// infrastructure
	ServiceUnavailable = New("SERVICE_UNAVAILABLE", "wallet service unavailable", http.StatusServiceUnavailable)
	Internal           = New("INTERNAL_ERROR", "internal error", http.StatusInternalServerError)
)
