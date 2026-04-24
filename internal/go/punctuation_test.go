package go_pkg

import "testing"

func TestExtractTrailingPunctuation(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantPunct  string
		wantRemain string
	}{
		// Basic punctuation extraction
		{
			name:       "semicolon at end",
			input:      "var x int = 1;",
			wantPunct:  ";",
			wantRemain: "var x int = 1",
		},
		{
			name:       "opening brace at start of line",
			input:      "{",
			wantPunct:  "{",
			wantRemain: "",
		},
		{
			name:       "closing brace at start of line",
			input:      "}",
			wantPunct:  "}",
			wantRemain: "",
		},
		{
			name:       "no punctuation",
			input:      "var x int",
			wantPunct:  "",
			wantRemain: "var x int",
		},
		// Punctuation inside strings should NOT be extracted
		{
			name:       "semicolon inside double-quoted string",
			input:      "fmt.Println(\"hello;\")",
			wantPunct:  "",
			wantRemain: "fmt.Println(\"hello;\")",
		},
		{
			name:       "semicolon inside rune literal",
			input:      "var c = ';'",
			wantPunct:  "",
			wantRemain: "var c = ';'",
		},
		{
			name:       "semicolon inside block comment",
			input:      "/* ; */",
			wantPunct:  "",
			wantRemain: "/* ; */",
		},
		{
			name:       "semicolon inside line comment",
			input:      "// ;",
			wantPunct:  "",
			wantRemain: "// ;",
		},
		{
			name:       "brace inside double-quoted string",
			input:      "fmt.Println(\"}\")",
			wantPunct:  "",
			wantRemain: "fmt.Println(\"}\")",
		},
		{
			name:       "brace inside rune literal",
			input:      "var c = '{'",
			wantPunct:  "",
			wantRemain: "var c = '{'",
		},
		// Multiple punctuation — last one extracted
		{
			name:       "multiple semicolons",
			input:      "x := 1; y := 2;",
			wantPunct:  ";",
			wantRemain: "x := 1; y := 2",
		},
		{
			name:       "brace and semicolon mixed",
			input:      "x := 1; y := {",
			wantPunct:  "{",
			wantRemain: "x := 1; y := ",
		},
		// Edge cases
		{
			name:       "empty line",
			input:      "",
			wantPunct:  "",
			wantRemain: "",
		},
		{
			name:       "whitespace only",
			input:      "   ",
			wantPunct:  "",
			wantRemain: "   ",
		},
		{
			name:       "punctuation only",
			input:      ";",
			wantPunct:  ";",
			wantRemain: "",
		},
		// Go raw string literal (backtick) — punctuation inside NOT extracted
		{
			name:       "semicolon inside raw string literal",
			input:      "`hello;`",
			wantPunct:  "",
			wantRemain: "`hello;`",
		},
		{
			name:       "brace inside raw string literal",
			input:      "`func main() {}`",
			wantPunct:  "",
			wantRemain: "`func main() {}`",
		},
		{
			name:       "raw string followed by semicolon",
			input:      "`hello;` + x;",
			wantPunct:  ";",
			wantRemain: "`hello;` + x",
		},
		// Go rune literal — punctuation inside NOT extracted
		{
			name:       "semicolon inside rune literal",
			input:      "var c = ';'",
			wantPunct:  "",
			wantRemain: "var c = ';'",
		},
		{
			name:       "brace inside rune literal",
			input:      "var c = '{'",
			wantPunct:  "",
			wantRemain: "var c = '{'",
		},
		// Go-specific: escaped characters in rune literals
		{
			name:       "escaped quote in rune literal",
			input:      "'\\'' ;",
			wantPunct:  ";",
			wantRemain: "'\\'' ",
		},
		{
			name:       "escaped backslash in rune literal",
			input:      "'\\\\' ;",
			wantPunct:  ";",
			wantRemain: "'\\\\' ",
		},
		// Go-specific: consecutive backticks (empty raw string)
		{
			name:       "consecutive backticks empty raw string",
			input:      "`` ;",
			wantPunct:  ";",
			wantRemain: "`` ",
		},
		// Go-specific: lone backtick is code
		{
			name:       "lone unclosed backtick treats rest as raw string",
			input:      "` }",
			wantPunct:  "",
			wantRemain: "` }",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPunct, gotRemain := extractTrailingPunctuation(tt.input)
			if gotPunct != tt.wantPunct {
				t.Errorf("extractTrailingPunctuation(%q) punctuation = %q, want %q", tt.input, gotPunct, tt.wantPunct)
			}
			if gotRemain != tt.wantRemain {
				t.Errorf("extractTrailingPunctuation(%q) remaining = %q, want %q", tt.input, gotRemain, tt.wantRemain)
			}
		})
	}
}
