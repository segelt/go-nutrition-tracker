package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client            *mongo.Client
	connectionTimeout time.Duration
	operationTimeout  time.Duration
	DSN               string
}

func NewDB(dsn string) *DB {
	db := &DB{
		DSN:               dsn,
		connectionTimeout: 30 * time.Second,
		operationTimeout:  10 * time.Second,
	}

	return db
}

func (db *DB) Open() error {
	// Ensure a DSN is set before attempting to open the database.
	if db.DSN == "" {
		return fmt.Errorf("dsn required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), db.connectionTimeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(db.DSN)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return errors.New("NewMongoClient")
	}

	db.client = client
	return nil
}

func (db *DB) Close() error {
	// Cancel background context.
	// db.cancel()

	// Close database.
	if db.client != nil {
		return db.client.Disconnect(context.TODO())
	}
	return nil
}
