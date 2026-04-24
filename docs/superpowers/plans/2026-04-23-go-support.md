# Go Support Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add Go language support to relinted by creating `internal/go/` package and wiring it into the CLI, so `relinted examples/linted-example-7.go` produces output identical to `examples/relinted-example-7.go`.

**Architecture:** Mirror the existing JS/Rust/Perl pattern — a self-contained `internal/go/` package with tokenizer, punctuation extractor, and formatter. The tokenizer adds Go-specific backtick string literal handling. The formatter is identical to C's since Go braces/semicolons behave the same for relinted's purposes.

**Tech Stack:** Go 1.22+, standard library `flag`, `strings`, `unicode`, `testing`

---

### Task 1: Create `internal/go/tokenize.go`

**Files:**
- Create: `internal/go/tokenize.go`

- [ ] **Step 1: Write the Go tokenizer**

Create `internal/go/tokenize.go` with the following content. This is identical to the C tokenizer with one addition: backtick raw string literal handling.

```go
package go_pkg

import "relinted/internal/tokenizer"

// Tokenize scans Go source code character by character and produces a list of segments.
// Go-specific rules:
//   - // starts a line comment (scan to newline)
//   - /* starts a block comment (scan to */)
//   - " starts a double-quoted string (scan to closing ", handling \" escapes)
//   - ' starts a rune literal (scan to closing ', handling \' escapes)
//   - ` starts a raw string literal (scan to closing `, NO escape processing)
//   - Everything else is Code
func Tokenize(input string) []tokenizer.Segment {
	var segments []tokenizer.Segment
	i := 0
	for i < len(input) {
		switch {
		case i+1 < len(input) && input[i] == '/' && input[i+1] == '/':
			// Line comment: scan to end of line (include trailing newline)
			j := i
			for j < len(input) && input[j] != '\n' {
				j++
			}
			if j < len(input) {
				j++
			}
			segments = append(segments, tokenizer.Segment{tokenizer.CommentLine, input[i:j]})
			i = j

		case i+1 < len(input) && input[i] == '/' && input[i+1] == '*':
			// Block comment: scan to */
			j := i + 2
			for j+1 < len(input) && !(input[j] == '*' && input[j+1] == '/') {
				j++
			}
			if j+1 < len(input) {
				j += 2
			}
			segments = append(segments, tokenizer.Segment{tokenizer.CommentBlock, input[i:j]})
			i = j

		case input[i] == '"':
			// Double-quoted string: scan to closing ", handling \" escapes
			j := i + 1
			for j < len(input) && input[j] != '"' {
				if input[j] == '\\' && j+1 < len(input) {
					j++
				}
				j++
			}
			if j < len(input) {
				j++
			}
			segments = append(segments, tokenizer.Segment{tokenizer.String, input[i:j]})
			i = j

		case input[i] == '\'':
			// Rune literal: scan to closing ', handling \' escapes
			j := i + 1
			for j < len(input) && input[j] != '\'' {
				if input[j] == '\\' && j+1 < len(input) {
					j++
				}
				j++
			}
			if j < len(input) {
				j++
			}
			segments = append(segments, tokenizer.Segment{tokenizer.Char, input[i:j]})
			i = j

		case input[i] == '`':
			// Raw string literal: scan to closing `, NO escape processing
			j := i + 1
			for j < len(input) && input[j] != '`' {
				j++
			}
			if j < len(input) {
				j++
			}
			segments = append(segments, tokenizer.Segment{tokenizer.String, input[i:j]})
			i = j

		default:
			// Code: scan to next special character
			j := i + 1
			for j < len(input) {
				ch := input[j]
				if ch == '"' || ch == '\'' || ch == '`' {
					break
				}
				if j+1 < len(input) && ch == '/' && input[j+1] == '/' {
					break
				}
				if j+1 < len(input) && ch == '/' && input[j+1] == '*' {
					break
				}
				j++
			}
			segments = append(segments, tokenizer.Segment{tokenizer.Code, input[i:j]})
			i = j
		}
	}
	return segments
}
```

- [ ] **Step 2: Commit**

```bash
git add internal/go/tokenize.go
git commit -m "feat(go): add Go tokenizer with backtick string support

Assisted-by: qwen3.6 via Bash"
```

---

### Task 2: Create `internal/go/punctuation.go`

**Files:**
- Create: `internal/go/punctuation.go`

- [ ] **Step 1: Write the Go punctuation extractor**

Create `internal/go/punctuation.go`. This is a copy of C's `extractTrailingPunctuation` with backtick awareness added to the state machine.

```go
package go_pkg

import (
	"strings"
	"unicode"

	"relinted/internal/tokenizer"
)

// extractTrailingPunctuation scans the line to find the last punctuation character
// (semicolon, opening brace, or closing brace) that appears in code context
// (not inside a string, rune literal, raw string literal, block comment, or line comment).
//
// Returns (punctuation, remaining).
func extractTrailingPunctuation(line string) (punctuation string, remaining string) {
	lastPunctPos := -1
	inString := false
	inRune := false
	inRawString := false
	inBlockComment := false
	inLineComment := false

	for i := 0; i < len(line); i++ {
		ch := line[i]

		if inLineComment {
			continue
		}

		if inRawString {
			if ch == '`' {
				inRawString = false
			}
			continue
		}

		if inString {
			if ch == '\\' && i+1 < len(line) {
				i++
				continue
			}
			if ch == '"' {
				inString = false
			}
			continue
		}

		if inRune {
			if ch == '\\' && i+1 < len(line) {
				i++
				continue
			}
			if ch == '\'' {
				inRune = false
			}
			continue
		}

		if inBlockComment {
			if ch == '*' && i+1 < len(line) && line[i+1] == '/' {
				inBlockComment = false
				i++
			}
			continue
		}

		// Code context
		if ch == '"' {
			inString = true
			continue
		}
		if ch == '\'' {
			inRune = true
			continue
		}
		if ch == '`' {
			inRawString = true
			continue
		}
		if i+1 < len(line) && ch == '/' && line[i+1] == '/' {
			inLineComment = true
			continue
		}
		if i+1 < len(line) && ch == '/' && line[i+1] == '*' {
			inBlockComment = true
			i++
			continue
		}

		if ch == '{' || ch == '}' || ch == ';' {
			lastPunctPos = i
		}
	}

	if lastPunctPos >= 0 {
		rest := line[lastPunctPos+1:]
		allSpaceOrComment := true
		for i := 0; i < len(rest); i++ {
			ch := rest[i]
			if ch == '/' && i+1 < len(rest) {
				if rest[i+1] == '/' {
					// Line comment
					break
				}
				if rest[i+1] == '*' {
					// Block comment - scan to end
					i++
					for i+1 < len(rest) && !(rest[i] == '*' && rest[i+1] == '/') {
						i++
					}
					if i+1 < len(rest) {
						i++
					}
					break
				}
			}
			if !unicode.IsSpace(rune(ch)) {
				allSpaceOrComment = false
				break
			}
		}
		if allSpaceOrComment {
			return string(line[lastPunctPos]), line[:lastPunctPos]
		}
	}

	return "", line
}

// expandTabs replaces each tab character with spaces to reach the next
// 4-column tab stop.
func expandTabs(s string) string {
	var sb strings.Builder
	col := 0
	for _, ch := range s {
		if ch == '\t' {
			spaces := 4 - (col % 4)
			sb.WriteString(strings.Repeat(" ", spaces))
			col += spaces
		} else {
			sb.WriteRune(ch)
			col++
		}
	}
	return sb.String()
}

// splitLines splits text into lines by '\n'.
func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	result := strings.Split(s, "\n")
	return result
}

// reconstructText joins all segment texts back together.
func reconstructText(segments []tokenizer.Segment) string {
	var sb strings.Builder
	for _, seg := range segments {
		sb.WriteString(seg.Text)
	}
	return sb.String()
}
```

- [ ] **Step 2: Commit**

```bash
git add internal/go/punctuation.go
git commit -m "feat(go): add Go punctuation extractor with backtick awareness

Assisted-by: qwen3.6 via Bash"
```

---

### Task 3: Create `internal/go/format.go`

**Files:**
- Create: `internal/go/format.go`

- [ ] **Step 1: Write the Go formatter**

Create `internal/go/format.go`. This is identical to C's `formatter.go` — no Go-specific formatting rules apply.

```go
package go_pkg

import (
	"strings"
	"unicode"
)

// Format takes Go source code and reformats it by relocating braces and
// semicolons to the far right, creating a Python-like visual structure.
func Format(input string) string {
	// Step 1: Expand tabs to 4-space tab stops
	input = expandTabs(input)

	// Step 2: Strip trailing whitespace from each line and the final newline
	lines := splitLines(input)
	for i, line := range lines {
		lines[i] = strings.TrimRightFunc(line, unicode.IsSpace)
	}
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	// Step 3: Tokenize entire input and reconstruct to normalize
	segments := Tokenize(input)
	fullText := reconstructText(segments)
	lines = splitLines(fullText)
	for i, line := range lines {
		lines[i] = strings.TrimRightFunc(line, unicode.IsSpace)
	}
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	// Step 4: Process each line — extract leading braces and trailing punctuation
	type lineData struct {
		content         string
		originallyEmpty bool
		ownPunct        string
		queuedBraces    []string
	}

	data := make([]lineData, len(lines))
	for i, line := range lines {
		stripped := strings.TrimSpace(line)
		originallyEmpty := stripped == ""

		var leadingBrace string
		content := line

		if len(stripped) > 0 && (stripped[0] == '{' || stripped[0] == '}') {
			leadingBrace = string(stripped[0])
			bracePos := -1
			for j, ch := range line {
				if !unicode.IsSpace(ch) {
					bracePos = j
					break
				}
			}
			if bracePos >= 0 {
				leadingSpaces := line[:bracePos]
				rest := strings.TrimLeftFunc(line[bracePos+1:], unicode.IsSpace)
				content = leadingSpaces + rest
				if strings.TrimSpace(content) == "" {
					content = ""
				}
			}
		}

		trimmed := strings.TrimRightFunc(content, unicode.IsSpace)
		punct, rest := extractTrailingPunctuation(trimmed)
		content = rest

		data[i] = lineData{
			content:         content,
			originallyEmpty: originallyEmpty,
			ownPunct:        punct,
			queuedBraces:    []string{leadingBrace},
		}
	}

	// Step 5: Apply queued leading braces to previous line's trailing queue
	for i := 1; i < len(data); i++ {
		if len(data[i].queuedBraces) > 0 {
			data[i-1].ownPunct += strings.Join(data[i].queuedBraces, "")
		}
	}

	// Step 6: Filter lines that became empty (propagate punctuation to previous line)
	var filtered []lineData
	for _, d := range data {
		if d.content == "" && !d.originallyEmpty {
			if len(filtered) > 0 {
				filtered[len(filtered)-1].ownPunct += d.ownPunct
			}
			continue
		}
		filtered = append(filtered, d)
	}

	// Step 7: Calculate max_len from content after punctuation extraction
	maxLen := 0
	for _, d := range filtered {
		if len(d.content) > maxLen {
			maxLen = len(d.content)
		}
	}

	// Step 8: Format output
	var sb strings.Builder
	for i, d := range filtered {
		if i > 0 {
			sb.WriteString("\n")
		}

		if d.ownPunct != "" {
			padded := d.content
			if len(padded) < maxLen {
				padded += strings.Repeat(" ", maxLen-len(padded))
			}
			sb.WriteString(padded)
			sb.WriteString(" ")
			sb.WriteString(d.ownPunct)
		} else {
			sb.WriteString(d.content)
		}
	}
	sb.WriteString("\n")

	return sb.String()
}
```

- [ ] **Step 2: Commit**

```bash
git add internal/go/format.go
git commit -m "feat(go): add Go formatter (identical to C formatter)

Assisted-by: qwen3.6 via Bash"
```

---

### Task 4: Create `internal/go/format_test.go`

**Files:**
- Create: `internal/go/format_test.go`

- [ ] **Step 1: Write the Go format test**

Create `internal/go/format_test.go` with basic unit tests and the integration test against the ground truth file.

```go
package go_pkg

import (
	"os"
	"testing"
)

func TestFormat_SimpleSemicolon(t *testing.T) {
	input := "var x int = 1\n"
	got := Format(input)
	if got == "" {
		t.Fatal("expected non-empty output")
	}
	if got[len(got)-2:] != ";\n" {
		t.Errorf("expected output ending with ';\\n', got %q", got[len(got)-2:])
	}
}

func TestFormat_BraceRelocation(t *testing.T) {
	input := "func main() {\n    fmt.Println(\"hi\")\n}\n"
	got := Format(input)
	if got == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormat_EmptyLinesPreserved(t *testing.T) {
	input := "var x int = 1\n\nvar y int = 2\n"
	got := Format(input)
	newlines := 0
	for _, ch := range got {
		if ch == '\n' {
			newlines++
		}
	}
	if newlines != 3 {
		t.Errorf("expected 3 newlines (2 content + 1 trailing), got %d", newlines)
	}
}

func TestFormat_StringBracesNotExtracted(t *testing.T) {
	input := `fmt.Println("}")`
	got := Format(input)
	if got == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormat_BlockCommentPreserved(t *testing.T) {
	input := "/* comment */\nfunc main() {}\n"
	got := Format(input)
	if got == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormat_RuneLiteralNotExtracted(t *testing.T) {
	input := "var c rune = ';'\n"
	got := Format(input)
	if got == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormat_SingleLineNoPunct(t *testing.T) {
	expected := "var x\n"
	got := Format("var x")
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestFormat_EmptyInput(t *testing.T) {
	got := Format("")
	if got != "\n" {
		t.Errorf("expected newline output, got %q", got)
	}
}

func TestFormat_BraceRelocationExact(t *testing.T) {
	input := "func main() {\n    fmt.Println(\"hi\")\n}\n"
	got := Format(input)
	expected := "func main()          {\n    fmt.Println(\"hi\") ;}\n"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestFormat_RightAlignedPunctuationExact(t *testing.T) {
	input := "var x int = 1\nvar y int = 2\n"
	got := Format(input)
	expected := "var x int = 1 ;\nvar y int = 2 ;\n"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestFormat_PunctuationNotInStringExact(t *testing.T) {
	input := `fmt.Println("}")`
	got := Format(input)
	expected := "fmt.Println(\"}\") ;\n"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestFormat_RawStringNotExtracted(t *testing.T) {
	input := "`fmt.Println(\"}\")`\n"
	got := Format(input)
	if got == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormat_Integration(t *testing.T) {
	input, err := os.ReadFile("examples/linted-example-7.go")
	if err != nil {
		t.Fatalf("failed to read input file: %v", err)
	}
	expected, err := os.ReadFile("examples/relinted-example-7.go")
	if err != nil {
		t.Fatalf("failed to read expected output: %v", err)
	}
	got := Format(string(input))
	if got != string(expected) {
		t.Errorf("output does not match expected\n--- expected ---\n%s\n--- got ---\n%s\n", expected, got)
	}
}
```

- [ ] **Step 2: Run the tests to verify they pass**

```bash
go test ./internal/go/... -v
```
Expected: All 13 tests PASS.

- [ ] **Step 3: Commit**

```bash
git add internal/go/format_test.go
git commit -m "test(go): add Go format tests including integration with ground truth

Assisted-by: qwen3.6 via Bash"
```

---

### Task 5: Wire Go into `cmd/relinted/main.go`

**Files:**
- Modify: `cmd/relinted/main.go`

- [ ] **Step 1: Add Go extension mapping**

Add `".go"` → `"go"` to the `extToLang` map in `cmd/relinted/main.go`:

```go
var extToLang = map[string]string{
	".c":   "c",
	".h":   "c",
	".cpp": "c",
	".cc":  "c",
	".pl":  "perl",
	".pm":  "perl",
	".js":  "js",
	".rs":  "rust",
	".go":  "go",
}
```

- [ ] **Step 2: Add Go import**

Add the import for the new Go package:

```go
	"relinted/internal/go_pkg"
```

- [ ] **Step 3: Add Go case to switch statement**

Add `"go"` case to the switch statement that dispatches to formatters:

```go
	case "go":
		output = go_pkg.Format(content)
```

- [ ] **Step 4: Commit**

```bash
git add cmd/relinted/main.go
git commit -m "feat(cmd): add Go language support to CLI

Add .go extension mapping, import internal/go_pkg, and dispatch
to go_pkg.Format() for Go source files.

Assisted-by: qwen3.6 via Bash"
```

---

### Task 6: Update `README.md`

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update banner sentence**

Change the banner to include Go:

```markdown
Relinted reformats C/C++, Perl, Rust, JavaScript, and Go source code to visually resemble Python by aligning braces (`{`, `}`) and semicolons (`;`) to the far right of the codebase. This creates a clean, Python-like visual structure while preserving the original syntax and semantics.
```

- [ ] **Step 2: Update Multi-Language feature bullet**

Change:
```markdown
- **Multi-Language**: Supports C/C++, Perl, Rust, and JavaScript with auto-detection from file extension
```
To:
```markdown
- **Multi-Language**: Supports C/C++, Perl, Rust, JavaScript, and Go with auto-detection from file extension
```

- [ ] **Step 3: Add .go to Language Detection table**

Add a row to the Language Detection table:

| Extension | Language |
|-----------|----------|
| `.c`, `.h`, `.cpp`, `.cc` | C |
| `.pl`, `.pm` | Perl |
| `.rs` | Rust |
| `.js` | JavaScript |
| `.go` | Go |

- [ ] **Step 4: Add Go to Usage Examples**

Add a Go example in the "Force a specific language" section:

```bash
./relinted -l go source.go
```

- [ ] **Step 5: Add Go to example file table**

Add a row:

| Linted source      | Relinted output      |
| ------------------ | -------------------- |
| `examples/linted-example-7.go` | `examples/relinted-example-7.go` |

- [ ] **Step 6: Commit**

```bash
git add README.md
git commit -m "docs: add Go references to README.md

Add Go to banner, feature list, language detection table, usage
examples, and example file table.

Assisted-by: qwen3.6 via Bash"
```

---

### Task 7: Update `Justfile`

**Files:**
- Modify: `Justfile`

- [ ] **Step 1: Add test-go target**

Add after `test-rust`:

```
test-go:
    go test ./internal/go/...
```

- [ ] **Step 2: Commit**

```bash
git add Justfile
git commit -m "chore: add test-go target to Justfile

Assisted-by: qwen3.6 via Bash"
```

---

### Task 8: Build, verify, and commit final binary

**Files:**
- Build: `relinted` binary

- [ ] **Step 1: Run all tests**

```bash
just test
```
Expected: All packages pass.

- [ ] **Step 2: Run Go-specific tests**

```bash
just test-go
```
Expected: All 13 Go tests pass.

- [ ] **Step 3: Build the binary**

```bash
go build -o relinted ./cmd/relinted/
```

- [ ] **Step 4: Verify Go formatting matches ground truth**

```bash
./relinted examples/linted-example-7.go | diff - examples/relinted-example-7.go
```
Expected: No diff output (files match).

- [ ] **Step 5: Verify help output includes Go**

```bash
./relinted --help
```
Expected: Banner mentions Go.

- [ ] **Step 6: Commit**

```bash
git add relinted
git commit -m "build: rebuild relinted binary with Go support

Assisted-by: qwen3.6 via Bash"
```

---

## Self-Review

**1. Spec coverage:**
- `internal/go/tokenize.go` — implements backtick string tokenization ✓
- `internal/go/punctuation.go` — implements backtick-aware punctuation extraction ✓
- `internal/go/format.go` — implements identical-to-C formatter ✓
- `internal/go/format_test.go` — implements unit tests + integration test ✓
- `cmd/relinted/main.go` — adds `.go` mapping, import, and dispatch case ✓
- `README.md` — updates banner, features, language table, examples ✓
- `Justfile` — adds `test-go` target ✓

**2. Placeholder scan:** No "TBD", "TODO", or vague requirements found. Every step has exact code, commands, and expected outcomes.

**3. Type consistency:** Package name `go_pkg` used consistently across all tasks. Function names (`Format`, `Tokenize`, `extractTrailingPunctuation`) match across files.

**4. Scope check:** Focused on single feature — Go support. Self-contained. One executable output file to match.
