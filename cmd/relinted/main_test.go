package main

import "testing"

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{"c extension", "foo.c", "c"},
		{"h extension", "foo.h", "c"},
		{"cpp extension", "foo.cpp", "c"},
		{"cc extension", "foo.cc", "c"},
		{"pl extension", "foo.pl", "perl"},
		{"pm extension", "foo.pm", "perl"},
		{"js extension", "foo.js", "js"},
		{"rs extension", "foo.rs", "rust"},
		{"uppercase C", "foo.C", "c"},
		{"unknown extension", "foo.unknown", "c"},
		{"path with directory", "/foo/bar.js", "js"},
		{"no extension", "Makefile", "c"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectLanguage(tt.path)
			if got != tt.expected {
				t.Errorf("detectLanguage(%q) = %q; want %q", tt.path, got, tt.expected)
			}
		})
	}
}
