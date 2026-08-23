package scheduler

import (
	"context"
	"predictive-maintenance/pkg/logger"
	"sync"
	"time"
)

type Job struct {
	Name     string
	Interval time.Duration
	Run      func(context.Context) error
}
type Runner struct {
	jobs   []Job
	log    *logger.Logger
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func New(log *logger.Logger) *Runner { return &Runner{log: log, jobs: []Job{}} }
func (r *Runner) Add(job Job) {
	if job.Interval <= 0 {
		job.Interval = time.Minute
	}
	if job.Run != nil {
		r.jobs = append(r.jobs, job)
	}
}
func (r *Runner) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	r.cancel = cancel
	for _, job := range r.jobs {
		j := job
		r.wg.Add(1)
		go r.loop(ctx, j)
	}
}
func (r *Runner) loop(ctx context.Context, j Job) {
	defer r.wg.Done()
	ticker := time.NewTicker(j.Interval)
	defer ticker.Stop()
	r.execute(ctx, j)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.execute(ctx, j)
		}
	}
}
func (r *Runner) execute(ctx context.Context, j Job) {
	started := time.Now()
	if e := j.Run(ctx); e != nil {
		r.log.Error("scheduler job failed", map[string]any{"job": j.Name, "error": e.Error()})
	} else {
		r.log.Debug("scheduler job completed", map[string]any{"job": j.Name, "duration": time.Since(started).String()})
	}
}
func (r *Runner) Stop() {
	if r.cancel != nil {
		r.cancel()
	}
	r.wg.Wait()
}
func (r *Runner) JobCount() int { return len(r.jobs) }
