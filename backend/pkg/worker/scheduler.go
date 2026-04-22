package worker

import (
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

// Scheduler provides a fluent API for task scheduling.
type Scheduler struct {
	instance *asynq.Scheduler
}

// NewScheduler creates a new Scheduler instance.
func NewScheduler(redisAddr string, redisPassword string, redisDB int) *Scheduler {
	s := asynq.NewScheduler(asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	}, &asynq.SchedulerOpts{
		Location: time.Local,
	})

	return &Scheduler{instance: s}
}

// JobBuilder helps build a scheduled job using a fluent API.
type JobBuilder struct {
	scheduler *Scheduler
	cron      string
	task      *asynq.Task
}

// Every starts building a new scheduled job.
// example: scheduler.Every(1).Hour().Do(task)
func (s *Scheduler) Every(interval int) *JobBuilder {
	return &JobBuilder{
		scheduler: s,
	}
}

// Minute sets the interval to minutes.
func (jb *JobBuilder) Minute() *JobBuilder {
	jb.cron = "*/1 * * * *" // Simplified for MVP, real impl would use interval
	return jb
}

// Hour sets the interval to hours.
func (jb *JobBuilder) Hour() *JobBuilder {
	jb.cron = "0 * * * *"
	return jb
}

// Day sets the interval to days.
func (jb *JobBuilder) Day() *JobBuilder {
	jb.cron = "0 0 * * *"
	return jb
}

// At sets a specific time for the daily job.
func (jb *JobBuilder) At(timeStr string) *JobBuilder {
	// timeStr example: "10:00"
	var hour, min int
	fmt.Sscanf(timeStr, "%d:%d", &hour, &min)
	jb.cron = fmt.Sprintf("%d %d * * *", min, hour)
	return jb
}

// Cron allows setting a raw cron expression.
func (jb *JobBuilder) Cron(expr string) *JobBuilder {
	jb.cron = expr
	return jb
}

// Do schedules the task to be executed.
func (jb *JobBuilder) Do(taskType string, payload []byte) (string, error) {
	jb.task = asynq.NewTask(taskType, payload)
	return jb.scheduler.instance.Register(jb.cron, jb.task)
}

// Run starts the scheduler.
func (s *Scheduler) Run() error {
	return s.instance.Run()
}

// Shutdown stops the scheduler.
func (s *Scheduler) Shutdown() {
	s.instance.Shutdown()
}
