package schedule

import (
	"github.com/cuigh/auxo/app/ioc"
	"time"

	"github.com/cuigh/auxo/app"
	"github.com/cuigh/auxo/config"
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/ext/times"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/skynet/lock"
	"github.com/cuigh/skynet/store"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ModeAuto int32 = iota
	ModeManual
)

// 调用器(Caller)：任务远程调用实现
// 执行器(Runner)：任务执行宿主程序
// 处理器(Handler)：业务处理逻辑
// 任务(Task)：任务执行排期
// 作业(Job)：一次具体的执行任务

type Job struct {
	oid     primitive.ObjectID
	fire    time.Time
	runner  string
	Id      string       `json:"id"`
	Task    string       `json:"task"`
	Handler string       `json:"handler"`
	Args    data.Options `json:"args"`
	Mode    int32        `json:"mode"` // 0-auto, 1-manual
	Fire    int64        `json:"fire"`
}

func NewJob(t *store.Task, args data.Options, mode int32, fire time.Time) *Job {
	id := primitive.NewObjectID()
	return &Job{
		oid:     id,
		fire:    fire,
		runner:  t.Runner,
		Id:      id.Hex(),
		Task:    t.Name,
		Handler: t.Handler,
		Mode:    mode,
		Fire:    times.ToUnixMilli(fire),
		Args:    mergeArgs(t.Args, args),
	}
}

// Scheduler dispatch task to executors to run by plan.
type Scheduler struct {
	node     string
	tf       *TaskFetcher
	th       *TaskHeap
	lock     lock.Lock
	resolver Resolver
	logger   log.Logger
	js       store.JobStore
	updater  chan *TaskHeap
	alerter  *Alerter
	closer   chan struct{}
	callers  map[string]Caller
}

func NewScheduler(lock lock.Lock, resolver Resolver, ts store.TaskStore, js store.JobStore, alerter *Alerter) *Scheduler {
	logger := log.Get("schedule")
	node := config.GetString("skynet.node")
	if node == "" {
		node = primitive.NewObjectID().Hex()[:8]
	}
	return &Scheduler{
		node: node,
		callers: map[string]Caller{
			"http": HTTPCaller{},
		},
		lock:     lock,
		resolver: resolver,
		tf:       NewTaskFetcher(ts, logger),
		js:       js,
		alerter:  alerter,
		updater:  make(chan *TaskHeap, 1),
		closer:   make(chan struct{}),
		logger:   logger,
	}
}

func (s *Scheduler) Start() {
	s.tf.Start(s.updater)

	var t Timer
	defer t.Stop()

	for {
		t.Reset(s.try())
		select {
		case <-t.C:
			continue
		case s.th = <-s.updater:
			s.logger.Info("update tasks")
			continue
		case <-s.closer:
			return
		}
	}
}

func (s *Scheduler) Stop() {
	close(s.closer)
	s.tf.Stop()
}

// Execute dispatches task immediately.
func (s *Scheduler) Execute(name string, args data.Options) error {
	task, err := s.tf.Find(name)
	if err != nil {
		return err
	}
	job := NewJob(task, args, ModeManual, time.Now())
	s.call(job, false)
	return nil
}

func (s *Scheduler) try() time.Duration {
	// sleep for one minute if tasks not fetched
	if s.th == nil {
		return time.Minute
	}

	now := time.Now()
	for {
		item := s.th.Peek()
		if item == nil {
			return time.Minute
		}

		if d := item.fire.Sub(now); d > 0 {
			return d
		}

		job := NewJob(item.task, nil, ModeAuto, item.fire)
		go s.call(job, false)

		// update next fire time of task
		item.next(now)
		s.th.Update(0)
	}
}

func (s *Scheduler) Retry(id string) error {
	j, err := s.js.Find(id)
	if err != nil {
		return err
	}

	t, err := s.tf.Find(j.Task)
	if err != nil {
		return err
	}

	job := &Job{
		//oid:
		fire:    time.Time(j.FireTime),
		runner:  t.Runner,
		Id:      j.Id.Hex(),
		Task:    j.Task,
		Handler: j.Handler,
		Mode:    j.Mode,
		Fire:    times.ToUnixMilli(time.Time(j.FireTime)),
		Args:    j.Args,
	}
	s.call(job, true)
	return nil
}

func (s *Scheduler) call(job *Job, retry bool) {
	if job.Mode == ModeAuto && !s.lock.Lock(job.Task, job.fire) {
		s.logger.Debugf("task {name: %s, fire: %s} was already dispatched by another node", job.Task, job.Fire)
		return
	}

	// save job info
	if !retry {
		err := s.js.Create(&store.Job{
			Id:        job.oid,
			Task:      job.Task,
			Handler:   job.Handler,
			Scheduler: s.node,
			Args:      job.Args,
			Mode:      job.Mode,
			FireTime:  store.Time(job.fire),
		})
		if err != nil {
			s.logger.Errorf("failed to save job to db: %s", err)
			// TODO: terminate dispatch or just ignore err?
			return
		}
	}

	// dispatch
	schema, addrs, err := s.resolver.Resolve(job.runner)
	if err != nil {
		s.logger.Errorf("failed to resolve runner: %s", err)
		return
	}
	caller := s.callers[schema]
	if caller == nil {
		s.logger.Errorf("caller not found: %s", schema)
		return
	}
	result := caller.Call(addrs, job)

	// update control info
	err = s.js.ModifyDispatch(job.Id, result.Success(), result.Info)
	if err != nil {
		s.logger.Errorf("failed to update job control info: %s", err)
	}

	if !result.Success() {
		go s.alerter.Alert(job.Id, result.Info)
	}
}

func mergeArgs(args1, args2 data.Options) data.Options {
	if len(args1) == 0 {
		return args2
	} else if len(args2) == 0 {
		return args1
	}

	var args data.Options
	args = append(args, args1...)
	for _, opt := range args2 {
		var exist bool
		for i := range args {
			if args[i].Name == opt.Name {
				args[i].Value = opt.Value
				exist = true
				break
			}
		}
		if !exist {
			args = append(args, opt)
		}
	}
	return args
}

type TaskFetcher struct {
	modify time.Time // last modify time of tasks
	th     *TaskHeap
	ts     store.TaskStore
	timer  *time.Timer
	logger log.Logger
}

func NewTaskFetcher(ts store.TaskStore, logger log.Logger) *TaskFetcher {
	return &TaskFetcher{
		ts:     ts,
		logger: logger,
	}
}

func (f *TaskFetcher) Find(name string) (*store.Task, error) {
	return f.ts.Find(name)
}

func (f *TaskFetcher) Start(c chan<- *TaskHeap) {
	if f.refresh() {
		c <- f.th
	}

	f.timer = time.AfterFunc(time.Minute, func() {
		if f.refresh() {
			c <- f.th
		}
		f.timer.Reset(time.Minute)
	})
}

func (f *TaskFetcher) Stop() {
	if f.timer != nil {
		f.timer.Stop()
	}
}

func (f *TaskFetcher) refresh() (updated bool) {
	modify, count, err := f.ts.GetState()
	if err != nil {
		f.logger.Error("failed to fetch task state: ", err)
		return false
	}
	if f.th != nil && f.th.Count() == int(count) && f.modify.Equal(modify) {
		return false
	}

	var tasks []*store.Task
	tasks, err = f.ts.FetchAll(true)
	if err != nil {
		f.logger.Error("failed to fetch tasks: ", err)
		return false
	}

	f.th, f.modify = NewTaskHeap(tasks), modify
	return true
}

type Timer struct {
	*time.Timer
}

func (t *Timer) Reset(d time.Duration) {
	if t.Timer == nil {
		t.Timer = time.NewTimer(d)
	} else {
		if !t.Timer.Stop() {
			select {
			case <-t.Timer.C: // try to drain from channel
			default:
			}
		}
		t.Timer.Reset(d)
	}
}

func (t *Timer) Stop() {
	if t.Timer != nil {
		t.Timer.Stop()
	}
}

func init() {
	app.OnInit(func() error {
		// register resolver service
		resolver := config.GetString("skynet.resolver")
		switch resolver {
		case "", "direct":
			ioc.Put(NewDirectResolver, ioc.Name("resolver"))
		default:
			return errors.Format("not supported resolver: %s", resolver)
		}

		ioc.Put(NewScheduler, ioc.Name("scheduler"))
		ioc.Put(NewAlerter, ioc.Name("alerter"))
		return nil
	})
}
