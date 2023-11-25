package db

import (
	"context"
	"finalproject/config"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoConn(cfg *config.Config) *mongo.Database {
	clientOptions := options.Client().ApplyURI(cfg.Database.Mongo.URI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(cfg.Database.Mongo.Db)
	return db
}
