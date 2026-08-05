package transactions

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	WalletID     string             `bson:"wallet_id"`
	Type         string             `bson:"type"`
	Amount       int64              `bson:"amount"`
	Currency     string             `bson:"currency"`
	BalanceBefore int64              `bson:"balance_before"`
	BalanceAfter int64              `bson:"balance_after"`
	RequestID    string             `bson:"request_id"`
	CreatedAt    time.Time          `bson:"created_at"`
}
