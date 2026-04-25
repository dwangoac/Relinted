package java

import (
	"os"
	"testing"
)

func TestFormat(t *testing.T) {
	input, err := os.ReadFile("../../examples/linted-example-8.java")
	if err != nil {
		t.Fatalf("Failed to read input: %v", err)
	}
	want, err := os.ReadFile("../../examples/relinted-example-8.java")
	if err != nil {
		t.Fatalf("Failed to read want: %v", err)
	}
	got := Format(string(input))
	if got != string(want) {
		t.Errorf("Format mismatch\ngot:\n%s\nwant:\n%s", got, want)
	}
}
