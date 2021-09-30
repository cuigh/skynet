package lock

import (
	"time"

	"github.com/cuigh/auxo/app"
	"github.com/cuigh/auxo/app/container"
	"github.com/cuigh/auxo/config"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/skynet/store"
)

type Lock interface {
	Lock(name string, fire time.Time) bool
	Unlock(name string, fire time.Time) bool
}

// NullLock is suitable for stand-alone mode.
type NullLock struct {
}

func NewNullLock() Lock {
	return NullLock{}
}

func (n NullLock) Lock(name string, fire time.Time) bool {
	return true
}

func (n NullLock) Unlock(name string, fire time.Time) bool {
	return true
}

func init() {
	app.OnInit(func() error {
		switch config.GetString("skynet.lock") {
		case "redis":
			container.Put(NewRedisLock, container.Name("lock"))
		case "mongo":
			container.Put(store.NewLockStore, container.Name("store.lock"))
			container.Put(NewMongoLock, container.Name("lock"))
		case "", "null":
			container.Put(NewNullLock, container.Name("lock"))
		default:
			return errors.NotSupported
		}
		return nil
	})
}
