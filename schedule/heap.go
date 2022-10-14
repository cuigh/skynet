package schedule

import (
	"time"

	"github.com/cuigh/auxo/log"
	"github.com/cuigh/skynet/store"
	"github.com/robfig/cron/v3"
)

var cronParser = cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

// An TaskItem is something we manage in a priority queue.
type TaskItem struct {
	fire     time.Time
	task     *store.Task
	triggers []cron.Schedule
}

func NewItem(task *store.Task) (*TaskItem, error) {
	item := &TaskItem{
		task: task,
	}
	for _, c := range task.Triggers {
		t, err := cronParser.Parse(c)
		if err != nil {
			return nil, err
		}
		item.triggers = append(item.triggers, t)
	}
	return item, nil
}

// update fire time
func (i *TaskItem) next(start time.Time) {
	var fire time.Time
	for _, t := range i.triggers {
		next := t.Next(start)
		if fire.IsZero() || next.Before(fire) {
			fire = next
		}
	}
	if fire.IsZero() {
		// avoid endless loop
		fire = start.AddDate(100, 0, 0)
	}
	i.fire = fire
}

// A TaskHeap implements minimum heap and holds tasks.
type TaskHeap struct {
	items []*TaskItem
}

func NewTaskHeap(tasks []*store.Task) *TaskHeap {
	now := time.Now()
	items := make([]*TaskItem, len(tasks))
	for i, task := range tasks {
		item, err := NewItem(task)
		if err != nil {
			log.Get("schedule").Errorf("failed to create TaskItem: %s", err)
			continue
		}

		item.next(now)
		items[i] = item
	}

	h := &TaskHeap{items: items}
	h.init()
	return h
}

func (h *TaskHeap) Count() int { return len(h.items) }

func (h *TaskHeap) Push(item *TaskItem) {
	n := len(h.items)
	h.items = append(h.items, item)
	h.up(n)
}

func (h *TaskHeap) Peek() *TaskItem {
	if len(h.items) == 0 {
		return nil
	}
	return h.items[0]
}

func (h *TaskHeap) Pop() *TaskItem {
	n := len(h.items) - 1
	h.swap(0, n)
	h.down(0, n)

	item := h.items[n]
	h.items[n] = nil
	h.items = h.items[:n]
	return item
}

func (h *TaskHeap) Update(i int) {
	if !h.down(i, h.Count()) {
		h.up(i)
	}
}

func (h *TaskHeap) init() {
	n := len(h.items)
	for i := n/2 - 1; i >= 0; i-- {
		h.down(i, n)
	}
}

func (h *TaskHeap) less(i, j int) bool {
	return h.items[i].fire.Before(h.items[j].fire)
}

func (h *TaskHeap) swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
}

func (h *TaskHeap) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.less(j, i) {
			break
		}
		h.swap(i, j)
		j = i
	}
}

func (h *TaskHeap) down(start, n int) bool {
	i := start
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.less(j, i) {
			break
		}
		h.swap(i, j)
		i = j
	}
	return i > start
}
