package transactions

import (
	"context"
	"errors"
	"svc-transactions/client/wallet"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
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
func (s *Service) ApplyTransaction(ctx context.Context, walletID string, changeType string, amount int64, requestID string) (*Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}
	if requestID == ""{
		return nil , errors.New("request_id is required")
	}


	transaction := &Transaction{
		WalletID:     walletID,
		Type:         changeType,
		Amount:       amount,
		Currency:     "EGP",
		Status:       StatusPending,
		RequestID:    requestID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := s.repo.Create(ctx,transaction)
	if err !=nil{
		if mongo.IsDuplicateKeyError(err) {
			existing , findErr := s.repo.FindByRequestID(ctx, requestID)
			if findErr != nil {
				return nil, findErr
			}
			return existing, nil
		}
		return nil , err
	}

	result, err := s.walletClient.ApplyBalanceChange(ctx, walletID, changeType,amount)
	if err != nil{
		_ = s.repo.MarkFailed(ctx ,transaction.ID, err.Error())
		return nil , err 
	}
	if err := s.repo.MarkCompleted(ctx ,transaction.ID,result.BalanceBefore ,result.BalanceAfter); err!=nil{
		return nil , err
	}
	transaction.Status = StatusCompleted
	transaction.BalanceBefore = result.BalanceBefore
	transaction.BalanceAfter = result.BalanceAfter
	return transaction , nil

}
func (s *Service) GetTransactions(ctx context.Context, walletID string) ([]Transaction, error) {

	transactions, err := s.repo.FindByWalletID(ctx, walletID)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
