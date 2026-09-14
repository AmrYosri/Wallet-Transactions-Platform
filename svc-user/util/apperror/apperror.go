package apperror

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AppError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
	Status  int                    `json:"-"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// WithDetails returns a copy. The vars below are shared across goroutines —
// mutating one is a data race.
func (e *AppError) WithDetails(d map[string]interface{}) *AppError {
	c := *e
	c.Details = d
	return &c
}

func New(code, message string, status int) *AppError {
	return &AppError{Code: code, Message: message, Status: status}
}

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
	UserNotFound       = New("USER_NOT_FOUND", "user not found", http.StatusNotFound)
	UserAlreadyExists  = New("USER_ALREADY_EXISTS", "user already exists", http.StatusConflict)
	InvalidRequestBody = New("INVALID_REQUEST_BODY", "invalid request body", http.StatusBadRequest)
	Internal           = New("INTERNAL_ERROR", "internal error", http.StatusInternalServerError)
)