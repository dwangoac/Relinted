package rust

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Format takes Rust source code and reformats it by relocating braces ({, })
// and semicolons (;) and commas (in match arm context) to the far right,
// creating a Python-like visual structure.
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

	// Step 3: Tokenize entire input and reconstruct to normalize
	// (re-tokenizing original input after tab expansion and trailing whitespace
	// stripping ensures string/char literals containing braces or punctuation
	// are handled correctly before line-level processing)
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
		trailStuff      string
		queuedBraces    []string
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
			// Find the position of the first non-space character in the original line
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

		// Check for trailing punctuation in code context
		trimmed := strings.TrimRightFunc(content, unicode.IsSpace)
		punct, rest := extractTrailingPunctuation(trimmed)
		content = rest
		if strings.TrimSpace(content) == "" {
			content = ""
		}

		data[i] = lineData{
			content:         content,
			originallyEmpty: originallyEmpty,
			trailStuff:      punct,
			queuedBraces:    []string{leadingBrace},
		}
	}

	// Step 5: Apply queued leading braces to previous line's trailing queue
	for i := 1; i < len(data); i++ {
		if len(data[i].queuedBraces) > 0 {
			data[i-1].trailStuff += strings.Join(data[i].queuedBraces, "")
		}
	}

	// Step 6: Filter lines that became empty (not originally empty, now empty)
	var filtered []lineData
	for _, d := range data {
		if d.content == "" && !d.originallyEmpty {
			if len(filtered) > 0 {
				filtered[len(filtered)-1].trailStuff += d.trailStuff
			}
			continue
		}
		filtered = append(filtered, d)
	}

	// Step 7: Calculate max_len from content after punctuation extraction
	maxLen := 0
	for _, d := range filtered {
		if utf8.RuneCountInString(d.content) > maxLen {
			maxLen = utf8.RuneCountInString(d.content)
		}
	}

	// Step 8: Format output
	var sb strings.Builder
	for i, d := range filtered {
		if i > 0 {
			sb.WriteString("\n")
		}

		if d.trailStuff != "" {
			padded := d.content
			contentRunes := utf8.RuneCountInString(padded)
			if contentRunes < maxLen {
				padded += strings.Repeat(" ", maxLen-contentRunes)
			}
			sb.WriteString(padded)
			sb.WriteString(" ")
			sb.WriteString(d.trailStuff)
		} else {
			sb.WriteString(d.content)
		}
	}
	sb.WriteString("\n")

	return sb.String()
}
