package java

import (
	"testing"

	"relinted/internal/tokenizer"
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

func TestTokenize_CodeOnly(t *testing.T) {
	checkSegments(t, "int x = 1;\n", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "int x = 1;\n"},
	})
}

func TestTokenize_LineComment(t *testing.T) {
	checkSegments(t, "int x = 1; // comment\n", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "int x = 1; "},
		{Type: tokenizer.CommentLine, Text: "// comment\n"},
	})
}

func TestTokenize_BlockComment(t *testing.T) {
	checkSegments(t, "int x = 1; /* comment */\n", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "int x = 1; "},
		{Type: tokenizer.CommentBlock, Text: "/* comment */"},
		{Type: tokenizer.Code, Text: "\n"},
	})
}

func TestTokenize_DoubleQuoteString(t *testing.T) {
	checkSegments(t, "System.out.println(\"hello\");\n", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "System.out.println("},
		{Type: tokenizer.String, Text: `"hello"`},
		{Type: tokenizer.Code, Text: ");\n"},
	})
}

func TestTokenize_DoubleQuoteStringWithEscape(t *testing.T) {
	checkSegments(t, "String s = \"hello\\n\";\n", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "String s = "},
		{Type: tokenizer.String, Text: `"hello\n"`},
		{Type: tokenizer.Code, Text: ";\n"},
	})
}

func TestTokenize_CharLiteral(t *testing.T) {
	checkSegments(t, "char c = 'a';\n", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "char c = "},
		{Type: tokenizer.Char, Text: `'a'`},
		{Type: tokenizer.Code, Text: ";\n"},
	})
}

func TestTokenize_CharLiteralWithEscape(t *testing.T) {
	checkSegments(t, "char c = '\\\\';\n", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "char c = "},
		{Type: tokenizer.Char, Text: `'\\'`},
		{Type: tokenizer.Code, Text: ";\n"},
	})
}

func TestTokenize_TextBlock(t *testing.T) {
	checkSegments(t, "String s = \"\"\"hello\"\"\";\n", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "String s = "},
		{Type: tokenizer.String, Text: `"""hello"""`},
		{Type: tokenizer.Code, Text: ";\n"},
	})
}

func TestTokenize_TextBlockMultiline(t *testing.T) {
	checkSegments(t, "String s = \"\"\"\nhello\nworld\n\"\"\";\n", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "String s = "},
		{Type: tokenizer.String, Text: `"""
hello
world
"""`},
		{Type: tokenizer.Code, Text: ";\n"},
	})
}

func TestTokenize_Mixed(t *testing.T) {
	checkSegments(t, "String s = \"hi\"; // comment\n", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "String s = "},
		{Type: tokenizer.String, Text: `"hi"`},
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
