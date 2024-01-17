package testr_test

import (
	"errors"
	"testing"

	"github.com/minizilla/testr"
)

func TestAssertEqual(t *testing.T) {
	tests := []struct {
		N string                        // Name.
		F func(assert *testr.Assertion) // Function.
		S testState                     // State.
	}{
		{
			N: "eq",
			F: func(assert *testr.Assertion) { assert.Equal(nil, nil) },
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
			N: "ne: bool",
			F: func(assert *testr.Assertion) { assert.Equal(false, nil) },
			S: fail("false != expected:<nil>"),
		},
		{
			N: "ne: int",
			F: func(assert *testr.Assertion) { assert.Equal(0, nil) },
			S: fail("0 != expected:<nil>"),
		},
		{
			N: "ne: float",
			F: func(assert *testr.Assertion) { assert.Equal(0.0, nil) },
			S: fail("0 != expected:<nil>"),
		},
		{
			N: "ne: complex",
			F: func(assert *testr.Assertion) { assert.Equal(0i, nil) },
			S: fail("(0+0i) != expected:<nil>"),
		},
		{
			N: "ne: array",
			F: func(assert *testr.Assertion) {
				assert.Equal([1]string{""}, nil)
			},
			S: fail("[1]string{\"\"} != expected:<nil>"),
		},
		{
			N: "ne: map",
			F: func(assert *testr.Assertion) {
				assert.Equal(map[bool]string{false: ""}, nil)
			},
			S: fail("map[bool]string{false:\"\"} != expected:<nil>"),
		},
		{
			N: "ne: slice",
			F: func(assert *testr.Assertion) {
				assert.Equal([]string{""}, nil)
			},
			S: fail("[]string{\"\"} != expected:<nil>"),
		},
		{
			N: "ne: string",
			F: func(assert *testr.Assertion) { assert.Equal("", nil) },
			S: fail("\"\" != expected:<nil>"),
		},
		{
			N: "ne: struct",
			F: func(assert *testr.Assertion) {
				type s struct{}
				assert.Equal(s{}, nil)
			},
			S: fail("testr_test.s{} != expected:<nil>"),
		},
		{
			N: "ne: with message",
			F: func(assert *testr.Assertion) {
				assert.Equal(false, nil, testr.WithMessage("message"))
			},
			S: fail("false != expected:<nil> // message"),
		},
		{
			N: "ne: with fail now",
			F: func(assert *testr.Assertion) {
				assert.Equal(false, nil, testr.WithFailNow())
			},
			S: failNow("false != expected:<nil>"),
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
			S: fail("error(foo) != expected:<nil>"),
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
			S: fail("error(foo) != expected:<nil> // message"),
		},
		{
			N: "ne: with fail now",
			F: func(assert *testr.Assertion) {
				assert.ErrorIs(errFoo, nil, testr.WithFailNow())
			},
			S: failNow("error(foo) != expected:<nil>"),
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

func TestMust(t *testing.T) {
	assert := testr.New(t)
	assert.Equal(testr.Must("ok", nil), "ok")
	assert.Panic(func() { testr.Must("panic", errors.New("intentional error")) })
}
