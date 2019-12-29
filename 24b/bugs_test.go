package main

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGrid_Rating(t *testing.T) {
	tests := []struct {
		name string
		grid *Grid
		want int
	}{
		{"given 1", Load(bytes.NewBufferString(`.....
.....
.....
#....
.#...`)), 2129920},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.grid.Rating(); got != tt.want {
				t.Errorf("Rating() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPoint_Adjacent(t *testing.T) {
	tests := []struct {
		x, y     int
		width    int
		expected int
	}{
		{0, 0, 5, 4},
		{1, 0, 5, 4},
		{2, 0, 5, 4},
		{3, 0, 5, 4},
		{4, 0, 5, 4},
		{0, 1, 5, 4},
		{1, 1, 5, 4},
		{2, 1, 5, 8},
		{3, 1, 5, 4},
		{4, 1, 5, 4},
		{0, 2, 5, 4},
		{1, 2, 5, 8},
		{3, 2, 5, 8},
		{4, 2, 5, 4},
		{0, 3, 5, 4},
		{1, 3, 5, 4},
		{2, 3, 5, 8},
		{3, 3, 5, 4},
		{4, 3, 5, 4},
		{0, 4, 5, 4},
		{1, 4, 5, 4},
		{2, 4, 5, 4},
		{3, 4, 5, 4},
		{4, 4, 5, 4},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d,%d", tt.x, tt.y), func(t *testing.T) {
			p := Point{
				X: tt.x,
				Y: tt.y,
			}
			assert.Equal(t, tt.expected, len(p.Adjacent(tt.width)))
		})
	}
}
