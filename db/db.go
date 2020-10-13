package db

import (
	"context"
	"time"

	"github.com/vitalyo61/atlant/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Database struct {
	client  *mongo.Client
	Timeout time.Duration
	Product *mongo.Collection
}

func New(cfg *config.DB) (*Database, error) {
	db := &Database{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
	}

	var err error

	if db.client, err = mongo.NewClient(options.Client().ApplyURI(cfg.Address)); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), db.Timeout)
	defer cancel()
	err = db.client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = db.client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	dBase := db.client.Database(cfg.Name)

	db.Product = dBase.Collection("product")

	if err := CreateIndexes(db.Product, db.Timeout); err != nil {
		return nil, err
	}

	return db, nil
}

func (d *Database) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout)
	defer cancel()
	return d.client.Disconnect(ctx)
}
