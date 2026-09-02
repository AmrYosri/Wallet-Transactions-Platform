package rest

import (
	"encoding/json"
	"net/http"

	"svc-notifications/internal/notification"
)

type Controller struct {
	service *notification.Service
}

func NewController(service *notification.Service) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) Notify(w http.ResponseWriter, r *http.Request) {
	var req notification.NotificationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	resp, err := c.service.Notify(ctx, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}