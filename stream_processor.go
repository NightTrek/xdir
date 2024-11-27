package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultBufferSize = 32 * 1024   // 32KB chunks
	maxBufferSize     = 1024 * 1024 // 1MB max
)

// StreamProcessor handles streaming file processing
type StreamProcessor struct {
	bufferSize int
	writer     io.Writer
}

// NewStreamProcessor creates a new processor with specified buffer size
func NewStreamProcessor(bufferSize int) *StreamProcessor {
	if bufferSize <= 0 {
		bufferSize = defaultBufferSize
	}
	if bufferSize > maxBufferSize {
		bufferSize = maxBufferSize
	}
	return &StreamProcessor{bufferSize: bufferSize}
}

// ProcessDirectory processes an entire directory streaming to XML
func (p *StreamProcessor) ProcessDirectory(dir string, w io.Writer) error {
	p.writer = w

	// Write XML header
	if _, err := fmt.Fprintf(w, "%s<files>\n", xml.Header); err != nil {
		return err
	}

	// Process all files
	if err := filepath.Walk(dir, p.processPath); err != nil {
		return err
	}

	// Write XML footer
	_, err := fmt.Fprintf(w, "</files>\n")
	return err
}

// processPath handles each file/directory during the walk
func (p *StreamProcessor) processPath(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil // Skip directories
	}

	if isExcluded(path) {
		return nil // Skip excluded files
	}

	return p.processFile(path, info)
}

// processFile handles a single file
func (p *StreamProcessor) processFile(path string, info os.FileInfo) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Start file element
	relPath, err := filepath.Rel(".", path)
	if err != nil {
		relPath = path
	}

	fmt.Fprintf(p.writer, `  <file path="%s" size="%d">`+"\n", relPath, info.Size())

	// Process context (imports, etc) while streaming
	ctx, err := p.analyzeContext(file)
	if err != nil {
		return err
	}

	// Write context if any found
	if len(ctx.Imports) > 0 {
		fmt.Fprintf(p.writer, "    <imports>\n")
		for _, imp := range ctx.Imports {
			fmt.Fprintf(p.writer, `      <import path="%s" />`, imp)
		}
		fmt.Fprintf(p.writer, "    </imports>\n")
	}

	// Reset file pointer for content
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	// Stream file content
	fmt.Fprintf(p.writer, "    <content><![CDATA[\n")

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, p.bufferSize), p.bufferSize)

	for scanner.Scan() {
		fmt.Fprintf(p.writer, "%s\n", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	fmt.Fprintf(p.writer, "    ]]></content>\n")
	fmt.Fprintf(p.writer, "  </file>\n")

	return nil
}

// Context holds file analysis results
type Context struct {
	Imports []string
}

// analyzeContext performs streaming analysis of file context
func (p *StreamProcessor) analyzeContext(r io.Reader) (Context, error) {
	var ctx Context
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, p.bufferSize), p.bufferSize)

	for scanner.Scan() {
		line := scanner.Text()

		// Simple import detection - extend as needed
		if strings.HasPrefix(line, "import ") {
			imp := strings.Trim(strings.TrimPrefix(line, "import "), `"' `)
			ctx.Imports = append(ctx.Imports, imp)
		}
	}

	return ctx, scanner.Err()
}

// isExcluded checks if a path should be excluded
func isExcluded(path string) bool {
	excludedPaths := []string{
		"node_modules",
		".git",
		".env",
		".DS_Store",
	}

	for _, excluded := range excludedPaths {
		if strings.Contains(path, excluded) {
			return true
		}
	}

	return false
}
