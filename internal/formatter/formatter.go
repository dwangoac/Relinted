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

	// Step 3: Tokenize entire input and reconstruct to normalize
	segments := tokenizer.Tokenize(input)
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

		data[i] = lineData{
			content:         content,
			originallyEmpty: originallyEmpty,
			ownPunct:        punct,
			queuedBraces:    []string{leadingBrace},
		}
	}

	// Step 6: Apply queued leading braces to previous line's trailing queue
	for i := 1; i < len(data); i++ {
		if len(data[i].queuedBraces) > 0 {
			data[i-1].ownPunct += strings.Join(data[i].queuedBraces, "")
		}
	}

	// Step 7: Filter lines that became empty (not originally empty, now empty)
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

	// Step 8: Calculate max_len from content after punctuation extraction
	maxLen := 0
	for _, d := range filtered {
		if len(d.content) > maxLen {
			maxLen = len(d.content)
		}
	}

	// Step 9: Format output
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

// extractTrailingPunctuation scans the line to find the last semicolon, brace,
// or closing brace that appears in code context (not inside a string, char
// literal, line comment, or block comment). If the punctuation is followed
// only by whitespace, it is extracted. Returns (remaining_content, punctuation).
func extractTrailingPunctuation(line string) (punctuation string, remaining string) {
	lastPunctPos := -1
	inString := false
	inChar := false
	inLineComment := false
	inBlockComment := false

	for i := 0; i < len(line); i++ {
		ch := line[i]

		if inLineComment {
			if ch == '\n' {
				inLineComment = false
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
		if i+1 < len(line) && ch == '/' && line[i+1] == '*' {
			inBlockComment = true
			i++
			continue
		}

		if ch == ';' || ch == '{' || ch == '}' {
			lastPunctPos = i
		}
	}

	if lastPunctPos >= 0 {
		rest := line[lastPunctPos+1:]
		allSpace := true
		for _, ch := range rest {
			if !unicode.IsSpace(ch) {
				allSpace = false
				break
			}
		}
		if allSpace {
			return string(line[lastPunctPos]), line[:lastPunctPos]
		}
	}

	return "", line
}
