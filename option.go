package testr

type option struct {
	t       T
	message string
	fail    func()
}

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

func WithMessage(message string) Option {
	return func(o *option) {
		if message != "" {
			o.message = " // " + message
		}
	}
}

func WithFailNow() Option {
	return func(o *option) {
		o.fail = o.t.FailNow
	}
}
