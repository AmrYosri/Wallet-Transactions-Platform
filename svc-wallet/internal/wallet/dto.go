package wallet


type CreateWalletRequest struct {
	OwnerPhone string `json:"owner_phone"`
	Currency   string `json:"currency"`
}

type WalletResponse struct {
	ID         string `json:"id"`
	OwnerPhone string `json:"owner_phone"`
	Currency   string `json:"currency"`
	Balance    int64  `json:"balance"`
	Status     string `json:"status"`
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