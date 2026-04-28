# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [7.0] - 2026-04-28

### Added
- `--version` flag to print version and exit
- Cross-platform binary releases (linux/amd64, linux/arm64, darwin/arm64, windows/amd64)
- MIT license

### Fixed
- C language support: added proper package wrapper and fixed punctuation extraction
- Added language-specific test targets in Justfile (test-c, test-perl, test-rust, test-js, test-go, test-java, test-php, test-swift)

## [6.0] - 2026-04-15

### Added
- Swift language support (tokenizer, formatter, punctuation extraction)
- PHP language support (tokenizer, formatter, punctuation extraction)
- C# language support (shares Java formatter)
- TypeScript language support (shares JavaScript formatter)
- Java language support (8-step pipeline)
- Go language support

