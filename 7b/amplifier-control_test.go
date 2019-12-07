package main

import (
	"github.com/pdbogen/aoc19/intcode"
	"github.com/stretchr/testify/assert"
	"testing"
)

var exampleA = intcode.MustLoadString("3,26,1001,26,-4,26,3,27,1002,27,2,27,1,27,26,27,4,27,1001,28,-1,28,1005,28,6,99,0,0,5")
var exampleB = intcode.MustLoadString("3,52,1001,52,-5,52,3,53,1,52,56,54,1007,54,5,55,1005,55,26,1001,54,-5,54,1105,1,12,1,53,54,53,1008,54,0,55,1001,55,1,55,2,53,55,53,4,53,1001,56,-1,56,1005,56,6,99,0,0,0,0,10")

func TestTry(t *testing.T) {
	tests := []struct {
		name      string
		program   []int
		expectSeq []int
		expectVal int
		err       bool
	}{
		{"example a", exampleA, []int{9, 8, 7, 6, 5}, 139629729, false},
		{"example b", exampleB, []int{9, 7, 8, 5, 6}, 18216, false},
		{"err example", []int{55}, []int{}, -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSeq, gotVal, err := try(tt.program, []int{}, []int{5, 6, 7, 8, 9})
			if tt.err {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tt.expectSeq, gotSeq)
			assert.Equal(t, tt.expectVal, gotVal)
		})
	}
}

func Test_runTry(t *testing.T) {
	tests := []struct {
		name    string
		program []int
		seq     []int
		expect  int
		err     bool
	}{
		{"example a", exampleA, []int{9, 8, 7, 6, 5}, 139629729, false},
		{"example b", exampleB, []int{9, 7, 8, 5, 6}, 18216, false},
		{"err example1", exampleA, []int{}, -1, true},
		{"err example2", exampleA, []int{1}, -1, true},
		{"err example3", []int{55}, []int{1, 2}, -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := runTry(tt.program, tt.seq)
			if tt.err {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tt.expect, got)
		})
	}
}
