package pflagsrc

import (
	"context"

	"github.com/micro/go-config/source"
)

type delimKey struct{}

// WithDelimiter sets the delimiter to use for splitting nested keys. Default is '-'
func WithDelimiter(delim string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, delimKey{}, delim)
	}
}
