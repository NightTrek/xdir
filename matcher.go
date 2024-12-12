package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

// cleanPattern ensures pattern starts with dot and is lowercase
func cleanPattern(pattern string) string {
	pattern = strings.TrimSpace(pattern)
	pattern = strings.ToLower(pattern)
	if !strings.HasPrefix(pattern, ".") {
		pattern = "." + pattern
	}
	return pattern
}

// isFileMatch checks if the file matches any of the patterns
func isFileMatch(path string, config Config) bool {
	ext := cleanPattern(filepath.Ext(path))
	fmt.Printf("Checking file %s with extension %s\n", path, ext)

	// Check file extensions first
	if len(config.filePatterns) > 0 {
		fmt.Printf("Checking against file patterns: %v\n", config.filePatterns)
		for _, pattern := range config.filePatterns {
			cleanedPattern := cleanPattern(pattern)
			if cleanedPattern == ext {
				fmt.Printf("File matched pattern %s\n", cleanedPattern)
				return true
			}
		}
		fmt.Printf("No file pattern match found\n")
		return false
	}

	// Then check glob patterns if specified
	if len(config.globPatterns) > 0 {
		fmt.Printf("Checking against glob patterns: %v\n", config.globPatterns)
		for _, pattern := range config.globPatterns {
			matched, err := filepath.Match(pattern, filepath.Base(path))
			if err == nil && matched {
				fmt.Printf("File matched glob pattern %s\n", pattern)
				return true
			}
		}
		fmt.Printf("No glob pattern match found\n")
		return false
	}

	// If no patterns specified, use default extensions
	fmt.Printf("Using default patterns\n")
	defaultPatterns := []string{".c", ".cpp", ".css", ".go", ".h",
		".hpp", ".html", ".java", ".js", ".json", ".jsx", ".md", ".mdx", ".php",
		".py", ".rb", ".rs", ".sql", ".swift", ".ts", ".tsx", ".txt", ".xml",
		".yaml", ".yml"}

	for _, pattern := range defaultPatterns {
		if cleanPattern(pattern) == ext {
			fmt.Printf("File matched default pattern %s\n", pattern)
			return true
		}
	}

	fmt.Printf("No default pattern match found\n")
	return false
}
