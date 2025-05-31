package main

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	debug  bool
	logger *log.Logger
}

func NewLogger(w io.Writer, debug bool) *Logger {
	return &Logger{
		debug:  debug,
		logger: log.New(w, "", log.LstdFlags),
	}
}

func (l *Logger) Debug(format string, v ...any) {
	if l.debug {
		l.logger.Printf("[DEBUG] "+format, v...)
	}
}

func (l *Logger) Info(format string, v ...any) {
	l.logger.Printf("[INFO] "+format, v...)
}

func (l *Logger) Fatal(format string, v ...any) {
	l.logger.Printf("[FATAL] "+format, v...)
	os.Exit(1)
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		debug:  l.debug,
		logger: log.New(l.logger.Writer(), l.logger.Prefix(), l.logger.Flags()),
	}
}
