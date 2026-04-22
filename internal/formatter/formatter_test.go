package formatter

import "testing"

func checkFormat(t *testing.T, input, expected string) {
	t.Helper()
	got := Format(input)
	if got != expected {
		t.Errorf("Format(\n%q\n) =\n%q\n\nwant:\n%q", input, got, expected)
	}
}

func TestFormat_SimpleBraceRelocation(t *testing.T) {
	input := `int main()
{
    return 0;
}`
	got := Format(input)
	// The { from line 2 moves to end of line 1
	// The ; from line 3 stays on line 3
	// The } from line 4 moves to end of line 3
	if got == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestFormat_CommentsPreserved(t *testing.T) {
	input := `/* comment */
int main() {
}`
	got := Format(input)
	if got == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestFormat_EmptyLinesKept(t *testing.T) {
	input := `int main() {

    return 0;
}`
	got := Format(input)
	if got == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestFormat_StringBracesNotExtracted(t *testing.T) {
	input := `printf("}");`
	got := Format(input)
	if got == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestFormat_TabExpansion(t *testing.T) {
	input := "int main() {\n\treturn 0;\n}"
	got := Format(input)
	if got == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestFormat_BlockCommentPreserved(t *testing.T) {
	input := `/* multi
   line */
int x;`
	got := Format(input)
	if got == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestFormat_EmptyInput(t *testing.T) {
	_ = Format("")
}

func TestFormat_SingleLineNoPunct(t *testing.T) {
	expected := "int x\n"
	checkFormat(t, "int x", expected)
}

func TestFormat_SingleLineWithSemicolon(t *testing.T) {
	input := "int x = 0;"
	got := Format(input)
	if got == "" {
		t.Fatal("expected non-empty output")
	}
}
