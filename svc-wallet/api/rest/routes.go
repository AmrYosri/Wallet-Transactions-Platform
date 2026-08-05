package rest

import "net/http"

func NewRouter(controller *Controller) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/wallets", controller.CreateWallet)
	mux.HandleFunc("GET /v1/wallets/{id}", controller.GetWallet)
	mux.HandleFunc("PATCH /v1/wallets/{id}/balance", controller.ChangeBalance)

	return mux
}
