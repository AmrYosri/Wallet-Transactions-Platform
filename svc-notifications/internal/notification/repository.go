package notification

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


type Repository struct {
	smsCollection  *mongo.Collection
	pushCollection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		smsCollection:  db.Collection("sms_notifications"),
		pushCollection: db.Collection("push_notifications"),
	}

}


func (r *Repository) CreateSMS(ctx context.Context, notification *SMSNotification) error{
	result,err := r.smsCollection.InsertOne(ctx,notification)
	if err != nil{
		return err
	}
	notification.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *Repository) MarkSMSSent(ctx context.Context, id primitive.ObjectID) error{
	update := bson.M{
		"$set":bson.M{
			"status": SMSStatusSent,
			"sent_at": time.Now(),
		},
	}
	_,err := r.smsCollection.UpdateOne(ctx ,bson.M{"_id":id},update)
	return err
}

func (r *Repository) MarkSMSFailed(ctx context.Context, id primitive.ObjectID, reason string) error{
	update := bson.M{
		"$set":bson.M{
			"status": SMSStatusFailed,
			"failure_reason": reason,
			"sent_at": time.Now(),
		},
	}
	_,err := r.smsCollection.UpdateOne(ctx,bson.M{"_id":id},update)
	return err
}



func (r *Repository) CreatePush(ctx context.Context, notification *PushNotification) error{
	result,err := r.pushCollection.InsertOne(ctx,notification)
	if err != nil{
		return err
	}
	notification.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *Repository) MarkPushSent(ctx context.Context, id primitive.ObjectID) error{
	update := bson.M{
		"$set":bson.M{
			"status": PushStatusSent,
			"sent_at": time.Now(),
		},
	}
	_,err := r.pushCollection.UpdateOne(ctx ,bson.M{"_id":id},update)
	return err
}

func (r *Repository) MarkPushFailed(ctx context.Context, id primitive.ObjectID, reason string) error{
	update := bson.M{
		"$set":bson.M{
			"status": PushStatusFailed,
			"failure_reason": reason,
			"sent_at": time.Now(),
		},
	}
	_,err := r.pushCollection.UpdateOne(ctx,bson.M{"_id":id},update)
	return err
}


