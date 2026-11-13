package main

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

const (
	stQueued   = "queued"
	stRunning  = "running"
	stDone     = "done"
	stFailed   = "failed"
	stCanceled = "canceled"

	maxErrors   = 200 // error messages kept per job
	maxFinished = 100 // finished jobs kept in the list
)

// Job is one queued operation: copy, move or delete.
type Job struct {
	ID       int64
	Op       string
	Sources  []string // virtual paths
	Dest     string   // virtual path
	Conflict string
	Created  time.Time

	mu       sync.Mutex
	state    string
	phase    string
	started  time.Time
	finished time.Time
	errors   []string
	errCount int
	speed    float64
	cancel   context.CancelFunc

	totalBytes, doneBytes atomic.Int64
	totalFiles, doneFiles atomic.Int64
	current               atomic.Pointer[string]
}

type jobView struct {
	ID         int64    `json:"id"`
	Op         string   `json:"op"`
	Sources    []string `json:"sources"`
	Dest       string   `json:"dest"`
	State      string   `json:"state"`
	Phase      string   `json:"phase"`
	TotalBytes int64    `json:"totalBytes"`
	DoneBytes  int64    `json:"doneBytes"`
	TotalFiles int64    `json:"totalFiles"`
	DoneFiles  int64    `json:"doneFiles"`
	Speed      float64  `json:"speed"`
	Current    string   `json:"current"`
	Errors     []string `json:"errors"`
	ErrCount   int      `json:"errCount"`
	Created    int64    `json:"created"`
	Started    int64    `json:"started"`
	Finished   int64    `json:"finished"`
	Elapsed    int64    `json:"elapsed"` // ms, measured by the server clock
}

func ms(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.UnixMilli()
}

func (j *Job) view() jobView {
	j.mu.Lock()
	defer j.mu.Unlock()
	v := jobView{
		ID: j.ID, Op: j.Op, Sources: j.Sources, Dest: j.Dest,
		State: j.state, Phase: j.phase, Speed: j.speed,
		TotalBytes: j.totalBytes.Load(), DoneBytes: j.doneBytes.Load(),
		TotalFiles: j.totalFiles.Load(), DoneFiles: j.doneFiles.Load(),
		Errors: append([]string{}, j.errors...), ErrCount: j.errCount,
		Created: ms(j.Created), Started: ms(j.started), Finished: ms(j.finished),
	}
	if !j.finished.IsZero() && !j.started.IsZero() {
		v.Elapsed = j.finished.Sub(j.started).Milliseconds()
	} else if !j.started.IsZero() {
		v.Elapsed = time.Since(j.started).Milliseconds()
	}
	if c := j.current.Load(); c != nil {
		v.Current = *c
	}
	return v
}

func (j *Job) setPhase(p string) {
	j.mu.Lock()
	j.phase = p
	j.mu.Unlock()
}

func (j *Job) addErr(msg string) {
	j.mu.Lock()
	j.errCount++
	if len(j.errors) < maxErrors {
		j.errors = append(j.errors, msg)
	}
	j.mu.Unlock()
}

// Queue runs jobs strictly one at a time, in the order they were added.
type Queue struct {
	app    *App
	mu     sync.Mutex
	jobs   []*Job
	nextID int64
	wake   chan struct{}
}

func newQueue(app *App) *Queue {
	return &Queue{app: app, wake: make(chan struct{}, 1)}
}

func (q *Queue) add(op string, sources []string, dest, conflict string) *Job {
	q.mu.Lock()
	q.nextID++
	j := &Job{ID: q.nextID, Op: op, Sources: sources, Dest: dest, Conflict: conflict,
		Created: time.Now(), state: stQueued}
	q.jobs = append(q.jobs, j)
	// Drop the oldest finished jobs beyond the limit.
	finished := 0
	for _, x := range q.jobs {
		if x.isFinished() {
			finished++
		}
	}
	if finished > maxFinished {
		kept := q.jobs[:0]
		for _, x := range q.jobs {
			if finished > maxFinished && x.isFinished() {
				finished--
				continue
			}
			kept = append(kept, x)
		}
		q.jobs = kept
	}
	q.mu.Unlock()
	select {
	case q.wake <- struct{}{}:
	default:
	}
	return j
}

func (j *Job) isFinished() bool {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.state != stQueued && j.state != stRunning
}

func (q *Queue) views() []jobView {
	q.mu.Lock()
	jobs := append([]*Job{}, q.jobs...)
	q.mu.Unlock()
	out := make([]jobView, len(jobs))
	for i, j := range jobs {
		out[i] = j.view()
	}
	return out
}

func (q *Queue) find(id int64) *Job {
	q.mu.Lock()
	defer q.mu.Unlock()
	for _, j := range q.jobs {
		if j.ID == id {
			return j
		}
	}
	return nil
}

func (q *Queue) cancel(id int64) {
	j := q.find(id)
	if j == nil {
		return
	}
	j.mu.Lock()
	defer j.mu.Unlock()
	switch j.state {
	case stQueued:
		j.state = stCanceled
		j.finished = time.Now()
	case stRunning:
		j.cancel()
	}
}

// remove deletes a finished job from the list; id 0 removes all finished jobs.
func (q *Queue) remove(id int64) {
	q.mu.Lock()
	defer q.mu.Unlock()
	kept := q.jobs[:0]
	for _, j := range q.jobs {
		if (id == 0 || j.ID == id) && j.isFinished() {
			continue
		}
		kept = append(kept, j)
	}
	q.jobs = kept
}

func (q *Queue) next() (*Job, context.Context) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for _, j := range q.jobs {
		j.mu.Lock()
		if j.state == stQueued {
			ctx, cancel := context.WithCancel(context.Background())
			j.state, j.started, j.cancel = stRunning, time.Now(), cancel
			j.mu.Unlock()
			return j, ctx
		}
		j.mu.Unlock()
	}
	return nil, nil
}

func (q *Queue) loop() {
	for {
		j, ctx := q.next()
		if j == nil {
			<-q.wake
			continue
		}
		q.exec(j, ctx)
	}
}

func (q *Queue) exec(j *Job, ctx context.Context) {
	stop := make(chan struct{})
	go j.sampleSpeed(stop)

	r := &runner{app: q.app, job: j, ctx: ctx}
	var err error
	if j.Op == "delete" {
		err = r.runDelete()
	} else {
		err = r.runTransfer()
	}
	close(stop)

	j.mu.Lock()
	switch {
	case ctx.Err() != nil:
		j.state = stCanceled
	case err != nil:
		j.state = stFailed
		j.errCount++
		j.errors = append(j.errors, err.Error())
	default:
		j.state = stDone
	}
	j.finished, j.speed, j.phase = time.Now(), 0, ""
	j.cancel()
	j.mu.Unlock()
	j.current.Store(nil)
}

// sampleSpeed keeps a moving average of the transfer rate over the last few seconds.
func (j *Job) sampleSpeed(stop <-chan struct{}) {
	const window = 5
	samples := []int64{j.doneBytes.Load()}
	t := time.NewTicker(time.Second)
	defer t.Stop()
	for {
		select {
		case <-stop:
			return
		case <-t.C:
			samples = append(samples, j.doneBytes.Load())
			if len(samples) > window+1 {
				samples = samples[1:]
			}
			sp := float64(samples[len(samples)-1]-samples[0]) / float64(len(samples)-1)
			j.mu.Lock()
			j.speed = sp
			j.mu.Unlock()
		}
	}
}
