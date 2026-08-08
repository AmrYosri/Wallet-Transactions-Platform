package user

import (
	"context"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateUser(ctx context.Context, firstName, lastName, nationalID string) (*User, error) {
	user := &User{
		FirstName: firstName,
		LastName: lastName,
		NationalID: nationalID,
		Status: "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err :=s.repo.Create(ctx,user)
	if err != nil {
		return  nil, err
	}
	return user,nil
}


func (s *Service) GetUserByNationalID(ctx context.Context, nationalID string) (*User, error) {
	user ,err := s.repo.FindByNationalID(ctx,nationalID)
	if err != nil {
		return nil , err
	}
	return user ,nil
}

func (s *Service) AddPhoneNumber(ctx context.Context, nationalID string, phone string) error {
	err := s.repo.AddPhoneNumber(ctx,nationalID,phone)
	if err != nil{
		return err
	}
	return nil 
}