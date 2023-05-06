package testr

type option struct {
	message string
}

type Option func(*option)

func newOption(opts ...Option) option {
	var opt option
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
