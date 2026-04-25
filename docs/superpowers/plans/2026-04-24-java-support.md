# Java Support Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add Java language support to relinted, matching the pattern established by Go, JS, Rust, and Perl.

**Architecture:** Create `internal/java/` package based on Go's implementation (closest C-family sibling). Java-specific additions: triple-quoted text block (`"""..."""`) handling in both tokenizer and punctuation extractor. Format pipeline is identical to Go/C.

**Tech Stack:** Go 1.22+, same tooling as rest of project.

---

## Task 1: Create Java example file

**Files:**
- Create: `examples/linted-example-8.java`

- [ ] **Step 1: Create linted-example-8.java**

Write a representative Java source file (~45 lines) covering: package declaration, imports, Javadoc comment, class with fields, methods, if/else, string literals, char literals, line comments, block comments, and Java text blocks (`"""..."""`).

Create `examples/linted-example-8.java` with this exact content:

```java
package com.example;

import java.util.Scanner;

/**
 * Simple calculator demo.
 * Demonstrates Java formatting.
 */
public class Calculator {
    private int result = 0;
    private String label = "Calculator";

    public void add(int n) {
        result += n;
        /* Update result */
    }

    public void subtract(int n) {
        result -= n;
    }

    public int getResult() {
        return result;
    }

    public char getSeparator() {
        char sep = '|';
        return sep;
    }

    public static void main(String[] args) {
        Calculator calc = new Calculator();
        calc.add(10);
        calc.add(5);

        String help = """
            Usage: add <n>, subtract <n>, result
            """;

        if (calc.getResult() > 0) {
            System.out.println("Positive!");
        } else {
            System.out.println("Zero or negative");
        }
    }
}
```

- [ ] **Step 2: Commit**

```bash
git add examples/linted-example-8.java
git commit -m "feat: add Java example file (linted-example-8.java)"
```

---

## Task 2: Create internal/java/ tokenize.go

**Files:**
- Create: `internal/java/tokenize.go`

- [ ] **Step 1: Create tokenize.go**

Create `internal/java/tokenize.go` with Go's tokenizer as the base, adding triple-quoted text block (`"""..."""`) handling. The text block is treated as an opaque `String` token to prevent brace/semicolon detection inside text block content.

```go
package java

import "relinted/internal/tokenizer"

// Tokenize scans Java source code character by character and produces a list of segments.
// Java-specific rules:
//   - // starts a line comment (scan to newline)
//   - /* starts a block comment (scan to */)
//   - " starts a double-quoted string (scan to closing ", handling \" escapes)
//   - ' starts a char literal (scan to closing ', handling \' escapes)
//   - """ starts a text block literal (scan to closing """, NO escape processing)
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

		case input[i] == '"' && i+2 < len(input) && input[i+1] == '"' && input[i+2] == '"':
			// Text block literal: scan to closing """, NO escape processing
			j := i + 3
			for j+2 < len(input) && !(input[j] == '"' && input[j+1] == '"' && input[j+2] == '"') {
				j++
			}
			if j+2 < len(input) {
				j += 3
			}
			segments = append(segments, tokenizer.Segment{tokenizer.String, input[i:j]})
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
			// Char literal: scan to closing ', handling \' escapes
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

Key differences from Go's tokenize.go:
- Added `"""` text block case before the regular `"` case (line 32-42)
- Removed backtick (`` ` ``) from the default case's break condition (line 89) — Java has no backtick strings

- [ ] **Step 2: Commit**

```bash
git add internal/java/tokenize.go
git commit -m "feat(java): add Java tokenizer with text block support"
```

---

## Task 3: Create internal/java/ punctuation.go

**Files:**
- Create: `internal/java/punctuation.go`

- [ ] **Step 1: Create punctuation.go**

Create `internal/java/punctuation.go` based on Go's punctuation.go with two changes:
1. Add `inTextBlock` state tracking (persisted across lines by format.go)
2. Text block end detection when `"""` is encountered while `inTextBlock` is true

```go
package java

import (
	"strings"
	"unicode"

	"relinted/internal/tokenizer"
)

// extractTrailingPunctuation scans the line to find the last punctuation character
// (semicolon, opening brace, or closing brace) that appears in code context
// (not inside a string, char literal, text block, block comment, or line comment).
//
// The inTextBlock parameter tracks whether we are inside a multi-line text block
// that started on a previous line. Returns (punctuation, remaining content, next inTextBlock state).
func extractTrailingPunctuation(line string, inTextBlock bool) (punctuation string, remaining string, nextInTextBlock bool) {
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

		if inTextBlock {
			// Check for text block end
			if ch == '"' && i+2 < len(line) && line[i+1] == '"' && line[i+2] == '"' {
				inTextBlock = false
				i += 2
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

		if inBlockComment {
			if ch == '*' && i+1 < len(line) && line[i+1] == '/' {
				inBlockComment = false
				i++
			}
			continue
		}

		// Code context
		if ch == '"' && i+2 < len(line) && line[i+1] == '"' && line[i+2] == '"' {
			inTextBlock = true
			i += 2
			continue
		}
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
			return string(line[lastPunctPos]), line[:lastPunctPos], inTextBlock
		}
	}

	return "", line, inTextBlock
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

Key differences from Go's punctuation.go:
- Function signature changed: `extractTrailingPunctuation(line string, inTextBlock bool) (punctuation, remaining string, nextInTextBlock bool)` — tracks text block state across lines
- Added `inTextBlock` handling before `inString` (lines 27-34) — when inside a text block, look for closing `"""`
- Added text block start detection in code context (lines 63-66) — checks for `"""` before `"`
- Removed `inRawString` / backtick handling (Java has no backtick strings)
- Function returns `inTextBlock` state for cross-line persistence

- [ ] **Step 2: Commit**

```bash
git add internal/java/punctuation.go
git commit -m "feat(java): add Java punctuation extractor with text block cross-line tracking"
```

---

## Task 4: Create internal/java/ format.go

**Files:**
- Create: `internal/java/format.go`

- [ ] **Step 1: Create format.go**

Create `internal/java/format.go` based on Go's format.go. The format pipeline is identical to Go/C. The only difference is the call to `extractTrailingPunctuation` now passes and receives the `inTextBlock` state.

```go
package java

import (
	"strings"
	"unicode"
)

// Format takes Java source code and reformats it by relocating braces and
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
	inTextBlock := false
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
		punct, rest, inTextBlock := extractTrailingPunctuation(trimmed, inTextBlock)
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

Key difference from Go's format.go:
- Line 57: `inTextBlock := false` — declares the cross-line state variable
- Line 77: `punct, rest, inTextBlock := extractTrailingPunctuation(trimmed, inTextBlock)` — passes and receives text block state

- [ ] **Step 2: Commit**

```bash
git add internal/java/format.go
git commit -m "feat(java): add Java formatter (8-step pipeline)"
```

---

## Task 5: Update main.go with Java support

**Files:**
- Modify: `cmd/relinted/main.go`

- [ ] **Step 1: Update extToLang mapping**

Add `.java` → `java` to the `extToLang` map. Find the existing map in `cmd/relinted/main.go` and add:

```go
var extToLang = map[string]string{
	".c":    "c",
	".h":    "c",
	".cpp":  "c",
	".cc":   "c",
	".pl":   "perl",
	".pm":   "perl",
	".js":   "js",
	".rs":   "rust",
	".go":   "go",
	".java": "java",
}
```

- [ ] **Step 2: Add import**

Add `"relinted/internal/java"` to the import block:

```go
import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"relinted/internal/formatter"
	"relinted/internal/go"
	"relinted/internal/io"
	"relinted/internal/java"
	"relinted/internal/js"
	"relinted/internal/perl"
	"relinted/internal/rust"
)
```

- [ ] **Step 3: Add switch case**

Add the `java` case to the switch statement in `main()`, right after the `go` case:

```go
	case "go":
		output = go_pkg.Format(content)
	case "java":
		output = java.Format(content)
	default:
```

- [ ] **Step 4: Update help text**

Update the `-l` flag help text and the banner to include Java:

Change line 40:
```go
flag.StringVar(&langFlag, "l", "", "Language to use (overrides extension detection) [c, perl, rust, js, go, java]")
```

Change line 42:
```go
fmt.Println("Relinted reformats C/C++, Perl, Rust, JavaScript, Go, and Java source code to visually resemble Python.")
```

Change line 80:
```go
fmt.Fprintf(os.Stderr, "Supported languages: c, perl, rust, js, go, java\n")
```

- [ ] **Step 5: Commit**

```bash
git add cmd/relinted/main.go
git commit -m "feat(cli): add Java language support to main"
```

---

## Task 6: Build relinted and generate relinted-example-8.java

**Files:**
- Create: `examples/relinted-example-8.java`

- [ ] **Step 1: Build relinted**

```bash
just build
```

- [ ] **Step 2: Generate ground truth output**

```bash
./relinted examples/linted-example-8.java examples/relinted-example-8.java
```

- [ ] **Step 3: Verify output looks correct**

```bash
cat examples/relinted-example-8.java
```

Check that:
- Braces are relocated to the end of the previous line
- Semicolons are right-aligned
- Text block content is preserved (braces/semicolons inside `"""..."""` are NOT extracted)
- Block comment content is preserved
- The output is visually Python-like

- [ ] **Step 4: Commit**

```bash
git add examples/relinted-example-8.java
git commit -m "feat: add Java ground truth file (relinted-example-8.java)"
```

---

## Task 7: Create format_test.go

**Files:**
- Create: `internal/java/format_test.go`

- [ ] **Step 1: Create format_test.go**

```go
package java

import (
	"os"
	"testing"
)

func TestFormat(t *testing.T) {
	input, err := os.ReadFile("../../examples/linted-example-8.java")
	if err != nil {
		t.Fatalf("Failed to read input: %v", err)
	}
	want, err := os.ReadFile("../../examples/relinted-example-8.java")
	if err != nil {
		t.Fatalf("Failed to read want: %v", err)
	}
	got := Format(string(input))
	if got != string(want) {
		t.Errorf("Format mismatch\ngot:\n%s\nwant:\n%s", got, want)
	}
}
```

- [ ] **Step 2: Run test to verify it passes**

```bash
just test-java
```

Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add internal/java/format_test.go
git commit -m "test(java): add format test comparing against ground truth"
```

---

## Task 8: Update README.md

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update banner**

Change line 3:
```
Relinted reformats C/C++, Perl, Rust, JavaScript, Go, and Java source code to visually resemble Python by aligning braces (`{`, `}`) and semicolons (`;`) to the far right of the codebase.
```

- [ ] **Step 2: Update Multi-Language bullet**

Change the "Multi-Language" bullet (around line 35):
```
- **Multi-Language**: Supports C/C++, Perl, Rust, JavaScript, Go, and Java with auto-detection from file extension
```

- [ ] **Step 3: Update -l flag argument table**

Change line 40:
```
| `-l` | *(Optional)* Language to use: `c`, `perl`, `rust`, `js`, `go`, or `java`. Overrides extension detection. |
```

- [ ] **Step 4: Add Java to Language Detection table**

Add a row to the Language Detection table (after the `.go` row):
```
| `.java` | Java |
```

- [ ] **Step 5: Add Java usage example**

Add to Usage Examples (after `./relinted program.go`):
```bash
./relinted program.java
./relinted -l java source.java
```

- [ ] **Step 6: Add Java before/after example**

Add after the Go example (around line 160):

```markdown
**Before (Java)**
```java
package com.example;

import java.util.Scanner;

/**
 * Simple calculator demo.
 * Demonstrates Java formatting.
 */
public class Calculator {
    private int result = 0;
    private String label = "Calculator";

    public void add(int n) {
        result += n;
        /* Update result */
    }

    public void subtract(int n) {
        result -= n;
    }

    public int getResult() {
        return result;
    }

    public char getSeparator() {
        char sep = '|';
        return sep;
    }

    public static void main(String[] args) {
        Calculator calc = new Calculator();
        calc.add(10);
        calc.add(5);

        String help = """
            Usage: add <n>, subtract <n>, result
            """;

        if (calc.getResult() > 0) {
            System.out.println("Positive!");
        } else {
            System.out.println("Zero or negative");
        }
    }
}
```

**After (Python-like Visual)**
```java
package com.example;
... (actual relinted output)
```

Use the actual output from `examples/relinted-example-8.java` for the "After" block.

- [ ] **Step 7: Add Java to example file table**

Add a row to the example file table (after the Go row):
```
| examples/linted-example-8.java | examples/relinted-example-8.java |
```

- [ ] **Step 8: Update "Tests" section description**

Change the Tests section description to mention Java:
```
The repository contains read-only source code test example files and their read-only ground truth counterparts; running Relinted on each linted-example file should produce output identical to the matching ground truth relinted-example file:
```
(No change needed — this is already generic.)

- [ ] **Step 9: Update "Notes" section**

Change the last bullet (around line 180):
```
- Supports C/C++, Perl, Rust, JavaScript, Go, and Java; other languages can be added via different parsers
```

- [ ] **Step 10: Commit**

```bash
git add README.md
git commit -m "docs: add Java references to README.md"
```

---

## Task 9: Update Justfile

**Files:**
- Modify: `Justfile`

- [ ] **Step 1: Add test-java target**

Add after the `test-go` target:

```
test-java:
    go test ./internal/java/...
```

- [ ] **Step 2: Commit**

```bash
git add Justfile
git commit -m "ci: add test-java target to Justfile"
```

---

## Task 10: Run all tests and verify

**Files:**
- Verify: all tests pass

- [ ] **Step 1: Run full test suite**

```bash
just test
```

Expected: All tests pass

- [ ] **Step 2: Run Java-specific test**

```bash
just test-java
```

Expected: PASS

- [ ] **Step 3: Run linter**

```bash
just lint
```

Expected: No warnings

- [ ] **Step 4: Run Go formatter**

```bash
just format
```

- [ ] **Step 5: Verify relinted binary works**

```bash
./relinted -l java examples/linted-example-8.java
```

Expected: Output matches `examples/relinted-example-8.java`

- [ ] **Step 6: Final commit**

```bash
git add -A
git status
```

If no untracked files remain, no commit needed. If there are changes (e.g., from `just format`), commit them.

---

## Self-Review

**1. Spec coverage:**
- `internal/java/` package with 3 source files + 1 test file ✓ (Tasks 2-4, 7)
- `tokenize.go` with `"""` text block handling ✓ (Task 2)
- `punctuation.go` with `"""` and cross-line state tracking ✓ (Task 3)
- `format.go` identical to Go's ✓ (Task 4)
- Example file `linted-example-8.java` + `relinted-example-8.java` ✓ (Tasks 1, 6)
- `main.go` changes: `.java` mapping, switch case, import ✓ (Task 5)
- `README.md` updates: banner, language table, example table, usage examples ✓ (Task 8)
- `Justfile` `test-java` target ✓ (Task 9)

**2. Placeholder scan:** No TBD, TODO, or vague requirements. All code is complete. All file paths are exact. All commands are specified.

**3. Type consistency:** Package name `java` used consistently. Function signature `extractTrailingPunctuation(line string, inTextBlock bool) (punctuation, remaining string, nextInTextBlock bool)` matches across punctuation.go and format.go.

**4. Ambiguity check:** Text block handling is explicit — `"""` is treated as an opaque string token in the tokenizer, and the punctuation extractor tracks `inTextBlock` state across lines to correctly handle multi-line text blocks.

**5. Edge cases:**
- Text block starting and ending on the same line: handled (both opening and closing `"""` detected in one pass)
- Text block spanning multiple lines: handled via cross-line `inTextBlock` state
- Text block content containing `""` (two quotes, not three): handled (only `"""` triggers text block state)
- Block comments `/* ... */`: handled (same as C/Go)
- Char literals with escaped characters: handled (same as C/Go)
