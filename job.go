package asyncq

import "sync"

type JobState string

const (
	ProcessingJobState JobState = "PROCESSING"
	SuccessJobState    JobState = "SUCCESS"
	FailureJobState    JobState = "FAILURE"
)

type Job struct {
	Params     []string
	retryCount int
	ID         string
	sync.Mutex
}

type Processor interface {
	Process(*Job) error
}
