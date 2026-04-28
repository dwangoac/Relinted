package rust

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
	punct, rest := extractTrailingPunctuation("fn main() {")
	if punct != "{" {
		t.Errorf("expected punct '{', got %q", punct)
	}
	if rest != "fn main() " {
		t.Errorf("expected rest 'fn main() ', got %q", rest)
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
	punct, rest := extractTrailingPunctuation(`println!("hello;");`)
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != `println!("hello;")` {
		t.Errorf("expected rest 'println!(\"hello;\")', got %q", rest)
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

func TestExtractTrailingPunctuation_CommaInMatchArm(t *testing.T) {
	punct, rest := extractTrailingPunctuation("Ok(num) => num,")
	if punct != "," {
		t.Errorf("expected punct ',', got %q", punct)
	}
	if rest != "Ok(num) => num" {
		t.Errorf("expected rest 'Ok(num) => num', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_CommaInMatchArmContinue(t *testing.T) {
	punct, rest := extractTrailingPunctuation("Err(_) => continue,")
	if punct != "," {
		t.Errorf("expected punct ',', got %q", punct)
	}
	if rest != "Err(_) => continue" {
		t.Errorf("expected rest 'Err(_) => continue', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_CommaNotInMatchArm(t *testing.T) {
	punct, rest := extractTrailingPunctuation("vec![1, 2, 3]")
	if punct != "" {
		t.Errorf("expected punct '', got %q", punct)
	}
	if rest != "vec![1, 2, 3]" {
		t.Errorf("expected rest 'vec![1, 2, 3]', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_CommaWithSemicolon(t *testing.T) {
	// Semicolon is the last punctuation character on the line
	punct, rest := extractTrailingPunctuation("let x = 1; // ,")
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != "let x = 1" {
		t.Errorf("expected rest 'let x = 1', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_CommaInStringNotExtracted(t *testing.T) {
	punct, rest := extractTrailingPunctuation(`println!("a, b")`)
	if punct != "" {
		t.Errorf("expected punct '', got %q", punct)
	}
	if rest != `println!("a, b")` {
		t.Errorf("expected rest 'println!(\"a, b\")', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_CommaInCharNotExtracted(t *testing.T) {
	input := "let c = ','"
	punct, rest := extractTrailingPunctuation(input)
	if punct != "" {
		t.Errorf("expected punct '', got %q", punct)
	}
	if rest != input {
		t.Errorf("expected rest %q, got %q", input, rest)
	}
}

func TestExtractTrailingPunctuation_ArrowAfterComma(t *testing.T) {
	punct, rest := extractTrailingPunctuation("Ok(n) => n, _ => 0;")
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != "Ok(n) => n, _ => 0" {
		t.Errorf("expected rest 'Ok(n) => n, _ => 0', got %q", rest)
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
		{Type: tokenizer.Code, Text: "let "},
		{Type: tokenizer.Code, Text: "x = "},
		{Type: tokenizer.Code, Text: "1;\n"},
	}
	got := reconstructText(segments)
	if got != "let x = 1;\n" {
		t.Errorf("got %q, want %q", got, "let x = 1;\n")
	}
}
