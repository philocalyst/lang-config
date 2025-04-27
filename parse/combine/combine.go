package main

import (
	"fmt"
	"io/fs"
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
	InputDir   string `arg:"-i,--input" help:"Directory containing language files" default:"language_files"`
	OutputFile string `arg:"-o,--output" help:"Output file path" default:"combined_languages.toml"`
	Format     string `arg:"-f,--format" help:"Output format: toml or yaml" default:"toml"`
	SourceExt  string `arg:"-s,--source" help:"Source files to process: toml, yaml, or both" default:"both"`
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
	if args.Format != "toml" && args.Format != "yaml" {
		log.Fatalf("Invalid output format: %s. Must be 'toml' or 'yaml'", args.Format)
	}

	// Validate source extension
	args.SourceExt = strings.ToLower(args.SourceExt)
	if args.SourceExt != "toml" && args.SourceExt != "yaml" && args.SourceExt != "both" {
		log.Fatalf("Invalid source extension: %s. Must be 'toml', 'yaml', or 'both'", args.SourceExt)
	}

	// Find all language files in the input directory
	languageFiles, err := findLanguageFiles(args.InputDir, args.SourceExt)
	if err != nil {
		log.Fatalf("Error finding language files: %v", err)
	}

	if len(languageFiles) == 0 {
		log.Fatalf("No language files found in %s", args.InputDir)
	}

	log.Infof("Found %d language files", len(languageFiles))

	// Collect all language entries
	var languages []map[string]interface{}

	for _, file := range languageFiles {
		lang, err := parseLanguageFile(file)
		if err != nil {
			log.Warnf("Error parsing %s: %v", file, err)
			continue
		}
		languages = append(languages, lang)
	}

	log.Infof("Successfully parsed %d language definitions", len(languages))

	// Create combined structure
	combined := map[string]interface{}{
		"language": languages,
	}

	// Create output file
	if err := os.MkdirAll(filepath.Dir(args.OutputFile), 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Generate output file in the specified format
	if args.Format == "toml" {
		writeTomlFile(args.OutputFile, combined)
	} else {
		writeYamlFile(args.OutputFile, combined)
	}

	fmt.Printf("Successfully generated %s with %d language entries\n",
		args.OutputFile, len(languages))
}

// findLanguageFiles returns a list of language files in the specified directory
func findLanguageFiles(dir, ext string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		lowerPath := strings.ToLower(path)
		if (ext == "toml" && strings.HasSuffix(lowerPath, ".toml")) ||
			(ext == "yaml" && (strings.HasSuffix(lowerPath, ".yaml") || strings.HasSuffix(lowerPath, ".yml"))) ||
			(ext == "both" && (strings.HasSuffix(lowerPath, ".toml") ||
				strings.HasSuffix(lowerPath, ".yaml") ||
				strings.HasSuffix(lowerPath, ".yml"))) {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// parseLanguageFile reads and parses a language file
func parseLanguageFile(filePath string) (map[string]interface{}, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	var parsed map[string]interface{}

	if ext == ".toml" {
		if err := toml.Unmarshal(data, &parsed); err != nil {
			return nil, fmt.Errorf("invalid TOML: %w", err)
		}
	} else if ext == ".yaml" || ext == ".yml" {
		if err := yaml.Unmarshal(data, &parsed); err != nil {
			return nil, fmt.Errorf("invalid YAML: %w", err)
		}
	} else {
		return nil, fmt.Errorf("unsupported file extension: %s", ext)
	}

	// Extract the language entry from the file
	langs, ok := parsed["language"].([]interface{})
	if !ok || len(langs) == 0 {
		return nil, fmt.Errorf("missing 'language' section or empty language array")
	}

	langMap, ok := langs[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid language entry format")
	}

	return langMap, nil
}

// writeTomlFile writes the combined structure to a TOML file
func writeTomlFile(filePath string, data map[string]interface{}) {
	tomlBytes, err := toml.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to marshal to TOML: %v", err)
	}

	if err := os.WriteFile(filePath, tomlBytes, 0644); err != nil {
		log.Fatalf("Failed to write TOML file: %v", err)
	}
}

// writeYamlFile writes the combined structure to a YAML file
func writeYamlFile(filePath string, data map[string]interface{}) {
	yamlBytes, err := yaml.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to marshal to YAML: %v", err)
	}

	if err := os.WriteFile(filePath, yamlBytes, 0644); err != nil {
		log.Fatalf("Failed to write YAML file: %v", err)
	}
}
