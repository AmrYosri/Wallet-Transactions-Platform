// Package transactions — data access for the transactions domain.
package transactions

import (
	"context"
	"errors"
	"time"

	"svc-transactions/util/apperror"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	collection *mongo.Collection // transactions
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		collection: db.Collection("transactions"),
	}
}

// Create inserts a transaction record and writes the generated ID back.
// A duplicate key means the unique index on request_id rejected it — this
// request has already been recorded, so the caller should look up the
// existing record rather than treating it as a failure.
func (r *Repository) Create(ctx context.Context, transaction *Transaction) error {
	result, err := r.collection.InsertOne(ctx, transaction)

	if mongo.IsDuplicateKeyError(err) {
		return apperror.DuplicateRequest
	}
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("request_id", transaction.RequestID).Msg("mongo insert transaction failed")
		return apperror.Internal
	}

	transaction.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *Repository) FindByWalletID(ctx context.Context, walletID string) ([]Transaction, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"wallet_id": walletID})
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("wallet_id", walletID).Msg("mongo find transactions failed")
		return nil, apperror.Internal
	}
	defer cursor.Close(ctx)

	var transactions []Transaction
	if err := cursor.All(ctx, &transactions); err != nil {
		log.Ctx(ctx).Error().Err(err).Str("wallet_id", walletID).Msg("mongo decode transactions failed")
		return nil, apperror.Internal
	}

	return transactions, nil
}

// FindByRequestID looks up an existing record by its idempotency key.
func (r *Repository) FindByRequestID(ctx context.Context, requestID string) (*Transaction, error) {
	var transaction Transaction

	err := r.collection.FindOne(ctx, bson.M{"request_id": requestID}).Decode(&transaction)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, apperror.TransactionNotFound
	}
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("request_id", requestID).Msg("mongo find by request_id failed")
		return nil, apperror.Internal
	}

	return &transaction, nil
}

// MarkCompleted moves a PENDING record to COMPLETED and stores the balances
// svc-wallet reported.
func (r *Repository) MarkCompleted(ctx context.Context, id primitive.ObjectID, balanceBefore, balanceAfter int64) error {
	update := bson.M{"$set": bson.M{
		"status":         StatusCompleted,
		"balance_before": balanceBefore,
		"balance_after":  balanceAfter,
		"updated_at":     time.Now(),
	}}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("transaction_id", id.Hex()).Msg("mongo mark completed failed")
		return apperror.Internal
	}

	return nil
}

// MarkFailed records why the wallet call failed. reason is the domain code
// (e.g. INSUFFICIENT_FUNDS), not prose — it's searchable that way.
func (r *Repository) MarkFailed(ctx context.Context, id primitive.ObjectID, reason string) error {
	update := bson.M{"$set": bson.M{
		"status":         StatusFailed,
		"failure_reason": reason,
		"updated_at":     time.Now(),
	}}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("transaction_id", id.Hex()).Msg("mongo mark failed failed")
		return apperror.Internal
	}

	return nil
}
