package rbo

import (
	"testing"
)

type Container struct {
	label   string
	members []string
}

func (c Container) Label() string {
	return c.label
}
func (c Container) Set() []string {
	return c.members
}
func (c Container) Length() int {
	return len(c.members)
}

func TestRBO(t *testing.T) {
	a := Container{
		label: "a",
		members: []string{
			"c",
			"a",
			"b",
			"d",
		},
	}
	b := Container{
		label: "b",
		members: []string{
			"a",
			"c",
			"b",
			"d",
		},
	}
	min, res, ext, err := RBO(a, b, .9)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	t.Logf("min: %f\n", min)
	t.Logf("res: %f\n", res)
	t.Logf("ext: %f\n", ext)
}

func TestAgreement(t *testing.T) {
	a := Container{
		label: "a",
		members: []string{
			"a",
			"b",
			"c",
			"d",
			"e",
		},
	}
	b := Container{
		label: "b",
		members: []string{
			"a",
			"b",
			"d",
			"c",
			"f",
		},
	}

	tt := []struct {
		name     string
		depth    int
		expected float64
	}{
		{
			name:     "at depth 1",
			depth:    1,
			expected: 1.0,
		},
		{
			name:     "at depth 3",
			depth:    3,
			expected: (2.0 / 3.0),
		},
		{
			name:     "at depth 4",
			depth:    4,
			expected: 1.0,
		},
		{
			name:     "at depth 5",
			depth:    5,
			expected: 0.8,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if actual := agreement(a, b, tc.depth); tc.expected != actual {
				t.Errorf("expected %f, got %f", tc.expected, actual)
			}
		})
	}
}

func TestOverlap(t *testing.T) {
	a := Container{
		label: "a",
		members: []string{
			"a",
			"b",
			"c",
			"d",
		},
	}
	b := Container{
		label: "b",
		members: []string{
			"a",
			"b",
			"c",
			"d",
		},
	}

	type testCase struct {
		name     string
		depth    int
		expected float64
	}
	tt := []testCase{
		{
			name:     "depth at 3",
			depth:    3,
			expected: 3.0,
		},
		{
			name:     "depth at 5",
			depth:    5,
			expected: 4.0,
		},
	}

	for _, tc := range tt {
		func(tc testCase) {
			t.Run(tc.name, func(t *testing.T) {
				if actual := overlap(a, b, tc.depth); tc.expected != actual {
					t.Errorf("expected %f, got %f", tc.expected, actual)
				}
			})
		}(tc)
	}
}

func TestIntersection(t *testing.T) {
	a := []string{
		"one",
		"two",
		"three",
	}
	b := []string{
		"two",
		"three",
		"four",
	}

	intersect := intersection(a, b)
	if l := len(intersect); l != 2 {
		t.Errorf("expected 2 entries, got %d", l)
	}

	for _, entry := range intersect {
		if !(entry == "two" || entry == "three") {
			t.Errorf("got unexpected entry: %s", entry)
		}
	}
}

func TestMin(t *testing.T) {
	tt := []struct {
		name     string
		input    []int
		expected int
	}{
		{
			name:     "positive ascending",
			input:    []int{1, 2, 3, 4, 5},
			expected: 1,
		},
		{
			name:     "positive descending",
			input:    []int{5, 4, 3, 2, 1},
			expected: 1,
		},
		{
			name:     "negative descending",
			input:    []int{-1, -2, -3, -4, -5},
			expected: -5,
		},
		{
			name:     "negative ascending",
			input:    []int{-5, -4, -3, -2, -1},
			expected: -5,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual := min(tc.input...)
			if actual != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, actual)
			}
		})
	}
}
