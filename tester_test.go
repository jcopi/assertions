package assertions

import "testing"

type TesterTB struct {
	testing.TB
	mustfail bool
	failed   bool
}

// Error implements testing.TB.
func (t *TesterTB) Error(args ...any) {
	t.failed = true
}

// Errorf implements testing.TB.
func (t *TesterTB) Errorf(format string, args ...any) {
	t.failed = true
}

// Fail implements testing.TB.
func (t *TesterTB) Fail() {
	t.failed = true
}

// FailNow implements testing.TB.
func (t *TesterTB) FailNow() {
	t.failed = true
}

// Failed implements testing.TB.
func (t *TesterTB) Failed() bool {
	return t.failed
}

// Fatal implements testing.TB.
func (t *TesterTB) Fatal(args ...any) {
	t.failed = true
}

// Fatalf implements testing.TB.
func (t *TesterTB) Fatalf(format string, args ...any) {
	t.failed = true
}

// Log implements testing.TB.
func (t *TesterTB) Log(args ...any) {
}

// Logf implements testing.TB.
func (t *TesterTB) Logf(format string, args ...any) {
}

func (t *TesterTB) AssertExpectation() {
	if t.failed != t.mustfail {
		t.TB.Logf("Failure was not as expected:\n > expected: %v\n < actual: %v\n", t.mustfail, t.failed)
		t.TB.FailNow()
	}
}

var _ testing.TB = &TesterTB{}

func NewTester(tb testing.TB, mustfail bool) *TesterTB {
	return &TesterTB{
		TB:       tb,
		mustfail: mustfail,
	}
}
