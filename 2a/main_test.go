package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdd(t *testing.T) {
	type test struct {
		input       []int
		expected    []int
		expectedPtr int
		err         bool
	}

	for i, test := range []test{
		{[]int{1, 1, 1, 5}, []int{1, 1, 1, 5, 0, 2}, 4, false},
		{[]int{1, 2, 3, 10}, []int{1, 2, 3, 10, 0, 0, 0, 0, 0, 0, 13}, 4, false},
		{[]int{2, 3, 4, 5}, nil, 0, true},
	} {
		actual, actualPtr, actualErr := Add(test.input, 0)
		if test.err {
			assert.Errorf(t, actualErr, "test %d: expected %+v to result in an error", i, test.input)
			assert.Nilf(t, actual, "test %d: expected %+v to result in nil output", i, test.input)
			assert.Equal(t, 0, actualPtr, "test %d: expected %+v to result in zero ptr", i, test.input)
		} else {
			assert.Nilf(t, actualErr, "test %d: expected %+v to _not_ result in an error, but got %v", i, test.input, actualErr)
			assert.Equalf(t, test.expected, actual, "test %d: output of %+v was not as expected", i, test.input)
			assert.Equal(t, test.expectedPtr, actualPtr, "test %d: result ptr of %+v was not as expected", i, test.input)
		}
	}
}

func TestMul(t *testing.T) {
	type test struct {
		input       []int
		expected    []int
		expectedPtr int
		err         bool
	}

	for i, test := range []test{
		{[]int{2, 1, 1, 5}, []int{2, 1, 1, 5, 0, 1}, 4, false},
		{[]int{2, 2, 3, 10}, []int{2, 2, 3, 10, 0, 0, 0, 0, 0, 0, 30}, 4, false},
		{[]int{1, 3, 4, 5}, nil, 0, true},
	} {
		actual, actualPtr, actualErr := Mul(test.input, 0)
		if test.err {
			assert.Errorf(t, actualErr, "test %d: expected %+v to result in an error", i, test.input)
			assert.Nilf(t, actual, "test %d: expected %+v to result in nil output", i, test.input)
			assert.Equal(t, 0, actualPtr, "test %d: expected %+v to result in zero ptr", i, test.input)
		} else {
			assert.Nilf(t, actualErr, "test %d: expected %+v to _not_ result in an error, but got %v", i, test.input, actualErr)
			assert.Equalf(t, test.expected, actual, "test %d: output of %+v was not as expected", i, test.input)
			assert.Equal(t, test.expectedPtr, actualPtr, "test %d: result ptr of %+v was not as expected", i, test.input)
		}
	}
}

func TestExecute(t *testing.T) {
	type test struct {
		input    []int
		expected []int
		err      bool
	}

	for i, test := range []test{
		{[]int{99}, []int{99}, false},
		{[]int{1, 5, 6, 7, 99, 1, 2, 0}, []int{1, 5, 6, 7, 99, 1, 2, 3}, false},
		{[]int{1, 9, 10, 3, 2, 3, 11, 0, 99, 30, 40, 50}, []int{3500, 9, 10, 70, 2, 3, 11, 0, 99, 30, 40, 50}, false},
		{[]int{1, 0, 0, 0, 99}, []int{2, 0, 0, 0, 99}, false},
		{[]int{2, 3, 0, 3, 99}, []int{2, 3, 0, 6, 99}, false},
		{[]int{2, 4, 4, 5, 99, 0}, []int{2, 4, 4, 5, 99, 9801}, false},
		{[]int{1, 1, 1, 4, 99, 5, 6, 0, 99}, []int{30, 1, 1, 4, 2, 5, 6, 0, 99}, false},
	} {
		actual, actualErr := execute(test.input)
		if test.err {
			assert.Errorf(t, actualErr, "test %d: expected %+v to result in an error", i, test.input)
			assert.Nilf(t, actual, "test %d: expected %+v to result in nil output", i, test.input)
		} else {
			assert.Nilf(t, actualErr, "test %d: expected %+v to _not_ result in an error, but got %v", i, test.input, actualErr)
			assert.Equalf(t, test.expected, actual, "test %d: output of %+v was not as expected", i, test.input)
		}
	}
}
