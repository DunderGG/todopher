// main.go provides the entry point for ToDopher.
// It orchestrates the high-level execution flow from flag parsing to report generation.
package main

import (
	"flag"
	"fmt"
)

func main() {
	parseFlags()
	config := initializeConfig()

	if !IsQuiet {
		printIntro()
	}

	// Get search directory from arguments, default to current directory
	searchDir := "."
	if args := flag.Args(); len(args) > 0 {
		searchDir = args[0]
	}

	// Walk the project directory and collect files to scan
	filesToScan, err := discoverFiles(searchDir, config)
	if err != nil {
		fmt.Printf("Error walking the path: %v\n", err)
		return
	}

	filesToScan = filterOutputFiles(filesToScan)

	if !IsQuiet {
		fmt.Printf("Found %d files to scan. Commencing concurrent audit...\n", len(filesToScan))
	}

	// Concurrent Scanning Engine with Regex Matcher
	findings := startWorkerPool(filesToScan, config)

	if !IsQuiet {
		fmt.Printf("\nAudit complete! Total findings across all files: %d\n", len(findings))
		skipCount := BinarySkipCount.Load()
		if skipCount > 0 {
			fmt.Printf("⚠️  Skipped %d binary/non-UTF8 files.\n", skipCount)
		}
	}

	generateReports(findings, config)
}
