package main

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedConfig *Config
	}{
		{
			name: "basic flag configuration",
			args: []string{
				"oasc",
				"-file", "file1.yaml",
				"-file", "file2.yaml",
				"-output", "merged.yaml",
				"-format", "yaml",
				"-debug",
			},
			expectedConfig: &Config{
				InputFiles:   []string{"file1.yaml", "file2.yaml"},
				OutputPath:   "merged.yaml",
				OutputFormat: "yaml",
				Debug:        true,
				Version:      false,
			},
		},
		{
			name: "version flag",
			args: []string{
				"oasc",
				"-version",
			},
			expectedConfig: &Config{
				Version: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			// Save original args
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			os.Args = tt.args

			config := ParseFlags()
			if tt.expectedConfig.Version {
				assert.Equal(t, tt.expectedConfig.Version, config.Version)
			} else {
				assert.Equal(t, tt.expectedConfig.InputFiles, config.InputFiles)
				assert.Equal(t, tt.expectedConfig.OutputPath, config.OutputPath)
				assert.Equal(t, tt.expectedConfig.OutputFormat, config.OutputFormat)
				assert.Equal(t, tt.expectedConfig.Debug, config.Debug)
				assert.Equal(t, tt.expectedConfig.Version, config.Version)
			}
		})
	}
}
