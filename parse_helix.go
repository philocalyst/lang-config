package main

import (
	"encoding/json"
	"flag"
	"io"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/pelletier/go-toml/v2"
)

func info(format string, a ...any) {
	color.Cyan(format+"\n", a...)
}

func magenta(format string, a ...any) {
	color.Magenta(format+"\n", a...)
}

func success(format string, a ...any) {
	color.Green(format+"\n", a...)
}

func warn(format string, a ...any) {
	color.New(color.FgYellow).Fprintf(os.Stderr, format+"\n", a...)
}

func fatal(format string, a ...any) {
	color.New(color.FgRed).Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(1)
}

func firstString(raw any) (string, bool) {
	switch v := raw.(type) {
	case string:
		return v, true
	case []any:
		if len(v) > 0 {
			if s, ok := v[0].(string); ok {
				return s, true
			}
		}
	}
	return "", false
}

func toStringList(raw any) ([]string, bool, int) {
	arr, ok := raw.([]any)
	if !ok {
		return nil, false, 0
	}
	var out []string
	for _, item := range arr {
		if s, ok := item.(string); ok {
			out = append(out, s)
		}
	}
	return out, true, len(arr) - len(out)
}

func readInput(path string) ([]byte, error) {
	if path == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(path)
}

func main() {
	var inputPath, outputPath string
	flag.StringVar(&inputPath, "i", "-", "TOML source path or '-' for stdin")
	flag.StringVar(&inputPath, "input", "-", "TOML source path or '-' for stdin")
	flag.StringVar(&outputPath, "o", "language_data.json", "Output JSON file path")
	flag.StringVar(&outputPath, "output", "language_data.json", "Output JSON file path")
	flag.Parse()

	parserName := "go-toml (github.com/pelletier/go-toml/v2)"
	magenta("Using TOML parser: %s", parserName)
	magenta("Reading configuration from: %s", inputPath)

	data, err := readInput(inputPath)
	if err != nil {
		fatal("Error reading %s: %v", inputPath, err)
	}

	info("Parsing TOML data...")
	var config map[string]interface{}
	if err := toml.Unmarshal(data, &config); err != nil {
		fatal("TOML parse error: %v", err)
	}

	languageData := make(map[string]interface{})
	rawLangs, exists := config["language"]
	if !exists {
		warn("No '[[language]]' sections found")
	} else if langs, ok := rawLangs.([]interface{}); !ok || len(langs) == 0 {
		warn("No '[[language]]' sections found")
	} else {
		success("Found %d language entries.", len(langs))
		for _, rawEntry := range rawLangs.([]interface{}) {
			entry, ok := rawEntry.(map[string]interface{})
			if !ok {
				warn("Skipping non-table entry in 'language'")
				continue
			}
			nameRaw, _ := entry["name"]
			name, okName := nameRaw.(string)
			if !okName || name == "" {
				warn("Skipping a language entry without a name")
				continue
			}

			dataMap := make(map[string]interface{})
			if tok, ok := firstString(entry["comment-token"]); ok {
				dataMap["comment_token"] = tok
			} else if tok, ok := firstString(entry["comment-tokens"]); ok {
				dataMap["comment_token"] = tok
			}
			if block, exists := entry["block-comment-tokens"]; exists && block != nil {
				dataMap["block_comment_tokens"] = block
			}
			if fts, exists := entry["file-types"]; exists {
				if list, ok, skipped := toStringList(fts); !ok {
					warn("Ignoring non-list file-types for '%s'", name)
				} else {
					if skipped > 0 {
						warn("Skipped %d non-string file-types for '%s'",
							skipped, name)
					}
					dataMap["file_types"] = list
				}
			}

			if len(dataMap) > 0 {
				languageData[name] = dataMap
			} else {
				warn("No data for '%s', skipping", name)
			}
		}
	}

	info("Writing JSON output to: %s", outputPath)
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		fatal("Error creating directories: %v", err)
	}
	outBytes, err := json.MarshalIndent(languageData, "", "    ")
	if err != nil {
		fatal("Error marshaling JSON: %v", err)
	}
	if err := os.WriteFile(outputPath, outBytes, 0644); err != nil {
		fatal("Error writing JSON: %v", err)
	}
	success("Done! JSON file generated successfully.")
}
