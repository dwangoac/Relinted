# JavaScript Support Design

## Goal

Add JavaScript language support to relinted by creating a new `internal/js/` package with tokenizer, punctuation extraction, and formatter. JavaScript uses `{`, `}`, and `;` as relocatable punctuation — same set as C and Perl.

## Architecture

Each language has its own `internal/<lang>/` package with 3 files: `tokenize.go`, `punctuation.go`, `format.go`. The CLI in `main.go` dispatches to the appropriate package via the `-l`/`--lang` flag.

## Components

### Tokenizer (`internal/js/tokenize.go`)

Identical to the C tokenizer. Handles:
- `//` line comments
- `/* */` block comments
- `"` string literals
- `'` char literals

No JavaScript-specific rules needed since JS regex literals require complex state tracking that doesn't affect brace/semicolon relocation.

### Punctuation (`internal/js/punctuation.go`)

Punctuation set: `{`, `}`, `;` — same as C and Perl. No special context rules (unlike Rust's match-arm comma). Extracts leading braces and trailing punctuation from each line.

### Formatter (`internal/js/format.go`)

Same algorithm as the C formatter:
1. Tokenize input into segments
2. For each line, extract leading braces and trailing punctuation
3. Relocate braces to their own lines
4. Remove trailing punctuation from code lines
5. Place them on a continuation line at column 80

### CLI Update (`main.go`)

- Add `"js"` case to language switch
- Add `.js` extension to `extToLang` map
- Import `internal/js`
- Update unsupported language error message

### Test Files

- `linted-example-6.js` — JavaScript "guess the number" game (input)
- `relinted-example-6.js` — reformatted output (ground truth)

### Justfile

- Add `test-js` target
- Update language detection table
- Update usage examples

### README

- Update supported languages list
- Update features table
- Add JS example transformation
- Update test table
