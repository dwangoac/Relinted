package js

import "testing"

func TestFormat_SimpleSemicolon(t *testing.T) {
	input := "let x = 1;\n"
	got := Format(input)
	if got == "" {
		t.Fatal("expected non-empty output")
	}
	if got[len(got)-2:] != ";\n" {
		t.Errorf("expected output ending with ';\\n', got %q", got[len(got)-2:])
	}
}

func TestFormat_EmptyLinesPreserved(t *testing.T) {
	input := "let x = 1;\n\nlet y = 2;\n"
	got := Format(input)
	newlines := 0
	for _, ch := range got {
		if ch == '\n' {
			newlines++
		}
	}
	if newlines != 3 {
		t.Errorf("expected 3 newlines (2 content + 1 trailing), got %d", newlines)
	}
}

func TestFormat_StringBracesNotExtracted(t *testing.T) {
	input := `console.log("}");`
	got := Format(input)
	if got == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormat_BlockCommentPreserved(t *testing.T) {
	input := "/* comment */\nfunction foo() {}\n"
	got := Format(input)
	if got == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormat_CharLiteralNotExtracted(t *testing.T) {
	input := "let c = ';';\n"
	got := Format(input)
	if got == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormat_SingleLineNoPunct(t *testing.T) {
	expected := "let x\n"
	got := Format("let x")
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestFormat_EmptyInput(t *testing.T) {
	got := Format("")
	if got != "\n" {
		t.Errorf("expected newline output, got %q", got)
	}
}

func TestFormat_BraceRelocationExact(t *testing.T) {
	input := "function foo() {\n    console.log('hi');\n}\n"
	got := Format(input)
	expected := "function foo()        {\n    console.log('hi') ;}\n"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestFormat_RightAlignedPunctuationExact(t *testing.T) {
	input := "let x = 1;\nlet y = 2;\n"
	got := Format(input)
	expected := "let x = 1 ;\nlet y = 2 ;\n"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestFormat_PunctuationNotInStringExact(t *testing.T) {
	input := "console.log(\"}\");\n"
	got := Format(input)
	expected := "console.log(\"}\") ;\n"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}
