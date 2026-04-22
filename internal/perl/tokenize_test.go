package perl

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
	checkSegments(t, "my $x = 1;\n", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "my $x = 1;\n"},
	})
}

func TestTokenize_HashComment(t *testing.T) {
	checkSegments(t, "my $x = 1; # comment\n", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "my $x = 1; "},
		{Type: tokenizer.CommentLine, Text: "# comment\n"},
	})
}

func TestTokenize_DoubleQuoteString(t *testing.T) {
	checkSegments(t, "print \"hello\";\n", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "print "},
		{Type: tokenizer.String, Text: "\"hello\""},
		{Type: tokenizer.Code, Text: ";\n"},
	})
}

func TestTokenize_DoubleQuoteStringWithEscape(t *testing.T) {
	checkSegments(t, "print \"hello\\n\";\n", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "print "},
		{Type: tokenizer.String, Text: "\"hello\\n\""},
		{Type: tokenizer.Code, Text: ";\n"},
	})
}

func TestTokenize_SingleQuoteString(t *testing.T) {
	checkSegments(t, "my $x = 'hello';\n", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "my $x = "},
		{Type: tokenizer.String, Text: "'hello'"},
		{Type: tokenizer.Code, Text: ";\n"},
	})
}

func TestTokenize_SingleQuoteStringWithEscape(t *testing.T) {
	checkSegments(t, "my $x = 'it\\'s';\n", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "my $x = "},
		{Type: tokenizer.String, Text: "'it\\'s'"},
		{Type: tokenizer.Code, Text: ";\n"},
	})
}

func TestTokenize_Regex(t *testing.T) {
	checkSegments(t, "if ($x =~ /pattern/) {\n", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "if ($x =~ "},
		{Type: tokenizer.String, Text: "/pattern/"},
		{Type: tokenizer.Code, Text: ") {\n"},
	})
}

func TestTokenize_RegexWithEscape(t *testing.T) {
	checkSegments(t, "if ($x =~ /pat\\/tern/) {\n", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "if ($x =~ "},
		{Type: tokenizer.String, Text: "/pat\\/tern/"},
		{Type: tokenizer.Code, Text: ") {\n"},
	})
}

func TestTokenize_Mixed(t *testing.T) {
	checkSegments(t, "my $x = \"hi\"; # comment\n", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "my $x = "},
		{Type: tokenizer.String, Text: "\"hi\""},
		{Type: tokenizer.Code, Text: "; "},
		{Type: tokenizer.CommentLine, Text: "# comment\n"},
	})
}

func TestTokenize_Empty(t *testing.T) {
	got := Tokenize("")
	if len(got) != 0 {
		t.Errorf("expected 0 segments, got %d", len(got))
	}
}
