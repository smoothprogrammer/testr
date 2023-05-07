package testr_test

import (
	"errors"
	"fmt"
	"testing"
)

type exampleT struct{}

func (m *exampleT) Helper()                         {}
func (m *exampleT) Logf(format string, args ...any) { fmt.Printf(format+"\n", args...) }
func (m *exampleT) Fail()                           {}
func (m *exampleT) FailNow()                        {}

var t = new(exampleT)

type mockT struct {
	t      *testing.T
	actual testState
}

func newMockT(t *testing.T) *mockT {
	return &mockT{
		t:      t,
		actual: pass(),
	}
}

func (m *mockT) Helper()                         { /* TODO: how to test helper? */ }
func (m *mockT) Logf(format string, args ...any) { m.actual.output = fmt.Sprintf(format, args...) }
func (m *mockT) Fail()                           { m.actual.state = fail("").state }
func (m *mockT) FailNow()                        { m.actual.state = failNow("").state }

// assert asserts the actual test state with the expected test state.
func (m *mockT) assert(expected testState) {
	m.t.Helper()
	if m.actual.state != expected.state {
		m.t.Errorf("%#v != %#v // state", m.actual.state, expected.state)
	}
	if m.actual.output != expected.output {
		m.t.Errorf("%#v != %#v // output", m.actual.output, expected.output)
	}
}

// testState represents the state of the test and it's output.
type testState struct {
	state, output string
}

// pass sets the test state as pass with empty output.
func pass() testState {
	return testState{"pass", ""}
}

// fail sets the test state as fail with an output.
func fail(output string) testState {
	return testState{"fail", output}
}

// failNow sets the test state as fail now with an output.
func failNow(output string) testState {
	return testState{"fail now", output}
}

var (
	errFoo     = errors.New("foo")
	errBar     = errors.New("bar")
	errWrapFoo = fmt.Errorf("wrap %w", errFoo)
)

type customError string

func (e customError) Error() string {
	return string(e)
}
