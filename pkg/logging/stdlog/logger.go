package stdlog

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/protomem/secrets-keeper/pkg/logging"
)

var ErrInvalidLogLevel = fmt.Errorf("invalid log level")

var _ logging.Logger = (*Logger)(nil)

type Logger struct {
	l *slog.Logger
}

func New(lvlStr string) (*Logger, error) {
	const op = "stdlog.New"
	var err error

	lvl, err := parseLogLeve(lvlStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, op)
	}

	l := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: lvl,
	}))

	return &Logger{
		l: l,
	}, nil
}

func (l *Logger) With(args ...any) logging.Logger {
	return &Logger{
		l: l.l.With(args...),
	}
}

func (l *Logger) Debug(msg string, args ...any) {
	l.l.Debug(msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.l.Info(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.l.Error(msg, args...)
}

func (l *Logger) Write(p []byte) (int, error) {
	logStr := string(p)
	logStr = strings.TrimSpace(logStr)
	logStr = strings.Trim(logStr, "\n")

	l.l.Info(logStr)

	return len(p), nil
}

func (l *Logger) Println(args ...any) {
	l.Debug(fmt.Sprintln(args...))
}

func (l *Logger) Sync(ctx context.Context) error {
	return nil
}

func parseLogLeve(lvlStr string) (slog.Level, error) {
	switch lvlStr {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, ErrInvalidLogLevel
	}
}
