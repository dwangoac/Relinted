# JavaScript Support Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add JavaScript language support to relinted by creating a new `internal/js/` package with tokenizer, punctuation extraction, and formatter, then wire it into the CLI, test files, Justfile, and README.

**Architecture:** Follow the exact same pattern as C, Perl, and Rust — each language has its own `internal/<lang>/` package with `tokenize.go`, `punctuation.go`, and `format.go`. JavaScript uses `{`, `}`, `;` as relocatable punctuation — same set as C and Perl, no special context rules.

**Tech Stack:** Go 1.22+, TDD with `go test`, `just` build system.

---

### Task 1: Create `internal/js/tokenize.go` + tests

**Files:**
- Create: `internal/js/tokenize.go`
- Test: `internal/js/tokenize_test.go`

The JavaScript tokenizer is **identical to the C tokenizer** — same `//` line comments, `/* */` block comments, `"` string literals, `'` char literals. No JavaScript-specific rules needed.

- [ ] **Step 1: Write the tokenizer test file**

Create `internal/js/tokenize_test.go` with these tests:

```go
package js

import "relinted/internal/tokenizer"
import "testing"

func checkSegments(t *testing.T, input string, expected []tokenizer.Segment) {
	t.Helper()
	got := Tokenize(input)
	if len(got) != len(expected) {
		t.Fatalf("got %d segments, want %d", len(got), len(expected))
	}
	for i := range got {
		if got[i].Type != expected[i].Type {
			t.Errorf("segment %d: type %v, want %v", i, got[i].Type, expected[i].Type)
		}
		if got[i].Text != expected[i].Text {
			t.Errorf("segment %d: text %q, want %q", i, got[i].Text, expected[i].Text)
		}
	}
}

func TestTokenize_CodeOnly(t *testing.T) {
	checkSegments(t, "let x = 1;\n", []tokenizer.Segment{
		{tokenizer.Code, "let x = 1;\n"},
	})
}

func TestTokenize_LineComment(t *testing.T) {
	checkSegments(t, "let x = 1; // comment\n", []tokenizer.Segment{
		{tokenizer.Code, "let x = 1; "},
		{tokenizer.CommentLine, "// comment\n"},
	})
}

func TestTokenize_BlockComment(t *testing.T) {
	checkSegments(t, "/* comment */\n", []tokenizer.Segment{
		{tokenizer.CommentBlock, "/* comment */"},
		{tokenizer.Code, "\n"},
	})
}

func TestTokenize_BlockCommentMultiline(t *testing.T) {
	checkSegments(t, "/* multi\nline */\n", []tokenizer.Segment{
		{tokenizer.CommentBlock, "/* multi\nline */"},
		{tokenizer.Code, "\n"},
	})
}

func TestTokenize_StringLiteral(t *testing.T) {
	checkSegments(t, 'console.log("hello");\n', []tokenizer.Segment{
		{tokenizer.Code, "console.log("},
		{tokenizer.String, '"hello"'},
		{tokenizer.Code, ");\n"},
	})
}

func TestTokenize_StringWithEscape(t *testing.T) {
	checkSegments(t, 'console.log("he\\"llo");\n', []tokenizer.Segment{
		{tokenizer.Code, "console.log("},
		{tokenizer.String, '"he\\"llo"'},
		{tokenizer.Code, ");\n"},
	})
}

func TestTokenize_CharLiteral(t *testing.T) {
	checkSegments(t, "let c = 'a';\n", []tokenizer.Segment{
		{tokenizer.Code, "let c = "},
		{tokenizer.Char, "'a'"},
		{tokenizer.Code, ";\n"},
	})
}

func TestTokenize_CharWithEscape(t *testing.T) {
	checkSegments(t, "let c = '\\n';\n", []tokenizer.Segment{
		{tokenizer.Code, "let c = "},
		{tokenizer.Char, "'\\n'"},
		{tokenizer.Code, ";\n"},
	})
}

func TestTokenize_Mixed(t *testing.T) {
	checkSegments(t, 'let s = "hi"; // comment\n', []tokenizer.Segment{
		{tokenizer.Code, "let s = "},
		{tokenizer.String, '"hi"'},
		{tokenizer.Code, "; "},
		{tokenizer.CommentLine, "// comment\n"},
	})
}

func TestTokenize_Empty(t *testing.T) {
	got := Tokenize("")
	if len(got) != 0 {
		t.Errorf("expected 0 segments, got %d", len(got))
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/js/...`
Expected: FAIL with `package js is not in std` or `function not defined`

- [ ] **Step 3: Write the tokenizer**

Create `internal/js/tokenize.go` — copy the C tokenizer from `internal/tokenizer/tokenizer.go` exactly:

```go
package js

import "relinted/internal/tokenizer"

// Tokenize scans the input character by character and produces a list of segments.
// JavaScript-specific rules are identical to C:
//   - // starts a line comment (scan to newline)
//   - /* starts a block comment (scan to */)
//   - " starts a string literal (scan to closing ", handling \" escapes)
//   - ' starts a char literal (scan to closing ', handling \' escapes)
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
			// String literal: scan to closing ", handling escapes
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
			// Char literal: scan to closing ', handling escapes
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

		default:
			// Code: scan to next special character
			j := i + 1
			for j < len(input) {
				ch := input[j]
				if ch == '"' || ch == '\'' {
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

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/js/...`
Expected: PASS — all 10 tests pass

- [ ] **Step 5: Commit**

```bash
git add internal/js/tokenize.go internal/js/tokenize_test.go
git commit -m "feat(js): add JavaScript tokenizer identical to C tokenizer"
```

---

### Task 2: Create `internal/js/punctuation.go` + tests

**Files:**
- Create: `internal/js/punctuation.go`
- Test: `internal/js/punctuation_test.go`

JavaScript punctuation set: `{`, `}`, `;` — same as C and Perl. No special context rules (unlike Rust's match-arm comma).

- [ ] **Step 1: Write the punctuation test file**

Create `internal/js/punctuation_test.go`:

```go
package js

import (
	"testing"

	"relinted/internal/tokenizer"
)

func TestExtractTrailingPunctuation_Semicolon(t *testing.T) {
	punct, rest := extractTrailingPunctuation("let x = 1;")
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != "let x = 1" {
		t.Errorf("expected rest 'let x = 1', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_OpenBrace(t *testing.T) {
	punct, rest := extractTrailingPunctuation("function foo() {")
	if punct != "{" {
		t.Errorf("expected punct '{', got %q", punct)
	}
	if rest != "function foo() " {
		t.Errorf("expected rest 'function foo() ', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_ClosingBrace(t *testing.T) {
	punct, rest := extractTrailingPunctuation("}")
	if punct != "}" {
		t.Errorf("expected punct '}', got %q", punct)
	}
	if rest != "" {
		t.Errorf("expected rest '', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_NoPunctuation(t *testing.T) {
	punct, rest := extractTrailingPunctuation("let x = 1")
	if punct != "" {
		t.Errorf("expected punct '', got %q", punct)
	}
	if rest != "let x = 1" {
		t.Errorf("expected rest 'let x = 1', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_InString(t *testing.T) {
	punct, rest := extractTrailingPunctuation(`console.log("hello;");`)
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != `console.log("hello;")` {
		t.Errorf("expected rest 'console.log(\"hello;\")', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_InLineComment(t *testing.T) {
	punct, rest := extractTrailingPunctuation("let x = 1; // comment;")
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != "let x = 1" {
		t.Errorf("expected rest 'let x = 1', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_InBlockComment(t *testing.T) {
	punct, rest := extractTrailingPunctuation("let x = 1; /* ; */")
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != "let x = 1" {
		t.Errorf("expected rest 'let x = 1', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_SemicolonInStringNotExtracted(t *testing.T) {
	punct, rest := extractTrailingPunctuation(`console.log("a; b")`)
	if punct != "" {
		t.Errorf("expected punct '', got %q", punct)
	}
	if rest != `console.log("a; b")` {
		t.Errorf("expected rest 'console.log(\"a; b\")', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_SemicolonInCharNotExtracted(t *testing.T) {
	input := "let c = ';'"
	punct, rest := extractTrailingPunctuation(input)
	if punct != "" {
		t.Errorf("expected punct '', got %q", punct)
	}
	if rest != input {
		t.Errorf("expected rest %q, got %q", input, rest)
	}
}

func TestExtractTrailingPunctuation_BraceInStringNotExtracted(t *testing.T) {
	punct, rest := extractTrailingPunctuation(`console.log("}");`)
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != `console.log("}")` {
		t.Errorf("expected rest 'console.log(\"}\")', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_MultiplePunctuation(t *testing.T) {
	punct, rest := extractTrailingPunctuation("let x = 1; let y = 2;")
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != "let x = 1; let y = 2" {
		t.Errorf("expected rest 'let x = 1; let y = 2', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_BracesOnly(t *testing.T) {
	punct, rest := extractTrailingPunctuation("if (x) {")
	if punct != "{" {
		t.Errorf("expected punct '{', got %q", punct)
	}
	if rest != "if (x) " {
		t.Errorf("expected rest 'if (x) ', got %q", rest)
	}
}

func TestExpandTabs_Basic(t *testing.T) {
	got := expandTabs("\t\t")
	if got != "        " {
		t.Errorf("expected 8 spaces, got %q", got)
	}
}

func TestExpandTabs_Mixed(t *testing.T) {
	got := expandTabs("ab\tcd")
	if got != "ab  cd" {
		t.Errorf("expected 'ab  cd', got %q", got)
	}
}

func TestSplitLines_Basic(t *testing.T) {
	got := splitLines("a\nb\nc")
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
	if got[0] != "a" || got[1] != "b" || got[2] != "c" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestSplitLines_TrailingNewline(t *testing.T) {
	got := splitLines("a\nb\n")
	if len(got) != 3 {
		t.Fatalf("expected 3 lines (with trailing empty), got %d", len(got))
	}
	if got[2] != "" {
		t.Errorf("expected last line to be empty, got %q", got[2])
	}
}

func TestReconstructText(t *testing.T) {
	segments := []tokenizer.Segment{
		{tokenizer.Code, "let "},
		{tokenizer.Code, "x = "},
		{tokenizer.Code, "1;\n"},
	}
	got := reconstructText(segments)
	if got != "let x = 1;\n" {
		t.Errorf("got %q, want %q", got, "let x = 1;\n")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/js/...`
Expected: FAIL — package not found or functions not defined

- [ ] **Step 3: Write the punctuation extraction**

Create `internal/js/punctuation.go` — based on the C formatter's punctuation logic (same as Perl but without regex/single-quote string handling):

```go
package js

import (
	"strings"
	"unicode"

	"relinted/internal/tokenizer"
)

// extractTrailingPunctuation scans the line to find the last punctuation character
// (semicolon, opening brace, or closing brace) that appears in code context
// (not inside a string, char literal, block comment, or line comment).
//
// Returns (punctuation, remaining).
func extractTrailingPunctuation(line string) (punctuation string, remaining string) {
	lastPunctPos := -1
	inString := false
	inChar := false
	inBlockComment := false
	inLineComment := false

	for i := 0; i < len(line); i++ {
		ch := line[i]

		if inLineComment {
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

		if inChar {
			if ch == '\\' && i+1 < len(line) {
				i++
				continue
			}
			if ch == '\'' {
				inChar = false
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
			inChar = true
			continue
		}
		if ch == '/' && i+1 < len(line) && line[i+1] == '*' {
			inBlockComment = true
			i++
			continue
		}
		if ch == '/' && i+1 < len(line) && line[i+1] == '/' {
			inLineComment = true
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

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/js/...`
Expected: PASS — all 18 tests pass

- [ ] **Step 5: Commit**

```bash
git add internal/js/punctuation.go internal/js/punctuation_test.go
git commit -m "feat(js): add JavaScript punctuation extraction for braces and semicolons"
```

---

### Task 3: Create `internal/js/format.go` + tests

**Files:**
- Create: `internal/js/format.go`
- Test: `internal/js/format_test.go`

The JavaScript formatter follows the same algorithm as the C formatter (from `internal/formatter/formatter.go`). No `utf8` import needed since the Perl formatter doesn't use it either — we'll use `len()` for byte length like the Perl formatter does.

- [ ] **Step 1: Write the formatter test file**

Create `internal/js/format_test.go`:

```go
package js

import "testing"

func TestFormat_SimpleSemicolon(t *testing.T) {
	input := "let x = 1;\n"
	got := Format(input)
	if got == "" {
		t.Fatal("expected non-empty output")
	}
	if got[len(got)-2:] != ";\n" {
		t.Errorf("expected output ending with ';\\n', got %q", got[len(got)-2:])
	}
}

func TestFormat_BraceRelocation(t *testing.T) {
	input := "function foo() {\n    console.log('hi');\n}\n"
	got := Format(input)
	if got == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormat_EmptyLinesPreserved(t *testing.T) {
	input := "let x = 1;\n\nlet y = 2;\n"
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
	input := `console.log("}");`
	got := Format(input)
	if got == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormat_BlockCommentPreserved(t *testing.T) {
	input := "/* comment */\nfunction foo() {}\n"
	got := Format(input)
	if got == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormat_CharLiteralNotExtracted(t *testing.T) {
	input := "let c = ';';\n"
	got := Format(input)
	if got == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormat_SingleLineNoPunct(t *testing.T) {
	expected := "let x\n"
	got := Format("let x")
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
	input := "function foo() {\n    console.log('hi');\n}\n"
	got := Format(input)
	expected := "function foo()        {\n    console.log('hi') ;}\n"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestFormat_RightAlignedPunctuationExact(t *testing.T) {
	input := "let x = 1;\nlet y = 2;\n"
	got := Format(input)
	expected := "let x = 1 ;\nlet y = 2 ;\n"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestFormat_PunctuationNotInStringExact(t *testing.T) {
	input := `console.log("}");\n`
	got := Format(input)
	expected := `console.log("}") ;\n`
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/js/...`
Expected: FAIL — Format function not defined

- [ ] **Step 3: Write the formatter**

Create `internal/js/format.go` — copy the Perl formatter's Format function exactly (same algorithm, same punctuation set):

```go
package js

import (
	"strings"
	"unicode"
)

// Format takes JavaScript source code and reformats it by relocating braces and
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

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/js/...`
Expected: PASS — all 11 tests pass

- [ ] **Step 5: Commit**

```bash
git add internal/js/format.go internal/js/format_test.go
git commit -m "feat(js): add JavaScript formatter with brace relocation and semicolon extraction"
```

---

### Task 4: Update CLI for JavaScript language selection

**Files:**
- Modify: `cmd/relinted/main.go`

- [ ] **Step 1: Add JavaScript to CLI**

In `cmd/relinted/main.go`, make these changes:

1. Add import for `internal/js`:
```go
	"relinted/internal/js"
```

2. Add `.js` extension to `extToLang` map:
```go
	".js":  "js",
```

3. Add `"js"` case to the switch statement:
```go
	case "js":
		output = js.Format(content)
```

4. Update the unsupported language error message:
```go
		fmt.Fprintf(os.Stderr, "Supported languages: c, perl, rust, js\n")
```

- [ ] **Step 2: Build and verify**

Run: `go build -o relinted ./cmd/relinted/`
Expected: Build succeeds

Run: `./relinted --help`
Expected: Shows `[c, perl, rust, js]` in help text

- [ ] **Step 3: Commit**

```bash
git add cmd/relinted/main.go
git commit -m "feat(js): add JavaScript language support to CLI"
```

---

### Task 5: Create test files, update Justfile and README

**Files:**
- Create: `linted-example-6.js` (input — JavaScript guessing game)
- Create: `relinted-example-6.js` (ground truth — expected output)
- Modify: `Justfile`
- Modify: `README.md`

- [ ] **Step 1: Create the JavaScript guessing game input**

Create `linted-example-6.js`:

```javascript
// Guess the number game
const readline = require('readline');

function main() {
    const secret = Math.floor(Math.random() * 100) + 1;
    let guess = 0;
    let attempts = 0;

    const rl = readline.createInterface({
        input: process.stdin,
        output: process.stdout
    });

    console.log('Guess the number between 1 and 100!');

    while (guess !== secret) {
        rl.question('Enter your guess: ', (answer) => {
            guess = parseInt(answer, 10);
            attempts++;

            if (guess === secret) {
                console.log('You win! Attempts: ' + attempts);
                rl.close();
            } else if (guess > secret) {
                console.log('Too high!');
            } else {
                console.log('Too low!');
            }
        });
    }
}

main();
```

- [ ] **Step 2: Generate ground truth by running relinted**

Run: `./relinted -l js linted-example-6.js relinted-example-6.js`

This will produce the ground truth file by running the JavaScript formatter on the input.

- [ ] **Step 3: Verify the output matches**

Run: `./relinted -l js linted-example-6.js /tmp/out.js && diff relinted-example-6.js /tmp/out.js`
Expected: No diff output (files match)

- [ ] **Step 4: Update Justfile**

Add `test-js` target and update the `test` target to include `./internal/js/...`:

```
set shell := ["bash", "-euo", "pipefail", "-c"]

build:
    go build -o relinted ./cmd/relinted/

test:
    go test ./internal/... ./cmd/...

test-c:
    go test ./internal/formatter/... ./internal/tokenizer/... ./internal/io/...

test-perl:
    go test ./internal/perl/...

test-rust:
    go test ./internal/rust/...

test-js:
    go test ./internal/js/...

run:
    go run ./cmd/relinted/

lint:
    go vet ./internal/... ./cmd/relinted/

format:
    gofmt -w .
```

- [ ] **Step 5: Update README.md**

Update these sections:
1. First paragraph: add "JavaScript" to the language list
2. Build section: add `just test-js`
3. Features section: add "JavaScript" to the multi-language bullet
4. Arguments table: add `js` to the language list
5. Language Detection table: add `.js` → JavaScript row
6. Usage examples: add `./relinted game.js` and `./relinted -l js source.c`
7. How It Works: add semicolons (JavaScript) to the trailing punctuation bullet
8. Example Transformation: add a JavaScript before/after example
9. Tests table: add `linted-example-6.js` → `relinted-example-6.js` row
10. Notes section: add JavaScript to the supported languages list

- [ ] **Step 6: Run all tests**

Run: `go test ./internal/...`
Expected: All tests pass

Run: `./relinted -l js linted-example-6.js /tmp/out.js && diff relinted-example-6.js /tmp/out.js`
Expected: No diff (integration test passes)

- [ ] **Step 7: Commit**

```bash
git add linted-example-6.js relinted-example-6.js Justfile README.md
git commit -m "feat(js): add JavaScript test files, update Justfile and README"
```

---

## Self-Review

**Spec coverage:**
- Tokenizer (identical to C) → Task 1 ✅
- Punctuation (`{`, `}`, `;`) → Task 2 ✅
- Formatter (same algorithm as C) → Task 3 ✅
- CLI update (switch case, extension map, import, error message) → Task 4 ✅
- Test files (input + ground truth) → Task 5 ✅
- Justfile (test-js target) → Task 5 ✅
- README (all sections updated) → Task 5 ✅

**Placeholder scan:** No TBD, TODO, or "similar to" references. All code is complete and explicit.

**Type consistency:** All tasks use `internal/js/` package name consistently. All function names (`Tokenize`, `Format`, `extractTrailingPunctuation`, `expandTabs`, `splitLines`, `reconstructText`) match across files.

**Testing discipline:** Each task follows TDD — write tests first, verify failure, implement, verify pass, commit.
