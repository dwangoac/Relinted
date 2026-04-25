# TypeScript, C#, PHP, Swift Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add TypeScript, C#, PHP, and Swift language support to relinted.

**Architecture:** TypeScript and C# reuse existing JS and Java tokenizers respectively (just extension mappings). PHP and Swift get new `internal/php/` and `internal/swift/` packages with language-specific tokenizers for heredocs, optionals, and other unique features. All languages share the same 8-step formatting pipeline.

**Tech Stack:** Go, relinted internal packages (tokenizer, formatter, language-specific packages).

---

### Task 1: TypeScript — Update main.go

**Files:**
- Modify: `cmd/relinted/main.go`

- [ ] **Step 1: Add `.ts` to `extToLang` map**

In `cmd/relinted/main.go`, add `.ts` → `"js"` mapping after the `.java` line:

```go
var extToLang = map[string]string{
	".c":    "c",
	".h":    "c",
	".cpp":  "c",
	".cc":   "c",
	".pl":   "perl",
	".pm":   "perl",
	".js":   "js",
	".ts":   "js",
	".rs":   "rust",
	".go":   "go",
	".java": "java",
}
```

- [ ] **Step 2: Add `.ts` to help text**

Change the `-l` flag help text:

From:
```go
flag.StringVar(&langFlag, "l", "", "Language to use (overrides extension detection) [c, perl, rust, js, go, java]")
```

To:
```go
flag.StringVar(&langFlag, "l", "", "Language to use (overrides extension detection) [c, perl, rust, js, ts, go, java]")
```

- [ ] **Step 3: Update Usage message**

Change:
```go
fmt.Println("Relinted reformats C/C++, Perl, Rust, JavaScript, Go, and Java source code to visually resemble Python.")
```

To:
```go
fmt.Println("Relinted reformats C/C++, Perl, Rust, JavaScript, TypeScript, Go, and Java source code to visually resemble Python.")
```

- [ ] **Step 4: Commit**

```bash
git add cmd/relinted/main.go
git commit -m "feat(cli): add TypeScript extension mapping"
```

### Task 2: TypeScript — Example files and test

**Files:**
- Create: `examples/linted-example-9.ts`
- Create: `examples/relinted-example-9.ts`
- Create: `internal/js/format_test.go` (add integration test)

- [ ] **Step 1: Create TypeScript example input**

Create `examples/linted-example-9.ts`:

```typescript
class Greeter {
    private greeting: string;

    constructor(message: string) {
        this.greeting = message;
    }

    greet(): string {
        return "Hello, " + this.greeting;
    }
}

function main() {
    const greeter = new Greeter("world");
    console.log(greeter.greet());

    const numbers: number[] = [1, 2, 3];
    for (const n of numbers) {
        console.log(n);
    }
}

main();
```

- [ ] **Step 2: Build relinted and generate ground truth**

```bash
go build -o relinted ./cmd/relinted/
./relinted examples/linted-example-9.ts examples/relinted-example-9.ts
```

- [ ] **Step 3: Add TypeScript integration test to JS format_test.go**

Append to `internal/js/format_test.go`:

```go
func TestFormat_TypeScript_Integration(t *testing.T) {
	input, err := os.ReadFile("../../examples/linted-example-9.ts")
	if err != nil {
		t.Fatalf("failed to read input: %v", err)
	}
	expected, err := os.ReadFile("../../examples/relinted-example-9.ts")
	if err != nil {
		t.Fatalf("failed to read expected: %v", err)
	}
	got := Format(string(input))
	if got != string(expected) {
		t.Errorf("TypeScript output mismatch\ngot:\n%s\nexpected:\n%s", got, expected)
	}
}
```

- [ ] **Step 4: Run tests**

```bash
go test ./internal/js/...
```

- [ ] **Step 5: Commit**

```bash
git add examples/linted-example-9.ts examples/relinted-example-9.ts internal/js/format_test.go
git commit -m "test(ts): add TypeScript integration test"
```

### Task 3: C# — Update main.go

**Files:**
- Modify: `cmd/relinted/main.go`

- [ ] **Step 1: Add `.cs` to `extToLang` map**

Add `.cs` → `"java"` mapping after `.java`:

```go
var extToLang = map[string]string{
	".c":    "c",
	".h":    "c",
	".cpp":  "c",
	".cc":   "c",
	".pl":   "perl",
	".pm":   "perl",
	".js":   "js",
	".ts":   "js",
	".rs":   "rust",
	".go":   "go",
	".java": "java",
	".cs":   "java",
}
```

- [ ] **Step 2: Add `.cs` to help text**

Change:
```go
flag.StringVar(&langFlag, "l", "", "Language to use (overrides extension detection) [c, perl, rust, js, ts, go, java]")
```

To:
```go
flag.StringVar(&langFlag, "l", "", "Language to use (overrides extension detection) [c, perl, rust, js, ts, go, java, cs]")
```

- [ ] **Step 3: Update Usage message**

Change:
```go
fmt.Println("Relinted reformats C/C++, Perl, Rust, JavaScript, TypeScript, Go, and Java source code to visually resemble Python.")
```

To:
```go
fmt.Println("Relinted reformats C/C++, Perl, Rust, JavaScript, TypeScript, Go, Java, and C# source code to visually resemble Python.")
```

- [ ] **Step 4: Commit**

```bash
git add cmd/relinted/main.go
git commit -m "feat(cli): add C# extension mapping"
```

### Task 4: C# — Example files and test

**Files:**
- Create: `examples/linted-example-10.cs`
- Create: `examples/relinted-example-10.cs`
- Create: `internal/java/format_test.go` (add integration test)

- [ ] **Step 1: Create C# example input**

Create `examples/linted-example-10.cs`:

```csharp
using System;
using System.Collections.Generic;

namespace Demo
{
    public class Program
    {
        public static void Main(string[] args)
        {
            var names = new List<string> { "Alice", "Bob", "Charlie" };

            foreach (var name in names)
            {
                Console.WriteLine("Hello, " + name);
            }

            int sum = 0;
            for (int i = 0; i < 5; i++)
            {
                sum += i;
            }

            Console.WriteLine("Sum: " + sum);
        }
    }
}
```

- [ ] **Step 2: Build relinted and generate ground truth**

```bash
go build -o relinted ./cmd/relinted/
./relinted examples/linted-example-10.cs examples/relinted-example-10.cs
```

- [ ] **Step 3: Add C# integration test to Java format_test.go**

Append to `internal/java/format_test.go`:

```go
func TestFormat_CSharp_Integration(t *testing.T) {
	input, err := os.ReadFile("../../examples/linted-example-10.cs")
	if err != nil {
		t.Fatalf("failed to read input: %v", err)
	}
	expected, err := os.ReadFile("../../examples/relinted-example-10.cs")
	if err != nil {
		t.Fatalf("failed to read expected: %v", err)
	}
	got := Format(string(input))
	if got != string(expected) {
		t.Errorf("C# output mismatch\ngot:\n%s\nexpected:\n%s", got, expected)
	}
}
```

- [ ] **Step 4: Run tests**

```bash
go test ./internal/java/...
```

- [ ] **Step 5: Commit**

```bash
git add examples/linted-example-10.cs examples/relinted-example-10.cs internal/java/format_test.go
git commit -m "test(cs): add C# integration test"
```

### Task 5: PHP — Create tokenizer

**Files:**
- Create: `internal/php/tokenize.go`

PHP tokenizer handles:
- Heredocs: `<<<'IDENT'\n...\nIDENT;` — treated as opaque `String` token, no escape processing
- Double-quoted strings: `"Hello $name"` — treated as opaque `String` token (no variable interpolation parsing)
- Single-quoted strings: `'literal'` — standard escape handling
- Backtick strings: `` `cmd` `` — treated as `String` token (no escape processing)
- PHP tags: `<?php`, `<?`, `<?=`, `?>` — treated as `Code` tokens
- Standard C-style comments: `//` and `/* */`
- Everything else is `Code`

Create `internal/php/tokenize.go`:

```go
package php

import "relinted/internal/tokenizer"

// Tokenize scans PHP source code character by character and produces a list of segments.
// PHP-specific rules:
//   - // starts a line comment (scan to newline)
//   - /* starts a block comment (scan to */)
//   - " starts a double-quoted string with variable interpolation (opaque token, no parsing)
//   - ' starts a single-quoted string (scan to closing ', handling \' escapes)
//   - ` starts a backtick string (opaque token, no escape processing)
//   - <<< starts a heredoc (scan to IDENTIFIER; on its own line, no escape processing)
//   - <?php, <?, <?=, ?> are PHP tags (treated as Code)
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
			// Double-quoted string with variable interpolation: opaque token, no escape processing
			j := i + 1
			for j < len(input) && input[j] != '"' {
				j++
			}
			if j < len(input) {
				j++
			}
			segments = append(segments, tokenizer.Segment{tokenizer.String, input[i:j]})
			i = j

		case input[i] == '\'':
			// Single-quoted string: scan to closing ', handling \' escapes
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
			// Backtick string: opaque token, no escape processing
			j := i + 1
			for j < len(input) && input[j] != '`' {
				j++
			}
			if j < len(input) {
				j++
			}
			segments = append(segments, tokenizer.Segment{tokenizer.String, input[i:j]})
			i = j

		case i+3 < len(input) && input[i] == '<' && input[i+1] == '?' && input[i+2] == '=':
			// Short echo tag: <?=
			j := i + 3
			for j < len(input) && !(j+1 < len(input) && input[j] == '?' && input[j+1] == '>') {
				j++
			}
			if j+1 < len(input) {
				j += 2
			}
			segments = append(segments, tokenizer.Segment{tokenizer.Code, input[i:j]})
			i = j

		case i+3 < len(input) && input[i] == '<' && input[i+1] == '?' && input[i+2] == 'p' && input[i+3] == 'h' && input[i+4] == 'p':
			// Full PHP open tag: <?php
			j := i + 5
			for j < len(input) && !(j+1 < len(input) && input[j] == '?' && input[j+1] == '>') {
				j++
			}
			if j+1 < len(input) {
				j += 2
			}
			segments = append(segments, tokenizer.Segment{tokenizer.Code, input[i:j]})
			i = j

		case i+1 < len(input) && input[i] == '<' && input[i+1] == '?':
			// Short open tag: <?
			j := i + 2
			for j < len(input) && !(j+1 < len(input) && input[j] == '?' && input[j+1] == '>') {
				j++
			}
			if j+1 < len(input) {
				j += 2
			}
			segments = append(segments, tokenizer.Segment{tokenizer.Code, input[i:j]})
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
				if ch == '<' && j+1 < len(input) && input[j+1] == '?' {
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
git add internal/php/tokenize.go
git commit -m "feat(php): add PHP tokenizer with heredoc and interpolation support"
```

### Task 6: PHP — Create punctuation.go, format.go, test, examples

**Files:**
- Create: `internal/php/punctuation.go`
- Create: `internal/php/format.go`
- Create: `examples/linted-example-11.php`
- Create: `examples/relinted-example-11.php`
- Create: `internal/php/format_test.go`

- [ ] **Step 1: Create PHP punctuation.go**

Same pattern as JS punctuation.go (PHP punctuation marks are `{`, `}`, `;` — no heredoc tracking needed since heredocs are already opaque tokens in the tokenizer).

Create `internal/php/punctuation.go`:

```go
package php

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
	inDoubleQuote := false
	inSingleQuote := false
	inBacktick := false
	inBlockComment := false
	inLineComment := false

	for i := 0; i < len(line); i++ {
		ch := line[i]

		if inLineComment {
			continue
		}

		if inDoubleQuote {
			if ch == '\\' && i+1 < len(line) {
				i++
				continue
			}
			if ch == '"' {
				inDoubleQuote = false
			}
			continue
		}

		if inSingleQuote {
			if ch == '\\' && i+1 < len(line) {
				i++
				continue
			}
			if ch == '\'' {
				inSingleQuote = false
			}
			continue
		}

		if inBacktick {
			if ch == '`' {
				inBacktick = false
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
			inDoubleQuote = true
			continue
		}
		if ch == '\'' {
			inSingleQuote = true
			continue
		}
		if ch == '`' {
			inBacktick = true
			continue
		}
		if i+1 < len(line) && ch == '/' && line[i+1] == '*' {
			inBlockComment = true
			i++
			continue
		}
		if i+1 < len(line) && ch == '/' && line[i+1] == '/' {
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
					break
				}
				if rest[i+1] == '*' {
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

- [ ] **Step 2: Create PHP format.go**

Identical to JS format.go (8-step pipeline, no PHP-specific logic).

Create `internal/php/format.go`:

```go
package php

import (
	"strings"
	"unicode"
)

// Format takes PHP source code and reformats it by relocating braces and
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

- [ ] **Step 3: Create PHP example input**

Create `examples/linted-example-11.php`:

```php
<?php

class Database {
    private $host;
    private $name;

    public function __construct($host, $name) {
        $this->host = $host;
        $this->name = $name;
    }

    public function connect() {
        $message = "Connecting to $this->name on $this->host";
        echo $message;
    }
}

function main() {
    $db = new Database('localhost', 'mydb');
    $db->connect();

    $items = ['apple', 'banana', 'cherry'];
    foreach ($items as $item) {
        echo "Item: $item\n";
    }
}

main();
?>
```

- [ ] **Step 4: Build relinted and generate ground truth**

```bash
go build -o relinted ./cmd/relinted/
./relinted examples/linted-example-11.php examples/relinted-example-11.php
```

- [ ] **Step 5: Create PHP format_test.go**

Create `internal/php/format_test.go`:

```go
package php

import (
	"os"
	"testing"
)

func TestFormat(t *testing.T) {
	input, err := os.ReadFile("../../examples/linted-example-11.php")
	if err != nil {
		t.Fatalf("Failed to read input: %v", err)
	}
	want, err := os.ReadFile("../../examples/relinted-example-11.php")
	if err != nil {
		t.Fatalf("Failed to read want: %v", err)
	}
	got := Format(string(input))
	if got != string(want) {
		t.Errorf("Format mismatch\ngot:\n%s\nwant:\n%s", got, want)
	}
}
```

- [ ] **Step 6: Run tests**

```bash
go test ./internal/php/...
```

- [ ] **Step 7: Commit**

```bash
git add internal/php/punctuation.go internal/php/format.go internal/php/format_test.go examples/linted-example-11.php examples/relinted-example-11.php
git commit -m "feat(php): add PHP formatter with punctuation extraction and format pipeline"
```

### Task 7: PHP — Update main.go

**Files:**
- Modify: `cmd/relinted/main.go`

- [ ] **Step 1: Add PHP import**

Add to imports:

```go
	"relinted/internal/php"
```

- [ ] **Step 2: Add `.php` to `extToLang` map**

```go
	".php": "php",
```

- [ ] **Step 3: Add PHP switch case**

```go
	case "php":
		output = php.Format(content)
```

- [ ] **Step 4: Update help text**

Change `-l` flag help:
```go
flag.StringVar(&langFlag, "l", "", "Language to use (overrides extension detection) [c, perl, rust, js, ts, go, java, cs, php]")
```

- [ ] **Step 5: Update Usage message**

Change:
```go
fmt.Println("Relinted reformats C/C++, Perl, Rust, JavaScript, TypeScript, Go, Java, and C# source code to visually resemble Python.")
```

To:
```go
fmt.Println("Relinted reformats C/C++, Perl, Rust, JavaScript, TypeScript, Go, Java, C#, and PHP source code to visually resemble Python.")
```

- [ ] **Step 6: Update default error message**

Change:
```go
fmt.Fprintf(os.Stderr, "Supported languages: c, perl, rust, js, go, java\n")
```

To:
```go
fmt.Fprintf(os.Stderr, "Supported languages: c, perl, rust, js, ts, go, java, cs, php\n")
```

- [ ] **Step 7: Commit**

```bash
git add cmd/relinted/main.go
git commit -m "feat(cli): add PHP language support to main"
```

### Task 8: Swift — Create tokenizer

**Files:**
- Create: `internal/swift/tokenize.go`

Swift tokenizer handles:
- Optional chaining: `?` after identifiers — treated as opaque `Code` token (not punctuation)
- Nil coalescing: `??` — treated as opaque `Code` token
- `@` attributes: `@objc`, `@escaping` — `@` followed by identifier is opaque `Code`
- String interpolation: `\(...)` — the `\` starts a sub-expression, `(...)` is skipped, the content inside is NOT parsed as code
- Standard C-style comments: `//` and `/* */`
- Standard strings: `"..."` with escape handling (no interpolation parsing, treated as opaque `String`)
- Standard char literals: `'...'` (single-char)
- Everything else is `Code`

Create `internal/swift/tokenize.go`:

```go
package swift

import "relinted/internal/tokenizer"

// Tokenize scans Swift source code character by character and produces a list of segments.
// Swift-specific rules:
//   - // starts a line comment (scan to newline)
//   - /* starts a block comment (scan to */)
//   - " starts a string literal (scan to closing ", handling \" escapes, opaque — no interpolation parsing)
//   - ' starts a char literal (scan to closing ', handling \' escapes)
//   - ? after identifier content is part of the Code token (optional chaining, not punctuation)
//   - ?? is part of Code token (nil coalescing, not punctuation)
//   - @ starts an attribute — @ followed by identifier chars is part of Code token
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
			// String literal: scan to closing ", handling \" escapes (opaque — no interpolation parsing)
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
			// Swift ? and @ are NOT punctuation — they stay in the code token
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

- [ ] **Step 2: Commit**

```bash
git add internal/swift/tokenize.go
git commit -m "feat(swift): add Swift tokenizer with optional and attribute support"
```

### Task 9: Swift — Create punctuation.go, format.go, test, examples

**Files:**
- Create: `internal/swift/punctuation.go`
- Create: `internal/swift/format.go`
- Create: `examples/linted-example-12.swift`
- Create: `examples/relinted-example-12.swift`
- Create: `internal/swift/format_test.go`

- [ ] **Step 1: Create Swift punctuation.go**

Same pattern as JS punctuation.go. Swift `?` and `@` are NOT punctuation marks — they stay in code tokens. Only `{`, `}`, `;` are punctuation.

Create `internal/swift/punctuation.go`:

```go
package swift

import (
	"strings"
	"unicode"

	"relinted/internal/tokenizer"
)

// extractTrailingPunctuation scans the line to find the last punctuation character
// (semicolon, opening brace, or closing brace) that appears in code context
// (not inside a string, char literal, block comment, or line comment).
//
// Swift ? and @ are NOT punctuation — they stay in code tokens.
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
		if i+1 < len(line) && ch == '/' && line[i+1] == '*' {
			inBlockComment = true
			i++
			continue
		}
		if i+1 < len(line) && ch == '/' && line[i+1] == '/' {
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
					break
				}
				if rest[i+1] == '*' {
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

- [ ] **Step 2: Create Swift format.go**

Identical to JS format.go (8-step pipeline).

Create `internal/swift/format.go`:

```go
package swift

import (
	"strings"
	"unicode"
)

// Format takes Swift source code and reformats it by relocating braces and
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

- [ ] **Step 3: Create Swift example input**

Create `examples/linted-example-12.swift`:

```swift
import Foundation

class User {
    var name: String?
    var age: Int?

    init(name: String, age: Int?) {
        self.name = name
        self.age = age
    }

    func greet() -> String {
        let displayName = name ?? "Anonymous"
        return "Hello, \(displayName)!"
    }
}

func main() {
    let users = [User(name: "Alice", age: 30), User(name: "Bob", age: nil)]

    for user in users {
        let message = user.greet()
        print(message)
    }

    let optionalValue: Int? = nil
    let result = optionalValue ?? 42
    print("Result: \(result)")
}

main()
```

- [ ] **Step 4: Build relinted and generate ground truth**

```bash
go build -o relinted ./cmd/relinted/
./relinted examples/linted-example-12.swift examples/relinted-example-12.swift
```

- [ ] **Step 5: Create Swift format_test.go**

Create `internal/swift/format_test.go`:

```go
package swift

import (
	"os"
	"testing"
)

func TestFormat(t *testing.T) {
	input, err := os.ReadFile("../../examples/linted-example-12.swift")
	if err != nil {
		t.Fatalf("Failed to read input: %v", err)
	}
	want, err := os.ReadFile("../../examples/relinted-example-12.swift")
	if err != nil {
		t.Fatalf("Failed to read want: %v", err)
	}
	got := Format(string(input))
	if got != string(want) {
		t.Errorf("Format mismatch\ngot:\n%s\nwant:\n%s", got, want)
	}
}
```

- [ ] **Step 6: Run tests**

```bash
go test ./internal/swift/...
```

- [ ] **Step 7: Commit**

```bash
git add internal/swift/punctuation.go internal/swift/format.go internal/swift/format_test.go examples/linted-example-12.swift examples/relinted-example-12.swift
git commit -m "feat(swift): add Swift formatter with optional and attribute support"
```

### Task 10: Swift — Update main.go

**Files:**
- Modify: `cmd/relinted/main.go`

- [ ] **Step 1: Add Swift import**

Add to imports:

```go
	"relinted/internal/swift"
```

- [ ] **Step 2: Add `.swift` to `extToLang` map**

```go
	".swift": "swift",
```

- [ ] **Step 3: Add Swift switch case**

```go
	case "swift":
		output = swift.Format(content)
```

- [ ] **Step 4: Update help text**

Change `-l` flag help:
```go
flag.StringVar(&langFlag, "l", "", "Language to use (overrides extension detection) [c, perl, rust, js, ts, go, java, cs, php, swift]")
```

- [ ] **Step 5: Update Usage message**

Change:
```go
fmt.Println("Relinted reformats C/C++, Perl, Rust, JavaScript, TypeScript, Go, Java, C#, and PHP source code to visually resemble Python.")
```

To:
```go
fmt.Println("Relinted reformats C/C++, Perl, Rust, JavaScript, TypeScript, Go, Java, C#, PHP, and Swift source code to visually resemble Python.")
```

- [ ] **Step 6: Update default error message**

Change:
```go
fmt.Fprintf(os.Stderr, "Supported languages: c, perl, rust, js, ts, go, java, cs, php\n")
```

To:
```go
fmt.Fprintf(os.Stderr, "Supported languages: c, perl, rust, js, ts, go, java, cs, php, swift\n")
```

- [ ] **Step 7: Commit**

```bash
git add cmd/relinted/main.go
git commit -m "feat(cli): add Swift language support to main"
```

### Task 11: Update documentation

**Files:**
- Modify: `README.md`
- Modify: `Justfile`

- [ ] **Step 1: Update README.md banner**

Change the banner to include TypeScript, C#, PHP, Swift.

- [ ] **Step 2: Update README.md detection table**

Add rows for `.ts`, `.cs`, `.php`, `.swift`.

- [ ] **Step 3: Update README.md usage examples**

Add usage examples for all 4 languages.

- [ ] **Step 4: Update README.md before/after examples**

Add TypeScript and C# before/after examples (reusing JS/Java patterns with TS/C# syntax).
Add PHP before/after example with heredoc and interpolation.
Add Swift before/after example with optionals and closures.

- [ ] **Step 5: Update README.md example file table**

Add rows for examples 9-12.

- [ ] **Step 6: Update README.md notes section**

- [ ] **Step 7: Add test targets to Justfile**

Add `test-php` and `test-swift` targets.

- [ ] **Step 8: Run all tests**

```bash
go test ./internal/... ./cmd/...
```

- [ ] **Step 9: Commit**

```bash
git add README.md Justfile
git commit -m "docs: add TypeScript, C#, PHP, Swift references to documentation"
```

---

## Self-Review

### Spec coverage
- TypeScript: Tasks 1-2 ✓ (ext mapping + examples + test)
- C#: Tasks 3-4 ✓ (ext mapping + examples + test)
- PHP: Tasks 5-7 ✓ (tokenizer + punctuation/format/examples/test + main.go)
- Swift: Tasks 8-10 ✓ (tokenizer + punctuation/format/examples/test + main.go)
- Docs: Task 11 ✓ (README + Justfile)

### Placeholder scan
- All code blocks contain complete, copyable code
- No "TBD", "TODO", "similar to" references
- All file paths are exact
- All commands with expected behavior

### Type consistency
- `Tokenize` returns `[]tokenizer.Segment` in all packages ✓
- `Format` returns `string` in all packages ✓
- `extractTrailingPunctuation` returns `(string, string)` in PHP/Swift (same as JS) ✓
- `extractTrailingPunctuation` returns `(string, string, bool)` in Java (with text block tracking) — PHP/Swift don't need this ✓

### Scope check
- Focused on 4 languages only ✓
- No over-engineering: PHP/Swift use opaque tokens for complex strings ✓
- TypeScript and C# reuse existing code (no new packages) ✓
- Example files are representative but simple ✓
