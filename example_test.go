package testr_test

import (
	"fmt"
	"io"
	"strings"

	"github.com/minizilla/testr"
)

func ExampleAssertion_Equal() {
	assert := testr.New(t) // using *testing.T

	assert.Equal(nil, nil)                    // PASS
	assert.Equal(false, true)                 // FAIL
	assert.Equal("nil", nil)                  // FAIL
	assert.Equal([]string{"hello\nworld"}, 0) // FAIL

	// Output:
	// false != expected:true
	// "nil" != expected:<nil>
	// []string{"hello\nworld"} != expected:0
}

func ExampleAssertion_ErrorIs() {
	assert := testr.New(t) // using *testing.T

	assert.ErrorIs(nil, nil)           // PASS
	assert.ErrorIs(errFoo, errFoo)     // PASS
	assert.ErrorIs(errWrapFoo, errFoo) // PASS
	assert.ErrorIs(errFoo, nil)        // FAIL
	assert.ErrorIs(errFoo, errBar)     // FAIL
	assert.ErrorIs(errWrapFoo, errBar) // FAIL

	// Output:
	// error(foo) != expected:<nil>
	// error(foo) != expected:error(bar)
	// error(wrap foo) != expected:error(bar)
}

func ExampleAssertion_ErrorAs() {
	assert := testr.New(t) // using *testing.T

	var e customError
	assert.ErrorAs(customError("err"), &e) // PASS
	assert.ErrorAs(errFoo, &e)             // FAIL

	// Output:
	// error(foo) != expected:as(*testr_test.customError)
}

func ExampleAssertion_Panic() {
	assert := testr.New(t) // using *testing.T

	assert.Panic(func() { panic("panic") }) // PASS
	assert.Panic(func() {})                 // FAIL

	// Output:
	// func() != expected:panic()
}

func ExampleWithMessage() {
	assert := testr.New(t) // using *testing.T

	assert.Equal(false, true, testr.WithMessage("assert equality"))
	assert.ErrorIs(errFoo, nil, testr.WithMessage("assert err is nil"))
	var e customError
	assert.ErrorAs(errFoo, &e, testr.WithMessage("assert err as customError"))
	assert.Panic(func() {}, testr.WithMessage("assert function is panic"))

	// Output:
	// false != expected:true // assert equality
	// error(foo) != expected:<nil> // assert err is nil
	// error(foo) != expected:as(*testr_test.customError) // assert err as customError
	// func() != expected:panic() // assert function is panic
}

func ExampleWithFailNow() {
	assert := testr.New(t) // using *testing.T

	assert.ErrorIs(errFoo, nil,
		testr.WithFailNow(),
		testr.WithMessage("using t.FailNow"),
	)

	// Output:
	// error(foo) != expected:<nil> // using t.FailNow
}

func ExampleMust() {
	r := strings.NewReader("testr")
	b := testr.Must(io.ReadAll(r))
	fmt.Println(string(b))

	// Output:
	// testr
}
