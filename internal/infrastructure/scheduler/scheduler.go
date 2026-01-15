package scheduler

import (
	"context"
	"log/slog"
	"time"
)

// Job represents a scheduled task
type Job struct {
	Name     string
	Interval time.Duration
	Run      func(ctx context.Context) error
}

// Scheduler manages and runs periodic jobs
type Scheduler struct {
	jobs   []Job
	logger *slog.Logger
}

// NewScheduler creates a new Scheduler instance
func NewScheduler(logger *slog.Logger) *Scheduler {
	return &Scheduler{
		jobs:   make([]Job, 0),
		logger: logger,
	}
}

// Register adds a job to the scheduler
func (s *Scheduler) Register(job Job) {
	s.jobs = append(s.jobs, job)
}

// Start begins running all registered jobs
// It blocks until the context is cancelled
func (s *Scheduler) Start(ctx context.Context) {
	for _, job := range s.jobs {
		go s.runJob(ctx, job)
	}

	<-ctx.Done()
	s.logger.Info("Scheduler stopped")
}

// runJob executes a single job on its interval
func (s *Scheduler) runJob(ctx context.Context, job Job) {
	ticker := time.NewTicker(job.Interval)
	defer ticker.Stop()

	s.logger.Info("Job started", "name", job.Name, "interval", job.Interval)

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Job stopped", "name", job.Name)
			return
		case <-ticker.C:
			if err := job.Run(ctx); err != nil {
				s.logger.Error("Job failed", "name", job.Name, "error", err)
			} else {
				s.logger.Debug("Job completed", "name", job.Name)
			}
		}
	}
}
