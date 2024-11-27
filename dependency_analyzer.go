package main

import (
	"bufio"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// DependencyAnalyzer handles dependency analysis for different file types
type DependencyAnalyzer struct {
	targetDir string
	fileMap   map[string]*FileContent // maps file paths to their content objects
}

// NewDependencyAnalyzer creates a new analyzer instance
func NewDependencyAnalyzer(targetDir string) *DependencyAnalyzer {
	return &DependencyAnalyzer{
		targetDir: targetDir,
		fileMap:   make(map[string]*FileContent),
	}
}

// RegisterFile adds a file to the dependency tracking system
func (da *DependencyAnalyzer) RegisterFile(path string, content *FileContent) {
	da.fileMap[path] = content
}

// AnalyzeDependencies analyzes dependencies for all registered files
func (da *DependencyAnalyzer) AnalyzeDependencies() error {
	for path, content := range da.fileMap {
		if err := da.analyzeFile(path, content); err != nil {
			return err
		}
	}
	return da.linkDependencies()
}

// analyzeFile determines the file type and calls the appropriate analyzer
func (da *DependencyAnalyzer) analyzeFile(path string, content *FileContent) error {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".go":
		return da.analyzeGoFile(path, content)
	case ".js", ".jsx", ".ts", ".tsx":
		return da.analyzeJSFile(path, content)
	case ".py":
		return da.analyzePythonFile(path, content)
	}
	return nil
}

// analyzeGoFile analyzes dependencies in Go files
func (da *DependencyAnalyzer) analyzeGoFile(path string, content *FileContent) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
	if err != nil {
		return err
	}

	content.Dependencies = &DependencyInfo{
		Imports: make([]ImportDependency, 0),
	}

	for _, imp := range f.Imports {
		// Remove quotes from import path
		importPath := strings.Trim(imp.Path.Value, "\"")

		depType := "standard"
		if strings.Contains(importPath, ".") || strings.Contains(importPath, "/") {
			depType = "external"
		}
		if strings.HasPrefix(importPath, da.targetDir) {
			depType = "local"
		}

		content.Dependencies.Imports = append(content.Dependencies.Imports, ImportDependency{
			Path: importPath,
			Type: depType,
		})
	}

	return nil
}

// analyzeJSFile analyzes dependencies in JavaScript/TypeScript files
func (da *DependencyAnalyzer) analyzeJSFile(path string, content *FileContent) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	content.Dependencies = &DependencyInfo{
		Imports: make([]ImportDependency, 0),
	}

	// Regular expressions for different import patterns
	importPatterns := []*regexp.Regexp{
		regexp.MustCompile(`import\s+.*\s+from\s+['"]([^'"]+)['"]`),
		regexp.MustCompile(`require\(['"]([^'"]+)['"]\)`),
		regexp.MustCompile(`import\s+['"]([^'"]+)['"]`),
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		for _, pattern := range importPatterns {
			matches := pattern.FindStringSubmatch(line)
			if len(matches) > 1 {
				importPath := matches[1]
				depType := "external"
				if strings.HasPrefix(importPath, ".") {
					depType = "local"
				} else if !strings.Contains(importPath, "/") {
					depType = "standard"
				}

				content.Dependencies.Imports = append(content.Dependencies.Imports, ImportDependency{
					Path: importPath,
					Type: depType,
				})
			}
		}
	}

	return scanner.Err()
}

// analyzePythonFile analyzes dependencies in Python files
func (da *DependencyAnalyzer) analyzePythonFile(path string, content *FileContent) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	content.Dependencies = &DependencyInfo{
		Imports: make([]ImportDependency, 0),
	}

	// Regular expressions for different import patterns
	importPatterns := []*regexp.Regexp{
		regexp.MustCompile(`^import\s+(\w+)`),
		regexp.MustCompile(`^from\s+([^\s]+)\s+import`),
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		for _, pattern := range importPatterns {
			matches := pattern.FindStringSubmatch(line)
			if len(matches) > 1 {
				importPath := matches[1]
				depType := "standard"
				if strings.Contains(importPath, ".") {
					depType = "local"
				}

				content.Dependencies.Imports = append(content.Dependencies.Imports, ImportDependency{
					Path: importPath,
					Type: depType,
				})
			}
		}
	}

	return scanner.Err()
}

// linkDependencies creates bidirectional relationships between files
func (da *DependencyAnalyzer) linkDependencies() error {
	// Reset all ImportedBy slices
	for _, content := range da.fileMap {
		if content.Dependencies != nil {
			content.Dependencies.ImportedBy = make([]ImportDependency, 0)
		}
	}

	// Build ImportedBy relationships
	for path, content := range da.fileMap {
		if content.Dependencies == nil {
			continue
		}

		for _, imp := range content.Dependencies.Imports {
			if imp.Type == "local" {
				// Convert import path to filesystem path
				importedPath := filepath.Join(da.targetDir, imp.Path)
				if importedContent, exists := da.fileMap[importedPath]; exists {
					if importedContent.Dependencies == nil {
						importedContent.Dependencies = &DependencyInfo{
							ImportedBy: make([]ImportDependency, 0),
						}
					}
					importedContent.Dependencies.ImportedBy = append(
						importedContent.Dependencies.ImportedBy,
						ImportDependency{
							Path: path,
							Type: "local",
						},
					)
				}
			}
		}
	}

	return nil
}
