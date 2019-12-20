package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestCanTraverse(t *testing.T) {
	tests := []struct {
		path        string
		constraints map[KeySet]KeySet
		want        bool
	}{
		{"ab", nil, true},
		{"ab", map[KeySet]KeySet{Keys['a']: Keys['b']}, false},
		{"ab", map[KeySet]KeySet{Keys['b']: Keys['c']}, false},
	}
	for _, tt := range tests {
		if got := CanTraverse([]byte(tt.path), tt.constraints); got != tt.want {
			t.Errorf("CanTraverse() = %v, want %v", got, tt.want)
		}
	}
}

func BenchmarkMap_Length(b *testing.B) {
	const givenD = `#################
#i.G..c...e..H.p#
########.########
#j.A..b...f..D.o#
########@########
#k.E..a...g..B.n#
########.########
#l.F..d...h..C.m#
#################`

	m, _ := LoadMap(bytes.NewBufferString(givenD))

	for i := 0; i < b.N; i++ {
		pathCache = sync.Map{}
		m.Length(nil, m.AllKeys())
	}
}

func TestMap_Path(t *testing.T) {
	const givenD = `#################
#i.G..c...e..H.p#
########.########
#j.A..b...f..D.o#
########@########
#k.E..a...g..B.n#
########.########
#l.F..d...h..C.m#
#################`

	m, err := LoadMap(bytes.NewBufferString(givenD))
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, m, "map")
	assert.NotNil(t, m.Start, "map.Start")
	assert.NotNil(t, m.Keys, "map.Keys")
	assert.NotNil(t, m.Keys[Keys['i']], "map.Keys['i']")
	assert.Equal(t,
		m.Path(m.Start.Point, m.Keys[Keys['i']].Point),
		[]Point{
			{8,4},
			{8,3},
			{8,2},
			{8,1},
			{7,1},
			{6,1},
			{5,1},
			{4,1},
			{3,1},
			{2,1},
		})
}
