package notifications

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/yourname/blog-kafka/models"
)

type NotificationJob struct {
	UserID    uuid.UUID
	BlogID    uuid.UUID
	BlogTitle string
	Type      models.NotificationType
}

type NotificationHandler interface {
	CreateSingleNotification(job NotificationJob) error
}

type WorkerPool struct {
	JobQueue      chan NotificationJob
	MaxWorkers    int
	IdleTimeout   time.Duration
	activeWorkers int
	wg            sync.WaitGroup
	mu            sync.Mutex
	handler       NotificationHandler
}

func NewWorkerPool(queueSize int, idleTimeout time.Duration, handler NotificationHandler) *WorkerPool {
	return &WorkerPool{
		JobQueue:    make(chan NotificationJob, queueSize),
		MaxWorkers:  5,
		IdleTimeout: idleTimeout,
		handler:     handler,
	}
}

func (wp *WorkerPool) Submit(job NotificationJob) {
	if wp == nil || wp.JobQueue == nil {
		fmt.Println("WorkerPool not initialized properly")
		return
	}
	wp.JobQueue <- job

	wp.mu.Lock()
	if wp.activeWorkers < wp.MaxWorkers {
		wp.activeWorkers++
		wp.wg.Add(1)
		go wp.worker(wp.activeWorkers)
	}
	wp.mu.Unlock()
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	fmt.Printf("Worker %d started\n", id)

	for {
		select {
		case job := <-wp.JobQueue:
			err := wp.handler.CreateSingleNotification(job)
			if err != nil {
				fmt.Printf("Worker %d error: %v\n", id, err)
			}
		case <-time.After(wp.IdleTimeout):
			fmt.Printf("Worker %d idle timeout — stopping\n", id)
			wp.mu.Lock()
			wp.activeWorkers--
			wp.mu.Unlock()
			return
		}
	}
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}
