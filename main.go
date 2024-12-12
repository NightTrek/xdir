package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func printHelp() {
	fmt.Println(`xdir - Directory to XML converter for AI context

Usage:
  xdir [flags] [sourcedir] [output.xml]
  xdir help    Show this help message

Flags:
  -patterns=<ext>   File extensions to include (comma-separated)
  -glob=<pattern>   Glob patterns to match files (comma-separated)
  -compress         Enable gzip compression for output
  -max-size=<bytes> Maximum file size in bytes (default: 10MB)
  -unsafe           Allow processing of normally excluded paths

Notes:
  - If no source directory is specified, the current directory is used
  - If no output file is specified, output.xml is used`)
}

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "help" || os.Args[1] == "-h" || os.Args[1] == "--help") {
		printHelp()
		os.Exit(0)
	}

	var config Config
	flag.Int64Var(&config.maxFileSize, "max-size", 10*1024*1024, "Maximum file size in bytes")
	flag.BoolVar(&config.compress, "compress", false, "Compress output with gzip")
	flag.BoolVar(&config.unsafeMode, "unsafe", false, "Allow processing of normally excluded paths")

	var patternsStr string
	var globPatternsStr string
	flag.StringVar(&patternsStr, "patterns", "", "File patterns to include (e.g., js,ts,md)")
	flag.StringVar(&globPatternsStr, "glob", "", "Glob patterns to include (e.g., '*.min.js')")

	flag.Parse()

	// Get current working directory as default
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// Set defaults and handle arguments
	config.targetDir = cwd           // Default to current directory
	config.outputFile = "output.xml" // Default output file

	// Override defaults with any provided arguments
	args := flag.Args()
	if len(args) > 0 {
		config.targetDir = args[0]
	}
	if len(args) > 1 {
		config.outputFile = args[1]
	}

	// Parse patterns
	if patternsStr != "" {
		rawPatterns := strings.Split(patternsStr, ",")
		config.filePatterns = make([]string, 0, len(rawPatterns))
		for _, p := range rawPatterns {
			if cleaned := cleanPattern(p); cleaned != "." {
				config.filePatterns = append(config.filePatterns, cleaned)
			}
		}
	}

	// Parse glob patterns
	if globPatternsStr != "" {
		config.globPatterns = strings.Split(globPatternsStr, ",")
	}

	writer, cleanup, err := setupOutput(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up output: %v\n", err)
		os.Exit(1)
	}

	stats, err := processFiles(config, writer)
	if err != nil {
		cleanup()
		fmt.Fprintf(os.Stderr, "Error processing files: %v\n", err)
		os.Exit(1)
	}

	// Close the file before renaming
	cleanup()

	// Rename the output file to include token count
	dir := filepath.Dir(config.outputFile)
	ext := filepath.Ext(config.outputFile)
	base := strings.TrimSuffix(filepath.Base(config.outputFile), ext)
	newName := fmt.Sprintf("%d-%s%s", stats.tokens, base, ext)
	newPath := filepath.Join(dir, newName)

	if err := os.Rename(config.outputFile, newPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error renaming output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nProcessing complete:\n")
	fmt.Printf("- Files processed: %d\n", stats.filesProc)
	fmt.Printf("- Total size: %.2f MB\n", float64(stats.bytesProc)/(1024*1024))
	fmt.Printf("- Total tokens: %d\n", stats.tokens)
	fmt.Printf("- Errors: %d\n", stats.errors)
	fmt.Printf("\nOutput written to: %s\n", newPath)
}
