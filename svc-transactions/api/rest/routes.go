package rest

import "net/http"

func NewRouter(controller *Controller) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/transactions/deposit", controller.Deposit)
	mux.HandleFunc("POST /v1/transactions/withdraw", controller.Withdraw)
	mux.HandleFunc("GET /v1/transactions", controller.GetTransactions)

	return mux
}
