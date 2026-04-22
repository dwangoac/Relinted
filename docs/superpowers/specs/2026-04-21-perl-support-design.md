# Perl Support Design

## Overview

Add Perl language support to relinted as a standalone `internal/perl/` package, parallel to the existing C implementation. The CLI gains extension auto-detection and a `-l`/`--lang` flag to select the language.

## Package Structure

```
internal/perl/
  tokenize.go      — Perl tokenizer (Code, String, CommentLine segments)
  punctuation.go   — extractTrailingPunctuation with Perl punctuation set
  format.go        — Format() entry point, mirrors C formatter algorithm
```

No changes to `internal/tokenizer/`, `internal/formatter/`, or any existing code.

## Tokenizer (`internal/perl/tokenize.go`)

Character-by-character tokenizer producing `tokenizer.Segment` values.

| Token type | Trigger | Rules |
|---|---|---|
| `CommentLine` | `#` | Scan to end of line (include trailing newline) |
| `String` | `"` | Scan to closing `"`, handle `\"` escapes |
| `String` | `'` | Scan to closing `'`, handle `\'` escapes (Perl has no char literals) |
| `String` | `/` | Regex literal: scan to closing `/`, handle `\/` escapes |
| `Code` | anything else | Scan to next special character |

Special characters that delimit Code segments: `"`, `'`, `/`, `#`.

Angle brackets `<...>` are treated as Code (they can be I/O operators or comparison operators in Perl).

## Punctuation (`internal/perl/punctuation.go`)

`extractTrailingPunctuation(line string) (punctuation string, remaining string)`

Perl punctuation set: `{`, `}`, `;`

Algorithm: scan left-to-right, track whether position is inside a string, regex, or `#` comment. Record the last punctuation character found in code context. If it is followed only by whitespace, extract it.

## Formatter (`internal/perl/format.go`)

`Format(input string) string` — mirrors the C formatter algorithm exactly:

1. Expand tabs to 4-space stops
2. Strip trailing whitespace from each line
3. Tokenize input, reconstruct to normalize
4. Extract leading `{`/`}` → queue for previous line
5. Extract trailing punctuation → right-align
6. Filter lines that became empty (preserve originally empty lines)
7. Calculate max_len from content after punctuation extraction
8. Pad and append punctuation

Only difference from C formatter: calls Perl tokenizer and Perl punctuation extraction.

## CLI Changes (`cmd/relinted/main.go`)

### Extension auto-detection

```go
var extToLang = map[string]string{
    ".c":  "c",
    ".h":  "c",
    ".cpp":"c",
    ".cc": "c",
    ".pl": "perl",
    ".pm": "perl",
}
```

### Language selection

- `-l` / `--lang` flag (via stdlib `flag` package) overrides extension detection
- If neither flag nor recognized extension, default to `c`
- Unrecognized language prints error and exits

### Flow

```
-l/--lang flag provided? → yes → use that language
                             → no → check file extension
                                     → yes → use detected language
                                     → no  → default to "c"
```

### Usage

```
relinted <input.pl> [output.pl]
relinted --lang perl <input> [output]
relinted -l perl <input> [output]
relinted --lang c <input.pl> [output]   # flag overrides extension
```

## Testing

- Unit tests for `internal/perl/tokenize.go` — verify segment types for strings, comments, regex, mixed code
- Unit tests for `internal/perl/punctuation.go` — verify punctuation extraction with strings/regex/comments
- Unit tests for `internal/perl/format.go` — verify full formatting output
- Integration test: `./relinted -l perl linted-example-4.pl` must match `relinted-example-4.pl` exactly

## Files to Create

1. `internal/perl/tokenize.go`
2. `internal/perl/tokenize_test.go`
3. `internal/perl/punctuation.go`
4. `internal/perl/punctuation_test.go`
5. `internal/perl/format.go`
6. `internal/perl/format_test.go`

## Files to Modify

1. `cmd/relinted/main.go` — add language detection and flag parsing
2. `Justfile` — add perl test target
3. `README.md` — update usage with language options
