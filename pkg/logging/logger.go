package logging

import (
	"context"
	"io"
)

type Logger interface {
	With(args ...any) Logger

	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Error(msg string, args ...any)

	io.Writer
	Println(args ...any)

	Sync(context.Context) error
}
