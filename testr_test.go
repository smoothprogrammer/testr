package testr_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/minizilla/testr"
)

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

func (m *mockT) assert(expected testState) {
	m.t.Helper()
	if m.actual.state != expected.state {
		m.t.Errorf("%#v != %#v // state", m.actual.state, expected.state)
	}
	if m.actual.output != expected.output {
		m.t.Errorf("%#v != %#v // output", m.actual.output, expected.output)
	}
}

type testState struct {
	state, output string
}

func pass() testState {
	return testState{"pass", ""}
}

func fail(output string) testState {
	return testState{"fail", output}
}

var (
	errFoo     = errors.New("foo")
	errBar     = errors.New("bar")
	errWrapFoo = fmt.Errorf("wrap %w", errFoo)
)

func TestAssertEqual(t *testing.T) {
	tests := []struct {
		name     string
		actual   any
		expected any
		testState
	}{
		{"eq: nil", nil, nil, pass()},
		{"eq: not nil", false, false, pass()},
		{"ne: diff val", false, true, fail("bool(false) != expected:bool(true)")},
		{"ne: diff type", nil, "nil", fail("nil() != expected:string(nil)")},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := newMockT(t)
			assert := testr.New(m)
			assert.Equal(tc.actual, tc.expected)
			m.assert(tc.testState)
		})
	}
}

func TestAssertErrorIs(t *testing.T) {
	tests := []struct {
		name     string
		actual   error
		expected error
		testState
	}{
		{"eq: nil", nil, nil, pass()},
		{"eq: not nil", errFoo, errFoo, pass()},
		{"eq: wrap", errWrapFoo, errFoo, pass()},
		{"ne: nil", errFoo, nil, fail("error(foo) != expected:nil()")},
		{"ne: not nil", errFoo, errBar, fail("error(foo) != expected:error(bar)")},
		{"ne: wrap", errWrapFoo, errBar, fail("error(wrap foo) != expected:error(bar)")},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := newMockT(t)
			assert := testr.New(m)
			assert.ErrorIs(tc.actual, tc.expected)
			m.assert(tc.testState)
		})
	}
}

func TestNilT(t *testing.T) {
	tests := []struct {
		name string
		f    func(assert *testr.Assertion)
	}{
		{"Equal", func(assert *testr.Assertion) { assert.Equal(nil, nil) }},
		{"ErrorIs", func(assert *testr.Assertion) { assert.ErrorIs(nil, nil) }},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				v := recover()
				if v != nil {
					return
				}
				t.Helper()
				t.Errorf("the function is not panic")
			}()
			assert := testr.New(nil)
			tc.f(assert)
		})
	}
}
