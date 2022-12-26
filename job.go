package asyncq

import "sync"

// JobState represent the Job state in queue
type JobState string

const (
	// ProcessingJobState represents that the Job is still in process in the queue
	ProcessingJobState JobState = "PROCESSING"
	// SuccessJobState represent that the job is successfully processed from the queue
	SuccessJobState JobState = "SUCCESS"
	// FailureJobState represent that the job is failed to process
	FailureJobState JobState = "FAILURE"
)

// Job structure hold the information for each message that needs to be processed from the queue
type Job struct {
	// Params are the parameters that is passed into the Processor method
	Params     []string // TODO allow interface type
	retryCount int
	// ID is unique identifier for a Job passed into the queue
	ID string
	sync.Mutex
}

// Processor is interface type to allow client to define Process method
// Queue will trigger process on a Job that is pushed to the queue
type Processor interface {
	Process(*Job) error
}
