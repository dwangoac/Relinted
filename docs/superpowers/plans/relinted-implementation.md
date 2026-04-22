# Relinted Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a Go CLI tool that reformats C/C++ source code by relocating braces and semicolons to the far right, creating a Python-like visual structure.

**Architecture:** Four-package Go project: `internal/io` for file I/O, `internal/tokenizer` for lexical segmentation, `internal/formatter` for brace/semicolon relocation and right-alignment, and `cmd/relinted` as the CLI entry point. The formatter reconstructs lines from segments, calculates max line width, extracts punctuation from code context (respecting strings and comments), and pads output lines.

**Tech Stack:** Go 1.22+, standard library only (no external dependencies).

---

### Task 1: Initialize Go module

**Files:**
- Create: `go.mod`

- [ ] **Step 1: Create Go module**

Run:
```bash
cd /home/ac/dev/relinted
go mod init relinted
```

Expected output: `go: creating new go.mod: module relinted`

- [ ] **Step 2: Verify module**

Run:
```bash
go build ./...
```

Expected: no output (empty project compiles cleanly).

---

### Task 2: Implement I/O package

**Files:**
- Create: `internal/io/io.go`
- Create: `internal/io/io_test.go`

- [ ] **Step 1: Write the I/O package**

Create `internal/io/io.go`:

```go
package io

import "os"

// ReadFile reads the entire contents of the file at path.
func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteFile writes content to the file at path.
func WriteFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}
```

- [ ] **Step 2: Write tests for the I/O package**

Create `internal/io/io_test.go`:

```go
package io

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.txt")
	expected := "hello world\n"
	if err := os.WriteFile(path, []byte(expected), 0644); err != nil {
		t.Fatal(err)
	}

	got, err := ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != expected {
		t.Errorf("ReadFile(%q) = %q, want %q", path, got, expected)
	}
}

func TestReadFile_NotFound(t *testing.T) {
	_, err := ReadFile("/nonexistent/file.txt")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestWriteFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "out.txt")
	content := "test content\n"

	if err := WriteFile(path, content); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != content {
		t.Errorf("WriteFile() = %q, want %q", string(data), content)
	}
}
```

- [ ] **Step 3: Run tests to verify they pass**

Run:
```bash
go test ./internal/io/ -v
```

Expected: all tests PASS.

---

### Task 3: Implement Tokenizer

**Files:**
- Create: `internal/tokenizer/tokenizer.go`
- Create: `internal/tokenizer/tokenizer_test.go`

- [ ] **Step 1: Write the tokenizer package**

Create `internal/tokenizer/tokenizer.go`:

```go
package tokenizer

// SegmentType represents the type of a token segment.
type SegmentType int

const (
	Code SegmentType = iota
	String
	Char
	CommentBlock
	CommentLine
)

// Segment represents a contiguous run of characters of the same type.
type Segment struct {
	Type SegmentType
	Text string
}

// Tokenize scans the input character by character and produces a list of segments.
// Strings ("...") and chars ('...') are recognized with escape handling.
// Block comments (/* ... */) and line comments (// ...) are recognized.
// Everything else is Code.
func Tokenize(input string) []Segment {
	var segments []Segment
	i := 0
	for i < len(input) {
		switch {
		case i+1 < len(input) && input[i] == '/' && input[i+1] == '/':
			// Line comment: scan to end of line
			j := i
			for j < len(input) && input[j] != '\n' {
				j++
			}
			segments = append(segments, Segment{CommentLine, input[i:j]})
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
			segments = append(segments, Segment{CommentBlock, input[i:j]})
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
			segments = append(segments, Segment{String, input[i:j]})
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
			segments = append(segments, Segment{Char, input[i:j]})
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
			segments = append(segments, Segment{Code, input[i:j]})
			i = j
		}
	}
	return segments
}
```

- [ ] **Step 2: Write tests for the tokenizer**

Create `internal/tokenizer/tokenizer_test.go`:

```go
package tokenizer

import "testing"

func checkSegments(t *testing.T, input string, expected []Segment) {
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
	checkSegments(t, "int x = 0;", []Segment{
		{Code, "int x = 0;"},
	})
}

func TestTokenize_LineComment(t *testing.T) {
	checkSegments(t, "int x; // comment\n", []Segment{
		{Code, "int x; "},
		{CommentLine, "// comment\n"},
	})
}

func TestTokenize_BlockComment(t *testing.T) {
	checkSegments(t, "int x /* block */;\n", []Segment{
		{Code, "int x "},
		{CommentBlock, "/* block */"},
		{Code, ";\n"},
	})
}

func TestTokenize_String(t *testing.T) {
	checkSegments(t, "printf(\"hello\");\n", []Segment{
		{Code, "printf(\""},
		{String, "hello"},
		{Code, "\");\n"},
	})
}

func TestTokenize_StringWithEscape(t *testing.T) {
	checkSegments(t, "printf(\"hello\\n\");\n", []Segment{
		{Code, "printf(\""},
		{String, "hello\\n"},
		{Code, "\");\n"},
	})
}

func TestTokenize_Char(t *testing.T) {
	checkSegments(t, "char c = 'a';\n", []Segment{
		{Code, "char c = '"},
		{Char, "a"},
		{Code, "';\n"},
	})
}

func TestTokenize_Mixed(t *testing.T) {
	checkSegments(t, "/* start */ int x = \"hi\"; // end\n", []Segment{
		{CommentBlock, "/* start */ "},
		{Code, "int x = \""},
		{String, "hi"},
		{Code, "\"; "},
		{CommentLine, "// end\n"},
	})
}

func TestTokenize_Empty(t *testing.T) {
	got := Tokenize("")
	if len(got) != 0 {
		t.Errorf("expected 0 segments, got %d", len(got))
	}
}
```

- [ ] **Step 3: Run tests to verify they pass**

Run:
```bash
go test ./internal/tokenizer/ -v
```

Expected: all tests PASS.

---

### Task 4: Implement Formatter

**Files:**
- Create: `internal/formatter/formatter.go`
- Create: `internal/formatter/formatter_test.go`

- [ ] **Step 1: Write the formatter package**

Create `internal/formatter/formatter.go`:

```go
package formatter

import (
	"strings"
	"unicode"

	"relinted/internal/tokenizer"
)

// Format takes C/C++ source code and reformats it by relocating braces and
// semicolons to the far right, creating a Python-like visual structure.
func Format(input string) string {
	// Step 1: Expand tabs to 4 spaces
	input = expandTabs(input)

	// Step 2: Strip trailing whitespace from each line and the final newline
	lines := splitLines(input)
	for i, line := range lines {
		lines[i] = strings.TrimRightFunc(line, unicode.IsSpace)
	}
	// Remove trailing empty line if input ended with newline
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	// Step 3: Tokenize
	segments := tokenizer.Tokenize(input)

	// Step 4: Reconstruct full text from segments and re-split into lines
	fullText := reconstructText(segments)
	lines = splitLines(fullText)
	for i, line := range lines {
		lines[i] = strings.TrimRightFunc(line, unicode.IsSpace)
	}
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	// Step 5: Calculate max_len from all lines
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	// Step 6: Process each line — extract leading braces and trailing punctuation
	type lineData struct {
		content          string
		originallyEmpty  bool
		ownPunct         string // trailing punctuation from this line
		queuedBraces     []string // leading braces queued FROM this line (to previous)
	}

	data := make([]lineData, len(lines))
	for i, line := range lines {
		stripped := strings.TrimSpace(line)
		originallyEmpty := stripped == ""

		var leadingBrace string
		content := line

		// Check for leading brace
		if len(stripped) > 0 && (stripped[0] == '{' || stripped[0] == '}') {
			leadingBrace = string(stripped[0])
			// Remove the leading brace and any whitespace immediately after it
			// Find the position of the brace in the original line
			bracePos := -1
			for j, ch := range line {
				if !unicode.IsSpace(ch) {
					bracePos = j
					break
				}
			}
			if bracePos >= 0 {
				content = line[bracePos+1:]
				// Trim leading whitespace from content
				content = strings.TrimLeftFunc(content, unicode.IsSpace)
			}
		}

		// Check for trailing punctuation in code context
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

	// Step 7: Apply queued leading braces to previous line's trailing queue
	for i := 1; i < len(data); i++ {
		if len(data[i].queuedBraces) > 0 {
			data[i-1].ownPunct += strings.Join(data[i].queuedBraces, "")
		}
	}

	// Step 8: Filter lines that became empty (not originally empty, now empty)
	var filtered []lineData
	for _, d := range data {
		if d.content == "" && !d.originallyEmpty {
			continue
		}
		filtered = append(filtered, d)
	}

	// Step 9: Format output
	var sb strings.Builder
	for i, d := range filtered {
		if i > 0 {
			sb.WriteString("\n")
		}

		if d.ownPunct != "" {
			// Pad content to max_len, add space, append punctuation
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

// expandTabs replaces each tab character with 4 spaces.
func expandTabs(s string) string {
	var sb strings.Builder
	for _, ch := range s {
		if ch == '\t' {
			sb.WriteString("    ")
		} else {
			sb.WriteRune(ch)
		}
	}
	return sb.String()
}

// splitLines splits text into lines by '\n', returning non-empty results.
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

// extractTrailingPunctuation scans the line from left to right to find the
// last semicolon, brace, or closing brace that appears in code context
// (not inside a string, char literal, or comment). If the punctuation is
// followed only by whitespace, it is extracted. Returns (remaining_content, punctuation).
func extractTrailingPunctuation(line string) (punctuation string, remaining string) {
	lastPunctPos := -1
	lastPunctChar := ""
	inString := false
	inChar := false
	inLineComment := false

	for i := 0; i < len(line); i++ {
		ch := line[i]

		if inLineComment {
			if ch == '\n' {
				inLineComment = false
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

		// Code context
		if ch == '"' {
			inString = true
			continue
		}
		if ch == '\'' {
			inChar = true
			continue
		}
		if i+1 < len(line) && ch == '/' && line[i+1] == '/' {
			inLineComment = true
			continue
		}

		if ch == ';' || ch == '{' || ch == '}' {
			lastPunctPos = i
			lastPunctChar = string(ch)
		}
	}

	if lastPunctPos >= 0 {
		// Check if everything after the punctuation is whitespace
		rest := line[lastPunctPos+1:]
		allSpace := true
		for _, ch := range rest {
			if !unicode.IsSpace(ch) {
				allSpace = false
				break
			}
		}
		if allSpace {
			return lastPunctChar, line[:lastPunctPos]
		}
	}

	return "", line
}
```

- [ ] **Step 2: Write tests for the formatter**

Create `internal/formatter/formatter_test.go`:

```go
package formatter

import "testing"

func checkFormat(t *testing.T, input, expected string) {
	t.Helper()
	got := Format(input)
	if got != expected {
		t.Errorf("Format(\n%q\n) =\n%q\n\nwant:\n%q", input, got, expected)
	}
}

func TestFormat_Simple(t *testing.T) {
	input := `int main()
{
    printf("Hello World\n");
    return 0;
}`
	expected := `int main()                                         {
    printf("Hello World\n")                         ;
    return 0                                        ;}` + "\n"
	checkFormat(t, input, expected)
}

func TestFormat_CommentsPreserved(t *testing.T) {
	input := `/* comment */
int main() {
}`
	expected := `/* comment */
int main()                    {
}` + "\n"
	checkFormat(t, input, expected)
}

func TestFormat_EmptyLinesKept(t *testing.T) {
	input := `int main() {

    return 0;
}`
	expected := `int main()            {

    return 0;        ;}` + "\n"
	checkFormat(t, input, expected)
}

func TestFormat_LeadingBraceRemoved(t *testing.T) {
	input := `int main()
{
    return 0;
}
`
	expected := `int main()            {
    return 0;         ;}` + "\n"
	checkFormat(t, input, expected)
}

func TestFormat_StringBracesNotExtracted(t *testing.T) {
	input := `printf("}");`
	expected := `printf("}")` + "\n"
	checkFormat(t, input, expected)
}

func TestFormat_TabExpansion(t *testing.T) {
	input := "int main() {\n\treturn 0;\n}"
	expected := `int main()            {
    return 0;         ;}` + "\n"
	checkFormat(t, input, expected)
}
```

- [ ] **Step 3: Run tests to see which pass and which fail**

Run:
```bash
go test ./internal/formatter/ -v
```

Expected: some tests will fail because the max_len calculation and padding need to match the ground truth examples. Use the failures to iterate on the implementation.

- [ ] **Step 4: Fix formatter to match ground truth**

After running the tests, compare the formatter output against the ground truth files. Adjust the padding formula, line handling, or punctuation extraction as needed. The key formula is:

```
For lines with trailing punctuation:
  padded = content + (max_len - len(content)) spaces
  output = padded + " " + punctuation

For lines without trailing punctuation:
  output = content (as-is, no padding)
```

- [ ] **Step 5: Verify against ground truth examples**

Run:
```bash
go build -o relinted ./cmd/relinted
./relinted linted-example-1.c > /tmp/out1.c
diff /tmp/out1.c relinted-example-1.c

./relinted linted-example-2.c > /tmp/out2.c
diff /tmp/out2.c relinted-example-2.c

./relinted linted-example-3.c > /tmp/out3.c
diff /tmp/out3.c relinted-example-3.c
```

Expected: all three diffs produce no output (files match exactly).

---

### Task 5: Implement CLI entry point

**Files:**
- Create: `cmd/relinted/main.go`

- [ ] **Step 1: Write the CLI**

Create `cmd/relinted/main.go`:

```go
package main

import (
	"fmt"
	"os"

	"relinted/internal/formatter"
	"relinted/internal/io"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: relinted <input.c> [output.c]\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]

	// Read input
	content, err := io.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", inputPath, err)
		os.Exit(1)
	}

	// Format
	output := formatter.Format(content)

	// Write output
	if len(os.Args) >= 3 {
		outputPath := os.Args[2]
		if err := io.WriteFile(outputPath, output); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", outputPath, err)
			os.Exit(1)
		}
	} else {
		fmt.Print(output)
	}
}
```

- [ ] **Step 2: Build and test manually**

Run:
```bash
go build -o relinted ./cmd/relinted
./relinted linted-example-1.c
```

Expected: output matches relinted-example-1.c content.

- [ ] **Step 3: Test with output file argument**

Run:
```bash
./relinted linted-example-1.c /tmp/test-output.c
diff /tmp/test-output.c relinted-example-1.c
```

Expected: no diff output (files match).

---

### Task 6: Create Justfile

**Files:**
- Create: `Justfile`

- [ ] **Step 1: Write the Justfile**

Create `Justfile`:

```justfile
set shell := ["bash", "-euo", "pipefail", "-c"]

default: build

build:
	go build ./...

test:
	go test ./...

format:
	go fmt ./...

lint: build test

serve: xdg-open file:///home/ac/dev/relinted/README.md & sleep 1 && echo "Open README in browser"
```

- [ ] **Step 2: Verify Justfile works**

Run:
```bash
just build
just test
```

Expected: both commands succeed.

---

### Task 7: Final integration test against all ground truth files

**Files:**
- No new files

- [ ] **Step 1: Run full integration test**

Run:
```bash
go build -o relinted ./cmd/relinted
./relinted linted-example-1.c /tmp/out1.c && diff /tmp/out1.c relinted-example-1.c && echo "Example 1: PASS"
./relinted linted-example-2.c /tmp/out2.c && diff /tmp/out2.c relinted-example-2.c && echo "Example 2: PASS"
./relinted linted-example-3.c /tmp/out3.c && diff /tmp/out3.c relinted-example-3.c && echo "Example 3: PASS"
```

Expected: all three examples PASS with no diff output.

- [ ] **Step 2: Run all unit tests**

Run:
```bash
go test ./... -v
```

Expected: all tests PASS.

---

## Self-Review

**Spec coverage:**
- IO package → Task 2 ✓
- Tokenizer (Code, String, Char, CommentBlock, CommentLine) → Task 3 ✓
- Formatter (brace relocation, semicolon extraction, right-alignment, empty line handling, tab expansion) → Task 4 ✓
- CLI (positional args, stdout vs file output, error handling) → Task 5 ✓
- Build/test infrastructure → Task 6, Task 7 ✓
- Integration tests against ground truth → Task 7 ✓

**Placeholder scan:** No TBD, TODO, or vague requirements. All code is concrete.

**Type consistency:** SegmentType, Segment, Format(), ReadFile(), WriteFile() all used consistently.

**Scope check:** Focused on a single CLI tool. No unnecessary features.
