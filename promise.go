package promise

import "sync"

// Promise ...
type Promise struct {
	mutex    *sync.Mutex
	executed bool
	task     func(*Promise)
	res      chan interface{}
	err      chan error
}

// Resolve ...
func (p *Promise) Resolve(res interface{}) {
	p.res <- res
}

// Reject ...
func (p *Promise) Reject(err error) {
	p.err <- err
}

// Await ...
func (p *Promise) Await() (interface{}, error) {
	return p.await(true)
}

func (p *Promise) await(execute bool) (interface{}, error) {
	if execute {
		p.execute()
	}

	select {
	case result := <-p.res:
		return result, nil
	case err := <-p.err:
		return nil, err
	}
}

func (p *Promise) execute() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.executed == false {
		go func() {
			p.task(p)
		}()
		p.executed = true
	}
}

// AwaitAll ...
func AwaitAll(plist []*Promise) ([]interface{}, error) {
	results := make([]interface{}, len(plist))

	// first kick off all promises
	for _, p := range plist {
		p.execute()
	}

	// then await each and collect the result or return immediately on first error
	for i, p := range plist {
		res, err := p.await(false)
		if err != nil {
			return nil, err
		}

		results[i] = res
	}

	return results, nil
}

// New ...
func New(task func(*Promise)) *Promise {
	return &Promise{
		mutex:    new(sync.Mutex),
		executed: false,
		task:     task,
		res:      make(chan interface{}),
		err:      make(chan error),
	}
}
