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
		{tokenizer.CommentLine, "// comment\n"},
	})
}

func TestPHP_BlockComment(t *testing.T) {
	checkSegments(t, "/* block */", []tokenizer.Segment{
		{tokenizer.CommentBlock, "/* block */"},
	})
}

func TestPHP_DoubleQuotedString(t *testing.T) {
	checkSegments(t, `"Hello $name"`, []tokenizer.Segment{
		{tokenizer.String, `"Hello $name"`},
	})
}

func TestPHP_SingleQuotedString(t *testing.T) {
	checkSegments(t, `'literal'`, []tokenizer.Segment{
		{tokenizer.Char, `'literal'`},
	})
}

func TestPHP_SingleQuotedStringWithEscape(t *testing.T) {
	checkSegments(t, `'it\'s'`, []tokenizer.Segment{
		{tokenizer.Char, `'it\'s'`},
	})
}

func TestPHP_BacktickString(t *testing.T) {
	checkSegments(t, "`cmd`", []tokenizer.Segment{
		{tokenizer.String, "`cmd`"},
	})
}

func TestPHP_EchoTag(t *testing.T) {
	checkSegments(t, "<?= $x ?>", []tokenizer.Segment{
		{tokenizer.Code, "<?= $x ?>"},
	})
}

func TestPHP_OpenTag(t *testing.T) {
	checkSegments(t, "<?php echo 1; ?>", []tokenizer.Segment{
		{tokenizer.Code, "<?php echo 1; ?>"},
	})
}

func TestPHP_ShortOpenTag(t *testing.T) {
	checkSegments(t, "<? echo 1; ?>", []tokenizer.Segment{
		{tokenizer.Code, "<? echo 1; ?>"},
	})
}

func TestPHP_Heredoc(t *testing.T) {
	checkSegments(t, "<<<EOT\nhello\nworld\nEOT;", []tokenizer.Segment{
		{tokenizer.String, "<<<EOT\nhello\nworld\nEOT;"},
	})
}

func TestPHP_HeredocQuoted(t *testing.T) {
	checkSegments(t, "<<<'EOT'\nhello\nworld\nEOT;", []tokenizer.Segment{
		{tokenizer.String, "<<<'EOT'\nhello\nworld\nEOT;"},
	})
}

func TestPHP_HeredocWithDollar(t *testing.T) {
	checkSegments(t, "<<<EOT\nhello\nEOT;", []tokenizer.Segment{
		{tokenizer.String, "<<<EOT\nhello\nEOT;"},
	})
}

func TestPHP_CodeOnly(t *testing.T) {
	checkSegments(t, "echo 1;", []tokenizer.Segment{
		{tokenizer.Code, "echo 1;"},
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
		{tokenizer.Code, "<?php echo "},
		{tokenizer.String, "\"hi\""},
		{tokenizer.Code, "; "},
		{tokenizer.CommentLine, "// end\n"},
	})
}
