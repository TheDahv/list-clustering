package graph

import (
	"fmt"
	"testing"

	"github.com/thedahv/list-clustering/pkg/rbo"
)

func TestComputeParallel(t *testing.T) {
	t.Skip()

	tt := []struct {
		Name          string
		PVal          float64
		Inputs        []rbo.Node
		ExpectedCount int
	}{
		{
			Name: "identical lists clustered together",
			PVal: .9,
			Inputs: []rbo.Node{
				rbo.SimpleContainer{
					Name: "a",
					Members: []string{
						"a",
						"b",
						"c",
						"d",
						"e",
					},
				},
				rbo.SimpleContainer{
					Name: "b",
					Members: []string{
						"a",
						"b",
						"c",
						"d",
						"e",
					},
				},
			},
			ExpectedCount: 1,
		},
	}

	for _, tc := range tt {
		edges, err := ComputeParallel(tc.PVal, tc.Inputs)
		if err != nil {
			t.Errorf("got error: %v", err)
		}

		if len(edges) != len(tc.Inputs) {
			t.Errorf("got %d edges, expected %d", len(edges), len(tc.Inputs))
		}
		//fmt.Println(edges)
	}
}

func TestComputePool(t *testing.T) {
	tt := []struct {
		Name               string
		PVal               float64
		InputCount         int
		ExpectedSimilarity float64
	}{
		{
			Name:               "identical lists clustered together",
			PVal:               .9,
			InputCount:         100,
			ExpectedSimilarity: 1.0,
		},
	}

	for _, tc := range tt {
		var inputs []rbo.Node
		for i := 0; i < tc.InputCount; i++ {
			inputs = append(
				inputs,
				rbo.SimpleContainer{
					Name: fmt.Sprintf("%d", i),
					Members: []string{
						"a",
						"b",
						"c",
						"d",
						"e",
					},
				},
			)
		}

		edges, err := ComputePool(tc.PVal, inputs)
		if err != nil {
			t.Errorf("got error: %v", err)
		}

		if l := len(edges); l == 0 {
			t.Errorf("got %d edges, expected more", l)
		}
		// fmt.Println(edges)
	}
}

func runBenchmark(size int, b *testing.B) {
	members := []string{
		"a",
		"b",
		"c",
		"d",
		"e",
	}
	labels := "abcdefghijklmnopqrstuvwxyz"

	var inputs []rbo.Node
	for i := 0; i < size; i++ {
		inputs = append(
			inputs,
			rbo.SimpleContainer{
				Name:    string(labels[i%len(labels)]),
				Members: members,
			},
		)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ComputePool(0.9, inputs)
	}

}

func BenchmarkCompute100(b *testing.B) {
	b.Skip()
	runBenchmark(100, b)
}
func BenchmarkCompute1000(b *testing.B) {
	runBenchmark(1000, b)
}

/*
func BenchmarkCompute5000(b *testing.B) {
	runBenchmark(5000, b)
}
*/
