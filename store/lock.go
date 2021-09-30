package store

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LockStore interface {
	Create(key, value string) error
	Delete(key string) error
	CreateIndexes(ctx context.Context) error
}

type lockStore struct {
	c *mongo.Collection
}

func NewLockStore(db *mongo.Database) LockStore {
	return &lockStore{
		c: db.Collection("lock"),
	}
}

func (s *lockStore) Create(key, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	doc := bson.M{
		"_id":   key,
		"value": value,
		"time":  time.Now(),
	}
	_, err := s.c.InsertOne(ctx, doc)
	return err
}

func (s *lockStore) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	filter := bson.M{
		"_id": key,
	}
	_, err := s.c.DeleteOne(ctx, filter)
	return err
}

func (s *lockStore) CreateIndexes(ctx context.Context) error {
	index := mongo.IndexModel{
		Keys:    bson.D{{"time", 1}},
		Options: options.Index().SetExpireAfterSeconds(300),
	}
	_, err := s.c.Indexes().CreateOne(ctx, index)
	return err
}