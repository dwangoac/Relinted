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
			// Line comment: scan to end of line (include trailing newline)
			j := i
			for j < len(input) && input[j] != '\n' {
				j++
			}
			if j < len(input) {
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
