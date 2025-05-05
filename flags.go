package main

import (
	"flag"
	"strings"
)

type Config struct {
	InputFiles   []string
	OutputPath   string
	OutputFormat string
}

type multiStringFlag []string

func (m *multiStringFlag) String() string {
	return strings.Join(*m, ", ")
}

func (m *multiStringFlag) Set(value string) error {
	*m = append(*m, value)
	return nil
}

func ParseFlags() *Config {
	var config Config
	var inputFiles multiStringFlag

	flag.Var(&inputFiles, "file", "Input OpenAPI file paths (can be specified multiple times)")
	flag.Var(&inputFiles, "f", "Input OpenAPI file paths (shorthand)")

	flag.StringVar(&config.OutputPath, "output", "merged.yaml", "Output file path")
	flag.StringVar(&config.OutputPath, "o", "merged.yaml", "Output file path (shorthand)")

	flag.StringVar(&config.OutputFormat, "format", "", "Output format (json or yaml)")
	flag.StringVar(&config.OutputFormat, "F", "", "Output format (shorthand)")

	flag.Parse()

	config.InputFiles = inputFiles
	return &config
}
