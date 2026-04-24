# Java Support Design

## Overview

Add Java language support to relinted by creating a dedicated `internal/java/` package. The formatter output for Java code will match the ground truth file `examples/relinted-example-8.java`.

## Architecture

### New Package: `internal/java/`

Three source files plus one test file, using Go's implementation as the base (closest C-family sibling).

#### `tokenize.go`

Based on Go's tokenizer with one addition: triple-quoted text block handling (`"""..."""`). When `"""` is encountered, scan to the next `"""` and treat the entire content as a single opaque `String` token. This prevents brace/semicolon detection inside text blocks (e.g., `"""fmt.Println("hello")"""`).

Tokenization rules:
- `//` starts a line comment (scan to newline)
- `/*` starts a block comment (scan to `*/`)
- `"` starts a double-quoted string (scan to closing `"`, handling `\"` escapes)
- `'` starts a char literal (scan to closing `'`, handling `\'` escapes)
- `"""` starts a text block literal (scan to closing `"""` — no escape processing)
- Everything else is Code

#### `punctuation.go`

Based on Go's punctuation extractor. Same state machine with `"""` awareness added. Text blocks are skipped just like double-quoted strings, so braces/semicolons inside `"""..."""` are not extracted as trailing punctuation.

Shared helper functions (duplicated from Go):
- `expandTabs` — replaces tabs with spaces to reach next 4-column tab stop
- `splitLines` — splits text into lines by `\n`
- `reconstructText` — joins tokenizer segments back together

#### `format.go`

Identical to Go's `format.go` (and C's `formatter.go`). No Java-specific formatting rules apply — Java braces and semicolons behave the same as C-family languages for relinted's purposes.

#### `format_test.go`

Table-driven test that runs `Format` on `examples/linted-example-8.java` and compares output against `examples/relinted-example-8.java`.

### Example File: `examples/linted-example-8.java`

~25-30 lines covering: package declaration, imports, class with fields, main method, if/else, for loop, string literals, char literals, line comments, block comments, and text blocks.

### Changes to `cmd/relinted/main.go`

- Add `".java"` → `"java"` mapping to `extToLang`
- Add `"java"` case to the switch statement in `main()`, calling `java_pkg.Format(content)`
- Import the new package

### Changes to `README.md`

- Add Java to the banner sentence
- Add `.java` to the Language Detection table
- Add Java row to the example file table
- Update "Multi-Language" feature bullet to include Java
- Add `./relinted -l java source.java` to Usage Examples

### Changes to `Justfile`

- Add `test-java` target running `go test ./internal/java/...`

## Testing

- Unit test: `examples/linted-example-8.java` formatted output matches `examples/relinted-example-8.java` exactly
- All existing tests continue to pass
