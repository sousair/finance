package workers

import (
	"context"
	"fmt"
	"time"
)

type (
	CronWorker interface {
		Register(job *CronJob)
		Start(ctx context.Context)
	}

	CronJob struct {
		Name     string
		Fn       func(ctx context.Context) error
		Interval time.Duration
	}
)

type cronWorker struct {
	jobs []*CronJob
}

func NewCronWorker() *cronWorker {
	return &cronWorker{
		jobs: []*CronJob{},
	}
}

func (cw *cronWorker) Register(job CronJob) {
	cw.jobs = append(cw.jobs, &job)
}

func (cw *cronWorker) Start(ctx context.Context) {
	for _, job := range cw.jobs {
		go func(job *CronJob) {
			for {
				select {
				case <-ctx.Done():
					return
				case <-time.After(job.Interval):
					fmt.Printf("[CronWorker] Running job (%s) \n", job.Name)

					if err := job.Fn(ctx); err != nil {
						fmt.Printf("[CronWorker] Job (%s) returned an error: %v \n", job.Name, err)
					}

					fmt.Printf("[CronWorker] Job (%s) finished \n", job.Name)
				}
			}
		}(job)
	}
}
