package store

import (
	"crypto/md5"
	"fmt"
	"github.com/cuigh/auxo/app/ioc"
	mongodb "github.com/cuigh/auxo/db/mongo"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/ext/times"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"strconv"
	"time"
)

func NewDatabase() *mongo.Database {
	return mongodb.MustOpen("skynet")
}

// generate 8-chars short id, only suitable for small dataset
func createId() string {
	id := [12]byte(primitive.NewObjectID())
	return fmt.Sprintf("%x", md5.Sum(id[:]))[:8]
}

type Time time.Time

func (t Time) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(time.Time(t))
}

func (t *Time) UnmarshalBSONValue(bt bsontype.Type, data []byte) error {
	if v, _, valid := bsoncore.ReadValue(data, bt); valid {
		*t = Time(v.Time())
		return nil
	}
	return errors.Format("unmarshal failed, type: %s, data:%s", bt, data)
}

func (t Time) MarshalJSON() (b []byte, err error) {
	return strconv.AppendInt(b, times.ToUnixMilli(time.Time(t)), 10), nil
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	i, err := strconv.ParseInt(string(data), 10, 64)
	if err == nil {
		*t = Time(times.FromUnixMilli(i))
	}
	return err
}

func (t Time) String() string {
	return time.Time(t).String()
}

func (t Time) Format(layout string) string {
	return time.Time(t).Format(layout)
}

func init() {
	ioc.Put(NewDatabase, ioc.Name("database"))
	ioc.Put(NewUserStore, ioc.Name("store.user"))
	ioc.Put(NewTaskStore, ioc.Name("store.task"))
	ioc.Put(NewJobStore, ioc.Name("store.job"))
	ioc.Put(NewRoleStore, ioc.Name("store.role"))
	ioc.Put(NewConfigStore, ioc.Name("store.config"))
}
