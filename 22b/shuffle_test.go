package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDoCut(t *testing.T) {
	for _, tt := range []struct {
		name     string
		n        int
		in       []int
		expected []int
	}{
		{"Cut(0)", 0, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{"Cut(1)", 1, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}},
		{"Cut(2)", 2, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, []int{2, 3, 4, 5, 6, 7, 8, 9, 0, 1}},
		{"Cut(-1)", -1, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, []int{9, 0, 1, 2, 3, 4, 5, 6, 7, 8}},
		{"Cut(-2)", -2, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, []int{8, 9, 0, 1, 2, 3, 4, 5, 6, 7}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			actual := make([]int, len(tt.in))
			for i, v := range tt.in {
				actual[DoCut(tt.n, i, len(tt.in))] = v
			}
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestUndoCut(t *testing.T) {
	assert.Equal(t, UndoCut(0, 0, 10), 0)
	assert.Equal(t, UndoCut(3, 7, 10), 0)
	assert.Equal(t, UndoCut(-3, 0, 10), 7)
}

func TestDeal(t *testing.T) {
	for _, tt := range []struct {
		name      string
		in        []int
		increment int
		expected  []int
	}{
		{
			"given",
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			3,
			[]int{0, 7, 4, 1, 8, 5, 2, 9, 6, 3},
		},
		{
			"given B",
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			7,
			[]int{0, 3, 6, 9, 2, 5, 8, 1, 4, 7},
		},
		{
			"inverse given",
			[]int{0, 7, 4, 1, 8, 5, 2, 9, 6, 3},
			7,
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			"inverse given B",
			[]int{0, 3, 6, 9, 2, 5, 8, 1, 4, 7},
			3,
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			actual := make([]int, len(tt.in))
			for i := range tt.in {
				actual[DoDeal(tt.increment, i, len(tt.in))] = tt.in[i]
			}
			assert.Equal(t, tt.expected, actual)
			undo := make([]int, len(actual))
			for i := range actual {
				undo[UndoDeal(tt.increment, i, len(tt.in))] = actual[i]
			}
			assert.Equal(t, tt.in, undo)
		})
	}
}

func TestStack_Flip(t *testing.T) {
	for _, tt := range []struct {
		name     string
		in       []int
		expected []int
	}{
		{
			"single flip",
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			[]int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			actual := make([]int, len(tt.in))
			for i := range tt.in {
				actual[DoFlip(i, len(tt.in))] = tt.in[i]
			}
			assert.Equal(t, tt.expected, actual)
			undo := make([]int, len(tt.in))
			for i := range tt.in {
				undo[UndoFlip(i, len(tt.in))] = actual[i]
			}
			assert.Equal(t, tt.in, undo)
		})
	}
}

func BenchmarkUndoShuffle(b *testing.B) {
	ops, err := LoadShuffle("input")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		UndoShuffle(ops, 2020, 119315717514047, 100)
	}
}

func Test_inverse(t *testing.T) {
	tests := []struct {
		name     string
		x, mod   int
		expected int
	}{
		{"inverse of 15 mod 26 is 7", 15, 26, 7},
		{"inverse of 3 mod 10 is 7", 3, 10, 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, inverse(tt.x, tt.mod))
		})
	}
}
