package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var givenExample = strings.Join([]string{"COM)B", "B)C", "C)D", "D)E", "E)F", "B)G", "G)H", "D)I", "E)J", "J)K", "K)L"}, "\n")
var givenResult = &Orbit{map[Body]*Orbit{
	"B": {map[Body]*Orbit{
		"C": {map[Body]*Orbit{
			"D": {map[Body]*Orbit{
				"E": {map[Body]*Orbit{
					"F": {map[Body]*Orbit{}},
					"J": {map[Body]*Orbit{
						"K": {map[Body]*Orbit{
							"L": {map[Body]*Orbit{}},
						}},
					}},
				}},
				"I": {map[Body]*Orbit{}},
			}},
		}},
		"G": {map[Body]*Orbit{
			"H": {map[Body]*Orbit{}},
		}},
	}},
}}

func TestParseMap(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *Orbit
		err      bool
	}{
		{"simple", "COM)B", &Orbit{Satellites: map[Body]*Orbit{"B": {map[Body]*Orbit{}}}}, false},
		{"three body", "COM)B\nB)C", &Orbit{map[Body]*Orbit{"B": &Orbit{map[Body]*Orbit{"C": {map[Body]*Orbit{}}}}}}, false},
		{
			"given example",
			givenExample,
			givenResult,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ParseMap(bytes.NewBufferString(tt.input))
			if tt.err {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestChecksum(t *testing.T) {
	tests := []struct {
		name     string
		orbit    *Orbit
		expected int
	}{
		{"simple", &Orbit{Satellites: map[Body]*Orbit{"B": {}}}, 1},
		{"example 1",
			&Orbit{
				map[Body]*Orbit{
					"B": {
						map[Body]*Orbit{
							"C": {
								map[Body]*Orbit{
									"D": {},
								},
							},
						},
					},
				},
			}, 6},
		{"given example", givenResult, 42},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := Checksum(tt.orbit)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
