// blame.go provides integration with Git history.
// It uses git blame and git log to identify the authors and dates of changes.
package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// gitBlame uses the 'git log' command to find the very first person who introduced
// the line containing the technical debt tag.
func gitBlame(filePath string, line int) (string, string) {
	absPath, _ := filepath.Abs(filePath)
	dir := filepath.Dir(absPath)

	rangeArg := fmt.Sprintf("%d,%d:%s", line, line, filepath.Base(absPath))
	cmd := exec.Command("git", "log", "-L", rangeArg, "--reverse", "--pretty=format:%an|%as", "--max-count=1", "--no-patch")
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fallbackBlame(filePath, line)
	}

	parts := strings.Split(strings.TrimSpace(string(output)), "|")
	if len(parts) >= 2 {
		return parts[0], parts[1]
	}

	return fallbackBlame(filePath, line)
}

// fallbackBlame is used when the deep history trace fails (e.g. for uncommitted local changes).
func fallbackBlame(filePath string, line int) (string, string) {
	absPath, _ := filepath.Abs(filePath)
	dir := filepath.Dir(absPath)

	cmd := exec.Command("git", "blame", "-L", fmt.Sprintf("%d,%d", line, line), "--porcelain", filepath.Base(absPath))
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "Unknown", ""
	}

	lines := strings.Split(string(output), "\n")
	author := "Unknown"
	date := ""

	for _, line := range lines {
		if strings.HasPrefix(line, "author ") {
			author = strings.TrimPrefix(line, "author ")
		}

		if strings.HasPrefix(line, "author-time ") {
			timestamp := strings.TrimPrefix(line, "author-time ")
			if ts, err := strconv.ParseInt(timestamp, 10, 64); err == nil {
				date = time.Unix(ts, 0).Format("2006-01-02")
			}
		}
	}
	return author, date
}
