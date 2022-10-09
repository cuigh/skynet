package store

import (
	"context"
	"time"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Task struct {
	Name        string       `json:"name" bson:"_id" valid:"required"`
	Runner      string       `json:"runner" bson:"runner" valid:"required"`
	Handler     string       `json:"handler,omitempty" bson:"handler,omitempty"`
	Args        data.Options `json:"args" bson:"args"`
	Triggers    []string     `json:"triggers" bson:"triggers"`
	Description string       `json:"desc,omitempty" bson:"desc,omitempty"`
	Enabled     bool         `json:"enabled" bson:"enabled"`
	Maintainers []string     `json:"maintainers" bson:"maintainers"`
	Alerts      []string     `json:"alerts" bson:"alerts"`
	ModifyTime  Time         `json:"modify_time" bson:"modify_time"`
}

type TaskStore interface {
	Find(name string) (*Task, error)
	Delete(name string) error
	Create(t *Task) error
	Modify(t *Task) error
	Search(name, runner string, pageIndex, pageSize int64) (tasks []*Task, total int64, err error)
	GetState() (modify time.Time, count int64, err error)
	FetchAll(enabled bool) ([]*Task, error)
	Count(ctx context.Context) (int64, error)
}

type taskStore struct {
	c *mongo.Collection
}

func NewTaskStore(db *mongo.Database) TaskStore {
	return &taskStore{
		c: db.Collection("task"),
	}
}

func (s *taskStore) Find(name string) (*Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r := s.c.FindOne(ctx, bson.M{"_id": name})
	t := &Task{}
	if err := r.Decode(&t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *taskStore) Delete(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r, err := s.c.DeleteOne(ctx, bson.M{"_id": name})
	if err == nil && r.DeletedCount == 0 {
		return errors.Format("can't find task '%s'", name)
	}
	return err
}

func (s *taskStore) Create(t *Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.ModifyTime = Time(time.Now())
	_, err := s.c.InsertOne(ctx, t)
	if mongo.IsDuplicateKeyError(err) {
		return errors.Format("任务名 %s 已经存在", t.Name)
	}
	return err
}

func (s *taskStore) Modify(t *Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.ModifyTime = Time(time.Now())
	r, err := s.c.UpdateByID(ctx, t.Name, bson.M{"$set": t})
	if err != nil {
		return err
	} else if r.MatchedCount == 0 {
		return errors.Format("can't find task '%s'", t.Name)
	}
	return nil
}

func (s *taskStore) Search(name, runner string, pageIndex, pageSize int64) (tasks []*Task, total int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	if name != "" {
		filter["_id"] = name
	}
	if runner != "" {
		filter["runner"] = runner
	}

	// fetch total count
	total, err = s.c.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().SetSkip(pageSize * (pageIndex - 1)).SetLimit(pageSize).SetSort(bson.M{"_id": 1})
	cur, err := s.c.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	err = cur.All(ctx, &tasks)
	if err != nil {
		return nil, 0, err
	}
	return tasks, total, nil
}

func (s *taskStore) GetState() (modify time.Time, count int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"enabled": true}

	// count
	count, err = s.c.CountDocuments(ctx, filter)
	if err != nil {
		return
	}

	// modify time
	r := s.c.FindOne(ctx, filter, options.FindOne().SetSort(bson.M{"modify_time": -1}))
	t := struct {
		ModifyTime time.Time `bson:"modify_time"`
	}{}
	if err = r.Decode(&t); err == mongo.ErrNoDocuments {
		err = nil
	}
	modify = t.ModifyTime
	return
}

func (s *taskStore) FetchAll(enabled bool) ([]*Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"enabled": enabled}
	cur, err := s.c.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var tasks []*Task
	err = cur.All(ctx, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *taskStore) Count(ctx context.Context) (int64, error) {
	filter := bson.M{}
	return s.c.CountDocuments(ctx, filter)
}
