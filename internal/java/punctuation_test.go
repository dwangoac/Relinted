package java

import (
	"testing"

	"relinted/internal/tokenizer"
)

func TestExtractTrailingPunctuation_Semicolon(t *testing.T) {
	punct, rest, nextInTextBlock := extractTrailingPunctuation("$x = 1;", false)
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != "$x = 1" {
		t.Errorf("expected rest '$x = 1', got %q", rest)
	}
	if nextInTextBlock {
		t.Errorf("expected nextInTextBlock false")
	}
}

func TestExtractTrailingPunctuation_OpenBrace(t *testing.T) {
	punct, rest, nextInTextBlock := extractTrailingPunctuation("function foo() {", false)
	if punct != "{" {
		t.Errorf("expected punct '{', got %q", punct)
	}
	if rest != "function foo() " {
		t.Errorf("expected rest 'function foo() ', got %q", rest)
	}
	if nextInTextBlock {
		t.Errorf("expected nextInTextBlock false")
	}
}

func TestExtractTrailingPunctuation_ClosingBrace(t *testing.T) {
	punct, rest, nextInTextBlock := extractTrailingPunctuation("}", false)
	if punct != "}" {
		t.Errorf("expected punct '}', got %q", punct)
	}
	if rest != "" {
		t.Errorf("expected rest '', got %q", rest)
	}
	if nextInTextBlock {
		t.Errorf("expected nextInTextBlock false")
	}
}

func TestExtractTrailingPunctuation_NoPunctuation(t *testing.T) {
	punct, rest, nextInTextBlock := extractTrailingPunctuation("$x = 1", false)
	if punct != "" {
		t.Errorf("expected punct '', got %q", punct)
	}
	if rest != "$x = 1" {
		t.Errorf("expected rest '$x = 1', got %q", rest)
	}
	if nextInTextBlock {
		t.Errorf("expected nextInTextBlock false")
	}
}

func TestExtractTrailingPunctuation_InString(t *testing.T) {
	punct, rest, nextInTextBlock := extractTrailingPunctuation(`System.out.println("hello;");`, false)
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != `System.out.println("hello;")` {
		t.Errorf("expected rest 'System.out.println(\"hello;\")', got %q", rest)
	}
	if nextInTextBlock {
		t.Errorf("expected nextInTextBlock false")
	}
}

func TestExtractTrailingPunctuation_InLineComment(t *testing.T) {
	punct, rest, nextInTextBlock := extractTrailingPunctuation("$x = 1; // comment;", false)
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != "$x = 1" {
		t.Errorf("expected rest '$x = 1', got %q", rest)
	}
	if nextInTextBlock {
		t.Errorf("expected nextInTextBlock false")
	}
}

func TestExtractTrailingPunctuation_InBlockComment(t *testing.T) {
	punct, rest, nextInTextBlock := extractTrailingPunctuation("$x = 1; /* ; */", false)
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != "$x = 1" {
		t.Errorf("expected rest '$x = 1', got %q", rest)
	}
	if nextInTextBlock {
		t.Errorf("expected nextInTextBlock false")
	}
}

func TestExtractTrailingPunctuation_SemicolonInStringNotExtracted(t *testing.T) {
	punct, rest, nextInTextBlock := extractTrailingPunctuation(`System.out.println("a; b")`, false)
	if punct != "" {
		t.Errorf("expected punct '', got %q", punct)
	}
	if rest != `System.out.println("a; b")` {
		t.Errorf("expected rest 'System.out.println(\"a; b\")', got %q", rest)
	}
	if nextInTextBlock {
		t.Errorf("expected nextInTextBlock false")
	}
}

func TestExtractTrailingPunctuation_BraceInStringNotExtracted(t *testing.T) {
	punct, rest, nextInTextBlock := extractTrailingPunctuation(`System.out.println("}");`, false)
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != `System.out.println("}")` {
		t.Errorf("expected rest 'System.out.println(\"}\")', got %q", rest)
	}
	if nextInTextBlock {
		t.Errorf("expected nextInTextBlock false")
	}
}

func TestExtractTrailingPunctuation_MultiplePunctuation(t *testing.T) {
	punct, rest, nextInTextBlock := extractTrailingPunctuation("$x = 1; $y = 2;", false)
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != "$x = 1; $y = 2" {
		t.Errorf("expected rest '$x = 1; $y = 2', got %q", rest)
	}
	if nextInTextBlock {
		t.Errorf("expected nextInTextBlock false")
	}
}

func TestExtractTrailingPunctuation_PunctWithTrailingComment(t *testing.T) {
	punct, rest, nextInTextBlock := extractTrailingPunctuation("$x = 1; // trailing comment", false)
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != "$x = 1" {
		t.Errorf("expected rest '$x = 1', got %q", rest)
	}
	if nextInTextBlock {
		t.Errorf("expected nextInTextBlock false")
	}
}

func TestExtractTrailingPunctuation_TextBlockStart(t *testing.T) {
	punct, rest, nextInTextBlock := extractTrailingPunctuation(`String s = """`, false)
	if punct != "" {
		t.Errorf("expected punct '', got %q", punct)
	}
	if rest != `String s = """` {
		t.Errorf("expected rest 'String s = \"\"\"', got %q", rest)
	}
	if !nextInTextBlock {
		t.Errorf("expected nextInTextBlock true")
	}
}

func TestExtractTrailingPunctuation_TextBlockContinue(t *testing.T) {
	punct, rest, nextInTextBlock := extractTrailingPunctuation(`hello; world;`, true)
	if punct != "" {
		t.Errorf("expected punct '', got %q", punct)
	}
	if rest != `hello; world;` {
		t.Errorf("expected rest 'hello; world;', got %q", rest)
	}
	// Text block continues (no closing """)
	if !nextInTextBlock {
		t.Errorf("expected nextInTextBlock true")
	}
}

func TestExtractTrailingPunctuation_TextBlockEnd(t *testing.T) {
	punct, rest, nextInTextBlock := extractTrailingPunctuation(`hello;"""`, true)
	if punct != "" {
		t.Errorf("expected punct '', got %q", punct)
	}
	if rest != `hello;"""` {
		t.Errorf("expected rest 'hello;\"\"\"', got %q", rest)
	}
	if nextInTextBlock {
		t.Errorf("expected nextInTextBlock false")
	}
}

func TestExtractTrailingPunctuation_CharLiteral(t *testing.T) {
	input := `char c = ';';`
	punct, rest, nextInTextBlock := extractTrailingPunctuation(input, false)
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != `char c = ';'` {
		t.Errorf("expected rest 'char c = ';'', got %q", rest)
	}
	if nextInTextBlock {
		t.Errorf("expected nextInTextBlock false")
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
		{Type: tokenizer.Code, Text: "System.out.println("},
		{Type: tokenizer.String, Text: `"hello"`},
		{Type: tokenizer.Code, Text: ");\n"},
	}
	got := reconstructText(segments)
	if got != "System.out.println(\"hello\");\n" {
		t.Errorf("got %q, want %q", got, "System.out.println(\"hello\");\n")
	}
}
