# Changelog

All notable changes to this project will be documented in this file.

This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] – 2025-04-27

### Added

* Introduced `justfile` for streamlined project tasks including building, testing, generating files, formatting, linting, and cleaning.
* Added a new `combine` CLI tool (`parse/combine`) capable of merging individual language definition files (TOML or YAML) into a single manifest file (TOML or YAML).
* Included `just` recipes (`build-combine`, `combine`, `run-*`) for the new `combine` tool.
* Established `just` recipes (`gen-json`, `split-files`) for generating a unified JSON from Helix TOML and splitting manifests into individual files, respectively, with sensible default output locations.
* Created a comprehensive `README.md` detailing project goals, features, prerequisites (Go, just), setup instructions, and usage via `just` commands.

### Changed

* Refactored internal Go package structure:
    * Renamed `parse/from_helix` to `parse/helix`.
    * Renamed `parse/to_all` to `parse/all`.
* Improved variable usage and structure within the `justfile` for better maintainability.
* Updated `README.md` significantly to align with `justfile` usage and improve clarity.
* Modified `gen-json` and `split-files` `just` recipes to default output to the directory containing the `justfile`.

### Deprecated

* The top-level `language_data.json` file as the primary, manually managed source format. Generated manifests or individual language files are now preferred.

### Removed

* Obsolete Lua parsing script (`parse/parse.lua`).
* Pre-compiled Go binaries from the repository (`parse/all/parse`, `parse/helix/parse_helix`); use `just build-all` instead.
* Unused copy of `languages.toml` from the `parse/helix` directory.

## [0.1.0] - 2025-04-26

### Added

* Comprehensive language configuration data (`language_data.json`) supporting a wide range of programming languages, markup languages, and file formats.
* Individual language configuration files generated in both TOML and YAML formats located in the `language_files/` directory.
* New Go tool (`parse/to_all`) for parsing the main language configuration file (TOML or JSON) and generating individual language files in TOML and/or YAML formats.
* Initial configuration support for Arduino (`.ino`, `.pde`), AutoHotkey (`.ahk`), and Jai (`.jai`).

### Changed

* Refactored the location of the language configuration parsing tool to `parse/from_helix/`.

[Unreleased]: https://github.com/philocalyst/lang-config/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/philocalyst/lang-config/compare/v0.1.0..v0.2.0
[0.1.0]: https://github.com/philocalyst/lang-config/compare/...v0.1.0
