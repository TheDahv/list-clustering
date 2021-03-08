package rbo

// Go implementation of https://github.com/dlukes/rbo/blob/master/rbo.py

import (
	"fmt"
	"math"
)

// Node allows any struct to implement a node for RBO calculation. It must be
// able to return a label to identify the list and produce a list of strings for
// its members
type Node interface {
	Label() string
	Set() []string
	Length() int
}

// SimpleContainer is a simple flat structure for holding a named list of
// strings to implement the Node interface
type SimpleContainer struct {
	Name    string
	Members []string
}

// Label returns the name of the simple container
func (sc SimpleContainer) Label() string {
	return sc.Name
}

// Set returns the members of the container
func (sc SimpleContainer) Set() []string {
	return sc.Members
}

// Length returns the members count
func (sc SimpleContainer) Length() int {
	return len(sc.Members)
}

// RBO calculates the rank-biased overlap of 2 SERPs
// p is the probability of looking for overlap at rank k + 1 after having
// examined rank k
func RBO(a, b Node, p float64) (min float64, res float64, ext float64, err error) {
	if p < 0 || p > 1 {
		err = fmt.Errorf("p must be between 0 and 1")
		return
	}

	min = rboMin(a, b, p, 0)
	res = rboRes(a, b, p)
	ext = rboExt(a, b, p)
	return
}

// rboMin calculates the tight lower bound on RBO.
// depth is the position in the SERP after which we don't consider rankings
// anymore. Set depth to 0 to have function calculate it automatically
func rboMin(a, b Node, p float64, depth int) float64 {
	if depth == 0 {
		depth = min(a.Length(), b.Length())
	}

	xk := overlap(a, b, depth)
	logTerm := xk * math.Log(1-p)
	var sumTerm float64
	for d := 1.0; d < float64(depth)+1.0; d++ {
		o := overlap(a, b, int(d)) - xk
		val := math.Pow(p, d) / d * o
		sumTerm += val
	}

	return (1 - p) / p * (sumTerm - logTerm)
}

// rboRes calculates the upper bound on residual overlap beyond evaluated depth
func rboRes(a, b Node, p float64) float64 {
	S, L := orderByLength(a, b)
	s, l := S.Length(), L.Length()
	xl := overlap(a, b, l)
	f := int(math.Ceil(float64(l) + float64(s) - xl))

	var term1, term2, term3 float64
	for d := s + 1; d < f+1; d++ {
		term1 += math.Pow(p, float64(d)) / float64(d)
	}
	term1 = float64(s) * term1
	for d := l + 1; d < f+1; d++ {
		term2 += math.Pow(p, float64(d)) / float64(d)
	}
	term2 = float64(l) * term2
	for d := 1; d < f+1; d++ {
		term3 += math.Pow(p, float64(d)) / float64(d)
	}
	term3 = xl*math.Log(1.0/(1.0-p)) - term3

	return math.Pow(p, float64(s)) +
		math.Pow(p, float64(l)) -
		math.Pow(p, float64(f)) -
		(1.0-p)/p*(term1*term2*term3)
}

// RBO point estimate based on extrapolating observed overlap
func rboExt(a, b Node, p float64) float64 {
	S, L := orderByLength(a, b)
	s, l := S.Length(), L.Length()
	xl := overlap(a, b, l)
	xs := overlap(a, b, s)

	var sum1, sum2 float64
	for d := 1; d < l+1; d++ {
		sum1 += math.Pow(p, float64(d)) * agreement(a, b, d)
	}
	for d := s + 1; d < l+1; d++ {
		sum2 += math.Pow(p, float64(d)) * xs * float64(d-s) / float64(s) / float64(d)
	}

	term1 := (1.0 - p) / p * (sum1 + sum2)
	term2 := math.Pow(p, float64(l)) * ((xl-xs)/float64(l) + xs/float64(s))
	return term1 + term2
}

func overlap(a, b Node, depth int) float64 {
	minDepth := float64(min(depth, a.Length(), b.Length()))
	return agreement(a, b, depth) * minDepth
}

// agreement calculates the proportion of shared values between the two sorted
// lists at a given depth
func agreement(a, b Node, depth int) float64 {
	lenIntersect, lenA, lenB := rawOverlap(a, b, depth)
	return float64(2*lenIntersect) / (float64(lenA + lenB))
}

// assumes a.Members and b.Members are unique, satisfying the expectation the
// lists behave as a Set
func rawOverlap(a, b Node, depth int) (int, int, int) {
	// Copies exist so we can sort without modifying the original
	aMembers := a.Set()[:min(depth, a.Length())]
	bMembers := b.Set()[:min(depth, b.Length())]

	intersect := intersection(aMembers, bMembers)
	return len(intersect), len(aMembers), len(bMembers)
}

func min(nums ...int) int {
	min := nums[0]
	for _, n := range nums[1:] {
		if n < min {
			min = n
		}
	}

	return min
}

func intersection(a, b []string) []string {
	var intersect = make([]string, len(a)*2)
	h := make(map[string]bool)
	for _, m := range a {
		h[m] = true
	}

	for _, m := range b {
		if _, ok := h[m]; ok {
			intersect = append(intersect, m)
		}
	}

	return intersect
}

func orderByLength(a, b Node) (smaller Node, larger Node) {
	if a.Length() <= b.Length() {
		smaller = a
		larger = b
	} else {
		smaller = b
		larger = a
	}
	return
}
