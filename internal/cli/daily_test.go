package cli

import (
	"os"
	"strings"
	"testing"
)

func TestResolveDailySaveText_FromArgs(t *testing.T) {
	text, err := resolveDailySaveText([]string{"hello", "world"}, strings.NewReader("ignored"), true, "")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if text != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", text)
	}
}

func TestResolveDailySaveText_FileAndArgs_Error(t *testing.T) {
	_, err := resolveDailySaveText([]string{"hi"}, strings.NewReader(""), true, "note.md")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestResolveDailySaveText_FromFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := dir + "/note.md"
	if err := os.WriteFile(path, []byte("note text\n"), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	text, err := resolveDailySaveText(nil, strings.NewReader("ignored"), true, path)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if text != "note text" {
		t.Fatalf("expected %q, got %q", "note text", text)
	}
}

func TestResolveDailySaveText_FromStdin_NotTTY(t *testing.T) {
	text, err := resolveDailySaveText(nil, strings.NewReader("piped\n"), false, "")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if text != "piped" {
		t.Fatalf("expected %q, got %q", "piped", text)
	}
}

func TestResolveDailySaveText_FromFileDash_UsesStdin(t *testing.T) {
	text, err := resolveDailySaveText(nil, strings.NewReader("from stdin\n"), true, "-")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if text != "from stdin" {
		t.Fatalf("expected %q, got %q", "from stdin", text)
	}
}

func TestResolveDailySaveText_NoInput_Error(t *testing.T) {
	_, err := resolveDailySaveText(nil, strings.NewReader(""), true, "")
	if err == nil {
		t.Fatalf("expected error")
	}
}
