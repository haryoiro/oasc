package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

func main() {
	config := ParseFlags()
	if len(config.InputFiles) < 2 {
		log.Fatal("At least two input files must be specified")
	}

	merged, err := ReadFile(config.InputFiles[0])
	if err != nil {
		log.Fatalf("Error reading first file: %v", err)
	}

	for i := 1; i < len(config.InputFiles); i++ {
		spec, err := ReadFile(config.InputFiles[i])
		if err != nil {
			log.Fatalf("Error reading file %s: %v", config.InputFiles[i], err)
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

	if err := WriteFile(merged, config.OutputPath, isJSON); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully merged %d files to %s in %s format\n",
		len(config.InputFiles), config.OutputPath, map[bool]string{true: "JSON", false: "YAML"}[isJSON])
}
