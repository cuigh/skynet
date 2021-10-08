package runner

import (
	"fmt"
	"sync"
	"time"

	"github.com/cuigh/auxo/app"
	"github.com/cuigh/auxo/app/container"
	"github.com/cuigh/auxo/ext/times"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/util/run"
	"github.com/cuigh/skynet/client"
	"github.com/cuigh/skynet/contract"
)

var (
	handlers = make(map[string]Handler)
	running  = sync.Map{}
)

type PreFilter func(job *contract.Job) error

type PostFilter func(job *contract.Job, err error)

type Handler interface {
	Handle(job *contract.Job) error
}

type ParallelHandler interface {
	Handler
	Split(job *contract.Job) ([]*contract.Batch, error)
	//Map()
	//Reduce()
	//Merge()
}

type HandlerFunc func(job *contract.Job) error

func (f HandlerFunc) Handle(job *contract.Job) error {
	return f(job)
}

func Register(name string, handler Handler) {
	handlers[name] = handler
}

func RegisterFunc(name string, handler func(job *contract.Job) error) {
	handlers[name] = HandlerFunc(handler)
}

func Serve(ws *web.Server) func(ctx *app.Context) error {
	ws.Post("/task/execute", HandleExecute)
	ws.Post("/task/split", HandleSplit)
	return func(ctx *app.Context) error {
		app.Run(ws)
		return nil
	}
}

func HandleExecute(ctx web.Context) error {
	var job contract.Job
	err := ctx.Bind(&job)
	if err != nil {
		return ctx.JSON(contract.Result{Code: contract.CodeFailed, Info: err.Error()})
	}

	go handle(&job)
	return ctx.JSON(contract.Result{})
}

func HandleSplit(ctx web.Context) error {
	var job contract.Job
	err := ctx.Bind(&job)
	if err != nil {
		return ctx.JSON(contract.Result{Code: contract.CodeFailed, Info: err.Error()})
	}

	result := split(&job)
	return ctx.JSON(result)
}

func handle(job *contract.Job) {
	log.Get("task").Debugf("handle job: %s", job)

	start := time.Now()

	handler := handlers[job.Handler]
	if handler == nil {
		notify(job, start, contract.CodeNotFound, "handler not found")
		return
	}

	if job.Mode == 0 {
		if last, exist := running.LoadOrStore(job.Task, job); exist {
			fire := times.FromUnixMilli(last.(*contract.Job).Fire)
			notify(job, start, contract.CodeTaskIsRunning, fmt.Sprintf("task is already running(fire: %s)",
				fire.Format("2006-01-02 15:04:05")))
			return
		}
		defer running.Delete(job.Task)
	}

	run.Safe(func() {
		err := handler.Handle(job)
		if err != nil {
			notify(job, start, contract.CodeFailed, err.Error())
		} else {
			notify(job, start, contract.CodeSuccess, "")
		}
	}, func(e interface{}) {
		notify(job, start, contract.CodeFailed, fmt.Sprint(e))
	})
}

func notify(job *contract.Job, start time.Time, code int32, info string) {
	param := contract.NotifyParam{
		Code:  code,
		Info:  info,
		Id:    job.Id,
		Start: times.ToUnixMilli(start),
		End:   times.ToUnixMilli(time.Now()),
	}
	err := container.Call(func(client *client.Client) error {
		return client.Notify(param)
	})
	if err != nil {
		log.Get("task").Errorf("failed to notify result of job(%s): %s", job.Id, err)
	}
}

func split(job *contract.Job) *contract.SplitResult {
	log.Get("task").Debugf("split job: %s", job)

	h := handlers[job.Handler]
	if h == nil {
		return &contract.SplitResult{Code: contract.CodeNotFound, Info: "handler not found"}
	}

	ph, ok := h.(ParallelHandler)
	if !ok {
		return &contract.SplitResult{Code: contract.CodeNotSupported, Info: "not supported"}
	}

	batches, err := ph.Split(job)
	if err != nil {
		return &contract.SplitResult{Code: contract.CodeFailed, Info: err.Error()}
	}
	return &contract.SplitResult{Code: contract.CodeSuccess, Batches: batches}
}
