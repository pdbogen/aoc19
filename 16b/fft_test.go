package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func BenchmarkFFT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FFTdigit(
			Repeat([]int{0, 3, 0, 3, 6, 7, 3, 2, 5, 7, 7, 2, 1, 2, 9, 4, 4, 0, 6, 3, 4, 9, 1, 5, 6, 5, 4, 7, 4, 6, 6, 4}, 10000),
			0,
			100,
		)
	}
}

func TestFFTdigit(t *testing.T) {
	tests := []struct {
		name     string
		in       []int
		n        int
		pass     int
		expected int
	}{
		{"given digit 0", []int{1, 2, 3, 4, 5, 6, 7, 8}, 0, 1, 4},
		{"given digit 1", []int{1, 2, 3, 4, 5, 6, 7, 8}, 1, 1, 8},
		{"given digit 2", []int{1, 2, 3, 4, 5, 6, 7, 8}, 2, 1, 2},
		{"given digit 3", []int{1, 2, 3, 4, 5, 6, 7, 8}, 3, 1, 2},
		{"given digit 4", []int{1, 2, 3, 4, 5, 6, 7, 8}, 4, 1, 6},
		{"given digit 5", []int{1, 2, 3, 4, 5, 6, 7, 8}, 5, 1, 1},
		{"given digit 6", []int{1, 2, 3, 4, 5, 6, 7, 8}, 6, 1, 5},
		{"given digit 7", []int{1, 2, 3, 4, 5, 6, 7, 8}, 7, 1, 8},
		{"given digit 0", []int{1, 2, 3, 4, 5, 6, 7, 8}, 0, 2, 3},
		{"given digit 1", []int{1, 2, 3, 4, 5, 6, 7, 8}, 1, 2, 4},
		{"given digit 2", []int{1, 2, 3, 4, 5, 6, 7, 8}, 2, 2, 0},
		{"given digit 3", []int{1, 2, 3, 4, 5, 6, 7, 8}, 3, 2, 4},
		{"given digit 4", []int{1, 2, 3, 4, 5, 6, 7, 8}, 4, 2, 0},
		{"given digit 5", []int{1, 2, 3, 4, 5, 6, 7, 8}, 5, 2, 4},
		{"given digit 6", []int{1, 2, 3, 4, 5, 6, 7, 8}, 6, 2, 3},
		{"given digit 7", []int{1, 2, 3, 4, 5, 6, 7, 8}, 7, 2, 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := FFTdigit(tt.in, tt.n, tt.pass)
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
		{"given 0", []int{1, 2, 3, 4, 5, 6, 7, 8}, []int{3, 4, 0, 4, 0, 4, 3, 8}, 2},
		//{"given 2",
		//	[]int{8, 0, 8, 7, 1, 2, 2, 4, 5, 8, 5, 9, 1, 4, 5, 4, 6, 6, 1, 9, 0, 8, 3, 2, 1, 8, 6, 4, 5, 5, 9, 5},
		//	[]int{2, 4, 1, 7, 6, 1, 7, 6, 4, 8, 0, 9, 1, 9, 0, 4, 6, 1, 1, 4, 0, 3, 8, 7, 6, 3, 1, 9, 5, 5, 9, 5},
		//	2},
	} {
		//cache = &sync.Map{}
		cache = map[cacheKey]int{}
		t.Run(tt.name, func(t *testing.T) {
			actual := FFT(tt.in, tt.phases)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

//
//func TestDecode(t *testing.T) {
//	for _, tt := range []struct {
//		name     string
//		in       []int
//		expected []int
//	}{
//		{"given 1",
//			[]int{0, 3, 0, 3, 6, 7, 3, 2, 5, 7, 7, 2, 1, 2, 9, 4, 4, 0, 6, 3, 4, 9, 1, 5, 6, 5, 4, 7, 4, 6, 6, 4},
//			[]int{8, 4, 4, 6, 2, 0, 2, 6}},
//	} {
//		t.Run(tt.name, func(t *testing.T) {
//			actual := Decode(tt.in)
//			assert.Equal(t, tt.expected, actual)
//		})
//	}
//}
