package store

import (
	"context"
	"time"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Job struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Task      string             `json:"task" bson:"task"`
	Handler   string             `json:"handler" bson:"handler"`
	Scheduler string             `json:"scheduler" bson:"scheduler"`
	Mode      int32              `json:"mode" bson:"mode"` // 0-Auto, 1-Manual
	Args      data.Options       `json:"args" bson:"args"`
	FireTime  Time               `json:"fire_time" bson:"fire_time"`
	Dispatch  struct {
		Status int32  `json:"status" bson:"status"` // 0-Unknown，1-Success，2-Failed
		Time   *Time  `json:"time,omitempty" bson:"time,omitempty"`
		Error  string `json:"error,omitempty" bson:"error,omitempty"`
	} `json:"dispatch" bson:"dispatch"`
	Execute struct {
		Status    int32  `json:"status" bson:"status"` // 0-Unknown，1-Success，2-Failed
		Error     string `json:"error,omitempty" bson:"error,omitempty"`
		StartTime *Time  `json:"start_time,omitempty" bson:"start_time,omitempty"`
		EndTime   *Time  `json:"end_time,omitempty" bson:"end_time,omitempty"`
	} `json:"execute" bson:"execute"`
}

type JobStore interface {
	Find(id string) (*Job, error)
	Search(task string, mode int32, dispatchStatus, executeStatus int32, pageIndex, pageSize int64) (jobs []*Job, total int64, err error)
	Create(job *Job) error
	ModifyDispatch(id string, success bool, error string) error
	ModifyExecute(id string, success bool, error string, start, end time.Time) error
	CreateIndexes(ctx context.Context) error
	Count(ctx context.Context) (int64, error)
}

func NewJobStore(db *mongo.Database) JobStore {
	return &jobStore{
		c: db.Collection("job"),
	}
}

type jobStore struct {
	c *mongo.Collection
}

func (s *jobStore) Search(task string, mode int32, dispatchStatus, executeStatus int32, pageIndex, pageSize int64) (jobs []*Job, total int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	if task != "" {
		filter["task"] = task
	}
	if mode != -1 {
		filter["mode"] = mode
	}
	if dispatchStatus != -1 {
		filter["dispatch.status"] = dispatchStatus
	}
	if executeStatus != -1 {
		filter["execute.status"] = executeStatus
	}

	// fetch total count
	total, err = s.c.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().SetSkip(pageSize * (pageIndex - 1)).SetLimit(pageSize).SetSort(bson.M{"_id": -1})
	cur, err := s.c.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	err = cur.All(ctx, &jobs)
	if err != nil {
		return nil, 0, err
	}
	return jobs, total, nil
}

func (s *jobStore) Find(id string) (*Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	r := s.c.FindOne(ctx, bson.M{"_id": oid})
	j := &Job{}
	if err := r.Decode(&j); err != nil {
		return nil, err
	}
	return j, nil
}

func (s *jobStore) Create(job *Job) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := s.c.InsertOne(ctx, job)
	return err
}

func (s *jobStore) ModifyDispatch(id string, success bool, error string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"dispatch.status": s.status(success),
		"dispatch.error":  error,
		"dispatch.time":   time.Now(),
	}
	r, err := s.c.UpdateByID(ctx, oid, bson.M{"$set": update})
	if err != nil {
		return err
	} else if r.MatchedCount == 0 {
		return errors.Format("can't find job '%s'", id)
	}
	return nil
}

func (s *jobStore) ModifyExecute(id string, success bool, error string, start, end time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"execute.status":     s.status(success),
		"execute.error":      error,
		"execute.start_time": start,
		"execute.end_time":   end,
	}
	r, err := s.c.UpdateByID(ctx, oid, bson.M{"$set": update})
	if err != nil {
		return err
	} else if r.MatchedCount == 0 {
		return errors.Format("can't find job '%s'", id)
	}
	return nil
}

func (s *jobStore) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{"task", 1}},
		},
		{
			Keys:    bson.D{{"fire_time", 1}},
			Options: options.Index().SetExpireAfterSeconds(3600 * 24 * 7),
		},
	}
	_, err := s.c.Indexes().CreateMany(ctx, indexes)
	return err
}

func (s *jobStore) Count(ctx context.Context) (int64, error) {
	filter := bson.M{}
	return s.c.CountDocuments(ctx, filter)
}

func (s *jobStore) status(success bool) int32 {
	if success {
		return 1
	}
	return 2
}
