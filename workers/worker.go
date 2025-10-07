package workers

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/yourname/blog-kafka/function"
	"github.com/yourname/blog-kafka/models"
)

type NotificationJob struct {
	ChannelID uuid.UUID
	AuthorID  uuid.UUID
	BlogID    uuid.UUID
	BlogTitle string
	Type      models.NotificationType
}

type WorkerPool struct {
	JobQueue      chan NotificationJob
	MaxWorkers    int
	IdleTimeout   time.Duration
	activeWorkers int
	wg            sync.WaitGroup
	mu            sync.Mutex
}

func NewWorkerPool(maxWorkers int, queueSize int, idleTimeout time.Duration) *WorkerPool {
	return &WorkerPool{
		JobQueue:    make(chan NotificationJob, queueSize),
		MaxWorkers:  maxWorkers,
		IdleTimeout: idleTimeout,
	}
}

func (wp *WorkerPool) Submit(job NotificationJob) {
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
			fmt.Printf("Worker %d processing job: %+v\n", id, job)
			err := function.CreateNotificationForChannelMembers(job.ChannelID, job.AuthorID, job.BlogID, job.BlogTitle, job.Type)
			if err != nil {
				fmt.Printf("Worker %d failed to create notification: %v\n", id, err)
			}

		case <-time.After(wp.IdleTimeout):
			fmt.Printf("Worker %d idle timeout, stopping\n", id)
			wp.mu.Lock()
			wp.activeWorkers--
			wp.mu.Unlock()
			return
		}
	}
}

// Wait waits for all active workers to finish
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}
