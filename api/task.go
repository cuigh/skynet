package api

import (
	"time"

	"github.com/cuigh/auxo/app/container"
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/ext/times"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/skynet/contract"
	"github.com/cuigh/skynet/schedule"
	"github.com/cuigh/skynet/store"
)

// TaskHandler is a controller of task.
type TaskHandler struct {
	Search  web.HandlerFunc `path:"/search" auth:"?" desc:"search tasks"`
	Find    web.HandlerFunc `path:"/find" auth:"?" desc:"find task by name"`
	Save    web.HandlerFunc `path:"/save" method:"post" auth:"task.edit" desc:"create or update task"`
	Delete  web.HandlerFunc `path:"/delete" method:"post" auth:"task.delete" desc:"delete task"`
	Execute web.HandlerFunc `path:"/execute" method:"post" auth:"task.exec" desc:"execute task"`
	Notify  web.HandlerFunc `path:"/notify" method:"post" auth:"*" desc:"notify execution result"`
}

// NewTask creates an instance of TaskHandler
func NewTask(store store.TaskStore, js store.JobStore) *TaskHandler {
	return &TaskHandler{
		Search:  taskSearch(store),
		Find:    taskFind(store),
		Save:    taskSave(store),
		Delete:  taskDelete(store),
		Execute: taskExecute(),
		Notify:  taskNotify(js),
	}
}

func taskSearch(ts store.TaskStore) web.HandlerFunc {
	type Args struct {
		Name      string `json:"name"`
		Runner    string `json:"runner"`
		PageIndex int64  `json:"page_index"`
		PageSize  int64  `json:"page_size"`
	}

	return func(ctx web.Context) error {
		args := &Args{}
		err := ctx.Bind(args)
		if err != nil {
			return err
		}

		tasks, total, err := ts.Search(args.Name, args.Runner, args.PageIndex, args.PageSize)
		if err != nil {
			return err
		}
		return success(ctx, data.Map{"items": tasks, "total": total})
	}
}

func taskFind(ts store.TaskStore) web.HandlerFunc {
	return func(ctx web.Context) error {
		name := ctx.Query("name")
		task, err := ts.Find(name)
		if err != nil {
			return err
		}
		return success(ctx, task)
	}
}

func taskSave(ts store.TaskStore) web.HandlerFunc {
	return func(ctx web.Context) error {
		t := &store.Task{}
		err := ctx.Bind(t, true)
		if err == nil {
			if time.Time(t.ModifyTime).IsZero() {
				err = ts.Create(t)
			} else {
				err = ts.Modify(t)
			}
		}
		return ajax(ctx, err)
	}
}

func taskDelete(ts store.TaskStore) web.HandlerFunc {
	return func(ctx web.Context) error {
		t := &store.Task{}
		err := ctx.Bind(t)
		if err == nil {
			err = ts.Delete(t.Name)
		}
		return ajax(ctx, err)
	}
}

func taskExecute() web.HandlerFunc {
	type Args struct {
		Name string       `json:"name"`
		Args data.Options `json:"args"`
	}

	return func(ctx web.Context) error {
		args := &Args{}
		err := ctx.Bind(args)
		if err == nil {
			err = container.Call(func(s *schedule.Scheduler) error {
				return s.Execute(args.Name, args.Args)
			})
		}
		return ajax(ctx, err)
	}
}

func taskNotify(js store.JobStore) web.HandlerFunc {
	type Args struct {
		Code  int32  `json:"code"`
		Info  string `json:"info,omitempty"`
		Id    string `json:"id"`
		Start int64  `json:"start,omitempty"` // unix milliseconds
		End   int64  `json:"end,omitempty"`   // unix milliseconds
	}

	return func(ctx web.Context) error {
		args := &Args{}
		err := ctx.Bind(args)
		if err != nil {
			return err
		}

		start := times.FromUnixMilli(args.Start)
		end := times.FromUnixMilli(args.End)
		if err = js.ModifyExecute(args.Id, args.Code == 0, args.Info, start, end); err != nil {
			return err
		}
		if args.Code != contract.CodeSuccess {
			err = container.Call(func(alerter *schedule.Alerter) {
				go alerter.Alert(args.Id, args.Info)
			})
			if err != nil {
				log.Get("api").Error(err)
			}
		}
		return success(ctx, nil)
	}
}
