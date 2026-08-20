package rest

import (
	"encoding/json"
	"net/http"
	"svc-transactions/internal/transactions"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
	var req transactions.TransactionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.RequestID == "" {
    req.RequestID = primitive.NewObjectID().Hex()
	}

	ctx := r.Context()

	result, err := c.service.ApplyTransaction(ctx, req.WalletID, "deposit", req.Amount, req.RequestID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

func (c *Controller) Withdraw(w http.ResponseWriter, r *http.Request) {
	var req transactions.TransactionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	
	
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.RequestID == "" {
    req.RequestID = primitive.NewObjectID().Hex()
	}

	ctx := r.Context()

	result, err := c.service.ApplyTransaction(ctx, req.WalletID, "withdraw", req.Amount, req.RequestID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	ctx := r.Context()

	transactions_, err := c.service.GetTransactions(ctx, walletID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var responseList []transactions.TransactionResponse
	for _, t := range transactions_ {
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
