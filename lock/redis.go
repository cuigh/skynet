package lock

import (
	"fmt"
	"time"

	"github.com/cuigh/auxo/data/guid"
	"github.com/cuigh/auxo/db/redis"
	"github.com/cuigh/auxo/log"
)

type RedisLock struct {
	client     redis.Client
	dispatcher string
	logger     log.Logger
}

func NewRedisLock() Lock {
	client, err := redis.Open("skynet")
	if err != nil {
		panic(err)
	}
	return &RedisLock{
		client:     client,
		dispatcher: guid.New().String(),
		logger:     log.Get("lock"),
	}
}

func (l *RedisLock) Lock(name string, fire time.Time) bool {
	key := l.key(name, fire)
	ok, err := l.client.SetNX(key, l.dispatcher, 5*time.Minute).Result()
	if err != nil {
		l.logger.Errorf("failed to lock job(%s:%s): %s", name, fire, err)
	}
	return ok
}

func (l *RedisLock) Unlock(name string, fire time.Time) bool {
	key := l.key(name, fire)
	err := l.client.Del(key).Err()
	if err != nil {
		l.logger.Errorf("failed to unlock job(%s:%s): %s", name, fire, err)
	}
	return err == nil
}

func (l *RedisLock) key(name string, fire time.Time) string {
	id := fire.Format("20060102150405")
	return fmt.Sprintf("task:%s:%s", name, id)
}
