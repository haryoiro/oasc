package main

import (
	"io"
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
	debug bool
}

func NewLogger(w io.Writer, debug bool) *Logger {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewTextHandler(w, opts)
	return &Logger{
		Logger: slog.New(handler),
		debug:  debug,
	}
}

func (l *Logger) Debug(format string, v ...any) {
	l.Logger.Debug(format, v...)
}

func (l *Logger) Info(format string, v ...any) {
	l.Logger.Info(format, v...)
}

func (l *Logger) Error(format string, v ...any) {
	l.Logger.Error(format, v...)
}

func (l *Logger) Fatal(format string, v ...any) {
	l.Error(format, v...)
	os.Exit(1)
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		Logger: l.Logger.With(args...),
		debug:  l.debug,
	}
}
