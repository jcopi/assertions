package assertions

import (
	"errors"
	"fmt"
	"math"
	"testing"
)

func TestErrorsMatch(t *testing.T) {
	cases := []struct {
		name     string
		expected error
		input    error
		mustFail bool
	}{
		{
			name:     "happy path",
			expected: nil,
			input:    nil,
			mustFail: false,
		},
		{
			name:     "match",
			expected: fmt.Errorf("matching"),
			input:    fmt.Errorf("matching"),
			mustFail: false,
		},
		{
			name:     "nil and non-nil",
			expected: fmt.Errorf("error"),
			input:    nil,
			mustFail: true,
		},
		{
			name:     "non-nil and nil",
			expected: nil,
			input:    fmt.Errorf("error"),
			mustFail: true,
		},
		{
			name:     "non-matching",
			expected: fmt.Errorf("different error"),
			input:    fmt.Errorf("error"),
			mustFail: true,
		},
		{
			name:     "matching errors.new",
			expected: errors.New("error"),
			input:    errors.New("error"),
			mustFail: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tb := NewTester(t, tc.mustFail)

			ErrorsMatch(tb, tc.expected, tc.input)
			tb.AssertExpectation()
		})
	}
}

func TestNoError(t *testing.T) {
	cases := []struct {
		name     string
		input    error
		mustFail bool
	}{
		{
			name:     "happy path",
			input:    nil,
			mustFail: false,
		},
		{
			name:     "fmt error",
			input:    fmt.Errorf("error"),
			mustFail: true,
		},
		{
			name:     "errors.new",
			input:    errors.New("error"),
			mustFail: true,
		},
		{
			name:     "empty fmt",
			input:    fmt.Errorf(""),
			mustFail: true,
		},
		{
			name:     "empty errors",
			input:    errors.New(""),
			mustFail: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tb := NewTester(t, tc.mustFail)

			NoError(tb, tc.input)
			tb.AssertExpectation()
		})
	}
}

func TestEqualAny(t *testing.T) {
	cases := []struct {
		name     string
		input    any
		expected any
		mustFail bool
	}{
		{name: "int equal", input: 42, expected: 42, mustFail: false},
		{name: "int not equal", input: 34, expected: 42, mustFail: true},
		{name: "float equal", input: 1.056, expected: 1.056, mustFail: false},
		{name: "float not equal", input: 3.1415926535, expected: 3.1, mustFail: true},
		{name: "float nan", input: math.NaN(), expected: math.NaN(), mustFail: true},
		{name: "string equal", input: "abc", expected: "abc", mustFail: false},
		{name: "string not equal", input: "def", expected: "abc", mustFail: true},
		{name: "bool equal", input: true, expected: true, mustFail: false},
		{name: "bool not equal", input: true, expected: false, mustFail: true},
		{name: "slice equal", input: []int{1, 2, 3, 4, 5}, expected: []int{1, 2, 3, 4, 5}, mustFail: false},
		{name: "slice not equal", input: []float64{1, 2, 3, 4, 5.5}, expected: []float64{1, 2, 3, 4, 5}, mustFail: true},
		{name: "map equal", input: map[string]int{"a": 2, "b": 3, "c": 255}, expected: map[string]int{"a": 2, "b": 3, "c": 255}, mustFail: false},
		{name: "map not equal", input: map[string]int{"a": 2, "b": 3, "c": 255}, expected: map[string]int{"a": 2, "b": 3, "c": 255, "d": 4}, mustFail: true},
		{name: "struct equal", input: struct {
			a int
			b float64
		}{a: 2, b: 2.1}, expected: struct {
			a int
			b float64
		}{a: 2, b: 2.1}, mustFail: false},
		{name: "struct not equal", input: struct {
			a int
			b float64
			c string
		}{a: 2, b: 2.1, c: " "}, expected: struct {
			a int
			b float64
		}{a: 2, b: 2.1}, mustFail: true},
		{name: "mismatched types", input: true, expected: 1, mustFail: true},
	}

	for _, tc := range cases {
		tb := NewTester(t, tc.mustFail)

		Equal(tb, tc.expected, tc.input)
		tb.AssertExpectation()
	}
}

func TestSlicesMatchAny(t *testing.T) {
	cases := []struct {
		name     string
		input    []any
		expected []any
		mustFail bool
	}{
		{
			name:     "int match",
			input:    []any{1, 2, 3, 4, 5},
			expected: []any{5, 1, 3, 4, 2},
			mustFail: false,
		},
		{
			name:     "int no match",
			input:    []any{1, 2, 3, 4, 5},
			expected: []any{5, 4, 1, 4, 2},
			mustFail: true,
		},
		{
			name:     "mismatched len",
			input:    []any{1, 2, 3, 4, 5},
			expected: []any{5, 1, 4, 2},
			mustFail: true,
		},
		{
			name:     "nil and len 0",
			input:    []any{},
			expected: nil,
			mustFail: false,
		},
		{
			name:     "string match",
			input:    []any{"abc", "def", "ghi", "jkl", "abc"},
			expected: []any{"abc", "def", "ghi", "jkl"},
			mustFail: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tb := NewTester(t, tc.mustFail)

			SlicesMatch(tb, tc.expected, tc.input)
			tb.AssertExpectation()
		})
	}
}

func TestMapsMatchAny(t *testing.T) {
	cases := []struct {
		name     string
		input    map[string]any
		expected map[string]any
		mustFail bool
	}{
		{
			name:     "nil and len 0",
			input:    nil,
			expected: map[string]any{},
			mustFail: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tb := NewTester(t, tc.mustFail)

			MapsMatch(tb, tc.expected, tc.input)
			tb.AssertExpectation()
		})
	}
}
