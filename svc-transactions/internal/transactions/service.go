package transactions

import (
	"context"
	"errors"
	"svc-transactions/client/wallet"
	"time"
)

type Service struct {
	repo         *Repository
	walletClient *wallet.Client
}

func NewService(repo *Repository, walletClient *wallet.Client) *Service {
	return &Service{
		repo:         repo,
		walletClient: walletClient,
	}
}
func (s *Service) ApplyTransaction(ctx context.Context, walletID string, changeType string, amount int64) (*Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	result, err := s.walletClient.ApplyBalanceChange(ctx, walletID, changeType, amount)
	if err != nil {
		return nil, err
	}

	transaction := &Transaction{
		WalletID:     walletID,
		Type:         changeType,
		Amount:       amount,
		Currency:     "EGP",
		BalanceBefore: result.BalanceBefore,
		BalanceAfter: result.BalanceAfter,
		CreatedAt:    time.Now(),
	}

	err = s.repo.Create(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
func (s *Service) GetTransactions(ctx context.Context, walletID string) ([]Transaction, error) {

	transactions, err := s.repo.FindByWalletID(ctx, walletID)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
