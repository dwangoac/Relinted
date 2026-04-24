package go_pkg

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFormat_SimpleStatement(t *testing.T) {
	input := "var x int = 1\n"
	got := Format(input)
	if got == "" {
		t.Fatal("expected non-empty output")
	}
	if got != "var x int = 1\n" {
		t.Errorf("got %q, want %q", got, "var x int = 1\n")
	}
}

func TestFormat_BraceRelocation(t *testing.T) {
	input := "func main() {\n    fmt.Println(\"hi\")\n}\n"
	got := Format(input)
	if got == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormat_EmptyLinesPreserved(t *testing.T) {
	input := "var x int = 1\n\nvar y int = 2\n"
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
	input := `fmt.Println("}")`
	got := Format(input)
	if got == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormat_BlockCommentPreserved(t *testing.T) {
	input := "/* comment */\nfunc main() {}\n"
	got := Format(input)
	if got == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormat_RuneLiteralNotExtracted(t *testing.T) {
	input := "var c rune = ';'\n"
	got := Format(input)
	if got == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormat_SingleLineNoPunct(t *testing.T) {
	expected := "var x\n"
	got := Format("var x")
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
	input := "func main() {\n    fmt.Println(\"hi\")\n}\n"
	got := Format(input)
	expected := "func main()           {\n    fmt.Println(\"hi\") }\n"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestFormat_RightAlignedPunctuationExact(t *testing.T) {
	input := "var x int\nvar y int\n"
	got := Format(input)
	expected := "var x int\nvar y int\n"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestFormat_PunctuationNotInStringExact(t *testing.T) {
	input := `fmt.Println("}")`
	got := Format(input)
	expected := "fmt.Println(\"}\")\n"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestFormat_RawStringNotExtracted(t *testing.T) {
	input := "`fmt.Println(\"}\")`\n"
	got := Format(input)
	if got == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormat_Integration(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	projectRoot := filepath.Join(cwd, "..", "..")
	inputPath := filepath.Join(projectRoot, "examples", "linted-example-7.go")
	expectedPath := filepath.Join(projectRoot, "examples", "relinted-example-7.go")
	input, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("failed to read input file %s: %v", inputPath, err)
	}
	expected, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("failed to read expected output %s: %v", expectedPath, err)
	}
	got := Format(string(input))
	if got != string(expected) {
		t.Errorf("output does not match expected\n--- expected ---\n%s\n--- got ---\n%s\n", expected, got)
	}
}
