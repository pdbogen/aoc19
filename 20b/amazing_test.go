package main

import (
	"os"
	"sync"
	"testing"
)

func BenchmarkMap_Path(b *testing.B) {
	f, err := os.Open("input")
	if err != nil {
		b.Fatal(err)
	}
	m, err := LoadMap(f)
	if err != nil {
		b.Fatal(err)
	}
	l := m.Layer(10)
	for i := 0; i < b.N; i++ {
		l.pathCache = &sync.Map{}
		l.Path(l.Start.Point, l.Goal.Point)
	}
}
