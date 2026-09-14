package wallet

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Wallet struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	OwnerPhone string             `bson:"owner_phone"`
	NationalID string             `bson:"national_id"`
	Currency   string             `bson:"currency"`
	Balance    int64              `bson:"balance"`
	Status     string             `bson:"status"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}

// BalanceChange is a receipt for one applied operation.
// _id is the request_id, so Mongo's built-in unique index gives us
// idempotency with no extra index to create.
type BalanceChange struct {
	RequestID     string             `bson:"_id"`
	WalletID      primitive.ObjectID `bson:"wallet_id"`
	Type          string             `bson:"type"`
	Amount        int64              `bson:"amount"`
	BalanceBefore int64              `bson:"balance_before"`
	BalanceAfter  int64              `bson:"balance_after"`
	CreatedAt     time.Time          `bson:"created_at"`
}
