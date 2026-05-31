package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/klauern/capacities-cli/internal/api"
)

const (
	formatTable = "table"
	formatJSON  = "json"
)

// validateFormat returns an error if format is not a recognized output format.
func validateFormat(format string) error {
	switch format {
	case formatTable, formatJSON:
		return nil
	default:
		return fmt.Errorf("unknown format %q: must be %q or %q", format, formatTable, formatJSON)
	}
}

func printJSON(out io.Writer, v any) error {
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func printTabTable(out io.Writer, header string, rows []string) {
	w := tabwriter.NewWriter(out, 0, 0, 3, ' ', 0)
	_, _ = fmt.Fprintln(w, header)
	for _, row := range rows {
		_, _ = fmt.Fprintln(w, row)
	}
	_ = w.Flush()
}

func printSpaces(out io.Writer, spaces []api.Space, format string) error {
	if format == formatJSON {
		return printJSON(out, spaces)
	}
	rows := make([]string, len(spaces))
	for i, s := range spaces {
		rows[i] = fmt.Sprintf("%s\t%s\t%s", s.ID, s.Title, s.Icon.Val)
	}
	printTabTable(out, "ID\tTITLE\tICON", rows)
	return nil
}

func printStructures(out io.Writer, structures []api.Structure, format string) error {
	if format == formatJSON {
		return printJSON(out, structures)
	}
	rows := make([]string, len(structures))
	for i, s := range structures {
		var collections []string
		for _, c := range s.Collections {
			collections = append(collections, c.Title)
		}
		cols := strings.Join(collections, ", ")
		if cols == "" {
			cols = "-"
		}
		rows[i] = fmt.Sprintf("%s\t%s\t%s", s.ID, s.Title, cols)
	}
	printTabTable(out, "STRUCTURE ID\tTITLE\tCOLLECTIONS", rows)
	return nil
}

func printLookupResults(out io.Writer, results []api.LookupResult, format string) error {
	if format == formatJSON {
		return printJSON(out, results)
	}
	rows := make([]string, len(results))
	for i, r := range results {
		rows[i] = fmt.Sprintf("%s\t%s\t%s", r.ID, r.StructureID, r.Title)
	}
	printTabTable(out, "ID\tSTRUCTURE ID\tTITLE", rows)
	return nil
}

func printSaveWebLinkResponse(out io.Writer, resp *api.SaveWebLinkResponse, format string) error {
	if format == formatJSON {
		return printJSON(out, resp)
	}
	tags := strings.Join(resp.Tags, ", ")
	if tags == "" {
		tags = "-"
	}
	rows := []string{fmt.Sprintf("%s\t%s\t%s\t%s", resp.ID, resp.Title, resp.Description, tags)}
	printTabTable(out, "ID\tTITLE\tDESCRIPTION\tTAGS", rows)
	return nil
}
