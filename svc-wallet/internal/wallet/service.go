package wallet

import (
	"context"
	"errors"
	"time"

	"svc-wallet/client/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	repo       *Repository
	userClient *user.Client
}

func NewService(repo *Repository, userClient *user.Client) *Service {
	return &Service{
		repo:       repo,
		userClient: userClient,
	}
}
func (s *Service) CreateWallet(ctx context.Context, ownerPhone, currency, nationalID string) (*Wallet, error) {
	_, err := s.userClient.GetUserByNationalID(ctx, nationalID)
	if err != nil {
		return nil, errors.New("user not found")
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

	err = s.repo.Create(ctx, wallet)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}
func (s *Service) GetWallet(ctx context.Context, id string) (*Wallet, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	wallet, err := s.repo.FindByID(ctx, objID)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *Service) ApplyBalanceChange(ctx context.Context, id string, changeType string, amount int64) (*Wallet, int64, error) {

	if amount > 60000 {
		return nil, 0, errors.New("exceeded the limit per transaction")
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, 0, err
	}
	var filter bson.M
	var update bson.M 
		switch changeType {
	case "withdraw":
		filter = bson.M{
			"_id": objID,
			"balance": bson.M{"$gte": amount},

		}
		update = bson.M{"$inc": bson.M{"balance": -amount}}

	case "deposit":
		filter = bson.M{
			"_id": objID,
			"balance": bson.M{"$lte":400000 - amount},
		}
		update = bson.M{"$inc": bson.M{"balance": amount}}


	default:
		return nil, 0, errors.New("invalid type")
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.Before)

	var wallet Wallet
	err = s.repo.collection.FindOneAndUpdate(ctx,filter,update,opts).Decode(&wallet)
	if err == nil {
		balanceBefore := wallet.Balance
		if changeType == "withdraw" {
			wallet.Balance -= amount
		}else if changeType == "deposit" {
			wallet.Balance += amount
		}
		return &wallet, balanceBefore, nil
	}
	if err != mongo.ErrNoDocuments {
		return nil, 0, err
	}
	var existing Wallet 
	lookupErr := s.repo.collection.FindOne(ctx , bson.M{"_id": objID}).Decode(&existing)

	if lookupErr == mongo.ErrNoDocuments {
		return nil, 0, errors.New("wallet not found")
	}
	if lookupErr != nil {
		return nil, 0, lookupErr
	}
	if changeType == "withdraw" {
		return nil, 0, errors.New("insufficient funds")
	}
	return nil , 0, errors.New("deposit would exceed balance limit")

}

