package schedule

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"sort"

	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/auxo/net/web"
)

// Caller format: http://abc, simple://
type Caller interface {
	// Call dispatches task to remote runner.
	Call(addrs []string, t *Job) *CallResult
	// Split asks remote runner how to split task.
	//Split(addrs []string, t *Job) SplitResult
}

type CallResult struct {
	Code int32  `json:"code"`
	Info string `json:"info,omitempty"`
}

func (r *CallResult) Success() bool {
	return r.Code == 0
}

type HTTPCaller struct {
}

func (c HTTPCaller) Call(addrs []string, j *Job) (r *CallResult) {
	addrs = shuffle(addrs)
	for _, addr := range addrs {
		r = c.call(addr+"/task/execute", j)
		if r.Success() {
			return
		}
		log.Get("schedule").Errorf("call with address '%s' failed: %s", addr, r.Info)
	}
	return
}

func (c HTTPCaller) call(addr string, j *Job) *CallResult {
	//return &CallResult{
	//	Code: 1,
	//	Info: fmt.Sprintf("not implemented: %s-%s", j.Task, times.Format(j.Fire, "yyyyMMddHHmmss")),
	//}

	var r CallResult
	if err := c.do(addr, j, &r); err != nil {
		r.Code = 1
		r.Info = err.Error()
	}
	return &r
}

func (c HTTPCaller) do(url string, args, result interface{}) error {
	data, err := json.Marshal(args)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, web.MIMEApplicationJSONCharsetUTF8, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	d := json.NewDecoder(resp.Body)
	return d.Decode(result)
}

func shuffle(addrs []string) []string {
	if l := len(addrs); l > 1 {
		arr := make([]string, len(addrs))
		copy(arr, addrs)
		sort.Slice(arr, func(i, j int) bool {
			return rand.Intn(2) > 0
		})
		return arr
	}
	return addrs
}
