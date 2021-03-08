package graph

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/thedahv/list-clustering/pkg/rbo"
)

// Edge is a relationship between two Nodes with a weight to indicate the
// RBO similarity of the sets of the two nodes
type Edge struct {
	Source     string
	Target     string
	similarity float64
}

// Compute calculates the RBO similarity among all nodes and returns weighted
// edges of the graph
func Compute(p float64, nodes []rbo.Node) ([]Edge, error) {
	var edges []Edge

	for i, source := range nodes {
		for j, target := range nodes {
			if i == j {
				continue
			}

			_, _, ext, err := rbo.RBO(source, target, p)
			if err != nil {
				return edges, fmt.Errorf("unable to compute edge for %s -> %s: %v",
					source.Label(), target.Label(), err)
			}

			edges = append(edges, Edge{Source: source.Label(), Target: target.Label(), similarity: ext})
		}
	}

	return edges, nil
}

type pair struct {
	source int
	target int
}

// ComputeParallel calculates the RBO similarity among all nodes and returns weighted
// edges of the graph
func ComputeParallel(p float64, nodes []rbo.Node) ([]Edge, error) {
	var edges []Edge
	maxWorkers := runtime.NumCPU()

	var wg sync.WaitGroup
	work := make(chan pair, 5)
	results := make(chan Edge, 5)

	for i := 0; i < maxWorkers; i++ {
		go func() {
			for w := range work {
				source := nodes[w.source]
				target := nodes[w.target]
				_, _, ext, err := rbo.RBO(source, target, p)

				if err != nil {
					fmt.Printf("unable to compute edge for %s -> %s: %v",
						source.Label(), target.Label(), err)
				}
				results <- Edge{
					Source:     source.Label(),
					Target:     target.Label(),
					similarity: ext,
				}
				wg.Add(1)
			}

		}()
	}

	var jobs int
	// TODO prevent mirrored comparisons (a->b === b->a)
	for i := range nodes {
		for j := i + 1; j < len(nodes); j++ {
			if i == j {
				continue
			}

			jobs++
			p := pair{source: i, target: j}
			work <- p
		}
	}
	close(work)

	/*
		for i := 0; i < jobs; i++ {
			edges = append(edges, <-results)
		}
	*/
	go func() {
		for edge := range results {
			edges = append(edges, edge)
			wg.Done()
		}
	}()

	wg.Wait()

	return edges, nil
}

// ComputePool computes with a pool
func ComputePool(p float64, nodes []rbo.Node) ([]Edge, error) {
	pool := NewPool(50)
	for i := range nodes {
		for j := range nodes {
			if i == j {
				continue
			}

			pool.Add(Task{
				Source: nodes[i],
				Target: nodes[j],
				pVal:   p,
			})
		}
	}
	pool.DoneAdding()

	return pool.Results()
}
