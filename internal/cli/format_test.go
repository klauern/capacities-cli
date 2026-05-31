package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/klauern/capacities-cli/internal/api"
)

func TestPrintJSON(t *testing.T) {
	var buf bytes.Buffer
	data := []string{"a", "b"}
	if err := printJSON(&buf, data); err != nil {
		t.Fatalf("printJSON error: %v", err)
	}
	var got []string
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Fatalf("unexpected JSON output: %v", got)
	}
}

func TestPrintSpaces(t *testing.T) {
	spaces := []api.Space{
		{ID: "s1", Title: "My Space", Icon: api.Icon{Val: "🚀"}},
		{ID: "s2", Title: "Work", Icon: api.Icon{Val: "💼"}},
	}

	t.Run("table", func(t *testing.T) {
		var buf bytes.Buffer
		if err := printSpaces(&buf, spaces, formatTable); err != nil {
			t.Fatalf("printSpaces error: %v", err)
		}
		out := buf.String()
		if !strings.Contains(out, "s1") || !strings.Contains(out, "My Space") {
			t.Errorf("table missing expected content: %q", out)
		}
		if !strings.Contains(out, "ID") || !strings.Contains(out, "TITLE") {
			t.Errorf("table missing header: %q", out)
		}
	})

	t.Run("json", func(t *testing.T) {
		var buf bytes.Buffer
		if err := printSpaces(&buf, spaces, formatJSON); err != nil {
			t.Fatalf("printSpaces error: %v", err)
		}
		var got []api.Space
		if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
			t.Fatalf("invalid JSON: %v\noutput: %s", err, buf.String())
		}
		if len(got) != 2 || got[0].ID != "s1" || got[1].Title != "Work" {
			t.Fatalf("unexpected JSON content: %+v", got)
		}
	})
}

func TestPrintStructures(t *testing.T) {
	structures := []api.Structure{
		{
			ID:    "str1",
			Title: "Notes",
			Collections: []api.Collection{
				{ID: "c1", Title: "Inbox"},
			},
		},
		{ID: "str2", Title: "Tasks"},
	}

	t.Run("table", func(t *testing.T) {
		var buf bytes.Buffer
		if err := printStructures(&buf, structures, formatTable); err != nil {
			t.Fatalf("printStructures error: %v", err)
		}
		out := buf.String()
		if !strings.Contains(out, "str1") || !strings.Contains(out, "Notes") {
			t.Errorf("table missing expected content: %q", out)
		}
		if !strings.Contains(out, "Inbox") {
			t.Errorf("table missing collection name: %q", out)
		}
	})

	t.Run("json", func(t *testing.T) {
		var buf bytes.Buffer
		if err := printStructures(&buf, structures, formatJSON); err != nil {
			t.Fatalf("printStructures error: %v", err)
		}
		var got []api.Structure
		if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
			t.Fatalf("invalid JSON: %v\noutput: %s", err, buf.String())
		}
		if len(got) != 2 || got[0].ID != "str1" {
			t.Fatalf("unexpected JSON content: %+v", got)
		}
	})
}

func TestPrintLookupResults(t *testing.T) {
	results := []api.LookupResult{
		{ID: "r1", StructureID: "str1", Title: "My Note"},
		{ID: "r2", StructureID: "str2", Title: "My Task"},
	}

	t.Run("table", func(t *testing.T) {
		var buf bytes.Buffer
		if err := printLookupResults(&buf, results, formatTable); err != nil {
			t.Fatalf("printLookupResults error: %v", err)
		}
		out := buf.String()
		if !strings.Contains(out, "r1") || !strings.Contains(out, "My Note") {
			t.Errorf("table missing expected content: %q", out)
		}
	})

	t.Run("json", func(t *testing.T) {
		var buf bytes.Buffer
		if err := printLookupResults(&buf, results, formatJSON); err != nil {
			t.Fatalf("printLookupResults error: %v", err)
		}
		var got []api.LookupResult
		if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
			t.Fatalf("invalid JSON: %v\noutput: %s", err, buf.String())
		}
		if len(got) != 2 || got[0].ID != "r1" || got[1].Title != "My Task" {
			t.Fatalf("unexpected JSON content: %+v", got)
		}
	})
}

func TestValidateFormat(t *testing.T) {
	tests := []struct {
		format  string
		wantErr bool
	}{
		{formatTable, false},
		{formatJSON, false},
		{"csv", true},
		{"", true},
		{"JSON", true},
		{"TABLE", true},
	}
	for _, tc := range tests {
		err := validateFormat(tc.format)
		if tc.wantErr && err == nil {
			t.Errorf("validateFormat(%q): expected error, got nil", tc.format)
		}
		if !tc.wantErr && err != nil {
			t.Errorf("validateFormat(%q): unexpected error: %v", tc.format, err)
		}
	}
}

func TestPrintSaveWebLinkResponse(t *testing.T) {
	resp := &api.SaveWebLinkResponse{
		SpaceID:     "sp1",
		ID:          "wl1",
		StructureID: "str1",
		Title:       "Go Blog",
		Description: "The Go Programming Language Blog",
		Tags:        []string{"go", "blog"},
	}

	t.Run("table", func(t *testing.T) {
		var buf bytes.Buffer
		if err := printSaveWebLinkResponse(&buf, resp, formatTable); err != nil {
			t.Fatalf("printSaveWebLinkResponse error: %v", err)
		}
		out := buf.String()
		if !strings.Contains(out, "wl1") || !strings.Contains(out, "Go Blog") {
			t.Errorf("table missing expected content: %q", out)
		}
		if !strings.Contains(out, "go") || !strings.Contains(out, "blog") {
			t.Errorf("table missing tags: %q", out)
		}
	})

	t.Run("table_no_tags", func(t *testing.T) {
		noTagResp := &api.SaveWebLinkResponse{ID: "wl2", Title: "Plain"}
		var buf bytes.Buffer
		if err := printSaveWebLinkResponse(&buf, noTagResp, formatTable); err != nil {
			t.Fatalf("printSaveWebLinkResponse error: %v", err)
		}
		out := buf.String()
		if !strings.Contains(out, "-") {
			t.Errorf("table should show '-' for empty tags: %q", out)
		}
	})

	t.Run("json", func(t *testing.T) {
		var buf bytes.Buffer
		if err := printSaveWebLinkResponse(&buf, resp, formatJSON); err != nil {
			t.Fatalf("printSaveWebLinkResponse error: %v", err)
		}
		var got api.SaveWebLinkResponse
		if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
			t.Fatalf("invalid JSON: %v\noutput: %s", err, buf.String())
		}
		if got.ID != "wl1" || got.Title != "Go Blog" || len(got.Tags) != 2 {
			t.Fatalf("unexpected JSON content: %+v", got)
		}
	})
}
