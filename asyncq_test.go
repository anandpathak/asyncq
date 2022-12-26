package asyncq

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
}

func (s *Suite) TestSetup() {
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestQueueSingleMessage() {
	mockJobProcess := &mockJobProcess{
		waitGroup: &sync.WaitGroup{},
	}
	q := New(2, 3, 1, 1*time.Second, mockJobProcess)
	q.Start()
	job := &Job{
		Params: nil,
	}
	mockJobProcess.waitGroup.Add(1)
	mockJobProcess.On("Process", job).Return(nil)
	expectedErr := q.Add(job)
	defer q.Close()
	assert.Equal(s.T(), nil, expectedErr)
	mockJobProcess.waitGroup.Wait()

}

func (s *Suite) TestShouldReturnErrorWhenChannelClosed() {
	mockJobProcess := &mockJobProcess{
		waitGroup: &sync.WaitGroup{},
	}
	q := New(2, 2, 1, 3*time.Second, mockJobProcess)
	q.Start()
	job := &Job{
		Params: nil,
	}
	q.Close()
	mockJobProcess.AssertNotCalled(s.T(), "Process", job)
	assert.Error(s.T(), q.Add(job))
	defer q.Close()
	mockJobProcess.waitGroup.Wait()
}
func (s *Suite) TestShouldProcessInBatchWhenQueueFull() {
	mockJobProcess := &mockJobProcess{
		waitGroup: &sync.WaitGroup{},
		delay:     1 * time.Second,
	}
	q := New(2, 2, 6, 3*time.Second, mockJobProcess)
	q.Start()

	timeStart := time.Now()
	for i := 0; i < 12; i++ {
		job := &Job{
			Params: nil,
		}
		mockJobProcess.waitGroup.Add(1)
		mockJobProcess.On("Process", job).Return(nil)
		assert.Equal(s.T(), nil, q.Add(job))
	}
	mockJobProcess.waitGroup.Wait()
	timeEnd := time.Now()
	fmt.Print(int64(timeEnd.Sub(timeStart) / time.Second))
	assert.True(s.T(), timeEnd.Sub(timeStart)*time.Second >= 2)
}

func (s *Suite) TestShouldRetryWhenProcessFails() {
	mockJobProcess := &mockJobProcess{
		waitGroup: &sync.WaitGroup{},
		delay:     1 * time.Second,
	}
	q := New(3, 2, 6, 1*time.Second, mockJobProcess)
	q.Start()
	job := &Job{
		Params: nil,
	}
	mockJobProcess.waitGroup.Add(3)
	mockJobProcess.On("Process", job).Return(errors.New("some error"))
	assert.Equal(s.T(), nil, q.Add(job))
	mockJobProcess.waitGroup.Wait()
	assert.Equal(s.T(), 2, job.retryCount)
}
