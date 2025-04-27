#!/usr/bin/env just

set shell := ["bash", "-eu", "-o", "pipefail", "-c"]

# ▰▰▰ Variables ▰▰▰ #
HELIX_PARSE_LOCATION := "parse/helix"
PARSE_ALL_LOCATION := "parse/all"
HELIX_PARSE_CLI := "parse_helix"
PARSE_ALL_CLI := "parse_all"

# ▰▰▰ Help ▰▰▰ #
help:
    @just --list

# ▰▰▰ Run ▰▰▰ #
run-parse-helix:
    @echo "Building parse_helix…"
    go run -C {{HELIX_PARSE_LOCATION}} {{HELIX_PARSE_CLI}}

run-parse-all:
    @echo "running parse…"
    go run -C {{PARSE_ALL_LOCATION}} {{PARSE_ALL_CLI}}

run-all: run-parse-helix run-parse-all

# ▰▰▰ Build ▰▰▰ #
build-parse-helix:
    @echo "Building parse_helix…"
    go build -C {{HELIX_PARSE_LOCATION}} -o {{justfile_directory()}}/{{HELIX_PARSE_CLI}}

build-parse-all:
    @echo "Building parse…"
    go build -C {{PARSE_ALL_LOCATION}} -o {{justfile_directory()}}/{{PARSE_ALL_CLI}}

build-all: build-parse-helix build-parse-all

# ▰▰▰ Install ▰▰▰ #
install-from-helix:
    @echo "Installing parse_helix…"
    go install github.com/you/langconfig/from_helix@latest

install-to-all:
    @echo "Installing parse…"
    go install github.com/you/langconfig/to_all@latest

install-all: install-from-helix install-to-all

# ▰▰▰ Test ▰▰▰ #
test-parse-helix:
    @echo "Testing parse_helix…"
    go test -C {{HELIX_PARSE_LOCATION}} {{HELIX_PARSE_CLI}}
        
test-parse-all:
    @echo "Testing parse…"
    go test -C {{HELIX_PARSE_LOCATION}} {PARSE_ALL_CLI}


test-all: test-parse-helix test-parse-all

# ▰▰▰ Format & Lint ▰▰▰ #
fmt:
    @echo "Formatting Go code and tidying modules…"
    go fmt ./...
    go mod tidy

lint:
    @echo "Running static analysis…"
    staticcheck ./from_helix/... ./to_all/...

# ▰▰▰ Generate Files ▰▰▰ #
# Generate a unified JSON from Helix's TOML.
#   just gen-json HELIX_TOML=path/to/languages.toml OUT=language_data.json
gen-json HELIX_TOML OUT: build-parse-helix
    @echo "→ {{HELIX_PARSE_CLI}} -i {{HELIX_TOML}} -o {{OUT}}"
    {{HELIX_PARSE_CLI}} -i {{HELIX_TOML}} -o {{OUT}}

# Split the JSON manifest into per-language files.
#   just split-files IN=language_data.json OUT_DIR=language_files FORMAT=both
split-files IN OUT_DIR FORMAT:
    @echo "→ {{PARSE_ALL_CLI}} -i {{IN}} -o {{OUT_DIR}} -f {{FORMAT}}"
    {{PARSE_ALL_CLI}} -i {{IN}} -o {{OUT_DIR}} -f {{FORMAT}}

# ▰▰▰ Clean ▰▰▰ #
clean:
    @echo "Cleaning artifacts…"
    rm -rf bin language_files *.json
