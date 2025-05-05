package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	Version   = "0.1.0"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	config := ParseFlags()

	if config.Version {
		fmt.Printf("oasc version %s (build: %s, commit: %s)\n", Version, BuildTime, GitCommit)
		return
	}

	logger := NewLogger(os.Stderr, config.Debug)
	logger.Debug("Debug mode enabled")

	if len(config.InputFiles) < 2 {
		logger.Fatal("At least two input files must be specified")
	}

	logger.Debug("Reading first file: %s", config.InputFiles[0])
	merged, err := ReadFile(config.InputFiles[0])
	if err != nil {
		logger.Fatal("Error reading first file: %v", err)
	}

	for i := 1; i < len(config.InputFiles); i++ {
		logger.Debug("Reading file: %s", config.InputFiles[i])
		spec, err := ReadFile(config.InputFiles[i])
		if err != nil {
			logger.Fatal("Error reading file %s: %v", config.InputFiles[i], err)
		}
		merged = MergeSpecs(merged, spec)
	}

	var isJSON bool
	if config.OutputFormat != "" {
		config.OutputFormat = strings.ToLower(config.OutputFormat)
		isJSON = config.OutputFormat == "json"
	} else {
		outExt := strings.ToLower(filepath.Ext(config.OutputPath))
		isJSON = outExt == ".json"
	}

	logger.Debug("Writing output to: %s", config.OutputPath)
	if err := WriteFile(merged, config.OutputPath, isJSON); err != nil {
		logger.Fatal("Error writing output file: %v", err)
	}

	logger.Info("Successfully merged %d files to %s in %s format",
		len(config.InputFiles), config.OutputPath, map[bool]string{true: "JSON", false: "YAML"}[isJSON])
}
