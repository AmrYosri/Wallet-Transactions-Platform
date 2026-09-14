// Package wallet — data access for the wallet domain.
// Two collections: wallets (balance of record), balance_changes (one receipt
// per applied operation, used for idempotency).
package wallet

import (
	"context"
	"errors"
	"svc-wallet/util/apperror"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository holds the Mongo collections this domain owns.
type Repository struct {
	collection     *mongo.Collection // wallets
	balanceChanges *mongo.Collection // balance_changes
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		collection:     db.Collection("wallets"),
		balanceChanges: db.Collection("balance_changes"),
	}
}

// Create inserts a wallet and writes the generated ID back onto the struct.
func (r *Repository) Create(ctx context.Context, wallet *Wallet) error {
	result, err := r.collection.InsertOne(ctx, wallet)

	// Only fires if a unique index exists on (national_id, currency).
	if mongo.IsDuplicateKeyError(err) {
		return apperror.WalletAlreadyExists
	}
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("mongo insert wallet failed")
		return apperror.Internal
	}

	wallet.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID returns one wallet.
// Two branches on purpose: "no document" is a 404, a timeout or dropped
// connection is a 500. Collapsing them would report a missing wallet when
// the database is actually unreachable.
func (r *Repository) FindByID(ctx context.Context, id primitive.ObjectID) (*Wallet, error) {
	var wallet Wallet

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&wallet)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, apperror.WalletNotFound
	}
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("wallet_id", id.Hex()).Msg("mongo find failed")
		return nil, apperror.Internal
	}

	return &wallet, nil
}

// ClaimRequest writes the receipt for an applied operation.
// The receipt's _id is the request_id, so Mongo's built-in unique index on
// _id rejects a second write for the same request. No extra index needed.
func (r *Repository) ClaimRequest(ctx context.Context, change *BalanceChange) error {
	_, err := r.balanceChanges.InsertOne(ctx, change)

	if mongo.IsDuplicateKeyError(err) {
		return apperror.DuplicateRequest
	}
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("request_id", change.RequestID).Msg("mongo claim request failed")
		return apperror.Internal
	}

	return nil
}

// FindBalanceChange looks for a receipt from a previous call.
// Returns (nil, nil) on a miss — unlike FindByID, "not there" is the normal
// case: it means this request hasn't been applied yet. A non-nil result means
// this is a retry, and the caller replays the stored numbers instead of
// moving money again.
func (r *Repository) FindBalanceChange(ctx context.Context, requestID string) (*BalanceChange, error) {
	var change BalanceChange

	err := r.balanceChanges.FindOne(ctx, bson.M{"_id": requestID}).Decode(&change)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("request_id", requestID).Msg("mongo find balance change failed")
		return nil, apperror.Internal
	}

	return &change, nil
}

// CompleteBalanceChange fills in the result on an existing receipt.
// Unused right now — the current flow writes the receipt complete. Kept for
// the claim-first flow that Part B requirement 2 will need.
func (r *Repository) CompleteBalanceChange(ctx context.Context, requestID string, before, after int64) error {
	update := bson.M{"$set": bson.M{
		"balance_before": before,
		"balance_after":  after,
		"completed_at":   time.Now(),
	}}

	_, err := r.balanceChanges.UpdateOne(ctx, bson.M{"_id": requestID}, update)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("request_id", requestID).Msg("mongo complete balance change failed")
		return apperror.Internal
	}

	return nil
}

// ReleaseRequest deletes a receipt so the request_id can be used again.
// Also unused for now — same reason as CompleteBalanceChange.
func (r *Repository) ReleaseRequest(ctx context.Context, requestID string) error {
	_, err := r.balanceChanges.DeleteOne(ctx, bson.M{"_id": requestID})
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("request_id", requestID).Msg("mongo release request failed")
		return apperror.Internal
	}

	return nil
}
