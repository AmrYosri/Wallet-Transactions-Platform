package notification

import (
	"time"	
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type SMSStatus string 
const (
	SMSStatusPending SMSStatus = "pending"
	SMSStatusSent  SMSStatus ="sent"
	SMSStatusFailed SMSStatus = "failed"
)


type SMSNotification struct {
	ID		  primitive.ObjectID `bson:"_id,omitempty"`
	TransactionID string `bson:"transaction_id"`
	Phone string `bson:"phone"`
	Message string `bson:"message"`

	Status SMSStatus `bson:"status"`
	FailureReason string `bson:"failed_reason,omitempty"`
	CreatedAt time.Time `bson:"created_at"`
	SentAt *time.Time `bson:"sent_at,omitempty"`

}

type PushStatus string
const (
	PushStatusPending PushStatus = "pending"
	PushStatusSent  PushStatus ="sent"
	PushStatusFailed PushStatus = "failed"
)

type PushNotification struct {
	ID		  primitive.ObjectID `bson:"_id,omitempty"`
	TransactionID string 		   `bson:"transaction_id"`
	DeviceToken string `bson:"device_token"`
	Title string `bson:"title"`
	Body string `bson:"body"`

	Status PushStatus `bson:"status"`
	FailureReason string `bson:"failed_reason,omitempty"`
	CreatedAt time.Time `bson:"created_at"`
	SentAt *time.Time `bson:"sent_at,omitempty"`
}
