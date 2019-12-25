package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStack_Cut(t *testing.T) {
	tests := []struct {
		name     string
		in       []int
		n        int
		expected []int
	}{
		{
			"given",
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			0,
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			"given",
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			3,
			[]int{3, 4, 5, 6, 7, 8, 9, 0, 1, 2},
		},
		{
			"given negative",
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			-3,
			[]int{7, 8, 9, 0, 1, 2, 3, 4, 5, 6},
		},
		{
			"cut 6",
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			6,
			[]int{6, 7, 8, 9, 0, 1, 2, 3, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := Stack{Cards: tt.in}.Cut(tt.n).Cards
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestStack_Deal(t *testing.T) {
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
	} {
		t.Run(tt.name, func(t *testing.T) {
			actual := Stack{tt.in}.Deal(tt.increment).Cards
			assert.Equal(t, tt.expected, actual)
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
			actual := Stack{tt.in}.Flip().Cards
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestStack_Shuffle(t *testing.T) {
	for _, tt := range []struct {
		name     string
		in       []int
		shuffle  []Operation
		expected []int
	}{
		{"given 1",
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			[]Operation{Deal(7), Flip(), Flip()},
			[]int{0, 3, 6, 9, 2, 5, 8, 1, 4, 7},
		},
		{
			"given 2",
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			[]Operation{Cut(6)},
			[]int{6, 7, 8, 9, 0, 1, 2, 3, 4, 5},
		},
		{
			"given 2",
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			[]Operation{Cut(6), Deal(7)},
			[]int{6, 9, 2, 5, 8, 1, 4, 7, 0, 3,},
		},
		{
			"given 2",
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			[]Operation{Cut(6), Deal(7), Flip()},
			[]int{3, 0, 7, 4, 1, 8, 5, 2, 9, 6},
		},
		{
			"given 3",
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			[]Operation{Flip(), Cut(-2), Deal(7), Cut(8), Cut(-4), Deal(7), Cut(3), Deal(9), Deal(3), Cut(-1)},
			[]int{9, 2, 5, 8, 1, 4, 7, 0, 3, 6},
		},
	} {
		actual := Stack{tt.in}.Shuffle(tt.shuffle).Cards
		assert.Equal(t, tt.expected, actual)
	}
}
