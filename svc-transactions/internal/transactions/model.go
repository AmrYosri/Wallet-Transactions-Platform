package transactions

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)
type Status string
const (
	StatusPending   Status = "pending"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
)

type Transaction struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	WalletID     string             `bson:"wallet_id"`
	Type         string             `bson:"type"`
	Amount       int64              `bson:"amount"`
	Currency     string             `bson:"currency"`
	BalanceBefore int64              `bson:"balance_before"`
	BalanceAfter int64              `bson:"balance_after"`
	FailureReason string             `bson:"failure_reason,omitempty"`
	RequestID    string             `bson:"request_id"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
	Status       Status             `bson:"status"`	
}
