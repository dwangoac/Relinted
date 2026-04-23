package rust

import "testing"

func TestFormat_SimpleSemicolon(t *testing.T) {
	input := "let x = 1;\n"
	got := Format(input)
	if got == "" {
		t.Fatal("expected non-empty output")
	}
	// Semicolon should be extracted and right-aligned
	if got[len(got)-2:] != ";\n" {
		t.Errorf("expected output ending with ';\\n', got %q", got[len(got)-2:])
	}
}

func TestFormat_BraceRelocation(t *testing.T) {
	input := "fn main() {\n    println!(\"hi\");\n}\n"
	got := Format(input)
	if got == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormat_CommaInMatchArm(t *testing.T) {
	input := "match x {\n    Ok(n) => n,\n    Err(_) => continue,\n}\n"
	got := Format(input)
	if got == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormat_CommaNotInMatchArm(t *testing.T) {
	input := "let v = vec![1, 2, 3];\n"
	got := Format(input)
	if got == "" {
		t.Error("expected non-empty output")
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
	input := `println!("}");`
	got := Format(input)
	if got == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormat_BlockCommentPreserved(t *testing.T) {
	input := "/* comment */\nfn main() {}\n"
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
	input := "fn main() {\n    println!(\"hi\");\n}\n"
	got := Format(input)
	expected := "fn main()          {\n    println!(\"hi\") ;}\n"
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

func TestFormat_CommaMatchArmExact(t *testing.T) {
	input := "match x {\n    Ok(n) => n,\n    Err(_) => continue,\n}\n"
	got := Format(input)
	expected := "match x                {\n    Ok(n) => n         ,\n    Err(_) => continue ,}\n"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestFormat_PunctuationNotInStringExact(t *testing.T) {
	input := "println!(\"}\");\n"
	got := Format(input)
	expected := "println!(\"}\") ;\n"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}
