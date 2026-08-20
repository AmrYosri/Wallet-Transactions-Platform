package transactions

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {

	return &Repository{
		collection: db.Collection("transactions"),
	}
}

func (r *Repository) Create(ctx context.Context, transaction *Transaction) error {
	result, err := r.collection.InsertOne(ctx, transaction)
	if err != nil {
		return err
	}
	transaction.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}
func (r *Repository) FindByWalletID(ctx context.Context, walletID string) ([]Transaction, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"wallet_id": walletID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []Transaction
	for cursor.Next(ctx) {
		var t Transaction
		if err := cursor.Decode(&t); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (r *Repository) MarkCompleted(ctx context.Context, id primitive.ObjectID , balanceBefore , balanceAfter int64) error {
	update := bson.M{
		"$set":bson.M{
			"status": StatusCompleted,
			"balance_before": balanceBefore,
			"balance_after": balanceAfter,
			"updated_at": time.Now(),
		},
	}
	_,err:= r.collection.UpdateOne(ctx, bson.M{"_id":id},update)
	return err
}

func (r *Repository) MarkFailed(ctx context.Context, id primitive.ObjectID , reason string) error {
	update := bson.M{
		"$set":bson.M{
			"status": StatusFailed,
			"failure_reason": reason,
			"updated_at": time.Now(),
		},
	}
	_,err:= r.collection.UpdateOne(ctx, bson.M{"_id":id},update)
	return err
}

func (r *Repository) FindByRequestID(ctx context.Context, requestID string) (*Transaction, error) {
	var transaction Transaction
	err := r.collection.FindOne(ctx, bson.M{"request_id": requestID}).Decode(&transaction)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}