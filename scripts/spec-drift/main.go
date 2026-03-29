package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

const defaultSpecURL = "https://api.capacities.io/openapi.json"

var validMethods = map[string]struct{}{
	"GET":     {},
	"POST":    {},
	"PUT":     {},
	"PATCH":   {},
	"DELETE":  {},
	"HEAD":    {},
	"OPTIONS": {},
	"TRACE":   {},
}

type openAPI struct {
	Paths map[string]map[string]json.RawMessage `json:"paths"`
}

func main() {
	specURL := flag.String("spec-url", defaultSpecURL, "OpenAPI spec URL")
	clientFile := flag.String("client-file", filepath.FromSlash("internal/api/client.go"), "Path to the API client source file")
	flag.Parse()

	specEndpoints, err := loadSpecEndpoints(*specURL)
	if err != nil {
		fatalf("load OpenAPI spec: %v", err)
	}

	clientEndpoints, err := loadClientEndpoints(*clientFile)
	if err != nil {
		fatalf("load client endpoints: %v", err)
	}

	missingFromClient := diffSorted(specEndpoints, clientEndpoints)
	extraInClient := diffSorted(clientEndpoints, specEndpoints)

	if len(missingFromClient) == 0 && len(extraInClient) == 0 {
		fmt.Printf("No endpoint drift detected between %s and the live OpenAPI spec.\n", *clientFile)
		return
	}

	fmt.Printf("Endpoint drift detected between %s and %s\n\n", *clientFile, *specURL)
	if len(missingFromClient) > 0 {
		fmt.Println("In spec, missing from client:")
		for _, endpoint := range missingFromClient {
			fmt.Printf("  %s\n", endpoint)
		}
		fmt.Println()
	}
	if len(extraInClient) > 0 {
		fmt.Println("In client, missing from spec:")
		for _, endpoint := range extraInClient {
			fmt.Printf("  %s\n", endpoint)
		}
		fmt.Println()
	}

	os.Exit(1)
}

func loadSpecEndpoints(specURL string) ([]string, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(http.MethodGet, specURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var spec openAPI
	if err := json.NewDecoder(resp.Body).Decode(&spec); err != nil {
		return nil, err
	}

	var endpoints []string
	for path, operations := range spec.Paths {
		for method := range operations {
			method = strings.ToUpper(method)
			if _, ok := validMethods[method]; !ok {
				continue
			}
			endpoints = append(endpoints, method+" "+path)
		}
	}

	slices.Sort(endpoints)
	return slices.Compact(endpoints), nil
}

func loadClientEndpoints(filename string) ([]string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return nil, err
	}

	var endpoints []string
	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok || len(call.Args) < 3 {
			return true
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel == nil || sel.Sel.Name != "doJSON" {
			return true
		}

		method, ok := methodLiteral(call.Args[1])
		if !ok {
			return true
		}

		path, ok := stringLiteral(call.Args[2])
		if !ok {
			return true
		}

		endpoints = append(endpoints, method+" "+path)
		return true
	})

	slices.Sort(endpoints)
	return slices.Compact(endpoints), nil
}

func diffSorted(left, right []string) []string {
	var diff []string
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		switch {
		case left[i] == right[j]:
			i++
			j++
		case left[i] < right[j]:
			diff = append(diff, left[i])
			i++
		default:
			j++
		}
	}

	diff = append(diff, left[i:]...)
	return diff
}

func methodLiteral(expr ast.Expr) (string, bool) {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return "", false
	}

	pkg, ok := sel.X.(*ast.Ident)
	if !ok || pkg.Name != "http" {
		return "", false
	}

	method := strings.ToUpper(strings.TrimPrefix(sel.Sel.Name, "Method"))
	if _, ok := validMethods[method]; !ok {
		return "", false
	}

	return method, true
}

func stringLiteral(expr ast.Expr) (string, bool) {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", false
	}

	value, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", false
	}

	return value, true
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
