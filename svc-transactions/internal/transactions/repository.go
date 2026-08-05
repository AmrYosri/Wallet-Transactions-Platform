package transactions

import (
	"context"

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
