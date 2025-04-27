#!/usr/bin/env just

set shell := ["bash", "-eu", "-o", "pipefail", "-c"]

# ▰▰▰ Variables ▰▰▰ #
HELIX_PARSE_CLI := "parse/helix"
PARSE_CLI := "parse/all"

# ▰▰▰ Help ▰▰▰ #
help:
    @just --list

# ▰▰▰ Run ▰▰▰ #
run-parse-helix:
    @echo "Building parse_helix…"
    go run -C {{HELIX_PARSE_CLI}} "parse_helix"

run-parse-all:
    @echo "running parse…"
    go run -C {{PARSE_CLI}} "parse_all"

run-all: run-parse-helix run-parse-all

# ▰▰▰ Build ▰▰▰ #
build-parse-helix:
    @echo "Building parse_helix…"
    go build -C {{HELIX_PARSE_CLI}} -o {{justfile_directory()}}/"parse_helix"

build-parse-all:
    @echo "Building parse…"
    go build -C {{PARSE_CLI}} -o {{justfile_directory()}}/"parse_all"

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
    go test -C {{HELIX_PARSE_CLI}} "parse_helix"
        
test-parse-all:
    @echo "Testing parse…"
    go test -C {{HELIX_PARSE_CLI}} "parse_helix"


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
gen-json HELIX_TOML OUT:
    @echo "→ {{HELIX_PARSE_CLI}} -i {{HELIX_TOML}} -o {{OUT}}"
    {{HELIX_PARSE_CLI}} -i {{HELIX_TOML}} -o {{OUT}}

# Split the JSON manifest into per-language files.
#   just split-files IN=language_data.json OUT_DIR=language_files FORMAT=both
split-files IN OUT_DIR FORMAT:
    @echo "→ {{PARSE_CLI}} -i {{IN}} -o {{OUT_DIR}} -f {{FORMAT}}"
    {{PARSE_CLI}} -i {{IN}} -o {{OUT_DIR}} -f {{FORMAT}}

# ▰▰▰ Clean ▰▰▰ #
clean:
    @echo "Cleaning artifacts…"
    rm -rf bin language_files *.json
