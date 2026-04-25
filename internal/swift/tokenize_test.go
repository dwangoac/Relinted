package swift

import (
	"testing"

	"relinted/internal/tokenizer"
)

func checkSegments(t *testing.T, input string, expected []tokenizer.Segment) {
	t.Helper()
	got := Tokenize(input)
	if len(got) != len(expected) {
		t.Fatalf("got %d segments, want %d", len(got), len(expected))
	}
	for i := range got {
		if got[i].Type != expected[i].Type {
			t.Errorf("segment %d: type %v, want %v", i, got[i].Type, expected[i].Type)
		}
		if got[i].Text != expected[i].Text {
			t.Errorf("segment %d: text %q, want %q", i, got[i].Text, expected[i].Text)
		}
	}
}

func TestTokenize_CodeOnly(t *testing.T) {
	checkSegments(t, "let x = 0", []tokenizer.Segment{
		{tokenizer.Code, "let x = 0"},
	})
}

func TestTokenize_OptionalChaining(t *testing.T) {
	checkSegments(t, "foo?.bar", []tokenizer.Segment{
		{tokenizer.Code, "foo?.bar"},
	})
}

func TestTokenize_NilCoalescing(t *testing.T) {
	checkSegments(t, "foo ?? bar", []tokenizer.Segment{
		{tokenizer.Code, "foo ?? bar"},
	})
}

func TestTokenize_Attribute(t *testing.T) {
	checkSegments(t, "@objc func hello()", []tokenizer.Segment{
		{tokenizer.Code, "@objc func hello()"},
	})
}

func TestTokenize_AttributeMultiple(t *testing.T) {
	checkSegments(t, "@escaping @MainActor func()", []tokenizer.Segment{
		{tokenizer.Code, "@escaping @MainActor func()"},
	})
}

func TestTokenize_StringLiteral(t *testing.T) {
	checkSegments(t, "let s = \"hello\"", []tokenizer.Segment{
		{tokenizer.Code, "let s = "},
		{tokenizer.String, "\"hello\""},
	})
}

func TestTokenize_StringWithEscape(t *testing.T) {
	checkSegments(t, "\"hello\\nworld\"", []tokenizer.Segment{
		{tokenizer.String, "\"hello\\nworld\""},
	})
}

func TestTokenize_StringInterpolationNotParsed(t *testing.T) {
	// String interpolation \(...) is inside a string literal, so the \ is
	// consumed by escape handling — the whole thing stays as one String token.
	checkSegments(t, "\"\\(foo)\"", []tokenizer.Segment{
		{tokenizer.String, "\"\\(foo)\""},
	})
}

func TestTokenize_CharLiteral(t *testing.T) {
	checkSegments(t, "let c = 'a'", []tokenizer.Segment{
		{tokenizer.Code, "let c = "},
		{tokenizer.Char, "'a'"},
	})
}

func TestTokenize_LineComment(t *testing.T) {
	checkSegments(t, "// comment\nfoo", []tokenizer.Segment{
		{tokenizer.CommentLine, "// comment\n"},
		{tokenizer.Code, "foo"},
	})
}

func TestTokenize_BlockComment(t *testing.T) {
	checkSegments(t, "/* block */", []tokenizer.Segment{
		{tokenizer.CommentBlock, "/* block */"},
	})
}

func TestTokenize_Mixed(t *testing.T) {
	checkSegments(t, "@objc let x = \"hi\" // end\n", []tokenizer.Segment{
		{tokenizer.Code, "@objc let x = "},
		{tokenizer.String, "\"hi\""},
		{tokenizer.Code, " "},
		{tokenizer.CommentLine, "// end\n"},
	})
}

func TestTokenize_Empty(t *testing.T) {
	got := Tokenize("")
	if len(got) != 0 {
		t.Errorf("expected 0 segments, got %d", len(got))
	}
}

func TestTokenize_OptionalChainingWithNilCoalescing(t *testing.T) {
	checkSegments(t, "foo?.bar ?? baz", []tokenizer.Segment{
		{tokenizer.Code, "foo?.bar ?? baz"},
	})
}

func TestTokenize_AttributeWithOptional(t *testing.T) {
	checkSegments(t, "@objc var x: String?", []tokenizer.Segment{
		{tokenizer.Code, "@objc var x: String?"},
	})
}

func TestTokenize_MultipleAttributes(t *testing.T) {
	checkSegments(t, "@State var count = 0", []tokenizer.Segment{
		{tokenizer.Code, "@State var count = 0"},
	})
}

func TestTokenize_CharWithEscape(t *testing.T) {
	checkSegments(t, "'\\n'", []tokenizer.Segment{
		{tokenizer.Char, "'\\n'"},
	})
}

func TestTokenize_CommentsAroundCode(t *testing.T) {
	checkSegments(t, "/* start */ let x = 5 // end\n", []tokenizer.Segment{
		{tokenizer.CommentBlock, "/* start */"},
		{tokenizer.Code, " let x = 5 "},
		{tokenizer.CommentLine, "// end\n"},
	})
}

func TestTokenize_JustQuestionMark(t *testing.T) {
	// A standalone ? in code stays as Code (not punctuation)
	checkSegments(t, "?", []tokenizer.Segment{
		{tokenizer.Code, "?"},
	})
}

func TestTokenize_DoubleQuestionMark(t *testing.T) {
	checkSegments(t, "??", []tokenizer.Segment{
		{tokenizer.Code, "??"},
	})
}

func TestTokenize_AttributeAtStart(t *testing.T) {
	checkSegments(t, "@property var name: String", []tokenizer.Segment{
		{tokenizer.Code, "@property var name: String"},
	})
}
