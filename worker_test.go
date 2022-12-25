package asyncq

import (
	"github.com/stretchr/testify/mock"
	"sync"
	"testing"
)

func TestWorkerPool(t *testing.T) {
	t.Run("should process job", func(t *testing.T) {
		w := newWorkerPool(2, 2)
		mockJob := &MockQueueExecute{wg: &sync.WaitGroup{}, channel: make(chan *Job, 2)}
		job := &Job{}
		mockJob.wg.Add(1)
		mockJob.On("exec", mock.Anything).Return(nil)
		w.start(mockJob)
		w.buffer <- job
		w.closeWorkers()
		mockJob.wg.Wait()
	})
}
