package perl

import "testing"

func TestExtractTrailingPunctuation_Semicolon(t *testing.T) {
	punct, rest := extractTrailingPunctuation("my $x = 1;")
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != "my $x = 1" {
		t.Errorf("expected rest 'my $x = 1', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_Brace(t *testing.T) {
	punct, rest := extractTrailingPunctuation("if ($x) {")
	if punct != "{" {
		t.Errorf("expected punct '{', got %q", punct)
	}
	if rest != "if ($x) " {
		t.Errorf("expected rest 'if ($x) ', got %q", rest)
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
	punct, rest := extractTrailingPunctuation("my $x = 1")
	if punct != "" {
		t.Errorf("expected punct '', got %q", punct)
	}
	if rest != "my $x = 1" {
		t.Errorf("expected rest 'my $x = 1', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_InString(t *testing.T) {
	punct, rest := extractTrailingPunctuation(`print "hello;";`)
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != `print "hello;"` {
		t.Errorf("expected rest 'print \"hello;\"', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_InComment(t *testing.T) {
	punct, rest := extractTrailingPunctuation("my $x = 1; # comment;")
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != "my $x = 1" {
		t.Errorf("expected rest 'my $x = 1', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_InRegex(t *testing.T) {
	punct, rest := extractTrailingPunctuation(`if ($x =~ /;/) {`)
	if punct != "{" {
		t.Errorf("expected punct '{', got %q", punct)
	}
	if rest != `if ($x =~ /;/) ` {
		t.Errorf("expected rest 'if ($x =~ /;/) ', got %q", rest)
	}
}

func TestExtractTrailingPunctuation_MultiplePunct(t *testing.T) {
	punct, rest := extractTrailingPunctuation("my $x = 1; # semicolon in comment;")
	if punct != ";" {
		t.Errorf("expected punct ';', got %q", punct)
	}
	if rest != "my $x = 1" {
		t.Errorf("expected rest 'my $x = 1', got %q", rest)
	}
}
