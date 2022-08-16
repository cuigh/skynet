package lock

import (
	"github.com/cuigh/auxo/app/ioc"
	"time"

	"github.com/cuigh/auxo/app"
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
			ioc.Put(NewRedisLock, ioc.Name("lock"))
		case "mongo":
			ioc.Put(store.NewLockStore, ioc.Name("store.lock"))
			ioc.Put(NewMongoLock, ioc.Name("lock"))
		case "", "null":
			ioc.Put(NewNullLock, ioc.Name("lock"))
		default:
			return errors.NotSupported
		}
		return nil
	})
}
