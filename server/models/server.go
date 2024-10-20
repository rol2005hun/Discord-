package models

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	OwnerID   primitive.ObjectID `bson:"owner_id"`
	Name      string             `bson:"name"`
	Image     string             `bson:"image,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
	Channels  []Channel          `bson:"channels,omitempty"`
	Roles     []Role             `bson:"roles,omitempty"`
}

type Role struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Permissions []string           `bson:"permissions,omitempty"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty"`
}

// ServerCollection is the MongoDB collection for servers
var ServerCollection *mongo.Collection

// InitializeServerCollection sets up the server collection
func InitializeServerCollection(dbClient *mongo.Client) {
	ServerCollection = dbClient.Database("your_database_name").Collection("servers")
}

// CreateServer hoz létre egy új szervert az adatbázisban
func CreateServer(server Server, dbClient *mongo.Client) (*mongo.InsertOneResult, error) {
	collection := dbClient.Database("mydb").Collection("servers")

	// Alapértelmezett kép beállítása, ha az image mező üres
	if server.Image == "" {
		server.Image = "default.png" // vagy a kívánt alapértelmezett URL
	}

	// Beállítjuk a létrehozás és frissítés dátumát
	server.CreatedAt = time.Now()
	server.UpdatedAt = time.Now()

	result, err := collection.InsertOne(context.TODO(), server)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetServers lekérdezi az összes szervert
func GetServers(dbClient *mongo.Client) ([]Server, error) {
	var servers []Server

	cursor, err := ServerCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find servers: %v", err)
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &servers); err != nil {
		return nil, fmt.Errorf("failed to decode servers: %v", err)
	}

	return servers, nil
}

func GetServersByOwnerID(ownerID primitive.ObjectID, dbClient *mongo.Client) ([]Server, error) {
	collection := dbClient.Database("mydb").Collection("servers")
	if collection == nil {
		return nil, fmt.Errorf("failed to get collection")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var servers []Server
	filter := bson.M{"owner_id": ownerID}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var server Server
		if err := cursor.Decode(&server); err != nil {
			return nil, err
		}
		servers = append(servers, server)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return servers, nil
}

// GetServerByID lekérdezi a szervert ID alapján
func GetServerByID(id string, dbClient *mongo.Client) (*Server, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid server ID: %v", err)
	}

	var server Server
	err = ServerCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&server)
	if err != nil {
		return nil, fmt.Errorf("server not found: %v", err)
	}

	return &server, nil
}

// UpdateServer frissíti a szerver adatait ID alapján
func UpdateServer(id string, updatedServer Server, dbClient *mongo.Client) (*mongo.UpdateResult, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid server ID: %v", err)
	}

	updatedServer.UpdatedAt = time.Now()

	update := bson.M{
		"$set": updatedServer,
	}
	result, err := ServerCollection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update server: %v", err)
	}
	return result, nil
}

// DeleteServer törli a szervert ID alapján
func DeleteServer(id string, dbClient *mongo.Client) (*mongo.DeleteResult, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid server ID: %v", err)
	}

	result, err := ServerCollection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		return nil, fmt.Errorf("failed to delete server: %v", err)
	}
	return result, nil
}
