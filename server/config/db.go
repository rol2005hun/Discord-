package config

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// ConnectDB csatlakozik a MongoDB adatbázishoz
func ConnectDB() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongourl")

	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ellenőrzés
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("[Wolimby] Az adatbázishoz való csatlakozás sikeres.")
	return client
}

// DisconnectDB bezárja az adatbázis kapcsolatot
func DisconnectDB(client *mongo.Client) {
	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[Wolimby] Az adatbáziskapcsolat bezárva.")
}
