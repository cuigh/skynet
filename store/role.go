package store

import (
	"context"
	"time"

	"github.com/cuigh/auxo/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Role struct {
	ID          string   `json:"id" bson:"_id"`
	Name        string   `json:"name" bson:"name"`
	Description string   `json:"desc" bson:"desc"`
	Perms       []string `json:"perms" bson:"perms"`
	CreateTime  Time     `json:"create_time" bson:"create_time"`
	ModifyTime  Time     `json:"modify_time" bson:"modify_time"`
}

type RoleStore interface {
	Create(r *Role) error
	Modify(role *Role) error
	Find(id string) (*Role, error)
	Search(name string) ([]*Role, error)
	Delete(id string) error
}

type roleStore struct {
	c *mongo.Collection
}

func NewRoleStore(db *mongo.Database) RoleStore {
	return &roleStore{db.Collection("role")}
}

func (s *roleStore) Create(r *Role) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r.CreateTime = Time(time.Now())
	r.ModifyTime = r.CreateTime
	_, err := s.c.InsertOne(ctx, r)
	if mongo.IsDuplicateKeyError(err) {
		return errors.Format("角色标志符 %s 已经存在", r.ID)
	}
	return err
}

func (s *roleStore) Modify(role *Role) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	role.ModifyTime = Time(time.Now())
	r, err := s.c.UpdateByID(ctx, role.ID, bson.M{"$set": role})
	if err != nil {
		return err
	} else if r.MatchedCount == 0 {
		return errors.Format("can't find role '%s'", role.ID)
	}
	return nil
}

func (s *roleStore) Find(id string) (*Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	role := &Role{}
	filter := bson.M{"_id": id}
	if err := s.c.FindOne(ctx, filter).Decode(&role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *roleStore) Search(name string) ([]*Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	if name != "" {
		filter["name"] = name
	}
	cur, err := s.c.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var roles []*Role
	err = cur.All(ctx, &roles)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (s *roleStore) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r, err := s.c.DeleteOne(ctx, bson.M{"_id": id})
	if err == nil && r.DeletedCount == 0 {
		return errors.Format("can't find role '%s'", id)
	}
	return err
}
