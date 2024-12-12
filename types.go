package main

import (
	"encoding/xml"
)

// FileContent represents a file's content in XML format
type FileContent struct {
	XMLName      xml.Name        `xml:"file"`
	Name         string          `xml:"name,attr"`
	Size         int64           `xml:"size,attr"`
	Content      string          `xml:",cdata"`
	Dependencies *DependencyInfo `xml:"dependencies,omitempty"`
}

// DependencyInfo represents file dependencies
type DependencyInfo struct {
	Imports    []ImportDependency `xml:"imports,omitempty"`
	ImportedBy []ImportDependency `xml:"imported_by,omitempty"`
}

// ImportDependency represents a single import relationship
type ImportDependency struct {
	Path     string `xml:"path,attr"`
	Type     string `xml:"type,attr"` // local, external, standard
	Location string `xml:"location,attr,omitempty"`
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

// Stats tracks processing statistics
type Stats struct {
	filesProc int64
	bytesProc int64
	errors    int64
	tokens    int64 // Added field for token count
}

// Default exclusion patterns
var ExcludedPaths = []string{
	"node_modules",
	".git",
	".env",
	".DS_Store",
}
