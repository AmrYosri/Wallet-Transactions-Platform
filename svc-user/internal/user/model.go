package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	FirstName    string             `bson:"first_name"`
	LastName     string             `bson:"last_name"`
	NationalID   string             `bson:"national_id"`
	PhoneNumbers []string           `bson:"phone_numbers"`
	Status       string             `bson:"status"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
}