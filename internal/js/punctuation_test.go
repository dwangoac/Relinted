package js

import (
	"testing"

	"relinted/internal/tokenizer"
)

func TestExtractTrailingPunctuation_Semicolon(t *testing.T) {
	punct, rest := extractTrailingPunctuation("let x = 1;")
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != "let x = 1" {
		t.Errorf("expected rest 'let x = 1', got %q", rest)
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
	punct, rest := extractTrailingPunctuation("let x = 1")
	if punct != "" {
		t.Errorf("expected punct '', got %q", punct)
	}
	if rest != "let x = 1" {
		t.Errorf("expected rest 'let x = 1', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_InString(t *testing.T) {
	punct, rest := extractTrailingPunctuation(`console.log("hello;");`)
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != `console.log("hello;")` {
		t.Errorf("expected rest 'console.log(\"hello;\")', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_InLineComment(t *testing.T) {
	punct, rest := extractTrailingPunctuation("let x = 1; // comment;")
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != "let x = 1" {
		t.Errorf("expected rest 'let x = 1', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_InBlockComment(t *testing.T) {
	punct, rest := extractTrailingPunctuation("let x = 1; /* ; */")
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != "let x = 1" {
		t.Errorf("expected rest 'let x = 1', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_SemicolonInStringNotExtracted(t *testing.T) {
	punct, rest := extractTrailingPunctuation(`console.log("a; b")`)
	if punct != "" {
		t.Errorf("expected punct '', got %q", punct)
	}
	if rest != `console.log("a; b")` {
		t.Errorf("expected rest 'console.log(\"a; b\")', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_SemicolonInCharNotExtracted(t *testing.T) {
	input := "let c = ';'"
	punct, rest := extractTrailingPunctuation(input)
	if punct != "" {
		t.Errorf("expected punct '', got %q", punct)
	}
	if rest != input {
		t.Errorf("expected rest %q, got %q", input, rest)
	}
}

func TestExtractTrailingPunctuation_BraceInStringNotExtracted(t *testing.T) {
	punct, rest := extractTrailingPunctuation(`console.log("}");`)
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != `console.log("}")` {
		t.Errorf("expected rest 'console.log(\"}\")', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_MultiplePunctuation(t *testing.T) {
	punct, rest := extractTrailingPunctuation("let x = 1; let y = 2;")
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != "let x = 1; let y = 2" {
		t.Errorf("expected rest 'let x = 1; let y = 2', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_BracesOnly(t *testing.T) {
	punct, rest := extractTrailingPunctuation("if (x) {")
	if punct != "{" {
		t.Errorf("expected punct '{', got %q", punct)
	}
	if rest != "if (x) " {
		t.Errorf("expected rest 'if (x) ', got %q", rest)
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
		{tokenizer.Code, "let "},
		{tokenizer.Code, "x = "},
		{tokenizer.Code, "1;\n"},
	}
	got := reconstructText(segments)
	if got != "let x = 1;\n" {
		t.Errorf("got %q, want %q", got, "let x = 1;\n")
	}
}
