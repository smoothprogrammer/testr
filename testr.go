package testr

import (
	"fmt"
	"reflect"
)

type T interface {
	Helper()
	Logf(format string, args ...any)
	Fail()
}

type Assertion struct {
	t T
}

func New(t T) *Assertion {
	return &Assertion{t}
}

func (assert *Assertion) checkNilT() {
	if assert.t == nil {
		panic("testr: T is nil")
	}
}

func (assert *Assertion) Equal(actual any, expected any) {
	assert.checkNilT()

	if reflect.DeepEqual(actual, expected) {
		return
	}
	defer assert.t.Fail()

	assert.t.Helper()
	assert.t.Logf("%s", ne(actual, expected))
}

func ne(actual, expected any) string {
	return fmt.Sprintf("%s != expected:%s", val(actual), val(expected))
}

func val(v any) string {
	if v == nil {
		return "nil()"
	}
	return fmt.Sprintf("%[1]T(%[1]v)", v)
}
