package testr_test

import "github.com/minizilla/testr"

func ExampleAssertion_Equal() {
	assert := testr.New(t) // using *testing.T

	assert.Equal(nil, nil)           // PASS
	assert.Equal(false, true)        // FAIL
	assert.Equal(int32(0), int64(0)) // FAIL

	// Output:
	// bool(false) != expected:bool(true)
	// int32(0) != expected:int64(0)
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
	// error(foo) != expected:nil()
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

func ExampleWithMessage() {
	assert := testr.New(t) // using *testing.T

	assert.Equal(int32(0), int64(0), testr.WithMessage("assert different type"))
	assert.ErrorIs(errFoo, nil, testr.WithMessage("assert err is nil"))
	var e customError
	assert.ErrorAs(errFoo, &e, testr.WithMessage("assert err as customError"))
	assert.Panic(func() {}, testr.WithMessage("assert function is panic"))

	// Output:
	// int32(0) != expected:int64(0) // assert different type
	// error(foo) != expected:nil() // assert err is nil
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
	// error(foo) != expected:nil() // using t.FailNow
}
