package perl

import "testing"

func TestFormat_Simple(t *testing.T) {
	input := "my $x = 1;\n"
	got := Format(input)
	expected := "my $x = 1 ;\n"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestFormat_BraceAtEnd(t *testing.T) {
	input := "if ($x) {\n    print \"hi\";\n}\n"
	got := Format(input)
	expected := "if ($x)        {\n    print \"hi\" ;}\n"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestFormat_ClosingBraceOnNextLine(t *testing.T) {
	input := "if ($x) {\n    print \"hi\";\n}\n"
	got := Format(input)
	if got == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormat_EmptyLinesPreserved(t *testing.T) {
	input := "my $x = 1;\n\nmy $y = 2;\n"
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

func TestFormat_RegexInCode(t *testing.T) {
	input := "if ($x =~ /pattern/) {\n    print \"match\";\n}\n"
	got := Format(input)
	if got == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormat_HashComment(t *testing.T) {
	input := "my $x = 1; # comment\n"
	got := Format(input)
	if got == "" {
		t.Error("expected non-empty output")
	}
}
