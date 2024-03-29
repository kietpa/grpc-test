package config

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func InitDB() (mongoClient *mongo.Client, err error) {
	ctx := context.TODO()
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	connectionOpts := options.Client().ApplyURI("mongodb://localhost:27017").SetServerAPIOptions(serverAPIOptions)

	mongoClient, err = mongo.Connect(ctx, connectionOpts)
	if err != nil {
		return nil, err
	}

	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return mongoClient, err
}
