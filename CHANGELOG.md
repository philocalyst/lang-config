# Changelog

All notable changes to this project will be documented in this file.

This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2025-04-26

### Added

* Comprehensive language configuration data (`language_data.json`) supporting a wide range of programming languages, markup languages, and file formats.
* Individual language configuration files generated in both TOML and YAML formats located in the `language_files/` directory.
* New Go tool (`parse/to_all`) for parsing the main language configuration file (TOML or JSON) and generating individual language files in TOML and/or YAML formats.
* Initial configuration support for Arduino (`.ino`, `.pde`), AutoHotkey (`.ahk`), and Jai (`.jai`).

### Changed

* Refactored the location of the language configuration parsing tool to `parse/from_helix/`.

[Unreleased]: https://github.com/philocalyst/infat/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/philocalyst/languages/compare/...v0.1.0
