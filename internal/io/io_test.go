package io

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.txt")
	expected := "hello world\n"
	if err := os.WriteFile(path, []byte(expected), 0644); err != nil {
		t.Fatal(err)
	}

	got, err := ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != expected {
		t.Errorf("ReadFile(%q) = %q, want %q", path, got, expected)
	}
}

func TestReadFile_NotFound(t *testing.T) {
	_, err := ReadFile("/nonexistent/file.txt")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestWriteFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "out.txt")
	content := "test content\n"

	if err := WriteFile(path, content); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != content {
		t.Errorf("WriteFile() = %q, want %q", string(data), content)
	}
}
