package lock

import (
	"fmt"
	"time"

	"github.com/cuigh/auxo/data/guid"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/skynet/store"
)

type MongoLock struct {
	s          store.LockStore
	dispatcher string
	logger     log.Logger
}

func NewMongoLock(s store.LockStore) Lock {
	return &MongoLock{
		s:          s,
		dispatcher: guid.New().String(),
		logger:     log.Get("lock"),
	}
}

func (l *MongoLock) Lock(name string, fire time.Time) bool {
	key := l.key(name, fire)
	err := l.s.Create(key, l.dispatcher)
	if err != nil {
		l.logger.Errorf("failed to lock job(%s:%s): %s", name, fire, err)
		return false
	}
	return true
}

func (l *MongoLock) Unlock(name string, fire time.Time) bool {
	key := l.key(name, fire)
	err := l.s.Delete(key)
	if err != nil {
		l.logger.Errorf("failed to unlock job(%s:%s): %s", name, fire, err)
		return false
	}
	return true
}

func (l *MongoLock) key(name string, fire time.Time) string {
	id := fire.Format("20060102150405")
	return fmt.Sprintf("%s-%s", name, id)
}
