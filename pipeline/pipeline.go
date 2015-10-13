package pipeline

import (
	"sync"
)

type (
	Pipeline interface {
		Err() <-chan error
		Output() <-chan interface{}
	}
	Handler  func(interface{}) (interface{}, error)
	pipeline struct {
		source <-chan interface{}
		dest   chan interface{}
		err    chan error
		handle Handler
		group  sync.WaitGroup
	}
)

func NewPipeline(source <-chan interface{}, workerCount int, handles ...Handler) Pipeline {
	var pipe Pipeline
	err := make(chan error)
	for _, handle := range handles {
		pipe = newPipeline(source, err, workerCount, handle)
		source = pipe.Output()
	}
	return pipe
}

func newPipeline(source <-chan interface{}, err chan error, workerCount int, handle Handler) Pipeline {
	pipe := &pipeline{
		source: source,
		dest:   make(chan interface{}),
		err:    err,
		handle: handle,
		group:  sync.WaitGroup{},
	}

	for i := 0; i < workerCount; i++ {
		pipe.group.Add(1)
		go pipe.work()
	}
	go pipe.waitForGroup()
	return pipe
}

func (pipe *pipeline) Err() <-chan error {
	return pipe.err
}

func (pipe *pipeline) Output() <-chan interface{} {
	return pipe.dest
}

func (pipe *pipeline) waitForGroup() {
	pipe.group.Wait()
	close(pipe.dest)
}

func (pipe *pipeline) work() {
	defer pipe.group.Done()
	for {
		job, more := <-pipe.source
		if !more {
			return
		}
		res, err := pipe.handle(job)
		if err != nil {
			pipe.err <- err
			continue
		}

		pipe.dest <- res
	}
}
