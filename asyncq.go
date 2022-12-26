package asyncq

import (
	"errors"
	"sync"
	"time"
)

// Queue is a structure represent async event Processor.
type Queue struct {
	// Retry represent the number of time each `Job` in the `Queue` will be retried, if the `Job` fails
	Retry int
	// BufferSize is the maximum `Job` that can be in queue for processing .
	BufferSize int
	// Func accept a method of Processor type, this method is performed on each Job.
	Func Processor
	once sync.Once
	// RetryDelay represent the delay between Job retries. Currently, Queue Only support linear `Retry` Delay
	RetryDelay time.Duration
	closed     bool
	pool       *workerPool
	sync.Mutex
}

// New returns a async queue instance
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

// Start initiate the Queue processor, It will allow Queue to accept `Job`
// Job can only be added to Queue Once Start is triggered
func (q *Queue) Start() {
	q.pool.start(q)
}

// Add adds a Job to the Queue for processing
// Once the Queue is closed Add to the Queue will not be allowed
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

// Close the queue process so that no new message will get processed
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
