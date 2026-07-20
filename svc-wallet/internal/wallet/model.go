package wallet

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Wallet struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	OwnerPhone string             `bson:"owner_phone"`
	Currency   string             `bson:"currency"`
	Balance    int64              `bson:"balance"`
	Status     string             `bson:"status"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}