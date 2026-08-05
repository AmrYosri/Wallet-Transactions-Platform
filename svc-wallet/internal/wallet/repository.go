package wallet

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
		collection: db.Collection("wallets"),
	}
}

func (r *Repository) Create(ctx context.Context, wallet *Wallet) error {
	result, err := r.collection.InsertOne(ctx, wallet)
	if err != nil {
		return err
	}
	wallet.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *Repository) FindByID(ctx context.Context, id primitive.ObjectID) (*Wallet, error) {
	var wallet Wallet
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&wallet)
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}
func (r *Repository) UpdateBalance(ctx context.Context, id primitive.ObjectID, newBalance int64) error {
	update := bson.M{
		"$set": bson.M{
			"balance":    newBalance,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}
