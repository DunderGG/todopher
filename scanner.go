// scanner.go contains the core scanning engine logic.
// It manages the worker pool, file traversal, regex extraction, and binary file detection.
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// discoverFiles traverses the file tree starting at searchDir and returns a list of files matching scan criteria.
func discoverFiles(searchDir string, config Config) ([]string, error) {
	var filesToScan []string
	err := filepath.WalkDir(searchDir, func(path string, dir os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Only consider files that are not in ignored folders and have allowed extensions
		if !dir.IsDir() && shouldScan(path, config) {
			filesToScan = append(filesToScan, path)
		}
		return nil
	})
	return filesToScan, err
}

// shouldScan evaluates if a file should be audited based on its path and the provided configuration.
func shouldScan(path string, config Config) bool {
	// Skip ignored folders
	for _, folder := range config.IgnoreFolders {
		if strings.Contains(path, string(os.PathSeparator)+folder+string(os.PathSeparator)) {
			return false
		}
	}

	// Filter by extension
	ext := filepath.Ext(path)
	for _, allowed := range config.AllowedExtensions {
		if ext == allowed {
			return true
		}
	}
	return false
}

// filterOutputFiles removes the report file from the list of files to scan to prevent self-scanning.
func filterOutputFiles(files []string) []string {
	// Filter out the report file itself to avoid self-scanning
	var filteredFiles []string
	for _, file := range files {
		base := filepath.Base(file)
		if base != OutputPath {
			filteredFiles = append(filteredFiles, file)
		}
	}
	return filteredFiles
}

// startWorkerPool initializes a pool of worker goroutines to scan files concurrently.
func startWorkerPool(filesToScan []string, config Config) []Finding {
	numWorkers := config.WorkerCount
	totalFiles := len(filesToScan)

	// Create channels for job distribution and result collection
	fileJobs := make(chan string, totalFiles)
	results := make(chan []Finding, totalFiles)
	var waitGroup sync.WaitGroup

	// Compile the case-insensitive regex once for performance using the global pattern.
	pattern := fmt.Sprintf(TagRegexPattern, strings.Join(config.SearchTags, "|"))
	regex := regexp.MustCompile(pattern)

	// Spawn workerIndex goroutines
	for workerIndex := 1; workerIndex <= numWorkers; workerIndex++ {
		// Increment the wait group counter which tracks active workers so we can wait for them to finish
		waitGroup.Add(1)

		// Each worker will read from the jobs channel, process the file, and send results back
		go func() {
			defer waitGroup.Done()
			for path := range fileJobs {
				findings := scanFile(path, regex)
				results <- findings
			}
		}()
	}

	// Send files to the worker pool through the jobs channel.
	for _, path := range filesToScan {
		fileJobs <- path
	}
	close(fileJobs)

	// Result collector and Progress Bar logic
	var allFindings []Finding
	count := 0
	done := make(chan bool)

	// This goroutine collects results from the workers and updates the progress bar.
	go func() {
		for findings := range results {
			allFindings = append(allFindings, findings...)
			count++
			if !IsQuiet {
				printProgressBar(count, totalFiles)
			}
		}
		done <- true
	}()

	// Wait for all workers to finish and close results channel
	waitGroup.Wait()
	close(results)

	// Wait for result collector to finish aggregating all data
	<-done

	return allFindings
}

// scanFile reads a file and identifies lines containing technical debt tags.
func scanFile(filePath string, regex *regexp.Regexp) []Finding {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filePath, err)
		return nil
	}
	defer file.Close()

	// Check if the file is binary before processing
	isBinary, err := isBinaryFile(file)
	if err != nil {
		fmt.Printf("Error checking if file is binary %s: %v\n", filePath, err)
		return nil
	}
	if isBinary {
		BinarySkipCount.Add(1)
		return nil
	}

	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		return nil
	}

	lines := strings.Split(string(content), "\n")
	var findings []Finding

	for lineCounter, line := range lines {
		matches := regex.FindStringSubmatch(line)
		if len(matches) > 0 {
			lineNum := lineCounter + 1
			// Determine context (up to 3 lines after the tag)
			var contextLines []string
			for contextCounter := 1; contextCounter <= 3 && lineCounter+contextCounter < len(lines); contextCounter++ {
				contextLines = append(contextLines, strings.TrimRight(lines[lineCounter+contextCounter], "\r"))
			}

			author := matches[2]
			when := ""

			// If no author extracted from tag, try git blame
			if author == "" {
				blameAuthor, blameDate := gitBlame(filePath, lineNum)
				if blameAuthor != "" {
					author = blameAuthor
				}
				when = blameDate
			}

			// Calculate Status based on age
			status := "Fresh"
			if when != "" {
				if commitTime, err := stringsToDate(when); err == nil {
					daysOld := timeSinceInDays(commitTime)
					if daysOld > 180 {
						status = "Ancient"
					} else if daysOld > 30 {
						status = "Stale"
					}
				}
			}

			finding := Finding{
				File:    filePath,
				Line:    lineNum,
				Tag:     strings.ToUpper(matches[1]),
				Author:  author,
				Content: strings.TrimSpace(matches[3]),
				Context: contextLines,
				When:    when,
				Status:  status,
			}
			findings = append(findings, finding)
		}
	}
	return findings
}

// isBinaryFile determines if a file is binary by checking for null bytes (\0) in the first 8KB of content.
func isBinaryFile(file *os.File) (bool, error) {
	// Reset file pointer to the beginning to ensure we read from the start
	_, err := file.Seek(0, 0)
	if err != nil {
		return false, err
	}

	// Read the first 8KB to check for binary content
	buffer := make([]byte, 8192)
	numBytesRead, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false, err
	}

	// Check for null bytes which usually indicate binary files
	for i := 0; i < numBytesRead; i++ {
		if buffer[i] == 0 {
			// Found a null byte, assume it's binary
			return true, nil
		}
	}

	// Reset file pointer to read the full content later
	_, err = file.Seek(0, 0)
	if err != nil {
		return false, err
	}

	return false, nil
}
