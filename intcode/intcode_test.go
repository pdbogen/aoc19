package intcode_test

import (
	"bytes"
	"fmt"
	"github.com/pdbogen/aoc19/intcode"
	. "github.com/pdbogen/aoc19/intcode/opcodes"
	"github.com/stretchr/testify/assert"
	"strconv"
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
		{[]int{1, 2, 3, 4}, []int{1, 2, 3, 4, 7}, 4, false},
		{[]int{101, 2, 3, 4}, []int{101, 2, 3, 4, 6}, 4, false},
		{[]int{1101, 2, 3, 4}, []int{1101, 2, 3, 4, 5}, 4, false},
		{[]int{1001, 2, 3, 4}, []int{1001, 2, 3, 4, 6}, 4, false},
	} {
		t.Run(fmt.Sprintf("add test #%d", i), func(t *testing.T) {
			actual, actualPtr, actualErr := intcode.Add(test.input, 0)
			if test.err {
				assert.Errorf(t, actualErr, "test %d: expected %+v to result in an error", i, test.input)
				assert.Nilf(t, actual, "test %d: expected %+v to result in nil output", i, test.input)
				assert.Equal(t, 0, actualPtr, "test %d: expected %+v to result in zero ptr", i, test.input)
				return
			}
			assert.Nilf(t, actualErr, "test %d: expected %+v to _not_ result in an error, but got %v", i, test.input, actualErr)
			assert.Equalf(t, test.expected, actual, "test %d: output of %+v was not as expected", i, test.input)
			assert.Equal(t, test.expectedPtr, actualPtr, "test %d: result ptr of %+v was not as expected", i, test.input)
		})
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
		{[]int{2, 2, 3, 3}, []int{2, 2, 3, 9}, 4, false},
		{[]int{102, 2, 3, 3}, []int{102, 2, 3, 6}, 4, false},
		{[]int{1102, 2, 3, 3}, []int{1102, 2, 3, 6}, 4, false},
		{[]int{1002, 2, 3, 3}, []int{1002, 2, 3, 9}, 4, false},
	} {
		actual, actualPtr, actualErr := intcode.Multiply(test.input, 0)
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
		program  []int
		input    string
		expected []int
		output   string
		err      bool
	}

	for i, test := range []test{
		{[]int{98}, "", nil, "", true},
		{[]int{99}, "", []int{99}, "", false},
		{[]int{1, 5, 6, 7, 99, 1, 2, 0}, "", []int{1, 5, 6, 7, 99, 1, 2, 3}, "", false},
		{[]int{1, 9, 10, 3, 2, 3, 11, 0, 99, 30, 40, 50}, "", []int{3500, 9, 10, 70, 2, 3, 11, 0, 99, 30, 40, 50}, "", false},
		{[]int{1, 0, 0, 0, 99}, "", []int{2, 0, 0, 0, 99}, "", false},
		{[]int{2, 3, 0, 3, 99}, "", []int{2, 3, 0, 6, 99}, "", false},
		{[]int{2, 4, 4, 5, 99, 0}, "", []int{2, 4, 4, 5, 99, 9801}, "", false},
		{[]int{1, 1, 1, 4, 99, 5, 6, 0, 99}, "", []int{30, 1, 1, 4, 2, 5, 6, 0, 99}, "", false},
		{[]int{OpInput, 5, OpOutput, 5, OpEnd}, "17", []int{OpInput, 5, OpOutput, 5, 99, 17}, "17\n", false},
		{[]int{1001, 1, 5, 3}, "", []int{1001, 1, 5, 6}, "", false},
		{[]int{1002, 2, 5, 3}, "", []int{1002, 2, 5, 25}, "", false},
		{[]int{1101, 100, -1, 4, 0}, "", []int{1101, 100, -1, 4, 99}, "", false},
		{[]int{
			1100 + OpJumpTrue, 1, 6, // 0 1 2
			100 + OpOutput, 0, // 3 4
			OpEnd,             // 5
			100 + OpOutput, 1, // 6 7
		}, "", nil, "1\n", false},
		{[]int{
			1100 + OpJumpFalse, 1, 6, // 0 1 2
			100 + OpOutput, 0, // 3 4
			OpEnd,             // 5
			100 + OpOutput, 1, // 6 7
		}, "", nil, "0\n", false},
		{[]int{1100 + OpLessThan, 1, 2, 5, 99}, "", []int{1100 + OpLessThan, 1, 2, 5, 99, 1}, "", false},
		{[]int{1100 + OpLessThan, 10, 2, 5, 99}, "", []int{1100 + OpLessThan, 10, 2, 5, 99, 0}, "", false},
		{[]int{1100 + OpLessThan, 2, 2, 5, 99}, "", []int{1100 + OpLessThan, 2, 2, 5, 99, 0}, "", false},
		{[]int{1100 + OpEquals, 2, 2, 5, 99}, "", []int{1100 + OpEquals, 2, 2, 5, 99, 1}, "", false},
		{[]int{1100 + OpEquals, 1, 2, 5, 99}, "", []int{1100 + OpEquals, 1, 2, 5, 99, 0}, "", false},
	} {
		if test.expected == nil {
			test.expected = test.program
		}

		t.Run(fmt.Sprintf("execute test #%d", i), func(t *testing.T) {
			actualOut := &bytes.Buffer{}
			actual, actualErr := intcode.Execute(test.program, bytes.NewBufferString(test.input), actualOut)
			if test.err {
				assert.Errorf(t, actualErr, "test %d: expected %+v to result in an error", i, test.program)
				return
			}
			assert.Nilf(t, actualErr, "test %d: expected %+v to _not_ result in an error, but got %v", i, test.program, actualErr)
			assert.Equalf(t, test.expected, actual, "test %d: output of %+v was not as expected", i, test.program)
			assert.Equal(t, test.output, actualOut.String())
		})
	}
}

func TestInput(t *testing.T) {
	type test struct {
		program  []int
		input    string
		expected []int
		err      bool
	}
	for i, test := range []test{
		{[]int{3, 0}, "99", []int{99, 0}, false},
		{[]int{3, 5, 99}, "99", []int{3, 5, 99, 0, 0, 99}, false},
		{[]int{3, 0}, "not a number", nil, true},
	} {
		t.Run(fmt.Sprintf("Input Test #%d", i), func(t *testing.T) {
			actual, actualErr := intcode.Execute(test.program, bytes.NewBufferString(test.input), nil)
			if test.err {
				assert.Error(t, actualErr)
				return
			}
			assert.Nil(t, actualErr)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestOutput(t *testing.T) {
	type test struct {
		program []int
		output  string
		err     bool
	}
	for i, test := range []test{
		{[]int{OpOutput, 0}, strconv.Itoa(OpOutput) + "\n", false},
		{[]int{OpOutput, 99}, "", true},
		{[]int{104, 17}, "17\n", false},
	} {
		t.Run(fmt.Sprintf("output test #%d", i), func(t *testing.T) {
			output := &bytes.Buffer{}
			actual, actualErr := intcode.Execute(test.program, nil, output)
			if test.err {
				assert.Error(t, actualErr)
				return
			}
			assert.Nil(t, actualErr)
			assert.Equal(t, test.output, output.String())
			assert.Equal(t, test.program, actual)
		})
	}
}

func TestReadOp(t *testing.T) {
	type test struct {
		program []int
		code    int
		args    []int
		modes   []int
		baked   []int
		err     bool
	}

	for i, test := range []test{
		{[]int{1, 0, 0, 0}, 1, []int{0, 0, 0}, []int{0, 0, 0}, []int{1, 1, 1}, false},
		{[]int{1, 1, 2, 3}, 1, []int{1, 2, 3}, []int{0, 0, 0}, []int{1, 2, 3}, false},
		{[]int{101, 0, 0, 0}, 1, []int{0, 0, 0}, []int{1, 0, 0}, []int{0, 101, 101}, false},
		{[]int{1101, 0, 0, 0}, 1, []int{0, 0, 0}, []int{1, 1, 0}, []int{0, 0, 1101}, false},
		{[]int{10101, 0, 0, 0}, 1, []int{0, 0, 0}, []int{1, 0, 1}, []int{0, 10101, 0}, false},
		{[]int{12301, 0, 0, 0}, 0, []int{0, 0, 0}, nil, nil, true},
	} {
		t.Run(fmt.Sprintf("read op test #%d", i), func(t *testing.T) {
			op, err := intcode.ReadOp(test.program, 0, len(test.args))
			if test.err {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, err)
			assert.NotNil(t, op)
			assert.Equal(t, test.code, op.Code)
			assert.Equal(t, test.args, op.Args)
			assert.Equal(t, test.modes, op.Modes)
			assert.Equal(t, len(op.Args), len(op.Modes), "# of args & modes didn't match")
		})
	}
}
