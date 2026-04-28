package js

import "relinted/internal/tokenizer"
import "testing"

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

func TestTokenize_CodeOnly(t *testing.T) {
	checkSegments(t, "let x = 1;\n", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "let x = 1;\n"},
	})
}

func TestTokenize_LineComment(t *testing.T) {
	checkSegments(t, "let x = 1; // comment\n", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "let x = 1; "},
		{Type: tokenizer.CommentLine, Text: "// comment\n"},
	})
}

func TestTokenize_BlockComment(t *testing.T) {
	checkSegments(t, "/* comment */\n", []tokenizer.Segment{
		{Type: tokenizer.CommentBlock, Text: "/* comment */"},
		{Type: tokenizer.Code, Text: "\n"},
	})
}

func TestTokenize_BlockCommentMultiline(t *testing.T) {
	checkSegments(t, "/* multi\nline */\n", []tokenizer.Segment{
		{Type: tokenizer.CommentBlock, Text: "/* multi\nline */"},
		{Type: tokenizer.Code, Text: "\n"},
	})
}

func TestTokenize_StringLiteral(t *testing.T) {
	checkSegments(t, "console.log(\"hello\");\n", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "console.log("},
		{Type: tokenizer.String, Text: "\"hello\""},
		{Type: tokenizer.Code, Text: ");\n"},
	})
}

func TestTokenize_StringWithEscape(t *testing.T) {
	checkSegments(t, "console.log(\"he\\\"llo\");\n", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "console.log("},
		{Type: tokenizer.String, Text: "\"he\\\"llo\""},
		{Type: tokenizer.Code, Text: ");\n"},
	})
}

func TestTokenize_CharLiteral(t *testing.T) {
	checkSegments(t, "let c = 'a';\n", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "let c = "},
		{Type: tokenizer.Char, Text: "'a'"},
		{Type: tokenizer.Code, Text: ";\n"},
	})
}

func TestTokenize_CharWithEscape(t *testing.T) {
	checkSegments(t, "let c = '\\n';\n", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "let c = "},
		{Type: tokenizer.Char, Text: "'\\n'"},
		{Type: tokenizer.Code, Text: ";\n"},
	})
}

func TestTokenize_Mixed(t *testing.T) {
	checkSegments(t, "let s = \"hi\"; // comment\n", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "let s = "},
		{Type: tokenizer.String, Text: "\"hi\""},
		{Type: tokenizer.Code, Text: "; "},
		{Type: tokenizer.CommentLine, Text: "// comment\n"},
	})
}

func TestTokenize_Empty(t *testing.T) {
	got := Tokenize("")
	if len(got) != 0 {
		t.Errorf("expected 0 segments, got %d", len(got))
	}
}
