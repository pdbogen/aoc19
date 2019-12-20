package main

import (
	"bytes"
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
		pathCache = &sync.Map{}
		m.Length(nil, m.AllKeys())
	}
}
