package rest

import (
	"encoding/json"
	"net/http"

	"svc-wallet/internal/wallet"
)

type Controller struct {
	service *wallet.Service
}

func NewController(service *wallet.Service) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) CreateWallet(w http.ResponseWriter, r *http.Request) {
	var req wallet.CreateWalletRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	newWallet, err := c.service.CreateWallet(ctx, req.OwnerPhone, req.Currency)
	if err != nil {
		http.Error(w, "failed to create wallet", http.StatusInternalServerError)
		return
	}

	resp := wallet.WalletResponse{
		ID:         newWallet.ID.Hex(),
		OwnerPhone: newWallet.OwnerPhone,
		Currency:   newWallet.Currency,
		Balance:    newWallet.Balance,
		Status:     newWallet.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (c *Controller) GetWallet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ctx := r.Context()

	foundWallet, err := c.service.GetWallet(ctx, id)
	if err != nil {
		http.Error(w, "not valid id", http.StatusNotFound)
		return
	}
	resp := wallet.WalletResponse{
		ID:         foundWallet.ID.Hex(),
		OwnerPhone: foundWallet.OwnerPhone,
		Balance:    foundWallet.Balance,
		Currency:   foundWallet.Currency,
		Status:     foundWallet.Status,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
func (c *Controller) ChangeBalance(w http.ResponseWriter, r *http.Request) {
	var req wallet.BalanceChangeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	id := r.PathValue("id")
	ctx := r.Context()
	updatedWallet, balanceBefore, err := c.service.ApplyBalanceChange(ctx, id, req.Type, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp := wallet.BalanceChangeResponse{
		WalletID:      updatedWallet.ID.Hex(),
		Type:          req.Type,
		Amount:        req.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  updatedWallet.Balance,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}
