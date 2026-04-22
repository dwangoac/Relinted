package tokenizer

import "testing"

func checkSegments(t *testing.T, input string, expected []Segment) {
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

func TestTokenize_CodeOnly(t *testing.T) {
	checkSegments(t, "int x = 0;", []Segment{
		{Code, "int x = 0;"},
	})
}

func TestTokenize_LineComment(t *testing.T) {
	checkSegments(t, "int x; // comment\n", []Segment{
		{Code, "int x; "},
		{CommentLine, "// comment\n"},
	})
}

func TestTokenize_BlockComment(t *testing.T) {
	checkSegments(t, "int x /* block */;\n", []Segment{
		{Code, "int x "},
		{CommentBlock, "/* block */"},
		{Code, ";\n"},
	})
}

func TestTokenize_String(t *testing.T) {
	checkSegments(t, "printf(\"hello\");\n", []Segment{
		{Code, "printf("},
		{String, "\"hello\""},
		{Code, ");\n"},
	})
}

func TestTokenize_StringWithEscape(t *testing.T) {
	checkSegments(t, "printf(\"hello\\n\");\n", []Segment{
		{Code, "printf("},
		{String, "\"hello\\n\""},
		{Code, ");\n"},
	})
}

func TestTokenize_Char(t *testing.T) {
	checkSegments(t, "char c = 'a';\n", []Segment{
		{Code, "char c = "},
		{Char, "'a'"},
		{Code, ";\n"},
	})
}

func TestTokenize_Mixed(t *testing.T) {
	checkSegments(t, "/* start */ int x = \"hi\"; // end\n", []Segment{
		{CommentBlock, "/* start */"},
		{Code, " int x = "},
		{String, "\"hi\""},
		{Code, "; "},
		{CommentLine, "// end\n"},
	})
}

func TestTokenize_Empty(t *testing.T) {
	got := Tokenize("")
	if len(got) != 0 {
		t.Errorf("expected 0 segments, got %d", len(got))
	}
}
