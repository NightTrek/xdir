package main

import (
	"encoding/xml"
)

// FileContent represents a file's content in XML format
type FileContent struct {
	XMLName xml.Name `xml:"file"`
	Name    string   `xml:"name,attr"`
	Size    int64    `xml:"size,attr"`
	Content string   `xml:",cdata"`
}

// Config holds the application configuration
type Config struct {
	maxFileSize  int64
	compress     bool
	globPatterns []string
	filePatterns []string
	targetDir    string
	outputFile   string
	unsafeMode   bool
}

// Default exclusion patterns
var ExcludedPaths = []string{
	"node_modules",
	".git",
	".env",
	".DS_Store",
}
