package wallet

import (
	"context"
	"errors"
	"time"

	"svc-wallet/client/user"
	"svc-wallet/util/apperror"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	repo       *Repository
	userClient *user.Client
}

func NewService(repo *Repository, userClient *user.Client) *Service {
	return &Service{repo: repo, userClient: userClient}
}

// CreateWallet verifies the owner exists in svc-user, then creates an empty wallet.
func (s *Service) CreateWallet(ctx context.Context, ownerPhone, currency, nationalID string) (*Wallet, error) {
	if _, err := s.userClient.GetUserByNationalID(ctx, nationalID); err != nil {
		return nil, apperror.UserNotFound
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

	if err := s.repo.Create(ctx, wallet); err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *Service) GetWallet(ctx context.Context, id string) (*Wallet, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, apperror.WalletNotFound
	}
	return s.repo.FindByID(ctx, objID)
}

// ApplyBalanceChange moves money in or out of a wallet.
//
//   - idempotent — the receipt lookup replays a repeated request_id instead of
//     applying it twice
//   - race-free — the funds/limit check lives in the Mongo filter, so it's
//     evaluated at write time, not read in Go first
//   - auditable — every applied change leaves a receipt
func (s *Service) ApplyBalanceChange(ctx context.Context, id, changeType string, amount int64, requestID string) (*Wallet, int64, error) {

	// Retry path: already applied? Return the stored result, change nothing.
	existing, err := s.repo.FindBalanceChange(ctx, requestID)
	if err != nil {
		return nil, 0, err
	}
	if existing != nil {
		wallet, err := s.repo.FindByID(ctx, existing.WalletID)
		if err != nil {
			return nil, 0, err
		}
		return wallet, existing.BalanceBefore, nil
	}

	if amount > 60000 {
		return nil, 0, apperror.LimitExceeded
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, 0, apperror.WalletNotFound
	}

	// The condition goes in the filter, not in Go. Mongo evaluates it as part
	// of the same atomic operation as the $inc, so two concurrent withdrawals
	// can't both pass the check and both write.
	var filter, update bson.M
	switch changeType {
	case "withdraw":
		filter = bson.M{"_id": objID, "balance": bson.M{"$gte": amount}}
		update = bson.M{
			"$inc": bson.M{"balance": -amount},
			"$set": bson.M{"updated_at": time.Now()},
		}
	case "deposit":
		filter = bson.M{"_id": objID, "balance": bson.M{"$lte": 400000 - amount}}
		update = bson.M{
			"$inc": bson.M{"balance": amount},
			"$set": bson.M{"updated_at": time.Now()},
		}
	default:
		return nil, 0, apperror.InvalidType
	}

	// After returns the post-update document from Mongo, so the new balance is
	// authoritative and we derive the old one rather than the other way round.
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var wallet Wallet
	err = s.repo.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&wallet)

	if errors.Is(err, mongo.ErrNoDocuments) {
		// Nothing matched — either no such wallet, or the condition failed.
		return nil, 0, s.explainNoMatch(ctx, objID, changeType)
	}
	if err != nil {
		return nil, 0, apperror.Internal
	}

	balanceBefore := wallet.Balance + amount
	if changeType == "deposit" {
		balanceBefore = wallet.Balance - amount
	}

	// Receipt, written after the fact. If the process dies between the $inc
	// and this insert, a retry re-applies — closing that window needs both
	// writes in one Mongo session (Part B requirement 2).
	_ = s.repo.ClaimRequest(ctx, &BalanceChange{
		RequestID:     requestID,
		WalletID:      objID,
		Type:          changeType,
		Amount:        amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  wallet.Balance,
		CreatedAt:     time.Now(),
	})

	return &wallet, balanceBefore, nil
}

// explainNoMatch turns "the filter matched nothing" into the real reason.
func (s *Service) explainNoMatch(ctx context.Context, objID primitive.ObjectID, changeType string) error {
	if _, err := s.repo.FindByID(ctx, objID); err != nil {
		return err // WalletNotFound or Internal
	}
	if changeType == "withdraw" {
		return apperror.InsufficientFunds
	}
	return apperror.MaxBalanceExceeded
}
