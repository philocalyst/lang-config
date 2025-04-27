# LangConfig

Centralized language definitions for editors, linters, and syntax-aware tools.

Manage filetypes, comment strings, extensions, and more in JSON, TOML, and YAML—all from a single source of truth. Kept up-to-date with the community and changing, shifting, times. It's a wild world out there. C++ might get a borrow checker. Don't venture into the veritable chaos alone. Let us lend you a hand.

---

## Why LangConfig?

Maintaining language metadata across different tools and editors can be a burden:

- Multiple formats (TOML, JSON, YAML)  
- Fragmented sources (editor-specific, plugin-specific)  
- Inconsistent comment strings, extensions, or embedded-filetype support  

LangConfig solves this by providing:

- A unified schema for language definitions  
- Easy conversion between TOML, JSON, and YAML  
- CLI utilities to import, generate, and split language data  
- An upcoming REST API for integration into dynamic applications.

Whether you’re authoring an editor, a plugin for an editor, or having fun with comment generation, lang-config allows you to operate with confidence on the languages you work with.

## Features

- Built off of Helix's `Languages.toml` -- all definitions in this early phase are editor-tested
- Well-structured data, lying in individual files or a single JSON manifest, so you can choose what you want to keep.
- Each language can hold a variety of options (comment tokens, block comment tokens, regex injections, active LSP's, diagonistic sources, roots, scopes, shebangs, supported auto-pairs)
- Easy to extend: add new languages, override defaults, or support new formats
- Support for a variety of configurations, and more to come. 

## Quickstart

### 1. Clone the repo

```shell
```

The project holds two CLI's
- `parse_helix` (imports Helix definitions)  
- `parse`       (generates per-language files)

### 2. Generate a Unified JSON

Convert Helix’s `languages.toml` into `language_data.json`:

```bash
parse_helix \
  -i path/to/languages.toml \
  -o language_data.json
```

### 3. Emit Per-Language Files

Split the JSON manifest into individual TOML and/or YAML files:

```bash
parse \
  -i language_data.json \
  -o language_files \
  -f both        # options: toml, yaml, both
```

You’ll end up with:
```
language_files/
├── python.toml
├── python.yaml
├── javascript.toml
├── javascript.yaml
└── ...
```

## Components

1. **parse_helix** (Go CLI)  
   - Reads Helix’s `languages.toml`  
   - Outputs a unified `language_data.json`

2. **parse** (Go CLI)  
   - Consumes `language_data.json` or a TOML manifest  
   - Emits per-language `{lang}.toml` and/or `{lang}.yaml`

## Design Principles

- **Unix-like**: small tools, single responsibility, pipeable workflows  
- **Modular**: clear separation between parsing, generation, and runtime use  
- **Extensible**: add new fields, formats, or override defaults with minimal code  
- **Accurate**: leverage Tree-sitter for embedded/compound filetype support

## Development

1. Clone this repo  
2. Ensure Go 1.24+ is installed  
3. Run tests:

   ```bash
   go test ./from_helix
   go test ./to_all
   ```

4. Build binaries with debug symbols:

   ```bash
   go build -o bin/parse_helix ./from_helix
   go build -o bin/parse        ./to_all
   ```

## Contributing

Contributions welcome! Whether it’s adding a new language, fixing a bug, or improving docs:
>[TIP!]
> If you're adding or updating a language configuration, JSON is the single source of truth going forward, so look there please.

1. Fork the repo  
2. Create a feature branch
3. Submit a PR against `main`  
4. Ensure tests pass and code is formatted (`go fmt`)

## Changelog

All notable changes and release notes can be found in [CHANGELOG.md](CHANGELOG.md).

## License

Distributed under the MIT License. See [LICENSE](LICENSE) for details.
