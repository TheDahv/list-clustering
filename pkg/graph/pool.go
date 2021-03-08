package graph

import (
	"fmt"
	"sync"

	"github.com/thedahv/list-clustering/pkg/rbo"
)

// Pool is a worker pool
type Pool struct {
	concurrency int
	errors      chan error
	wg          sync.WaitGroup
	work        chan Task
	results     chan Edge

	edges []Edge
	errs  []error
}

// Task is some work going into the pool
type Task struct {
	Source rbo.Node
	Target rbo.Node
	pVal   float64
}

// NewPool initializes a pool
func NewPool(concurrency int) *Pool {
	p := &Pool{
		concurrency: concurrency,
		errors:      make(chan error, 10),
		work:        make(chan Task, 100),
		results:     make(chan Edge, 50),
	}

	for i := 0; i < concurrency; i++ {
		go p.startWorker()
	}

	go p.collectResults()
	go p.collectErrors()

	return p
}

func (p *Pool) startWorker() {
	for t := range p.work {
		_, _, ext, err := rbo.RBO(t.Source, t.Target, t.pVal)
		if err != nil {
			p.errors <- fmt.Errorf("could not compute RBO for %s -> %s",
				t.Source.Label(), t.Target.Label())
		} else {
			p.results <- Edge{
				Source:     t.Source.Label(),
				Target:     t.Target.Label(),
				similarity: ext,
			}
		}
	}
}

func (p *Pool) collectResults() {
	for edge := range p.results {
		p.edges = append(p.edges, edge)
		p.wg.Done()
	}
}

func (p *Pool) collectErrors() {
	for err := range p.errors {
		p.errs = append(p.errs, err)
		p.wg.Done()
	}
}

// Add adds work to the pool
func (p *Pool) Add(task Task) {
	p.wg.Add(1)
	p.work <- task
}

// DoneAdding indicates all the work is in
func (p *Pool) DoneAdding() {
	close(p.work)
}

// Results blocks until all results are collected or the first error comes in
func (p *Pool) Results() ([]Edge, error) {
	var err error
	p.wg.Wait()
	if len(p.errors) > 0 {
		err = p.errs[0]
	}
	close(p.results)
	return p.edges, err
}
