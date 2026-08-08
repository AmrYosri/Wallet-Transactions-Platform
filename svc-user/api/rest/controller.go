package rest

import (
	"encoding/json"
	"net/http"
	"svc-user/internal/user"
)

type Controller struct {
	service *user.Service
}

func NewController(service *user.Service) *Controller {
	return &Controller{
		service: service,
	}
}


func (c *Controller) CreateUser(w http.ResponseWriter, r *http.Request) {

	var req user.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil{
		http.Error(w,"invalid request body" , http.StatusBadRequest)
		return
	}

	ctx :=r.Context()
	newUser , err := c.service.CreateUser(ctx , req.FirstName , req.LastName ,req.NationalID)
	if err != nil {
		http.Error(w,"failed to create user" ,http.StatusInternalServerError)
		return
	}

	resp := user.UserResponse{
		ID: newUser.ID.Hex(),
		FirstName: newUser.FirstName,
		LastName: newUser.LastName,
		NationalID: newUser.NationalID,
		Status: newUser.Status,
		PhoneNumbers: newUser.PhoneNumbers,
	}
	w.Header().Set("Content-Type" , "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
	
}

func (c *Controller) GetUser (w http.ResponseWriter , r *http.Request){
	nationalID := r.PathValue("national_id")
	ctx := r.Context()

	foundUser , err := c.service.GetUserByNationalID(ctx,nationalID)
	if err != nil{
		http.Error(w,"not valid national id " , http.StatusNotFound)
		return
	}
	resp := user.UserResponse{
		ID: foundUser.ID.Hex(),
		FirstName: foundUser.FirstName,
		LastName: foundUser.LastName,
		NationalID: foundUser.NationalID,
		PhoneNumbers: foundUser.PhoneNumbers,
		Status: foundUser.Status,

	}

	w.Header().Set("Content-Type" , "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}