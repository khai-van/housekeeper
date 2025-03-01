package mongox

import (
	"context"
	"fmt"
	"log"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectMongoDB(ctx context.Context, uri, databaseName string) error {
	clientOptions := options.Client().ApplyURI(uri)

	err := mgm.SetDefaultConfig(&mgm.Config{}, databaseName, clientOptions)
	if err != nil {
		return fmt.Errorf("can't connect to mongoDB: %w", err)
	}

	_, client, _, err := mgm.DefaultConfigs()
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Check the connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Connected to MongoDB!")
	return nil
}
