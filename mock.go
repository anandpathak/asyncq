package asyncq

import (
	"fmt"
	"github.com/stretchr/testify/mock"
	"sync"
	"time"
)

type MockQueueExecute struct {
	mock.Mock
	wg      *sync.WaitGroup
	channel chan *Job
}

func (mqe *MockQueueExecute) exec(job *Job) {
	fmt.Print("exec is called")
	_ = mqe.Called(job)
	mqe.wg.Done()
	return
}

type MockJobProcess struct {
	mock.Mock
	waitGroup *sync.WaitGroup
	delay     time.Duration
}

func (mjp *MockJobProcess) Process(j *Job) error {
	if mjp.delay > 0 {
		time.Sleep(mjp.delay)
	}
	args := mjp.Called(j)
	mjp.waitGroup.Done()
	return args.Error(0)
}
