# asyncq
Async Job processor


## Quick start
 - Install package using
```go get -u github.com/anandpathak/asynq```
 - Below is the example of creating client for async processing 
```go
    import (
		"github.com/anandpathak/asyncq"
)

 func main(){
	 consumer := eventConsumer{}
    // this will create a queue 
	 queue := asyncq.New(retryCount, bufferSize, poolSize , delay, consumer)
	 // start the queue
	 queue.Start()
	 
}

type eventConsumer struct {}
func (p eventConsumer)Process(j *asyncq.Job ) error {
	// add your consumer logic in Process
	return nil
}
```

 - You can create jobs and pass to the async queue. async queue will pull message, and process
```go
     // create the job to process
	 job := asyncq.Job{
        Params:[]string{"hello", "world"},
        ID: "testJob1",
    }
	// Add job to queue
	queue.Add(&job)
```
 - You can close the consumer to stop receiving the message
```go
    queue.Close()
```
