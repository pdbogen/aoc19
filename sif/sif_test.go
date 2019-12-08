package sif

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToLayers(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		width  int
		height int
		want   [][][]int
	}{
		{"given", "123456789012", 3, 2, [][][]int{
			{
				{1, 7},
				{4, 0},
			},
			{
				{2, 8},
				{5, 1},
			},
			{
				{3, 9},
				{6, 2},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToLayers(tt.data, tt.width, tt.height)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFlattenLayers(t *testing.T) {
	tests := []struct {
		name   string
		layers [][][]int
		expect [][]int
	}{
		{"given", [][][]int{
			{
				{0, 1, 2, 0},
				{2, 2, 1, 0},
			},
			{
				{2, 1, 2, 0},
				{2, 2, 2, 0},
			},
		}, [][]int{
			{0, 1},
			{1, 0},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FlattenLayers(tt.layers)
			assert.Equal(t, tt.expect, got)
		})
	}
}
