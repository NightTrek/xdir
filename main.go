package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func printHelp() {
	fmt.Println(`xdir - Directory to XML converter for AI context

Usage:
  xdir [flags] sourcedir [output.xml]
  xdir help    Show this help message

Flags:
  -patterns=<ext>   File extensions to include (comma-separated)
  -glob=<pattern>   Glob patterns to match files (comma-separated)
  -compress         Enable gzip compression for output
  -max-size=<bytes> Maximum file size in bytes (default: 10MB)
  -unsafe           Allow processing of normally excluded paths (node_modules or .env files etc)`)
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

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Error: No source directory specified")
		printHelp()
		os.Exit(1)
	}

	config.targetDir = args[0]
	if len(args) > 1 {
		config.outputFile = args[1]
	}

	if patternsStr != "" {
		rawPatterns := strings.Split(patternsStr, ",")
		config.filePatterns = make([]string, 0, len(rawPatterns))
		for _, p := range rawPatterns {
			if cleaned := cleanPattern(p); cleaned != "." {
				config.filePatterns = append(config.filePatterns, cleaned)
			}
		}
	}

	if globPatternsStr != "" {
		config.globPatterns = strings.Split(globPatternsStr, ",")
	}

	writer, cleanup, err := setupOutput(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up output: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	stats, err := processFiles(config, writer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing files: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nProcessing complete:\n")
	fmt.Printf("- Files processed: %d\n", stats.filesProc)
	fmt.Printf("- Total size: %.2f MB\n", float64(stats.bytesProc)/(1024*1024))
	fmt.Printf("- Errors: %d\n", stats.errors)
}
