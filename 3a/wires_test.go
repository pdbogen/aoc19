package main

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

func TestParseLine(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		args    args
		wantRet []Move
		wantErr bool
	}{
		{name: "basic 1", args: args{line: "U8"}, wantRet: []Move{{Up, 8}}, wantErr: false},
		{name: "basic 2", args: args{line: "U8,R5"}, wantRet: []Move{
			{Up, 8},
			{Right, 5},
		}, wantErr: false},
		{name: "basic 4", args: args{line: "U8,R5,L5,D3"}, wantRet: []Move{
			{Up, 8},
			{Right, 5},
			{Left, 5},
			{Down, 3},
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRet, err := ParseLine(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("ParseLine() gotRet = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestFindCrossings(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		name    string
		args    args
		want    []Point
		wantErr bool
	}{
		{"crossing 1",
			args{[]string{"R8,U5,L5,D3", "U7,R6,D4,L4"}},
			[]Point{{3, 3}, {6, 5}},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grid, err := ToGrid(tt.args.lines)
			assert.Nil(t, err, "conversion to grid should have worked")

			got, err := FindCrossings(grid)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindCrossings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Subset(t, got, tt.want, "result did not contain some desired intersections")
			assert.Subset(t, tt.want, got, "result contained extra intersections")
		})
	}
}

func TestDistance(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"basic example", args{[]string{"R8,U5,L5,D3", "U7,R6,D4,L4"}}, 6, false},
		{"simple 1", args{[]string{"R10,U10", "U10,R10"}}, 20, false},
		{"simple 2", args{[]string{"U10,R10", "R10,U10"}}, 20, false},
		{"more example 1",
			args{[]string{
				"R75,D30,R83,U83,L12,D49,R71,U7,L72",
				"U62,R66,U55,R34,D71,R55,D58,R83",
			}},
			159,
			false,
		},
		{"self crossing",
			args{[]string{
				"U5,R5,D2,L2,U6,R5",
				"L1,U7,R30",
			}},
			10,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grid, err := ToGrid(tt.args.lines)
			assert.Nil(t, err, "ToGrid should not produce an error")

			got, err := ClosestCrossingDistance(grid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClosestCrossingDistance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ClosestCrossingDistance() got = %v, want %v", got, tt.want)
				for _, line := range strings.Split(grid.String(), "\n") {
					t.Log(line)
				}
			}
		})
	}
}

func Test_dist2d(t *testing.T) {
	type args struct {
		a Point
		b Point
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"1", args{Point{0, 0}, Point{1, 1}}, 2},
		{"2", args{Point{2, 2}, Point{1, 1}}, 2},
		{"3", args{Point{2, 0}, Point{1, 1}}, 2},
		{"4", args{Point{0, 2}, Point{1, 1}}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dist2d(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("dist2d() = %v, want %v", got, tt.want)
			}
		})
	}
}
