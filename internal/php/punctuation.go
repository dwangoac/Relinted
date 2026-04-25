package php

import (
	"strings"
	"unicode"

	"relinted/internal/tokenizer"
)

// extractTrailingPunctuation scans the line to find the last punctuation character
// (semicolon, opening brace, or closing brace) that appears in code context
// (not inside a string, char literal, backtick string, block comment, or line comment).
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
		if ch == '#' {
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
			if ch == '#' {
				break
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
