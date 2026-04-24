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
		{tokenizer.CommentLine, "// comment\n"},
	})
}

func TestBlockComment(t *testing.T) {
	checkSegments(t, "/* block comment */", []tokenizer.Segment{
		{tokenizer.CommentBlock, "/* block comment */"},
	})
}

func TestDoubleQuotedString(t *testing.T) {
	checkSegments(t, "\"hello\"", []tokenizer.Segment{
		{tokenizer.String, "\"hello\""},
	})
}

func TestRuneLiteral(t *testing.T) {
	checkSegments(t, "'a'", []tokenizer.Segment{
		{tokenizer.Char, "'a'"},
	})
}

func TestRawStringLiteral(t *testing.T) {
	checkSegments(t, "`hello`", []tokenizer.Segment{
		{tokenizer.String, "`hello`"},
	})
}

func TestRawStringWithSpecialChars(t *testing.T) {
	checkSegments(t, "`fmt.Println(\"}\")`", []tokenizer.Segment{
		{tokenizer.String, "`fmt.Println(\"}\")`"},
	})
}

func TestRawStringMultiline(t *testing.T) {
	input := "`line1\nline2\nline3`"
	checkSegments(t, input, []tokenizer.Segment{
		{tokenizer.String, "`line1\nline2\nline3`"},
	})
}

func TestCodeSegment(t *testing.T) {
	checkSegments(t, "func main() {", []tokenizer.Segment{
		{tokenizer.Code, "func main() {"},
	})
}

func TestMixedInput(t *testing.T) {
	input := "func main() {\n\t// print hello\n\tfmt.Println(\"hello\")\n}"
	checkSegments(t, input, []tokenizer.Segment{
		{tokenizer.Code, "func main() {\n\t"},
		{tokenizer.CommentLine, "// print hello\n"},
		{tokenizer.Code, "\tfmt.Println("},
		{tokenizer.String, "\"hello\""},
		{tokenizer.Code, ")\n}"},
	})
}
