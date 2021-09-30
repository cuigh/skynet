package contract

import (
	"fmt"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/ext/times"
)

const (
	CodeSuccess int32 = iota
	CodeFailed
	CodeNotFound
	CodeNotSupported
	CodeTaskIsRunning
)

type Result struct {
	Code int32  `json:"code"`
	Info string `json:"info,omitempty"`
}

type Job struct {
	Id      string       `json:"id"`
	Task    string       `json:"task"`
	Handler string       `json:"handler"`
	Args    data.Options `json:"args"`
	Mode    int32        `json:"mode"` // 0-auto, 1-manual
	Fire    int64        `json:"fire"` // unix milliseconds
}

func (j *Job) String() string {
	return fmt.Sprintf("{mode: %v, Id: %s, task: %s, handler: %s, args: %v, fire: %s}",
		j.Mode, j.Id, j.Task, j.Handler, j.Args, times.FromUnixMilli(j.Fire).Format("2006-01-02 15:04:05.999"))
}

type ExecuteParam struct {
	Name string       `json:"name"`
	Args data.Options `json:"args,omitempty"`
}

type NotifyParam struct {
	Code  int32  `json:"code"`
	Info  string `json:"info,omitempty"`
	Id    string `json:"id,omitempty"`
	Start int64  `json:"start,omitempty"` // unix milliseconds
	End   int64  `json:"end,omitempty"`   // unix milliseconds
}

type SplitResult struct {
	Code    int32  `json:"code"` // 0-成功, 1-失败, 2-不支持拆分
	Info    string `json:"info,omitempty"`
	Batches []*Batch
}

type Batch struct {
	Id   string
	Args data.Options
}
