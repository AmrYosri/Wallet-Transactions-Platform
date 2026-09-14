package rest

import (
	"encoding/json"
	"net/http"

	"svc-transactions/internal/transactions"
	"svc-transactions/util/apperror"
)

type Controller struct {
	service *transactions.Service
}

func NewController(service *transactions.Service) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) Deposit(w http.ResponseWriter, r *http.Request) {
	c.applyTransaction(w, r, "deposit")
}

func (c *Controller) Withdraw(w http.ResponseWriter, r *http.Request) {
	c.applyTransaction(w, r, "withdraw")
}

// applyTransaction is the shared body of Deposit and Withdraw — they only
// differ by the type string.
func (c *Controller) applyTransaction(w http.ResponseWriter, r *http.Request, changeType string) {
	var req transactions.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.Write(w, apperror.InvalidRequestBody)
		return
	}

	// The idempotency key must come from the client. Generating one here would
	// look like it works but silently break retries: every retry would get a
	// fresh key and be applied as a new transaction.
	if req.RequestID == "" {
		apperror.Write(w, apperror.InvalidRequestBody)
		return
	}

	result, err := c.service.ApplyTransaction(r.Context(), req.WalletID, changeType, req.Amount, req.RequestID)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	resp := transactions.TransactionResponse{
		ID:            result.ID.Hex(),
		WalletID:      result.WalletID,
		Type:          result.Type,
		Amount:        result.Amount,
		Currency:      result.Currency,
		BalanceBefore: result.BalanceBefore,
		BalanceAfter:  result.BalanceAfter,
		CreatedAt:     result.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (c *Controller) GetTransactions(w http.ResponseWriter, r *http.Request) {
	walletID := r.URL.Query().Get("wallet_id")

	found, err := c.service.GetTransactions(r.Context(), walletID)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	responseList := make([]transactions.TransactionResponse, 0, len(found))
	for _, t := range found {
		responseList = append(responseList, transactions.TransactionResponse{
			ID:            t.ID.Hex(),
			WalletID:      t.WalletID,
			Type:          t.Type,
			Amount:        t.Amount,
			Currency:      t.Currency,
			BalanceBefore: t.BalanceBefore,
			BalanceAfter:  t.BalanceAfter,
			CreatedAt:     t.CreatedAt,
		})
	}

	resp := transactions.TransactionListResponse{
		WalletID:     walletID,
		Transactions: responseList,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
