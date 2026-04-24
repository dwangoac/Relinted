# Go Support Design

## Overview

Add Go language support to relinted by creating a dedicated `internal/go/` package. The formatter output for Go code will match the ground truth file `examples/relinted-example-7.go`.

## Architecture

### New Package: `internal/go/`

Three source files plus one test file, mirroring the JS/Rust/Perl pattern.

#### `tokenize.go`

Go-specific tokenizer. Identical to the C tokenizer with one addition: backtick string literal handling. When a `` ` `` is encountered, scan to the next `` ` `` and treat the entire content as a single opaque `String` token. This prevents brace/semicolon detection inside raw string literals (e.g., `` `fmt.Println("hello")` ``).

Tokenization rules:
- `//` starts a line comment (scan to newline)
- `/*` starts a block comment (scan to `*/`)
- `"` starts a double-quoted string (scan to closing `"`, handling `\"` escapes)
- `'` starts a character literal (scan to closing `'`, handling `\'` escapes)
- `` ` `` starts a raw string literal (scan to closing `` ` `` — no escape processing)
- Everything else is Code

#### `punctuation.go`

Copy of C's `extractTrailingPunctuation` with backtick awareness added to the state machine. Backtick strings are skipped just like double-quoted strings, so braces/semicolons inside `` `...` `` are not extracted as trailing punctuation.

Shared helper functions (duplicated from C):
- `expandTabs` — replaces tabs with spaces to reach next 4-column tab stop
- `splitLines` — splits text into lines by `\n`
- `reconstructText` — joins tokenizer segments back together

#### `format.go`

Identical to C's `formatter.go`. The 8-step format pipeline (expand tabs, strip trailing whitespace, tokenize and reconstruct, extract leading braces, extract trailing punctuation, filter empty lines, calculate max width, format output) is duplicated here. No Go-specific formatting rules apply — Go braces and semicolons behave the same as C's for relinted's purposes.

#### `format_test.go`

Table-driven test that runs `Format` on `examples/linted-example-7.go` and compares output against `examples/relinted-example-7.go`.

### Changes to `cmd/relinted/main.go`

- Add `".go"` → `"go"` mapping to `extToLang`
- Add `"go"` case to the switch statement in `main()`, calling `go_pkg.Format(content)`
- Import the new package

### Changes to `README.md`

- Add `.go` to the Language Detection table
- Add `./relinted -l go source.go` to Usage Examples
- Add Go row to the example file table
- Update "Multi-Language" feature bullet to include Go
- Update the banner sentence to include Go

### Changes to `Justfile`

- Add `test-go` target running `go test ./internal/go/...`

## Testing

- Unit test: `examples/linted-example-7.go` formatted output matches `examples/relinted-example-7.go` exactly
- All existing tests continue to pass
