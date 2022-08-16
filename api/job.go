package api

import (
	"github.com/cuigh/auxo/app/ioc"
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/skynet/schedule"
	"github.com/cuigh/skynet/store"
)

// JobHandler encapsulates job related handlers.
type JobHandler struct {
	Search web.HandlerFunc `path:"/search" auth:"?" desc:"search jobs"`
	Find   web.HandlerFunc `path:"/find" auth:"?" desc:"find job by id"`
	Retry  web.HandlerFunc `path:"/retry" method:"post" auth:"job.exec" desc:"retry job"`
}

// NewJob creates an instance of JobHandler
func NewJob(store store.JobStore) *JobHandler {
	return &JobHandler{
		Search: jobSearch(store),
		Find:   jobFind(store),
		Retry:  jobRetry(),
	}
}

func jobSearch(s store.JobStore) web.HandlerFunc {
	type Args struct {
		Task           string `json:"task"`
		Mode           int32  `json:"mode"`
		DispatchStatus int32  `json:"dispatch_status"`
		ExecuteStatus  int32  `json:"execute_status"`
		PageIndex      int64  `json:"page_index"`
		PageSize       int64  `json:"page_size"`
	}

	return func(ctx web.Context) error {
		args := &Args{}
		err := ctx.Bind(args)
		if err != nil {
			return err
		}

		jobs, total, err := s.Search(args.Task, args.Mode, args.DispatchStatus, args.ExecuteStatus, args.PageIndex, args.PageSize)
		if err != nil {
			return err
		}
		return success(ctx, data.Map{"items": jobs, "total": total})
	}
}

func jobFind(s store.JobStore) web.HandlerFunc {
	return func(ctx web.Context) error {
		id := ctx.Query("id")
		job, err := s.Find(id)
		if err != nil {
			return err
		}

		return success(ctx, job)
	}
}

func jobRetry() web.HandlerFunc {
	type Args struct {
		Id string `json:"id"`
	}

	return func(ctx web.Context) error {
		args := &Args{}
		err := ctx.Bind(args)
		if err == nil {
			err = ioc.Call(func(s *schedule.Scheduler) error {
				return s.Retry(args.Id)
			})
		}
		return ajax(ctx, err)
	}
}
