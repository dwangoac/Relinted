package php

import (
	"testing"

	"relinted/internal/tokenizer"
)

func TestExtractTrailingPunctuation_Semicolon(t *testing.T) {
	punct, rest := extractTrailingPunctuation("$x = 1;")
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != "$x = 1" {
		t.Errorf("expected rest '$x = 1', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_OpenBrace(t *testing.T) {
	punct, rest := extractTrailingPunctuation("function foo() {")
	if punct != "{" {
		t.Errorf("expected punct '{', got %q", punct)
	}
	if rest != "function foo() " {
		t.Errorf("expected rest 'function foo() ', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_ClosingBrace(t *testing.T) {
	punct, rest := extractTrailingPunctuation("}")
	if punct != "}" {
		t.Errorf("expected punct '}', got %q", punct)
	}
	if rest != "" {
		t.Errorf("expected rest '', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_NoPunctuation(t *testing.T) {
	punct, rest := extractTrailingPunctuation("$x = 1")
	if punct != "" {
		t.Errorf("expected punct '', got %q", punct)
	}
	if rest != "$x = 1" {
		t.Errorf("expected rest '$x = 1', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_InString(t *testing.T) {
	punct, rest := extractTrailingPunctuation(`echo "hello;";`)
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != `echo "hello;"` {
		t.Errorf("expected rest 'echo \"hello;\"', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_InLineComment(t *testing.T) {
	punct, rest := extractTrailingPunctuation("$x = 1; // comment;")
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != "$x = 1" {
		t.Errorf("expected rest '$x = 1', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_InHashComment(t *testing.T) {
	punct, rest := extractTrailingPunctuation("$x = 1; # comment;")
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != "$x = 1" {
		t.Errorf("expected rest '$x = 1', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_InBlockComment(t *testing.T) {
	punct, rest := extractTrailingPunctuation("$x = 1; /* ; */")
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != "$x = 1" {
		t.Errorf("expected rest '$x = 1', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_SemicolonInStringNotExtracted(t *testing.T) {
	punct, rest := extractTrailingPunctuation(`echo "a; b"`)
	if punct != "" {
		t.Errorf("expected punct '', got %q", punct)
	}
	if rest != `echo "a; b"` {
		t.Errorf("expected rest 'echo \"a; b\"', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_SemicolonInSingleQuoteNotExtracted(t *testing.T) {
	input := `$x = ';';`
	punct, rest := extractTrailingPunctuation(input)
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != `$x = ';'` {
		t.Errorf("expected rest %q, got %q", `$x = ';'`, rest)
	}
}

func TestExtractTrailingPunctuation_BraceInStringNotExtracted(t *testing.T) {
	punct, rest := extractTrailingPunctuation(`echo "}";`)
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != `echo "}"` {
		t.Errorf("expected rest 'echo \"}\"', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_MultiplePunctuation(t *testing.T) {
	punct, rest := extractTrailingPunctuation("$x = 1; $y = 2;")
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != "$x = 1; $y = 2" {
		t.Errorf("expected rest '$x = 1; $y = 2', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_Heredoc(t *testing.T) {
	input := `echo <<<EOT
hello;
EOT;`
	punct, rest := extractTrailingPunctuation(input)
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != `echo <<<EOT
hello;
EOT` {
		t.Errorf("unexpected rest: %q", rest)
	}
}

func TestExtractTrailingPunctuation_PunctWithTrailingComment(t *testing.T) {
	punct, rest := extractTrailingPunctuation("$x = 1; // trailing comment")
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != "$x = 1" {
		t.Errorf("expected rest '$x = 1', got %q", rest)
	}
}

func TestExpandTabs_Basic(t *testing.T) {
	got := expandTabs("\t\t")
	if got != "        " {
		t.Errorf("expected 8 spaces, got %q", got)
	}
}

func TestExpandTabs_Mixed(t *testing.T) {
	got := expandTabs("ab\tcd")
	if got != "ab  cd" {
		t.Errorf("expected 'ab  cd', got %q", got)
	}
}

func TestSplitLines_Basic(t *testing.T) {
	got := splitLines("a\nb\nc")
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
	if got[0] != "a" || got[1] != "b" || got[2] != "c" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestSplitLines_TrailingNewline(t *testing.T) {
	got := splitLines("a\nb\n")
	if len(got) != 3 {
		t.Fatalf("expected 3 lines (with trailing empty), got %d", len(got))
	}
	if got[2] != "" {
		t.Errorf("expected last line to be empty, got %q", got[2])
	}
}

func TestReconstructText(t *testing.T) {
	segments := []tokenizer.Segment{
		{tokenizer.Code, "echo "},
		{tokenizer.Code, "\"hi\";"},
		{tokenizer.CommentLine, "\n"},
	}
	got := reconstructText(segments)
	if got != "echo \"hi\";\n" {
		t.Errorf("got %q, want %q", got, "echo \"hi\";\n")
	}
}
