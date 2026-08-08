package rest

import "net/http"

func NewRouter (controller *Controller) *http.ServeMux{
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/users" , controller.CreateUser)
	mux.HandleFunc("GET /v1/users/{national_id}",controller.GetUser)
	return mux
}