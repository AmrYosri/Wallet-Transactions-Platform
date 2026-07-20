package wallet

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
		collection: db.Collection("wallets"),
	}
}

func (r *Repository) Create(ctx context.Context, wallet *Wallet) error {
	_, err := r.collection.InsertOne(ctx, wallet)
	return err
}

func (r *Repository) FindByID(ctx context.Context, id primitive.ObjectID) (*Wallet, error) {
	var wallet Wallet
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&wallet)
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}
