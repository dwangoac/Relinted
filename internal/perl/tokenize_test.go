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
		{tokenizer.Code, "my $x = 1;\n"},
	})
}

func TestTokenize_HashComment(t *testing.T) {
	checkSegments(t, "my $x = 1; # comment\n", []tokenizer.Segment{
		{tokenizer.Code, "my $x = 1; "},
		{tokenizer.CommentLine, "# comment\n"},
	})
}

func TestTokenize_DoubleQuoteString(t *testing.T) {
	checkSegments(t, "print \"hello\";\n", []tokenizer.Segment{
		{tokenizer.Code, "print "},
		{tokenizer.String, "\"hello\""},
		{tokenizer.Code, ";\n"},
	})
}

func TestTokenize_DoubleQuoteStringWithEscape(t *testing.T) {
	checkSegments(t, "print \"hello\\n\";\n", []tokenizer.Segment{
		{tokenizer.Code, "print "},
		{tokenizer.String, "\"hello\\n\""},
		{tokenizer.Code, ";\n"},
	})
}

func TestTokenize_SingleQuoteString(t *testing.T) {
	checkSegments(t, "my $x = 'hello';\n", []tokenizer.Segment{
		{tokenizer.Code, "my $x = "},
		{tokenizer.String, "'hello'"},
		{tokenizer.Code, ";\n"},
	})
}

func TestTokenize_SingleQuoteStringWithEscape(t *testing.T) {
	checkSegments(t, "my $x = 'it\\'s';\n", []tokenizer.Segment{
		{tokenizer.Code, "my $x = "},
		{tokenizer.String, "'it\\'s'"},
		{tokenizer.Code, ";\n"},
	})
}

func TestTokenize_Regex(t *testing.T) {
	checkSegments(t, "if ($x =~ /pattern/) {\n", []tokenizer.Segment{
		{tokenizer.Code, "if ($x =~ "},
		{tokenizer.String, "/pattern/"},
		{tokenizer.Code, ") {\n"},
	})
}

func TestTokenize_RegexWithEscape(t *testing.T) {
	checkSegments(t, "if ($x =~ /pat\\/tern/) {\n", []tokenizer.Segment{
		{tokenizer.Code, "if ($x =~ "},
		{tokenizer.String, "/pat\\/tern/"},
		{tokenizer.Code, ") {\n"},
	})
}

func TestTokenize_Mixed(t *testing.T) {
	checkSegments(t, "my $x = \"hi\"; # comment\n", []tokenizer.Segment{
		{tokenizer.Code, "my $x = "},
		{tokenizer.String, "\"hi\""},
		{tokenizer.Code, "; "},
		{tokenizer.CommentLine, "# comment\n"},
	})
}

func TestTokenize_Empty(t *testing.T) {
	got := Tokenize("")
	if len(got) != 0 {
		t.Errorf("expected 0 segments, got %d", len(got))
	}
}
