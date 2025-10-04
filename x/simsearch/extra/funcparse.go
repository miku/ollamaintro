package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// FunctionInfo represents the extracted function metadata
type FunctionInfo struct {
	ID       int    `json:"id"`
	Text     string `json:"text"`
	Filename string `json:"filename,omitempty"`
	Line     int    `json:"line,omitempty"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]

	functions, err := extractFunctions(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(functions, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	// Write to file
	outputFile := "functions.json"
	err = ioutil.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Extracted %d functions to %s\n", len(functions), outputFile)
}

func extractFunctions(dir string) ([]FunctionInfo, error) {
	var functions []FunctionInfo
	idCounter := 1

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-Go files
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Parse the Go file
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not parse %s: %v\n", path, err)
			return nil
		}

		// Read file content for extracting function bodies
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		// Inspect AST for function declarations
		ast.Inspect(node, func(n ast.Node) bool {
			var funcDecl *ast.FuncDecl

			switch x := n.(type) {
			case *ast.FuncDecl:
				funcDecl = x
			}

			if funcDecl != nil {
				// Get position information
				pos := fset.Position(funcDecl.Pos())
				endPos := fset.Position(funcDecl.End())

				// Extract function text
				funcText := string(content[pos.Offset:endPos.Offset])

				funcInfo := FunctionInfo{
					ID:       idCounter,
					Text:     funcText,
					Filename: path,
					Line:     pos.Line,
				}

				functions = append(functions, funcInfo)
				idCounter++
			}

			return true
		})

		return nil
	})

	return functions, err
}
