package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alexflint/go-arg"
	"github.com/pelletier/go-toml/v2"
	log "github.com/sirupsen/logrus"
)

// CLI defines command line arguments
type CLI struct {
	Input  string `arg:"-i,--input" help:"TOML file path" default:"languages.toml"`
	Output string `arg:"-o,--output" help:"Output JSON file path" default:"language_data.json"`
}

func main() {
	var args CLI
	arg.MustParse(&args)

	// Setup logging
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})

	// Read the input file
	data, err := os.ReadFile(args.Input)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	// Parse TOML into a map (using an interface type to handle unknown structure)
	var tomlData map[string]interface{}
	err = toml.Unmarshal(data, &tomlData)
	if err != nil {
		log.Fatalf("Error parsing TOML: %v", err)
	}

	// Extract language definitions
	languages, ok := tomlData["language"].([]interface{})
	if !ok {
		log.Fatal("Could not find or parse 'language' section in TOML")
	}

	log.Infof("Found %d language definitions", len(languages))

	// Transform to JSON structure
	langData := make(map[string]map[string]interface{})

	for _, lang := range languages {
		langMap, ok := lang.(map[string]interface{})
		if !ok {
			log.Warn("Skipping invalid language entry")
			continue
		}

		name, ok := langMap["name"].(string)
		if !ok || name == "" {
			log.Warn("Skipping language entry without a valid name")
			continue
		}

		// Create a new map for each language
		entry := make(map[string]interface{})

		// Collect all data for this language
		for key, value := range langMap {
			// Convert hyphens to underscores for JSON field names
			jsonKey := key
			if key == "comment-token" {
				jsonKey = "comment_token"
			} else if key == "comment-tokens" {
				jsonKey = "comment_tokens"
			} else if key == "block-comment-tokens" {
				jsonKey = "block_comment_tokens"
			} else if key == "file-types" {
				jsonKey = "file_types"
			}

			// Skip the name field as it's used as the map key
			if key != "name" {
				entry[jsonKey] = value
			}
		}

		langData[name] = entry
	}

	// Write output file
	if err := os.MkdirAll(filepath.Dir(args.Output), 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	jsonOutput, err := json.MarshalIndent(langData, "", "  ")
	if err != nil {
		log.Fatalf("Error generating JSON: %v", err)
	}

	if err := os.WriteFile(args.Output, jsonOutput, 0644); err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	fmt.Printf("Successfully generated %s with %d language entries\n",
		args.Output, len(langData))
}
