package rest

import (
	"encoding/json"
	"net/http"

	"svc-wallet/internal/wallet"
	"svc-wallet/util/apperror"
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.Write(w, apperror.InvalidRequestBody)
		return
	}

	newWallet, err := c.service.CreateWallet(r.Context(), req.OwnerPhone, req.Currency, req.NationalID)
	if err != nil {
		// The service already returns an AppError, so the status comes from
		// the error itself rather than being guessed here.
		apperror.Write(w, err)
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

	foundWallet, err := c.service.GetWallet(r.Context(), id)
	if err != nil {
		apperror.Write(w, err)
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.Write(w, apperror.InvalidRequestBody)
		return
	}

	// Idempotency key. Required — without one, a retried request applies the
	// balance change a second time.
	requestID := r.Header.Get("Idempotency-Key")
	if requestID == "" {
		apperror.Write(w, apperror.InvalidRequestBody)
		return
	}

	id := r.PathValue("id")

	updatedWallet, balanceBefore, err := c.service.ApplyBalanceChange(r.Context(), id, req.Type, req.Amount, requestID)
	if err != nil {
		apperror.Write(w, err)
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
