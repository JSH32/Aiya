package database

import (
	"context"
	"log"
	"time"

	"nawp.com/util/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitDB : Mongo Initializer
func InitDB(config config.Config) *mongo.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.DB.URL))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	return client.Database(config.DB.DBNAME)
}
