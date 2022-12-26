package asyncq

import (
	"fmt"
	"sync"
	"time"

	"github.com/stretchr/testify/mock"
)

type mockQueueExecute struct {
	mock.Mock
	wg      *sync.WaitGroup
	channel chan *Job
}

func (mqe *mockQueueExecute) exec(job *Job) {
	fmt.Print("exec is called")
	_ = mqe.Called(job)
	mqe.wg.Done()
	return
}

type mockJobProcess struct {
	mock.Mock
	waitGroup *sync.WaitGroup
	delay     time.Duration
}

func (mjp *mockJobProcess) Process(j *Job) error {
	if mjp.delay > 0 {
		time.Sleep(mjp.delay)
	}
	args := mjp.Called(j)
	mjp.waitGroup.Done()
	return args.Error(0)
}
