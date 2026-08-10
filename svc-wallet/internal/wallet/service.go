package wallet

import (
	"context"
	"errors"
	"time"

	"svc-wallet/client/user"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo       *Repository
	userClient *user.Client
}

func NewService(repo *Repository, userClient *user.Client) *Service {
	return &Service{
		repo:       repo,
		userClient: userClient,
	}
}
func (s *Service) CreateWallet(ctx context.Context, ownerPhone, currency, nationalID string) (*Wallet, error) {
	_, err := s.userClient.GetUserByNationalID(ctx, nationalID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	wallet := &Wallet{
		OwnerPhone: ownerPhone,
		NationalID: nationalID,
		Currency:   currency,
		Balance:    0,
		Status:     "active",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = s.repo.Create(ctx, wallet)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}
func (s *Service) GetWallet(ctx context.Context, id string) (*Wallet, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	wallet, err := s.repo.FindByID(ctx, objID)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *Service) ApplyBalanceChange(ctx context.Context, id string, changeType string, amount int64) (*Wallet, int64, error) {

	if amount > 60000 {
		return nil, 0, errors.New("exceeded the limit per transaction")
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, 0, err
	}

	wallet, err := s.repo.FindByID(ctx, objID)
	if err != nil {
		return nil, 0, err
	}

	balanceBefore := wallet.Balance
	var newBalance int64

	switch changeType {
	case "withdraw":
		if wallet.Balance-amount < 0 {
			return nil, 0, errors.New("insufficient funds")
		}
		newBalance = wallet.Balance - amount

	case "deposit":
		if wallet.Balance+amount > 400000 {
			return nil, 0, errors.New("this deposit would push the balance over the limit")
		}
		newBalance = wallet.Balance + amount

	default:
		return nil, 0, errors.New("invalid type")
	}

	err = s.repo.UpdateBalance(ctx, objID, newBalance)
	if err != nil {
		return nil, 0, err
	}

	wallet.Balance = newBalance

	return wallet, balanceBefore, nil
}
