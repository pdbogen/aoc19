package main

import (
	"bytes"
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
