package argsx

type options[T any] struct {
	delimiter string
	defaultV  []T
}

type Option[T any] func(*options[T])

// getOpts parse option the default delimiter is ,
func getOpts[T any](opts []Option[T]) *options[T] {
	op := &options[T]{delimiter: ","}
	for _, opt := range opts {
		opt(op)
	}
	return op
}

// getDefault get the default value of T type
func (opts options[T]) getDefault() (t [][]T) {
	if opts.defaultV == nil {
		return
	}
	list := make([][]T, 0)
	return append(list, opts.defaultV)
}

// WithDelimiter specify slice delimiter default is ","
//
//	WithDelimiter[string]("-")
//	WithDelimiter[float64](".")
func WithDelimiter[T any](delimiter string) Option[T] {
	return func(opts *options[T]) {
		opts.delimiter = delimiter
	}
}

// WithDefault specify slice default value if payload is empty
//
//	WithDefault[string]("First", "Second")
//	WithDefault[float32](float32(1.2), float32(3.4))
func WithDefault[T any](dv ...T) Option[T] {
	return func(opts *options[T]) {
		opts.defaultV = dv
	}
}
