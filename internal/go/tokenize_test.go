package go_pkg

import (
	"relinted/internal/tokenizer"
	"testing"
)

func checkSegments(t *testing.T, input string, expected []tokenizer.Segment) {
	t.Helper()
	got := Tokenize(input)
	if len(got) != len(expected) {
		t.Fatalf("got %d segments, want %d", len(got), len(expected))
	}
	for i := range got {
		if got[i].Type != expected[i].Type {
			t.Errorf("segment %d: type %v, want %v", i, got[i].Type, expected[i].Type)
		}
		if got[i].Text != expected[i].Text {
			t.Errorf("segment %d: text %q, want %q", i, got[i].Text, expected[i].Text)
		}
	}
}

func TestLineComment(t *testing.T) {
	checkSegments(t, "// comment\n", []tokenizer.Segment{
		{Type: tokenizer.CommentLine, Text: "// comment\n"},
	})
}

func TestBlockComment(t *testing.T) {
	checkSegments(t, "/* block comment */", []tokenizer.Segment{
		{Type: tokenizer.CommentBlock, Text: "/* block comment */"},
	})
}

func TestDoubleQuotedString(t *testing.T) {
	checkSegments(t, "\"hello\"", []tokenizer.Segment{
		{Type: tokenizer.String, Text: "\"hello\""},
	})
}

func TestRuneLiteral(t *testing.T) {
	checkSegments(t, "'a'", []tokenizer.Segment{
		{Type: tokenizer.Char, Text: "'a'"},
	})
}

func TestRawStringLiteral(t *testing.T) {
	checkSegments(t, "`hello`", []tokenizer.Segment{
		{Type: tokenizer.String, Text: "`hello`"},
	})
}

func TestRawStringWithSpecialChars(t *testing.T) {
	checkSegments(t, "`fmt.Println(\"}\")`", []tokenizer.Segment{
		{Type: tokenizer.String, Text: "`fmt.Println(\"}\")`"},
	})
}

func TestRawStringMultiline(t *testing.T) {
	input := "`line1\nline2\nline3`"
	checkSegments(t, input, []tokenizer.Segment{
		{Type: tokenizer.String, Text: "`line1\nline2\nline3`"},
	})
}

func TestCodeSegment(t *testing.T) {
	checkSegments(t, "func main() {", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "func main() {"},
	})
}

func TestMixedInput(t *testing.T) {
	input := "func main() {\n\t// print hello\n\tfmt.Println(\"hello\")\n}"
	checkSegments(t, input, []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "func main() {\n\t"},
		{Type: tokenizer.CommentLine, Text: "// print hello\n"},
		{Type: tokenizer.Code, Text: "\tfmt.Println("},
		{Type: tokenizer.String, Text: "\"hello\""},
		{Type: tokenizer.Code, Text: ")\n}"},
	})
}
