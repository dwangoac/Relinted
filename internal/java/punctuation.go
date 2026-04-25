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
