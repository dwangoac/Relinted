package rust

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
		{tokenizer.Code, "let x = 1;\n"},
	})
}

func TestTokenize_LineComment(t *testing.T) {
	checkSegments(t, "let x = 1; // comment\n", []tokenizer.Segment{
		{tokenizer.Code, "let x = 1; "},
		{tokenizer.CommentLine, "// comment\n"},
	})
}

func TestTokenize_BlockComment(t *testing.T) {
	checkSegments(t, "/* comment */\n", []tokenizer.Segment{
		{tokenizer.CommentBlock, "/* comment */"},
		{tokenizer.Code, "\n"},
	})
}

func TestTokenize_BlockCommentMultiline(t *testing.T) {
	checkSegments(t, "/* multi\nline */\n", []tokenizer.Segment{
		{tokenizer.CommentBlock, "/* multi\nline */"},
		{tokenizer.Code, "\n"},
	})
}

func TestTokenize_StringLiteral(t *testing.T) {
	checkSegments(t, "println!(\"hello\");\n", []tokenizer.Segment{
		{tokenizer.Code, "println!("},
		{tokenizer.String, "\"hello\""},
		{tokenizer.Code, ");\n"},
	})
}

func TestTokenize_StringWithEscape(t *testing.T) {
	checkSegments(t, "println!(\"he\\\"llo\");\n", []tokenizer.Segment{
		{tokenizer.Code, "println!("},
		{tokenizer.String, "\"he\\\"llo\""},
		{tokenizer.Code, ");\n"},
	})
}

func TestTokenize_CharLiteral(t *testing.T) {
	checkSegments(t, "let c = 'a';\n", []tokenizer.Segment{
		{tokenizer.Code, "let c = "},
		{tokenizer.Char, "'a'"},
		{tokenizer.Code, ";\n"},
	})
}

func TestTokenize_CharWithEscape(t *testing.T) {
	checkSegments(t, "let c = '\\n';\n", []tokenizer.Segment{
		{tokenizer.Code, "let c = "},
		{tokenizer.Char, "'\\n'"},
		{tokenizer.Code, ";\n"},
	})
}

func TestTokenize_Mixed(t *testing.T) {
	checkSegments(t, "let s = \"hi\"; // comment\n", []tokenizer.Segment{
		{tokenizer.Code, "let s = "},
		{tokenizer.String, "\"hi\""},
		{tokenizer.Code, "; "},
		{tokenizer.CommentLine, "// comment\n"},
	})
}

func TestTokenize_Attribute(t *testing.T) {
	checkSegments(t, "#[derive(Debug)]\n", []tokenizer.Segment{
		{tokenizer.Code, "#[derive(Debug)]\n"},
	})
}

func TestTokenize_Empty(t *testing.T) {
	got := Tokenize("")
	if len(got) != 0 {
		t.Errorf("expected 0 segments, got %d", len(got))
	}
}
