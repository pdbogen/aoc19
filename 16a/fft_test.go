package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFFTdigit(t *testing.T) {
	tests := []struct {
		name     string
		in       []int
		n        int
		expected int
	}{
		{"given digit 0", []int{1, 2, 3, 4, 5, 6, 7, 8}, 0, 4},
		{"given digit 1", []int{1, 2, 3, 4, 5, 6, 7, 8}, 1, 8},
		{"given digit 2", []int{1, 2, 3, 4, 5, 6, 7, 8}, 2, 2},
		{"given digit 3", []int{1, 2, 3, 4, 5, 6, 7, 8}, 3, 2},
		{"given digit 4", []int{1, 2, 3, 4, 5, 6, 7, 8}, 4, 6},
		{"given digit 5", []int{1, 2, 3, 4, 5, 6, 7, 8}, 5, 1},
		{"given digit 6", []int{1, 2, 3, 4, 5, 6, 7, 8}, 6, 5},
		{"given digit 7", []int{1, 2, 3, 4, 5, 6, 7, 8}, 7, 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := FFTdigit(tt.in, tt.n)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestFFT(t *testing.T) {
	for _, tt := range []struct {
		name     string
		in       []int
		expected []int
		phases   int
	}{
		{"given 0", []int{1, 2, 3, 4, 5, 6, 7, 8}, []int{4, 8, 2, 2, 6, 1, 5, 8}, 1},
		{"given 1", []int{4, 8, 2, 2, 6, 1, 5, 8}, []int{3, 4, 0, 4, 0, 4, 3, 8}, 1},
		{"given 2",
			[]int{8, 0, 8, 7, 1, 2, 2, 4, 5, 8, 5, 9, 1, 4, 5, 4, 6, 6, 1, 9, 0, 8, 3, 2, 1, 8, 6, 4, 5, 5, 9, 5},
			[]int{2, 4, 1, 7, 6, 1, 7, 6, 4, 8, 0, 9, 1, 9, 0, 4, 6, 1, 1, 4, 0, 3, 8, 7, 6, 3, 1, 9, 5, 5, 9, 5},
			100},
	} {
		t.Run(tt.name, func(t *testing.T) {
			actual := FFT(tt.in, tt.phases)
			assert.Equal(t, tt.expected, actual)
		})
	}
}