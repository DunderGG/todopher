// types.go defines the core data structures and global variables used throughout ToDopher.
// It includes the Finding and Config structs, as well as the embedded assets filesystem.
package main

import (
	"embed"
	"sync/atomic"
)

//go:embed assets/*
var content embed.FS

const (
	// DefaultOutputPath is the filename for the generated technical debt report.
	DefaultOutputPath = "report.html"
	// TagRegexPattern is core regex for identifying technical debt tags.
	// 		Prefix: Supports // (C++/Go), # (Python/Shell), /* (Block Start), or * (Block Middle)
	// 		Group 1: Placeholder for search tags (e.g., TODO, FIXME)
	//	 	Group 2: An optional author name following a hyphen (e.g., TODO-Dunder)
	// 		Group 3: A colon followed by any amount of whitespace
	// 		Group 4: The remaining text on the line as the description/content
	TagRegexPattern = `(?i)(?://|#|/\*|\*)\s*(%s)(?:-([A-Z]+))?:\s*(.*)`
)

var (
	// IsQuiet mode flag for suppressing non-essential output
	IsQuiet bool
	// OutputPath is the destination for the HTML report
	OutputPath string
	// BinarySkipCount tracks how many files were identified as binary and skipped
	BinarySkipCount atomic.Int32
	// CustomTags is a comma-separated list of search tags
	CustomTags string
	// CustomExtensions is a comma-separated list of file extensions
	CustomExtensions string
	// CustomIgnore is a comma-separated list of folders to ignore
	CustomIgnore string
	// JsonPath is the destination for the JSON export
	JsonPath string
	// CsvPath is the destination for the CSV export
	CsvPath string
	// MdPath is the destination for the Markdown export
	MdPath string
	// WorkerCount is the number of concurrent worker goroutines
	WorkerCount int
)

// Config holds the scanner settings
type Config struct {
	SearchTags        []string
	IgnoreFolders     []string
	AllowedExtensions []string
	WorkerCount       int
}

// Finding represents a single technical debt entry found in the source code.
// The 'json:"file"' parts tell the JSON encoder how to name these fields when converting to JSON, which is useful for the API endpoint.
type Finding struct {
	File    string   `json:"file"`
	Line    int      `json:"line"`
	Tag     string   `json:"tag"`
	Author  string   `json:"author"`
	Content string   `json:"content"`
	Context []string `json:"context"`
	When    string   `json:"when"`
	Status  string   `json:"status"` // "Fresh", "Stale", or "Ancient"
}
