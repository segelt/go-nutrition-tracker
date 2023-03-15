package mongo

import (
	"context"
	"errors"
	"fmt"
	"nutritiontracker/configs"
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

func NewDB(conf configs.DBConfig) *DB {
	db := &DB{
		DSN:               conf.Dsn,
		connectionTimeout: conf.TimeoutDbconn,
		operationTimeout:  conf.TimeoutRead,
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

	ctx, cancel = context.WithTimeout(context.Background(), db.connectionTimeout)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		return errors.New("MongoDB.Ping")
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
