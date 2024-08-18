package assertions

import (
	"cmp"
	"reflect"
	"runtime/debug"
	"testing"
)

func errorfNow(tb testing.TB, format string, args ...any) {
	tb.Logf(format, args...)
	tb.FailNow()
}

// NoError asserts that the input error is nil
func NoError(tb testing.TB, input error) {
	const failureFormat = "Unexpected error occurred\n > Error: %v\n"

	if input != nil {
		errorfNow(tb, failureFormat, input)
		return
	}
}

// Error asserts that the input error is non-nil
func Error(tb testing.TB, input error) {
	const failureFormat = "expected error did not occur\n"

	if input == nil {
		errorfNow(tb, failureFormat)
		return
	}
}

// ErrorsMatch asserts that the input and expected error either are both nil
// or both have the same string returned by Error. This is to facilitate table
// driven test with a single expected error field.
func ErrorsMatch(tb testing.TB, expected, input error) {
	const failureFormat = "Errors do not match\n > expected: %v\n < input:    %v\n"
	// If the errors are equal by direct comparison they must match, either both nil or equivalent errors
	if expected != input {
		// Don't call .Error() on a nil error
		if expected == nil || input == nil {
			errorfNow(tb, failureFormat, expected, input)
			return
		}

		if expected.Error() != input.Error() {
			errorfNow(tb, failureFormat, expected, input)
			return
		}
	}
}

// Equal asserts that 2 values of the same type are equal using reflect.DeepEqual
func Equal[T any](tb testing.TB, expected, input T) {
	const failureFormat = "Values are not equal\n > expected: %v\n < input:    %v\n"
	if !reflect.DeepEqual(expected, input) {
		errorfNow(tb, failureFormat, expected, input)
	}
}

func nonMatchingMaps[K comparable, E any, T ~map[K]E](a T, b T) (T, T) {
	aout := make(T)
	bout := make(T)

	for ak, av := range a {
		bv, ok := b[ak]
		if !ok {
			aout[ak] = av
		}
		if !reflect.DeepEqual(av, bv) {
			aout[ak] = av
			bout[ak] = bv
		}
	}

	for bk, bv := range b {
		// We've already compared all elements that exist in both maps
		_, ok := a[bk]
		if !ok {
			bout[bk] = bv
		}
	}

	return aout, bout
}

func nonMatchingSlices[E any, T ~[]E](a T, b T) (T, T) {
	nonMatchedA := make(T, 0)
	nonMatchedB := make(T, 0)

	visited := make([]bool, len(a))

	for _, elementA := range a {
		found := false
		for i, elementB := range b {
			if !visited[i] && reflect.DeepEqual(elementA, elementB) {
				visited[i] = true
				found = true
				break
			}
		}
		if !found {
			nonMatchedA = append(nonMatchedA, elementA)
		}
	}

	clear(visited)

	for _, elementB := range b {
		found := false
		for i, elementA := range a {
			if !visited[i] && reflect.DeepEqual(elementA, elementB) {
				visited[i] = true
				found = true
				break
			}
		}
		if !found {
			nonMatchedB = append(nonMatchedB, elementB)
		}
	}

	return nonMatchedA, nonMatchedB
}

// SlicesMatch asserts that both expected and input have the same members regardless of order
// elements in expected and input are compared using reflect.DeepEqual.
// Failing results will only print the non-matching elements
func SlicesMatch[E any, T ~[]E](tb testing.TB, expected, input T) {
	if len(expected) != len(input) {
		errorfNow(tb, "Elements do not match, slices have different lengths\n > expected length: %v\n, < input length:    %v\n", len(expected), len(input))
		return
	}

	const failureFormat = "Elements do not match\n > expected: %#v\n < input:    %#v\n"
	// Slices will required a different approach
	// since we're not requiring that elements be orderable we can't easily sort the elements
	// we'll take the n^2 approach for simplicity
	expectedNoMatch, inputNoMatch := nonMatchingSlices(expected, input)
	if len(expectedNoMatch) > 0 || len(inputNoMatch) > 0 {
		errorfNow(tb, failureFormat, expectedNoMatch, inputNoMatch)
		return
	}
}

// MapsMatch asserts that both expected and input have the same members regardless of order
// elements in expected and input are compared using reflect.DeepEqual.
// Failing results will only print the non-matching elements
func MapsMatch[K comparable, E any, T ~map[K]E](tb testing.TB, expected, input T) {
	const failureFormat = "Elements do not match\n > expected: %#v\n < input:    %#v\n"

	expectedNoMatch, inputNoMatch := nonMatchingMaps(expected, input)
	if len(expectedNoMatch) > 0 || len(inputNoMatch) > 0 {
		errorfNow(tb, failureFormat, expectedNoMatch, inputNoMatch)
		return
	}
}

// Within asserts that input is within the range [minT, maxT]
// The assertion will pass while input is >= minT and input is <= maxT
func Within[T cmp.Ordered](tb testing.TB, minT, maxT, input T) {
	const failureFormat = "value is not in the expected range\n > expected: [%v, %v]\n < input: %v\n"

	if input < minT || input > maxT {
		errorfNow(tb, failureFormat, minT, maxT, input)
		return
	}
}

func panicHandler(fn func()) (panicked bool, msg any, stack string) {
	panicked = true

	defer func() {
		msg = recover()
		if panicked {
			stack = string(debug.Stack())
		}
	}()

	fn()
	panicked = false
	return
}

// Panics asserts that the provided function panics during execution
func Panics(tb testing.TB, fn func()) {
	const failureFormat = "function %#v did not panic\n > revcovered value: %#v\n"

	panicked, recovered, _ := panicHandler(fn)
	if !panicked {
		errorfNow(tb, failureFormat, fn, recovered)
		return
	}
}

// NotPanics asserts that the provided function does not panic durion execution
func NotPanics(tb testing.TB, fn func()) {
	const failureFormat = "function %#v panic\n > revcovered value: %#v\n > stack: %v\n"

	panicked, recovered, stack := panicHandler(fn)
	if panicked {
		errorfNow(tb, failureFormat, fn, recovered, stack)
		return
	}
}
