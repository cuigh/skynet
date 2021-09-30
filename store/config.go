package store

import (
	"context"
	"time"

	"github.com/cuigh/auxo/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Id         string       `json:"id" bson:"_id"`
	Options    data.Options `json:"options" bson:"options"`
	ModifyTime Time         `json:"modify_time" bson:"modify_time"`
}

//type Config struct {
//	Id  string `json:"id" bson:"_id"`
//	SMS struct {
//		Body string `json:"body" bson:"body"`
//	} `json:"sms" bson:"sms"`
//	Email struct {
//		Title    string `json:"title" bson:"title"` // 标题模版
//		Body     string `json:"body" bson:"body"`   // 内容模版
//		Server   string
//		Username string
//		Password string
//	} `json:"email" bson:"email"`
//	ModifyTime Time `json:"modify_time" bson:"modify_time"`
//}

type ConfigStore interface {
	Find(id string) (data.Options, error)
	Save(id string, opts data.Options) error
}

func NewConfigStore(db *mongo.Database) ConfigStore {
	return &configStore{
		c: db.Collection("config"),
	}
}

type configStore struct {
	c *mongo.Collection
}

func (s *configStore) Find(id string) (data.Options, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r := s.c.FindOne(ctx, bson.M{"_id": id})
	c := &Config{}
	if err := r.Decode(&c); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return c.Options, nil
}

func (s *configStore) Save(id string, opts data.Options) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"options":     opts,
		"modify_time": time.Now(),
	}
	_, err := s.c.UpdateByID(ctx, id, bson.M{"$set": update}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}
