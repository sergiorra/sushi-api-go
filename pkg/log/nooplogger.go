package log

import "context"

type noop struct {
}

// NewNoopLogger use for test
func NewNoopLogger() Logger {
	return &noop{}
}

func (l *noop) UnexpectedError(ctx context.Context, err error) {
	// nothing to do here, use for test
}