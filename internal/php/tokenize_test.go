package php

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
	for i, exp := range expected {
		if got[i].Type != exp.Type {
			t.Errorf("segment %d: type %v, want %v", i, got[i].Type, exp.Type)
		}
		if got[i].Text != exp.Text {
			t.Errorf("segment %d: text %q, want %q", i, got[i].Text, exp.Text)
		}
	}
}

func TestPHP_LineComment(t *testing.T) {
	checkSegments(t, "// comment\n", []tokenizer.Segment{
		{Type: tokenizer.CommentLine, Text: "// comment\n"},
	})
}

func TestPHP_BlockComment(t *testing.T) {
	checkSegments(t, "/* block */", []tokenizer.Segment{
		{Type: tokenizer.CommentBlock, Text: "/* block */"},
	})
}

func TestPHP_DoubleQuotedString(t *testing.T) {
	checkSegments(t, `"Hello $name"`, []tokenizer.Segment{
		{Type: tokenizer.String, Text: `"Hello $name"`},
	})
}

func TestPHP_SingleQuotedString(t *testing.T) {
	checkSegments(t, `'literal'`, []tokenizer.Segment{
		{Type: tokenizer.Char, Text: `'literal'`},
	})
}

func TestPHP_SingleQuotedStringWithEscape(t *testing.T) {
	checkSegments(t, `'it\'s'`, []tokenizer.Segment{
		{Type: tokenizer.Char, Text: `'it\'s'`},
	})
}

func TestPHP_BacktickString(t *testing.T) {
	checkSegments(t, "`cmd`", []tokenizer.Segment{
		{Type: tokenizer.String, Text: "`cmd`"},
	})
}

func TestPHP_EchoTag(t *testing.T) {
	checkSegments(t, "<?= $x ?>", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "<?= $x ?>"},
	})
}

func TestPHP_OpenTag(t *testing.T) {
	checkSegments(t, "<?php echo 1; ?>", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "<?php echo 1; ?>"},
	})
}

func TestPHP_ShortOpenTag(t *testing.T) {
	checkSegments(t, "<? echo 1; ?>", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "<? echo 1; ?>"},
	})
}

func TestPHP_Heredoc(t *testing.T) {
	checkSegments(t, "<<<EOT\nhello\nworld\nEOT;", []tokenizer.Segment{
		{Type: tokenizer.String, Text: "<<<EOT\nhello\nworld\nEOT;"},
	})
}

func TestPHP_HeredocQuoted(t *testing.T) {
	checkSegments(t, "<<<'EOT'\nhello\nworld\nEOT;", []tokenizer.Segment{
		{Type: tokenizer.String, Text: "<<<'EOT'\nhello\nworld\nEOT;"},
	})
}

func TestPHP_HeredocWithDollar(t *testing.T) {
	checkSegments(t, "<<<EOT\nhello\nEOT;", []tokenizer.Segment{
		{Type: tokenizer.String, Text: "<<<EOT\nhello\nEOT;"},
	})
}

func TestPHP_CodeOnly(t *testing.T) {
	checkSegments(t, "echo 1;", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "echo 1;"},
	})
}

func TestPHP_Empty(t *testing.T) {
	got := Tokenize("")
	if len(got) != 0 {
		t.Errorf("expected 0 segments, got %d", len(got))
	}
}

func TestPHP_Mixed(t *testing.T) {
	checkSegments(t, "<?php echo \"hi\"; // end\n", []tokenizer.Segment{
		{Type: tokenizer.Code, Text: "<?php echo "},
		{Type: tokenizer.String, Text: "\"hi\""},
		{Type: tokenizer.Code, Text: "; "},
		{Type: tokenizer.CommentLine, Text: "// end\n"},
	})
}
