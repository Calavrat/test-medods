package store

import (
	"context"
	"fmt"

	"github.com/Calavrat/TestMedods/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthMongo struct {
	db         *mongo.Database
	collection *mongo.Collection
}

var ctx context.Context

func NewAuthorization(db *mongo.Database, collection string) *AuthMongo {
	return &AuthMongo{
		db:         db,
		collection: db.Collection(collection),
	}
}

func (r *AuthMongo) GetUser(userid string) (model.User, error) {
	var user model.User
	oid, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		return user, err
	}

	filter := bson.M{"_id": oid}
	result := r.collection.FindOne(ctx, filter)

	if err := result.Decode(&user); err != nil {
		return user, err
	}

	return user, nil
}

func (r *AuthMongo) AddRefresh(user model.User) error {
	oid, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}

	userBytes, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("error to marshal user: %v", err)
	}

	var updateUserObj bson.M
	err = bson.Unmarshal(userBytes, &updateUserObj)
	if err != nil {
		return fmt.Errorf("error to unmarshal user: %v", err)
	}
	delete(updateUserObj, "_id")

	update := bson.M{"$set": updateUserObj}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error to update user: %v", err)
	}

	return nil
}
