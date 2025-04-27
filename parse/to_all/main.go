package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/pelletier/go-toml/v2"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// CLI defines command line arguments
type CLI struct {
	Input     string `arg:"-i,--input" help:"Input file (TOML or JSON)" default:"languages.toml"`
	OutputDir string `arg:"-o,--outdir" help:"Output directory for language files" default:"language_files"`
	Format    string `arg:"-f,--format" help:"Output format: toml, yaml, or both" default:"both"`
}

func main() {
	var args CLI
	arg.MustParse(&args)

	// Setup logging
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})

	// Validate format option
	args.Format = strings.ToLower(args.Format)
	if args.Format != "toml" && args.Format != "yaml" && args.Format != "both" {
		log.Fatalf("Invalid format: %s. Must be 'toml', 'yaml', or 'both'", args.Format)
	}

	// Read the input file
	data, err := os.ReadFile(args.Input)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	// Determine file type and parse accordingly
	var languages map[string]map[string]interface{}
	ext := strings.ToLower(filepath.Ext(args.Input))

	if ext == ".toml" {
		// Parse TOML to intermediate map
		var tomlData map[string]interface{}
		if err := toml.Unmarshal(data, &tomlData); err != nil {
			log.Fatalf("Error parsing TOML: %v", err)
		}

		// Extract language entries
		langs, ok := tomlData["language"].([]interface{})
		if !ok {
			log.Fatal("Could not find or parse 'language' section in TOML")
		}

		languages = make(map[string]map[string]interface{}, len(langs))
		for _, lang := range langs {
			langMap, ok := lang.(map[string]interface{})
			if !ok {
				continue
			}

			name, ok := langMap["name"].(string)
			if !ok || name == "" {
				continue
			}

			// Create language entry
			entry := make(map[string]interface{})
			for key, value := range langMap {
				if key != "name" {
					entry[key] = value
				}
			}

			languages[name] = entry
		}
	} else if ext == ".json" {
		// Parse JSON directly
		if err := json.Unmarshal(data, &languages); err != nil {
			log.Fatalf("Error parsing JSON: %v", err)
		}
	} else {
		log.Fatalf("Unsupported input file format: %s", ext)
	}

	// Create output directory
	if err := os.MkdirAll(args.OutputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Create files for each language
	for name, data := range languages {
		// Sanitize filename
		filename := sanitizeFilename(name)

		if args.Format == "toml" || args.Format == "both" {
			// Create TOML file
			outputToml(filepath.Join(args.OutputDir, filename+".toml"), name, data)
		}

		if args.Format == "yaml" || args.Format == "both" {
			// Create YAML file
			outputYaml(filepath.Join(args.OutputDir, filename+".yaml"), name, data)
		}
	}

	fmt.Printf("Successfully generated %d language files in %s\n",
		len(languages), args.OutputDir)
}

// outputToml writes a language definition to a TOML file
func outputToml(filepath string, name string, data map[string]interface{}) {
	// Create a wrapper structure for TOML
	tomlData := map[string]interface{}{
		"language": []map[string]interface{}{
			{
				"name": name,
			},
		},
	}

	// Add all properties
	for key, value := range data {
		tomlData["language"].([]map[string]interface{})[0][key] = value
	}

	// Marshal to TOML
	tomlBytes, err := toml.Marshal(tomlData)
	if err != nil {
		log.Warnf("Failed to marshal %s to TOML: %v", name, err)
		return
	}

	// Write to file
	if err := os.WriteFile(filepath, tomlBytes, 0644); err != nil {
		log.Warnf("Failed to write TOML file for %s: %v", name, err)
	}
}

// outputYaml writes a language definition to a YAML file
func outputYaml(filepath string, name string, data map[string]interface{}) {
	// Create a wrapper structure for YAML
	yamlData := map[string]interface{}{
		"language": []map[string]interface{}{
			{
				"name": name,
			},
		},
	}

	// Add all properties
	for key, value := range data {
		yamlData["language"].([]map[string]interface{})[0][key] = value
	}

	// Marshal to YAML
	yamlBytes, err := yaml.Marshal(yamlData)
	if err != nil {
		log.Warnf("Failed to marshal %s to YAML: %v", name, err)
		return
	}

	// Write to file
	if err := os.WriteFile(filepath, yamlBytes, 0644); err != nil {
		log.Warnf("Failed to write YAML file for %s: %v", name, err)
	}
}

// sanitizeFilename makes sure the language name is a valid filename
func sanitizeFilename(name string) string {
	// Replace characters that are problematic in filenames
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
		" ", "_",
		"#", "_",
	)

	return replacer.Replace(name)
}
