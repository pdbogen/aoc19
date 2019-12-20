package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestBox(t *testing.T) {
	type test struct {
		lines    []string
		height   int
		expected int
	}

	for _, tt := range []test{
		{
			[]string{
				"#",
				"##",
				"###",
				"####",
				"#####",
			},
			5,
			1,
		},
		{
			[]string{
				".....#########",
				"......#########",
				".......#########",
				"........#########",
				".........#########",
			},
			5,
			5,
		},
		{
			strings.Split(`.....................######
......................######
......................#######
.......................######
........................######
........................#######
.........................#######
.........................#######`, "\n"), 8, 2,
		},
	} {
		var lines []Line
		for _, l := range tt.lines {
			lines = append(lines, Line{
				offset: strings.IndexByte(l, '#'),
				width:  len(l) - strings.IndexByte(l, '#'),
			})
		}
		assert.Equal(t, tt.expected, box(lines, tt.height))
	}
}
