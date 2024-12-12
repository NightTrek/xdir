package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// isExcludedPath checks if the path should be excluded
func isExcludedPath(path string, config Config) bool {
	if config.unsafeMode {
		return false
	}

	parts := strings.Split(filepath.ToSlash(path), "/")
	for _, part := range parts {
		for _, excluded := range ExcludedPaths {
			if part == excluded {
				fmt.Printf("Skipping excluded path: %s (matched %s)\n", path, excluded)
				return true
			}
		}
	}

	base := filepath.Base(path)
	if strings.HasPrefix(base, ".") && base != "." && base != ".." {
		fmt.Printf("Skipping hidden file/directory: %s\n", path)
		return true
	}

	return false
}

// processFiles walks through the directory and processes files
func processFiles(config Config, writer io.Writer) (Stats, error) {
	var stats Stats
	absTargetDir, err := filepath.Abs(config.targetDir)
	if err != nil {
		return stats, fmt.Errorf("error resolving target directory: %v", err)
	}

	fmt.Printf("Processing directory: %s\n", absTargetDir)

	// Create dependency analyzer and token counter
	analyzer := NewDependencyAnalyzer(absTargetDir)
	tokenCounter := NewTokenCounter()
	fileContents := make(map[string]*FileContent)

	// First pass: Collect all files and their contents
	err = filepath.Walk(absTargetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			stats.errors++
			return nil
		}

		if info.IsDir() {
			if isExcludedPath(path, config) {
				fmt.Printf("Skipping excluded directory: %s\n", path)
				return filepath.SkipDir
			}
			return nil
		}

		if isExcludedPath(path, config) {
			return nil
		}

		if !isFileMatch(path, config) {
			return nil
		}

		// Check file size before processing
		if config.maxFileSize > 0 && info.Size() > config.maxFileSize {
			fmt.Printf("Skipping file exceeding size limit: %s (%d bytes)\n", path, info.Size())
			stats.errors++
			return nil
		}

		relPath, err := filepath.Rel(absTargetDir, path)
		if err != nil {
			fmt.Printf("Error getting relative path for %s: %v\n", path, err)
			stats.errors++
			return nil
		}

		fmt.Printf("Processing: %s\n", relPath)

		// Read file content
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("Error opening file %s: %v\n", path, err)
			stats.errors++
			return nil
		}
		defer file.Close()

		// Use a buffer with reasonable size
		buf := bytes.NewBuffer(make([]byte, 0, 32*1024))
		_, err = io.Copy(buf, file)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", path, err)
			stats.errors++
			return nil
		}

		// Create file content
		content := &FileContent{
			Name:    relPath,
			Size:    info.Size(),
			Content: buf.String(),
		}

		// Register file with dependency analyzer
		analyzer.RegisterFile(path, content)
		fileContents[path] = content

		stats.filesProc++
		stats.bytesProc += info.Size()
		fmt.Printf("Processed: %s (%.2f KB)\n", relPath, float64(info.Size())/1024)
		return nil
	})

	if err != nil {
		return stats, err
	}

	// Analyze dependencies
	fmt.Println("Analyzing dependencies...")
	if err := analyzer.AnalyzeDependencies(); err != nil {
		fmt.Printf("Error analyzing dependencies: %v\n", err)
		stats.errors++
	}

	// Create a buffer to hold the XML output for token counting
	xmlBuffer := &bytes.Buffer{}

	// Write XML output to buffer
	xmlBuffer.WriteString(xml.Header)
	xmlBuffer.WriteString("<files>\n")

	// Write each file with its dependencies
	encoder := xml.NewEncoder(xmlBuffer)
	encoder.Indent("", "  ")

	for _, content := range fileContents {
		if err := encoder.Encode(content); err != nil {
			fmt.Printf("Error encoding %s: %v\n", content.Name, err)
			stats.errors++
			continue
		}
	}

	xmlBuffer.WriteString("</files>\n")

	// Count tokens in the complete XML output
	xmlContent := xmlBuffer.String()
	tokenCount := tokenCounter.CountTokens(xmlContent)
	stats.tokens = int64(tokenCount)
	fmt.Printf("Token count: %d\n", tokenCount)

	// Write the XML to the actual output
	if _, err := xmlBuffer.WriteTo(writer); err != nil {
		return stats, fmt.Errorf("error writing XML output: %v", err)
	}

	return stats, nil
}
