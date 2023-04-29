package utils

import "sync"

type worker interface {
	Task()
}

type workerFunc struct {
	fn func()
}

func (w workerFunc) Task() {
	w.fn()
}

type pool struct {
	works chan worker
	wg    sync.WaitGroup
}

func NewWorkerPool(n int) *pool {
	p := pool{
		works: make(chan worker),
	}

	p.wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			for w := range p.works {
				w.Task()
			}
			p.wg.Done()
		}()
	}

	return &p
}

func (p *pool) Run(w worker) {
	p.works <- w
}

func (p *pool) RunFunc(fn func()) {
	p.works <- workerFunc{fn: fn}
}

func (p *pool) Stop() {
	close(p.works)
	p.wg.Wait()
}
