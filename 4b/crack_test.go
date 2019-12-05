package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAtLeastOneExactDouble(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  bool
	}{
		{"example 1", "112233", true},
		{"example 2", "123444", false},
		{"example 3", "111122", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.out, AtLeastOneExactDouble(0, tt.in))
		})
	}
}
