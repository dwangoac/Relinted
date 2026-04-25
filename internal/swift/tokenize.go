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
