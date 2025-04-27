# LangConfig

Centralized language definitions for editors, linters, and syntax-aware tools.

Manage filetypes, comment strings, extensions, and more in JSON, TOML, and YAML—all from a single source of truth. Kept up-to-date with the community and changing, shifting times. It’s a wild world out there. C++ might get a borrow checker. Don’t venture into the veritable chaos alone—let us lend you a hand.

---

## Why LangConfig?

Maintaining language metadata across different tools and editors can be a burden:

- Multiple formats (TOML, JSON, YAML)  
- Fragmented sources (editor-specific, plugin-specific)  
- Inconsistent comment strings, extensions, or embedded-filetype support  

**LangConfig** solves this by providing:

- A unified schema for language definitions  
- Easy conversion between TOML, JSON, and YAML  
- CLI utilities to import, generate, and split language data  
- An upcoming REST API for dynamic integrations  

Whether you’re authoring an editor, a plugin, or having fun with comment generation, LangConfig lets you operate with confidence.

---

## Features

- Based on Helix’s `Languages.toml`—all definitions in this early phase are editor-tested  
- Individual files or single JSON manifest—pick your workflow  
- Rich per-language options (comment tokens, block tokens, injections, LSPs, diagnostics, roots, scopes, shebangs, auto-pairs…)  
- Easy to extend: add or override formats, fields, defaults  
- Accurate Tree-sitter-powered support for embedded/compound filetypes  

---

## Prerequisites

- Go 1.24+  
- [just](https://github.com/casey/just) (a modern task runner)

---

## Quickstart

1. Clone the repo

   ```bash
   git clone https://github.com/you/langconfig.git
   cd langconfig
   ```

2. Explore available tasks

   ```bash
   just --list
   ```

3. Generate a unified JSON from Helix’s TOML:

   ```bash
   just gen-json HELIX_TOML=path/to/languages.toml OUT=language_data.json
   ```

4. Split that JSON into per-language TOML/YAML:

   ```bash
   just split-files IN=language_data.json OUT_DIR=language_files FORMAT=both
   ```

   You’ll end up with:

   ```
   language_files/
   ├── python.toml
   ├── python.yaml
   ├── javascript.toml
   ├── javascript.yaml
   └── …
   ```

---

## Justfile Commands

Use `just <recipe> [ARGS…]` to streamline your workflow:

• help  
    List all available tasks.

• gen-json HELIX_TOML OUT  
    Convert Helix’s `languages.toml` to a single JSON manifest.

• split-files IN OUT_DIR FORMAT  
    Split a JSON manifest into `{lang}.toml` and/or `{lang}.yaml`.

• run-all  
    Runs both parse steps (for ad-hoc testing).

• build-all  
    Build both `parse_helix` and `parse_all` binaries into the project root.

• install-all  
    `go install` the CLIs into your `GOPATH`.

• test-all  
    Run all unit tests.

• fmt  
    `go fmt` + `go mod tidy`.

• lint  
    Static analysis via `staticcheck`.

• clean  
    Remove binaries, JSON outputs, and generated language files.

---

## Components

1. **parse_helix** (Go CLI)  
   Reads Helix’s `languages.toml` → outputs `language_data.json`

2. **parse** (Go CLI)  
   Reads `language_data.json` (or TOML) → emits per-language TOML/YAML

---

## Development

1. Fork & clone  
2. Install prerequisites (Go, just)  
3. Run tests

   ```bash
   just test-all
   ```

4. Build binaries

   ```bash
   just build-all
   ```

5. Generate or split language data

   ```bash
   just gen-json HELIX_TOML=… OUT=…
   just split-files IN=… OUT_DIR=… FORMAT=…
   ```

---

## Contributing

Contributions welcome! When updating or adding new languages, remember:
> JSON is the single source of truth going forward—please author there.

1. Fork the repo  
2. Create a feature branch  
3. Submit a PR against `main`  
4. Ensure tests pass and code is formatted (`just fmt`)

---

## Changelog

See [CHANGELOG.md](CHANGELOG.md).

---

## License

Distributed under the MIT License. See [LICENSE](LICENSE).
