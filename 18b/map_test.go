package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMap_Path(t *testing.T) {

	m, err := LoadMap(bytes.NewBufferString(givenD))
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, m, "map")
	assert.Len(t, m.Starts, 2)
	assert.NotNil(t, m.Keys, "map.Keys")
	assert.NotNil(t, m.Keys[Keys['i']], "map.Keys['i']")
	assert.Nil(t, m.Path(m.Starts[0].Point, m.Keys[Keys['g']].Point))
	assert.Equal(t,
		m.Path(m.Starts[0].Point, m.Keys[Keys['i']].Point),
		[]Point{
			{8, 3},
			{8, 2},
			{8, 1},
			{7, 1},
			{6, 1},
			{5, 1},
			{4, 1},
			{3, 1},
			{2, 1},
		})
}

func TestMap_Length(t *testing.T) {
	type test struct {
		Name    string
		Map     string
		Lengths []int32
	}
	for _, tt := range []test{
		{"given A", givenA, []int32{2, 2, 2, 2}},
		{"given B", givenB, []int32{6, 6, 6, 6}},
		{"given C", givenC, []int32{7, 9, 9, 7}},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			m, err := LoadMap(bytes.NewBufferString(tt.Map))
			assert.NoError(t, err)
			assert.NotNil(t, m)
			if !assert.Equal(t, len(m.Starts), len(tt.Lengths)) {
				return
			}
			for s := range m.Starts {
				l := m.Length(nil, m.KeysInSection(s))
				assert.Equal(t, tt.Lengths[s], l)
			}
		})
	}
}
