package user


type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	NationalID string `json:"national_id"`

}

type UserResponse struct {
	ID string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	NationalID string `json:"national_id"`
	PhoneNumbers []string `json:"phone_numbers"`
	Status string `json:"status"`

}