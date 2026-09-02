package rest

import "net/http"

func NewRouter(controller *Controller) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/notify", controller.Notify)

	return mux
}