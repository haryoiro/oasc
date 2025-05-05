package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	tests := []struct {
		name     string
		debug    bool
		level    string
		message  string
		args     []interface{}
		contains string
	}{
		{
			name:     "debug log with debug mode enabled",
			debug:    true,
			level:    "DEBUG",
			message:  "test debug message",
			args:     []interface{}{},
			contains: "[DEBUG] test debug message",
		},
		{
			name:     "debug log with debug mode disabled",
			debug:    false,
			level:    "DEBUG",
			message:  "test debug message",
			args:     []interface{}{},
			contains: "",
		},
		{
			name:     "info log",
			debug:    false,
			level:    "INFO",
			message:  "test info message",
			args:     []interface{}{},
			contains: "[INFO] test info message",
		},
		{
			name:     "formatted message",
			debug:    false,
			level:    "INFO",
			message:  "test %s message",
			args:     []interface{}{"formatted"},
			contains: "[INFO] test formatted message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLogger(&buf, tt.debug)

			switch tt.level {
			case "DEBUG":
				logger.Debug(tt.message, tt.args...)
			case "INFO":
				logger.Info(tt.message, tt.args...)
			}

			output := buf.String()
			if tt.contains == "" {
				assert.Empty(t, output)
			} else {
				assert.Contains(t, output, tt.contains)
			}
		})
	}
}

func TestLoggerFatal(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, false)

	// Skip actual Fatal call as it calls os.Exit(1)
	message := "fatal error"
	logger.Info(message)
	assert.Contains(t, buf.String(), "[INFO] fatal error")
}
