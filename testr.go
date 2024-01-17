// Package testr provides a minimal extension to the standard library's testing.
package testr

import (
	"errors"
	"fmt"
	"reflect"
)

// T represents testing.T.
type T interface {
	Helper()
	Logf(format string, args ...any)
	Fail()
	FailNow()
}

// Assertion provides assertion methods around T.
type Assertion struct{ t T }

// New returns a new Assertion given by T.
func New(t T) *Assertion {
	return &Assertion{t}
}

func (assert *Assertion) checkNilT() {
	if assert.t == nil {
		panic("testr: T is nil")
	}
}

// Equal asserts that the actual object are equal to the expected object.
// It can take options.
func (assert *Assertion) Equal(actual any, expected any, options ...Option) {
	assert.checkNilT()
	opt := newOption(assert.t, options...)

	if reflect.DeepEqual(actual, expected) {
		return
	}
	defer opt.fail()

	assert.t.Helper()
	assert.t.Logf("%s%s", ne(actual, expected), opt.message)
}

// ErrorIs asserts that the actual error are equal to the expected error.
// ErrorIs uses errors.Is so it can use any perks that errors.Is provides.
// It can take options.
func (assert *Assertion) ErrorIs(actual error, expected error, options ...Option) {
	assert.checkNilT()
	opt := newOption(assert.t, options...)

	if errors.Is(actual, expected) {
		return
	}
	defer opt.fail()

	assert.t.Helper()
	assert.t.Logf("%s%s", ne(actual, expected), opt.message)
}

// ErrorAs asserts that the actual error as the expected target.
// ErrorAs uses errors.As so it can use any perks that errors.As provides.
// It can take options.
func (assert *Assertion) ErrorAs(actual error, expected any, options ...Option) {
	assert.checkNilT()
	opt := newOption(assert.t, options...)

	if errors.As(actual, expected) {
		return
	}
	defer opt.fail()

	assert.t.Helper()
	assert.t.Logf("%s%s", ne(actual, raw(fmt.Sprintf("as(%T)", expected))), opt.message)
}

// Panic assert that the function is panic.
// It can take options.
func (assert *Assertion) Panic(f func(), options ...Option) {
	assert.checkNilT()
	opt := newOption(assert.t, options...)

	defer func() {
		v := recover()
		if v != nil {
			return
		}
		defer opt.fail()

		assert.t.Helper()
		assert.t.Logf("%s%s", ne(raw("func()"), raw("panic()")), opt.message)
	}()

	assert.t.Helper()
	f()
}

// Must is a helper that wraps a call to a function returning (T, error)
// and panics if the error is non-nil. It is intended for use in variable
// initializations such as
//
//	t := testr.Must(template.New("name").Parse("text"))
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// ne returns a formatted string that tells the two objects are not equal.
func ne(actual, expected any) string {
	return fmt.Sprintf("%s != expected:%s", val(actual), val(expected))
}

// raw represents a raw string that intended to be printed as is.
type raw string

// val returns a string of v with it's type or returns as is if v is raw.
func val(v any) string {
	if err, ok := v.(error); ok {
		return fmt.Sprintf("error(%v)", err)
	}

	if s, ok := v.(raw); ok {
		return string(s)
	}

	return fmt.Sprintf("%#v", v)
}
