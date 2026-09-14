package apperror

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AppError is the error shape every layer returns and the REST layer writes out.
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
// Copy, not mutation — the vars below are shared across every goroutine,
// so writing to one would be a data race.
func (e *AppError) WithDetails(d map[string]interface{}) *AppError {
	c := *e
	c.Details = d
	return &c
}

func New(code, message string, status int) *AppError {
	return &AppError{Code: code, Message: message, Status: status}
}

// Write renders any error as JSON. Anything that isn't an AppError becomes
// a generic 500, so no driver message ever reaches the client.
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
	// not found
	WalletNotFound = New("WALLET_NOT_FOUND", "wallet not found", http.StatusNotFound)
	UserNotFound   = New("USER_NOT_FOUND", "user not found", http.StatusNotFound)

	// conflict
	WalletAlreadyExists = New("WALLET_ALREADY_EXISTS", "wallet already exists for this user and currency", http.StatusConflict)
	DuplicateRequest    = New("DUPLICATE_REQUEST", "request already processed", http.StatusConflict)

	// business rules
	InsufficientFunds  = New("INSUFFICIENT_FUNDS", "insufficient funds", http.StatusUnprocessableEntity)
	LimitExceeded      = New("LIMIT_EXCEEDED", "amount exceeds per-transaction limit", http.StatusUnprocessableEntity)
	MaxBalanceExceeded = New("MAX_BALANCE_EXCEEDED", "resulting balance exceeds maximum", http.StatusUnprocessableEntity)

	// bad input
	InvalidType        = New("INVALID_TYPE", "invalid transaction type", http.StatusBadRequest)
	InvalidRequestBody = New("INVALID_REQUEST_BODY", "invalid request body", http.StatusBadRequest)

	// fallback
	Internal = New("INTERNAL_ERROR", "internal error", http.StatusInternalServerError)
)
