package php

import (
	"relinted/internal/tokenizer"
)

// Tokenize scans PHP source code character by character and produces a list of segments.
func Tokenize(input string) []tokenizer.Segment {
	var segments []tokenizer.Segment
	i := 0
	for i < len(input) {
		switch {
		case i+1 < len(input) && input[i] == '/' && input[i+1] == '/':
			j := i
			for j < len(input) && input[j] != '\n' {
				j++
			}
			if j < len(input) {
				j++
			}
			segments = append(segments, tokenizer.Segment{Type: tokenizer.CommentLine, Text: input[i:j]})
			i = j

		case i+1 < len(input) && input[i] == '/' && input[i+1] == '*':
			j := i + 2
			for j+1 < len(input) && !(input[j] == '*' && input[j+1] == '/') {
				j++
			}
			if j+1 < len(input) {
				j += 2
			}
			segments = append(segments, tokenizer.Segment{Type: tokenizer.CommentBlock, Text: input[i:j]})
			i = j

		case input[i] == '"':
			j := i + 1
			for j < len(input) && input[j] != '"' {
				j++
			}
			if j < len(input) {
				j++
			}
			segments = append(segments, tokenizer.Segment{Type: tokenizer.String, Text: input[i:j]})
			i = j

		case input[i] == '\'':
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
			segments = append(segments, tokenizer.Segment{Type: tokenizer.Char, Text: input[i:j]})
			i = j

		case input[i] == '`':
			j := i + 1
			for j < len(input) && input[j] != '`' {
				j++
			}
			if j < len(input) {
				j++
			}
			segments = append(segments, tokenizer.Segment{Type: tokenizer.String, Text: input[i:j]})
			i = j

		case i+2 < len(input) && input[i] == '<' && input[i+1] == '<' && input[i+2] == '<':
			j := i + 3
			var quote byte
			if j < len(input) && (input[j] == '"' || input[j] == '\'') {
				quote = input[j]
				j++
			}
			identStart := j
			for j < len(input) && (input[j] == '_' || input[j] == '$' ||
				(input[j] >= 'a' && input[j] <= 'z') ||
				(input[j] >= 'A' && input[j] <= 'Z') ||
				(input[j] >= '0' && input[j] <= '9')) {
				j++
			}
			if j < len(input) && input[j] == quote {
				j++
			}
			if j == identStart {
				segments = append(segments, tokenizer.Segment{Type: tokenizer.Code, Text: input[i : i+1]})
				i++
				continue
			}
			ident := input[identStart:j]
			found := false
			for j < len(input) {
				if input[j] == '\n' {
					k := j + 1
					if k+len(ident) <= len(input) && input[k:k+len(ident)] == ident {
						pos := k + len(ident)
						if pos < len(input) && input[pos] == ';' {
							j = pos + 1
							segments = append(segments, tokenizer.Segment{Type: tokenizer.String, Text: input[i:j]})
							i = j
							found = true
							break
						}
					}
				}
				j++
			}
			if !found {
				segments = append(segments, tokenizer.Segment{Type: tokenizer.String, Text: input[i:]})
				i = len(input)
			}

		case input[i] == '<' && i+1 < len(input) && input[i+1] == '?':
			tagStart := i
			j := i + 2
			for j < len(input) {
				if j+1 < len(input) && input[j] == '?' && input[j+1] == '>' {
					j += 2
					break
				}
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
				if ch == '<' && j+1 < len(input) && j+2 < len(input) && input[j+1] == '<' && input[j+2] == '<' {
					break
				}
				j++
			}
			segments = append(segments, tokenizer.Segment{Type: tokenizer.Code, Text: input[tagStart:j]})
			i = j

		default:
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
				if ch == '<' && j+1 < len(input) && j+2 < len(input) && input[j+1] == '<' && input[j+2] == '<' {
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
