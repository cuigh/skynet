package api

import (
	"context"
	"runtime"
	"time"

	"github.com/cuigh/auxo/app"
	"github.com/cuigh/auxo/app/ioc"
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/security/passwd"
	"github.com/cuigh/auxo/util/run"
	"github.com/cuigh/skynet/store"
)

//var ErrSystemInitialized = errors.New("system was already initialized")

// SystemHandler encapsulates system related handlers.
type SystemHandler struct {
	CheckState web.HandlerFunc `path:"/check-state" auth:"*" desc:"check system state"`
	InitDB     web.HandlerFunc `path:"/init-db" method:"post" auth:"*" desc:"initialize database"`
	InitUser   web.HandlerFunc `path:"/init-user" method:"post" auth:"*" desc:"initialize administrator account"`
	Summarize  web.HandlerFunc `path:"/summarize" auth:"?" desc:"fetch statistics data"`
}

// NewSystem creates an instance of SystemHandler
func NewSystem() *SystemHandler {
	return &SystemHandler{
		CheckState: systemCheckState,
		InitDB:     systemInitDB,
		InitUser:   systemInitUser,
		Summarize:  systemSummarize,
	}
}

func systemCheckState(c web.Context) error {
	var count int64
	err := ioc.Call(func(us store.UserStore) (err error) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		count, err = us.Count(ctx)
		return err
	})
	if err != nil {
		return err
	}
	return success(c, data.Map{"fresh": count == 0})
}

func systemInitDB(ctx web.Context) error {
	return ajax(ctx, ioc.Call(func(js store.JobStore, ls store.LockStore, us store.UserStore) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		return run.Pipeline(
			func() error { return js.CreateIndexes(ctx) },
			func() error { return ls.CreateIndexes(ctx) },
			func() error { return us.CreateIndexes(ctx) },
		)
	}))
}

func systemInitUser(c web.Context) error {
	type Args struct {
		Name      string `json:"name" valid:"required"`
		LoginName string `json:"login_name" valid:"required"`
		Email     string `json:"email,omitempty" valid:"email"`
		Phone     string `json:"phone,omitempty"`
		Password  string `json:"password"`
	}

	return ajax(c, ioc.Call(func(us store.UserStore) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		count, err := us.Count(ctx)
		if err != nil {
			return err
		} else if count > 0 {
			return errors.New("system was already initialized")
		}

		args := &Args{}
		err = c.Bind(args, true)
		if err == nil {
			u := &store.User{
				Name:      args.Name,
				LoginName: args.LoginName,
				Email:     args.Email,
				Phone:     args.Phone,
				Admin:     true,
				Status:    store.UserStatusNormal,
			}
			u.Password, u.Salt, err = passwd.Generate(args.Password)
			if err == nil {
				err = us.Create(u)
			}
		}
		return err
	}))
}

func systemSummarize(c web.Context) error {
	summary := struct {
		Version   string `json:"version"`
		GoVersion string `json:"goVersion"`
		TaskCount int64  `json:"taskCount"`
		JobCount  int64  `json:"jobCount"`
		UserCount int64  `json:"userCount"`
	}{
		Version:   app.Version,
		GoVersion: runtime.Version(),
	}

	err := ioc.Call(func(ts store.TaskStore, js store.JobStore, us store.UserStore) (err error) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		summary.TaskCount, err = ts.Count(ctx)
		if err != nil {
			return
		}

		summary.JobCount, err = js.Count(ctx)
		if err != nil {
			return
		}

		summary.UserCount, err = us.Count(ctx)
		return
	})
	if err != nil {
		return err
	}
	return success(c, summary)
}
