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
func (m *mockT) FailNow()                        { m.actual.state = failNow("").state }

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

func TestAssertEqual(t *testing.T) {
	tests := []struct {
		N string                        // Name.
		F func(assert *testr.Assertion) // Function.
		S testState                     // State.
	}{
		{
			N: "eq: nil",
			F: func(assert *testr.Assertion) { assert.Equal(nil, nil) },
			S: pass(),
		},
		{
			N: "eq: not nil",
			F: func(assert *testr.Assertion) { assert.Equal(false, false) },
			S: pass(),
		},
		{
			N: "eq: with message",
			F: func(assert *testr.Assertion) {
				assert.Equal(nil, nil, testr.WithMessage("message"))
			},
			S: pass(),
		},
		{
			N: "eq: with fail now",
			F: func(assert *testr.Assertion) {
				assert.Equal(nil, nil, testr.WithFailNow())
			},
			S: pass(),
		},
		{
			N: "ne: diff val",
			F: func(assert *testr.Assertion) { assert.Equal(false, true) },
			S: fail("bool(false) != expected:bool(true)"),
		},
		{
			N: "ne: diff type",
			F: func(assert *testr.Assertion) { assert.Equal(nil, "nil") },
			S: fail("nil() != expected:string(nil)"),
		},
		{
			N: "ne: with message",
			F: func(assert *testr.Assertion) {
				assert.Equal(false, true, testr.WithMessage("message"))
			},
			S: fail("bool(false) != expected:bool(true) // message"),
		},
		{
			N: "ne: with fail now",
			F: func(assert *testr.Assertion) {
				assert.Equal(nil, "nil", testr.WithFailNow())
			},
			S: failNow("nil() != expected:string(nil)"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.N, func(t *testing.T) {
			m := newMockT(t)
			tc.F(testr.New(m))
			m.assert(tc.S)
		})
	}
}

func TestAssertErrorIs(t *testing.T) {
	tests := []struct {
		N string                        // Name.
		F func(assert *testr.Assertion) // Function.
		S testState                     // State.
	}{
		{
			N: "eq: nil",
			F: func(assert *testr.Assertion) { assert.ErrorIs(nil, nil) },
			S: pass(),
		},
		{
			N: "eq: not nil",
			F: func(assert *testr.Assertion) { assert.ErrorIs(errFoo, errFoo) },
			S: pass(),
		},
		{
			N: "eq: wrap",
			F: func(assert *testr.Assertion) { assert.ErrorIs(errWrapFoo, errFoo) },
			S: pass(),
		},
		{
			N: "eq: with message",
			F: func(assert *testr.Assertion) {
				assert.ErrorIs(nil, nil, testr.WithMessage("message"))
			},
			S: pass(),
		},
		{
			N: "eq: with fail now",
			F: func(assert *testr.Assertion) {
				assert.ErrorIs(nil, nil, testr.WithFailNow())
			},
			S: pass(),
		},
		{
			N: "ne: nil",
			F: func(assert *testr.Assertion) { assert.ErrorIs(errFoo, nil) },
			S: fail("error(foo) != expected:nil()"),
		},
		{
			N: "ne: not nil",
			F: func(assert *testr.Assertion) { assert.ErrorIs(errFoo, errBar) },
			S: fail("error(foo) != expected:error(bar)"),
		},
		{
			N: "ne: wrap",
			F: func(assert *testr.Assertion) { assert.ErrorIs(errWrapFoo, errBar) },
			S: fail("error(wrap foo) != expected:error(bar)"),
		},
		{
			N: "ne: with message",
			F: func(assert *testr.Assertion) {
				assert.ErrorIs(errFoo, nil, testr.WithMessage("message"))
			},
			S: fail("error(foo) != expected:nil() // message"),
		},
		{
			N: "ne: with fail now",
			F: func(assert *testr.Assertion) {
				assert.ErrorIs(errFoo, nil, testr.WithFailNow())
			},
			S: failNow("error(foo) != expected:nil()"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.N, func(t *testing.T) {
			m := newMockT(t)
			tc.F(testr.New(m))
			m.assert(tc.S)
		})
	}
}

func TestAssertErrorAs(t *testing.T) {
	tests := []struct {
		N string                        // Name.
		F func(assert *testr.Assertion) // Function.
		S testState                     // State.
	}{
		{
			N: "eq",
			F: func(assert *testr.Assertion) {
				var e customError
				assert.ErrorAs(customError("err"), &e)
			},
			S: pass(),
		},
		{
			N: "eq: with message",
			F: func(assert *testr.Assertion) {
				var e customError
				assert.ErrorAs(customError("err"), &e, testr.WithMessage("message"))
			},
			S: pass(),
		},
		{
			N: "eq: with fail now",
			F: func(assert *testr.Assertion) {
				var e customError
				assert.ErrorAs(customError("err"), &e, testr.WithFailNow())
			},
			S: pass(),
		},
		{
			N: "ne",
			F: func(assert *testr.Assertion) {
				var e customError
				assert.ErrorAs(errFoo, &e)
			},
			S: fail("error(foo) != expected:as(*testr_test.customError)"),
		},
		{
			N: "ne: with message",
			F: func(assert *testr.Assertion) {
				var e customError
				assert.ErrorAs(errFoo, &e, testr.WithMessage("message"))
			},
			S: fail("error(foo) != expected:as(*testr_test.customError) // message"),
		},
		{
			N: "ne: with fail now",
			F: func(assert *testr.Assertion) {
				var e customError
				assert.ErrorAs(errFoo, &e, testr.WithFailNow())
			},
			S: failNow("error(foo) != expected:as(*testr_test.customError)"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.N, func(t *testing.T) {
			m := newMockT(t)
			tc.F(testr.New(m))
			m.assert(tc.S)
		})
	}
}

func TestAssertPanic(t *testing.T) {
	tests := []struct {
		N string                        // Name.
		F func(assert *testr.Assertion) // Function.
		S testState                     // State.
	}{
		{
			N: "panic",
			F: func(assert *testr.Assertion) {
				assert.Panic(func() { panic("panic") })
			},
			S: pass(),
		},
		{
			N: "panic: with message",
			F: func(assert *testr.Assertion) {
				assert.Panic(func() { panic("panic") }, testr.WithMessage("message"))
			},
			S: pass(),
		},
		{
			N: "panic: with fail now",
			F: func(assert *testr.Assertion) {
				assert.Panic(func() { panic("panic") }, testr.WithFailNow())
			},
			S: pass(),
		},
		{
			N: "not panic",
			F: func(assert *testr.Assertion) { assert.Panic(func() {}) },
			S: fail("func() != expected:panic()"),
		},
		{
			N: "not panic: with message",
			F: func(assert *testr.Assertion) {
				assert.Panic(func() {}, testr.WithMessage("message"))
			},
			S: fail("func() != expected:panic() // message"),
		},
		{
			N: "not panic: with fail now",
			F: func(assert *testr.Assertion) {
				assert.Panic(func() {}, testr.WithFailNow())
			},
			S: failNow("func() != expected:panic()"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.N, func(t *testing.T) {
			m := newMockT(t)
			tc.F(testr.New(m))
			m.assert(tc.S)
		})
	}
}

func TestNilT(t *testing.T) {
	tests := []struct {
		N string                        // Name.
		F func(assert *testr.Assertion) // Function.
	}{
		{
			N: "Equal",
			F: func(assert *testr.Assertion) { assert.Equal(nil, nil) },
		},
		{
			N: "ErrorIs",
			F: func(assert *testr.Assertion) { assert.ErrorIs(nil, nil) },
		},
		{
			N: "ErrorAs",
			F: func(assert *testr.Assertion) {
				var e customError
				assert.ErrorAs(customError("err"), &e)
			},
		},
		{
			N: "Panic",
			F: func(assert *testr.Assertion) {
				assert.Panic(func() { panic("panic") })
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.N, func(t *testing.T) {
			defer func() {
				v := recover()
				if v != nil {
					return
				}
				t.Errorf("func() != panic()")
			}()
			assert := testr.New(nil)
			tc.F(assert)
		})
	}
}
