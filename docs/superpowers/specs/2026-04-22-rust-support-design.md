# Rust Support Design

## Overview

Add Rust language support to relinted as a standalone `internal/rust/` package, parallel to the existing C implementation. The key difference from C/Perl: Rust uses `,` (comma) as a statement-terminating punctuation that should be extracted as trailing punctuation — but only in match arm contexts (after `=>`).

## Package Structure

```
internal/rust/
  tokenize.go      — Rust tokenizer (identical to C tokenizer)
  punctuation.go   — extractTrailingPunctuation with Rust punctuation set
  format.go        — Format() entry point, mirrors C formatter algorithm
```

No changes to `internal/tokenizer/`, `internal/formatter/`, or any existing code.

## Tokenizer (`internal/rust/tokenize.go`)

Character-by-character tokenizer producing `tokenizer.Segment` values. **Identical to C tokenizer.**

| Token type | Trigger | Rules |
|---|---|---|
| `CommentLine` | `//` | Scan to end of line (include trailing newline) |
| `CommentBlock` | `/*` | Scan to `*/` |
| `String` | `"` | Scan to closing `"`, handle `\"` escapes |
| `String` | `'` | Scan to closing `'`, handle `\'` escapes (char literals) |
| `Code` | anything else | Scan to next special character |

Special characters that delimit Code segments: `"`, `'`, `//`, `/*`.

**No Rust-specific handling:**
- `#` in attributes like `#[derive(Debug)]` tokenizes as a regular operator token — fine for punctuation extraction
- Raw strings `r#"..."#` tokenize as `r` (identifier) + `"..."` (string) — fine for punctuation extraction
- `=>` operator tokenizes as two separate `Code` segments (`=>`) — fine for punctuation extraction

## Punctuation (`internal/rust/punctuation.go`)

`extractTrailingPunctuation(line string) (punctuation string, remaining string)`

Rust punctuation set: `{`, `}`, `;`, `,`

Algorithm: scan left-to-right, track whether position is inside a string, block comment, or line comment. Record the last punctuation character found in code context. If it is followed only by whitespace, extract it.

**Comma-specific rule:** `,` is only extracted when `=>` appears somewhere before the `,` on the same line (match arm context). This means:
- `Ok(num) => num,` → `,` is extracted (`=>` appears before `,` on same line)
- `Err(_) => continue,` → `,` is extracted (`=>` appears before `,` on same line)
- `Ordering::Less => println!("Too small!"),` → `,` is extracted (`=>` appears before `,` on same line)
- `vec![1, 2, 3]` → `,` is NOT extracted (no `=>` on line)
- `(a, b)` → `,` is NOT extracted (no `=>` on line)
- `.gen_range(1, 101);` → `,` is NOT extracted (no `=>` on line)

The `=>` detection: for each `,` found in code context, scan the line from start to the comma position. If `=>` appears anywhere before it, extract the comma.

## Formatter (`internal/rust/format.go`)

`Format(input string) string` — mirrors the C formatter algorithm exactly:

1. Expand tabs to 4-space stops
2. Strip trailing whitespace from each line
3. Tokenize input, reconstruct to normalize
4. Extract leading `{`/`}` → queue for previous line
5. Extract trailing punctuation → right-align
6. Filter lines that became empty (preserve originally empty lines)
7. Calculate max_len from content after punctuation extraction
8. Pad and append punctuation

Only difference from C formatter: calls Rust tokenizer and Rust punctuation extraction (which includes comma extraction in match arm context).

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
    ".rs": "rust",
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
relinted <input.rs> [output.rs]
relinted --lang rust <input> [output]
relinted -l rust <input> [output]
relinted --lang c <input.rs> [output]   # flag overrides extension
```

## Testing

- Unit tests for `internal/rust/tokenize.go` — verify segment types for strings, comments, mixed code
- Unit tests for `internal/rust/punctuation.go` — verify punctuation extraction with strings/comments, comma extraction in match arm vs non-match-arm contexts
- Unit tests for `internal/rust/format.go` — verify full formatting output
- Integration test: `./relinted -l rust linted-example-5.rs` must match `relinted-example-5.rs` exactly

## Files to Create

1. `internal/rust/tokenize.go`
2. `internal/rust/tokenize_test.go`
3. `internal/rust/punctuation.go`
4. `internal/rust/punctuation_test.go`
5. `internal/rust/format.go`
6. `internal/rust/format_test.go`

## Files to Modify

1. `cmd/relinted/main.go` — add `"rust"` case and `.rs` extension
2. `Justfile` — add rust test target
3. `README.md` — update usage with Rust language option
