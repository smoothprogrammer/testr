package testr

type option struct {
	t       T
	message string
	fail    func()
}

// Option represents options for the Assertion's methods.
type Option func(*option)

func newOption(t T, opts ...Option) option {
	opt := option{
		t:       t,
		message: "",
		fail:    t.Fail,
	}
	for _, f := range opts {
		f(&opt)
	}
	return opt
}

// WithMessage appends the message to the end of the output
// if the assertion is fail.
func WithMessage(message string) Option {
	return func(o *option) {
		if message != "" {
			o.message = " // " + message
		}
	}
}

// WithFailNow marks the test state as fail and stop the execution
// if the assertion is fail.
func WithFailNow() Option {
	return func(o *option) {
		o.fail = o.t.FailNow
	}
}
