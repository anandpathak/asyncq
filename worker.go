package asyncq

type worker struct {
	close chan bool
	pool  chan *Job
}

type workerPool struct {
	workers []*worker
	size    int
	buffer  chan *Job
}

func newWorkerPool(size int, channelBufferSize int) *workerPool {
	var ws []*worker
	for i := 0; i < size; i++ {
		w := &worker{close: make(chan bool)}
		ws = append(ws, w)
	}
	return &workerPool{
		workers: ws,
		size:    size,
		buffer:  make(chan *Job, channelBufferSize),
	}
}

func (wp *workerPool) start(e executor) {
	// start the workers
	for i := 0; i < wp.size; i++ {
		go wp.workers[i].start(e)(wp.buffer)
	}
	log.Infof("started worker pool of size %d", wp.size)
}

func (w *worker) start(e executor) func(Queue <-chan *Job) {
	return func(q <-chan *Job) {
		for {
			select {
			case job, ok := <-q:
				if !ok {
					break
				}
				e.exec(job)
			case <-w.close:
				return
			}
		}
	}

}

func (wp *workerPool) closeWorkers() {
	for i := range wp.workers {
		wp.workers[i].close <- true
	}
	log.Infof("[async] closing worker pool")
}

type executor interface {
	exec(job *Job)
}
