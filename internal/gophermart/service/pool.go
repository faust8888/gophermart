package service

import (
	"context"
	"github.com/faust8888/gophermart/internal/middleware/logger"
	"sync"
	"time"
)

type WorkerPool struct {
	registerAccrualJobs    chan FetchOrderAccrualStatusJob
	wg                     sync.WaitGroup
	poolSize               int
	fixedScheduleInSeconds int
	selectLimit            int
}

func (p *WorkerPool) Run(jobFunc func(FetchOrderAccrualStatusJob)) {
	for w := 1; w <= p.poolSize; w++ {
		p.wg.Add(1)
		go p.worker(p.registerAccrualJobs, jobFunc)
	}
	ticker := time.NewTicker(time.Duration(p.fixedScheduleInSeconds) * time.Second)
	defer ticker.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("worker pool shutting down")
			cancel()
			close(p.registerAccrualJobs)
			p.wg.Wait()
			logger.Log.Info("all workers don")
		case <-ticker.C:
			for i := 1; i <= p.poolSize; i++ {
				p.startFetchingOrderAccrualStatusAsync(p.selectLimit, i)
			}
		}
	}
}

func (p *WorkerPool) startFetchingOrderAccrualStatusAsync(selectLimit, taskNumber int) {
	p.registerAccrualJobs <- FetchOrderAccrualStatusJob{selectLimit: selectLimit, taskNumber: taskNumber}
}

func NewWorkerPool(poolSize, fixedScheduleInSeconds, selectLimit int) *WorkerPool {
	return &WorkerPool{
		registerAccrualJobs:    make(chan FetchOrderAccrualStatusJob, 100),
		wg:                     sync.WaitGroup{},
		poolSize:               poolSize,
		fixedScheduleInSeconds: fixedScheduleInSeconds,
		selectLimit:            selectLimit,
	}
}

func (p *WorkerPool) worker(jobs <-chan FetchOrderAccrualStatusJob, jobFunc func(FetchOrderAccrualStatusJob)) {
	defer p.wg.Done()
	for job := range jobs {
		jobFunc(job)
	}
}

type JobFunc func(FetchOrderAccrualStatusJob)

type FetchOrderAccrualStatusJob struct {
	taskNumber  int
	selectLimit int
}
