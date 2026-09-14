package transactions

import (
	"context"
	"fmt"
	"time"

	"svc-transactions/client/notification"
	"svc-transactions/client/wallet"
	"svc-transactions/util/apperror"

	"google.golang.org/grpc/status"
)

type Service struct {
	repo               *Repository
	walletClient       *wallet.Client
	notificationClient *notification.Client
}

func NewService(repo *Repository, walletClient *wallet.Client, notificationClient *notification.Client) *Service {
	return &Service{
		repo:               repo,
		walletClient:       walletClient,
		notificationClient: notificationClient,
	}
}

// ApplyTransaction records a money movement, then asks svc-wallet to move it.
//
// The record is written PENDING first, so a crash mid-flow leaves a trace
// rather than nothing. requestID is carried through to svc-wallet so a retry
// is recognised on both sides.
func (s *Service) ApplyTransaction(ctx context.Context, walletID string, changeType string, amount int64, requestID string) (*Transaction, error) {
	if amount <= 0 {
		return nil, apperror.InvalidAmount
	}
	if requestID == "" {
		return nil, apperror.InvalidRequestBody
	}

	transaction := &Transaction{
		WalletID:  walletID,
		Type:      changeType,
		Amount:    amount,
		Currency:  "EGP",
		Status:    StatusPending,
		RequestID: requestID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

		// Our own idempotency check: a duplicate request_id means we've seen this
	// before, so return the existing record instead of starting again.
	err := s.repo.Create(ctx, transaction)
	if err != nil {
		if err == apperror.DuplicateRequest {
			existing, findErr := s.repo.FindByRequestID(ctx, requestID)
			if findErr != nil {
				return nil, findErr
			}
			return existing, nil
		}
		return nil, err
	}

	result, err := s.walletClient.ApplyBalanceChange(ctx, walletID, changeType, amount, requestID)
	if err != nil {
		appErr := fromGRPC(err)
		_ = s.repo.MarkFailed(ctx, transaction.ID, appErr.Code)
		return nil, appErr
	}

	if err := s.repo.MarkCompleted(ctx, transaction.ID, result.BalanceBefore, result.BalanceAfter); err != nil {
		return nil, apperror.Internal
	}

	transaction.Status = StatusCompleted
	transaction.BalanceBefore = result.BalanceBefore
	transaction.BalanceAfter = result.BalanceAfter

	// Fire-and-forget: the money has already moved, so a failed notification
	// must not fail the transaction. Phase 4 replaces this with Kafka, which
	// is what makes it actually reliable — right now a process restart loses it.
	go func() {
		_ = s.notificationClient.Notify(notification.NotifyRequest{
			TransactionID: transaction.ID.Hex(),
			Phone:         "01000000000",
			DeviceToken:   "stub-device-token",
			SMSMessage:    fmt.Sprintf("Your %s of %d %s was successful", changeType, amount, transaction.Currency),
			PushTitle:     "Transaction Successful",
			PushBody:      fmt.Sprintf("%s of %d %s completed", changeType, amount, transaction.Currency),
		})
	}()

	return transaction, nil
}

func (s *Service) GetTransactions(ctx context.Context, walletID string) ([]Transaction, error) {
	return s.repo.FindByWalletID(ctx, walletID)
}

// fromGRPC turns an error from svc-wallet back into an AppError.
// svc-wallet puts the domain code in the status message, so the caller sees
// the real reason instead of a generic transport failure.
func fromGRPC(err error) *apperror.AppError {
	st, ok := status.FromError(err)
	if !ok {
		return apperror.ServiceUnavailable
	}

	switch st.Message() {
	case "INSUFFICIENT_FUNDS":
		return apperror.InsufficientFunds
	case "WALLET_NOT_FOUND":
		return apperror.WalletNotFound
	case "LIMIT_EXCEEDED":
		return apperror.LimitExceeded
	case "MAX_BALANCE_EXCEEDED":
		return apperror.MaxBalanceExceeded
	case "INVALID_TYPE":
		return apperror.InvalidType
	case "DUPLICATE_REQUEST":
		return apperror.DuplicateRequest
	default:
		// svc-wallet unreachable, or an error we don't have a mapping for.
		return apperror.ServiceUnavailable
	}
}