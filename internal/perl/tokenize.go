package perl

import "relinted/internal/tokenizer"

// Tokenize scans the input character by character and produces a list of segments.
// Perl-specific rules:
//   - # starts a line comment (scan to newline)
//   - " starts a string literal (scan to closing ", handling \" escapes)
//   - ' starts a string literal (scan to closing ', handling \' escapes)
//   - / starts a regex literal (scan to closing /, handling \/ escapes)
//   - Everything else is Code
func Tokenize(input string) []tokenizer.Segment {
	var segments []tokenizer.Segment
	i := 0
	for i < len(input) {
		switch {
		case input[i] == '#':
			// Line comment: scan to end of line (include trailing newline)
			j := i
			for j < len(input) && input[j] != '\n' {
				j++
			}
			if j < len(input) {
				j++
			}
			segments = append(segments, tokenizer.Segment{Type: tokenizer.CommentLine, Text: input[i:j]})
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
			segments = append(segments, tokenizer.Segment{Type: tokenizer.String, Text: input[i:j]})
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
			segments = append(segments, tokenizer.Segment{Type: tokenizer.String, Text: input[i:j]})
			i = j

		case input[i] == '/':
			// Regex literal: scan to closing /, handling \/ escapes
			j := i + 1
			for j < len(input) && input[j] != '/' {
				if input[j] == '\\' && j+1 < len(input) {
					j++
				}
				j++
			}
			if j < len(input) {
				j++
			}
			segments = append(segments, tokenizer.Segment{Type: tokenizer.String, Text: input[i:j]})
			i = j

		default:
			// Code: scan to next special character
			j := i + 1
			for j < len(input) {
				ch := input[j]
				if ch == '"' || ch == '\'' || ch == '/' || ch == '#' {
					break
				}
				j++
			}
			segments = append(segments, tokenizer.Segment{Type: tokenizer.Code, Text: input[i:j]})
			i = j
		}
	}
	return segments
}
