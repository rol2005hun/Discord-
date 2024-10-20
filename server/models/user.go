package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name,omitempty"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

func GetUserByUsernameOrEmail(usernamOrEmail string, dbClient *mongo.Client) (User, error) {
	collection := dbClient.Database("mydb").Collection("users")

	var user User
	err := collection.FindOne(context.TODO(), bson.M{"username": usernamOrEmail}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		err = collection.FindOne(context.TODO(), bson.M{"email": usernamOrEmail}).Decode(&user)
		if err != nil {
			return User{}, err
		}
	} else if err != nil {
		return User{}, err
	}

	return user, nil
}

// CreateUser létrehoz egy új felhasználót
func CreateUser(user User, dbClient *mongo.Client) (*mongo.InsertOneResult, error) {
	collection := dbClient.Database("mydb").Collection("users")

	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetUsers visszaadja az összes felhasználót
func GetUsers(dbClient *mongo.Client) ([]User, error) {
	collection := dbClient.Database("mydb").Collection("users")

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	var users []User
	if err = cursor.All(context.TODO(), &users); err != nil {
		return nil, err
	}

	return users, nil
}
