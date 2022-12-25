package asyncq

import (
	"errors"
	"sync"
	"time"
)

type Queue struct {
	Retry      int
	BufferSize int
	Func       Processor
	once       sync.Once
	RetryDelay time.Duration
	closed     bool
	pool       *workerPool
	sync.Mutex
}

func New(retry int, bufferSize int, poolSize int, delay time.Duration, processor Processor) *Queue {
	queue := &Queue{
		Retry:      retry,
		BufferSize: bufferSize,
		Func:       processor,
		RetryDelay: delay,
		closed:     false,
		pool:       newWorkerPool(poolSize, bufferSize),
	}
	return queue
}

func (q *Queue) Start() {
	q.pool.start(q)
}

func (q *Queue) Add(j *Job) error {
	q.Lock()
	defer q.Unlock()
	if !q.closed {
		go func(j *Job) {
			q.pool.buffer <- j
		}(j)
		return nil
	}
	return errors.New("channel already closed")
}

func (q *Queue) reQueue(j *Job) {
	if j.retryCount < q.Retry {
		go func(j *Job) {
			select {
			case <-time.After(q.RetryDelay):
				j.Mutex.Lock()
				defer j.Mutex.Unlock()
				j.retryCount++
				q.pool.buffer <- j
			}
		}(j)
	}
}

func (q *Queue) Close() {
	q.once.Do(func() {
		q.Mutex.Lock()
		defer q.Mutex.Unlock()
		q.closed = true
		q.pool.closeWorkers()

		log.Info("[async] closing the queue worker")
	})
}

func (q *Queue) exec(job *Job) {
	log.Infof("[async] job triggered with id=%s params=%v", job.ID, job.Params)
	err := q.Func.Process(job)
	if err != nil {
		log.Errorf("[async] job failed, id=%s , params=%v\n err=%s", job.ID, job.Params, err.Error())
		q.reQueue(job)
		return
	}
	log.Infof("[async] job success for id=%s", job.ID)
}
