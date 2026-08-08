package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	collection *mongo.Collection

}

func NewRepository (db *mongo.Database) *Repository{
	return &Repository{
		collection: db.Collection("users"),
	}
}


func (r *Repository) Create(ctx context.Context, user *User) error {

	result , err := r.collection.InsertOne(ctx,user)
	if err != nil {
		return err
	}
	user.ID =result.InsertedID.(primitive.ObjectID)
	return nil
}


func(r *Repository) FindByID(ctx context.Context , id primitive.ObjectID) (*User,error){
	var user User
	err :=r.collection.FindOne(ctx , bson.M{"_id":id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user,nil 
}

func (r *Repository) FindByNationalID(ctx context.Context, nationalID string) (*User,error){
	var user User
	err := r.collection.FindOne(ctx,bson.M{"national_id":nationalID}).Decode(&user)
	if err != nil {
		return nil ,err
	}
	return &user , nil
}


func(r *Repository) AddPhoneNumber(ctx context.Context , nationalID string , phone string) error{
	update := bson.M{
		"$push":bson.M{
			"phone_numbers":phone,
		},
	}
	_,err := r.collection.UpdateOne(ctx ,bson.M{"national_id":nationalID},update)
	return err
}