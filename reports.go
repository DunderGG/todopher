// reports.go manages the generation and export of audit findings.
// It handles HTML dashboard creation, asset extraction, and JSON/CSV/Markdown exports.
package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// generateReports coordinates the generation of all requested audit reports.
func generateReports(findings []Finding, config Config) {
	// Generate static HTML report
	err := generateHtmlReport(findings, config, OutputPath)
	if err != nil {
		fmt.Printf("Error generating report: %v\n", err)
	}

	// Export to JSON if requested
	if JsonPath != "" {
		exportToJson(findings)
	}

	// Export to CSV if requested
	if CsvPath != "" {
		exportToCsv(findings)
	}

	// Export to Markdown if requested
	if MdPath != "" {
		exportToMarkdown(findings)
	}

	if absPath, err := filepath.Abs(OutputPath); err == nil {
		if !IsQuiet {
			fmt.Printf("📊 Report generated successfully at:\n %s\n", absPath)
		}
	} else {
		if !IsQuiet {
			fmt.Println("📊 Report generated successfully at:\n " + OutputPath)
		}
	}
}

// exportToJson serializes the findings into a pretty-printed JSON file.
func exportToJson(findings []Finding) {
	jsonData, err := json.MarshalIndent(findings, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	err = os.WriteFile(JsonPath, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing JSON file: %v\n", err)
	} else if !IsQuiet {
		if absJson, err := filepath.Abs(JsonPath); err == nil {
			fmt.Printf("📄 JSON data exported to:\n %s\n", absJson)
		} else {
			fmt.Printf("📄 JSON data exported to:\n %s\n", JsonPath)
		}
	}
}

// exportToMarkdown serializes the findings into a clean Markdown table.
func exportToMarkdown(findings []Finding) {
	file, err := os.Create(MdPath)
	if err != nil {
		fmt.Printf("Error creating Markdown file: %v\n", err)
		return
	}
	defer file.Close()

	var builder strings.Builder
	builder.WriteString("# ToDopher Audit Report\n\n")
	builder.WriteString(fmt.Sprintf("Total Issues Found: **%d**\n\n", len(findings)))

	if len(findings) > 0 {
		builder.WriteString("| File | Line | Tag | Author | Content | Date | Status |\n")
		builder.WriteString("| --- | --- | --- | --- | --- | --- | --- |\n")

		for _, finding := range findings {
			content := strings.ReplaceAll(finding.Content, "|", "\\|")
			content = strings.ReplaceAll(content, "\n", " ")

			row := fmt.Sprintf("| %s | %d | **%s** | %s | %s | %s | %s |\n",
				finding.File, finding.Line, finding.Tag, finding.Author, content, finding.When, finding.Status)
			builder.WriteString(row)
		}
	} else {
		builder.WriteString("*No issues found.*\n")
	}

	_, err = file.WriteString(builder.String())
	if err != nil {
		fmt.Printf("Error writing Markdown file: %v\n", err)
		return
	}

	if !IsQuiet {
		if absMd, err := filepath.Abs(MdPath); err == nil {
			fmt.Printf("📝 Markdown data exported to:\n %s\n", absMd)
		} else {
			fmt.Printf("📝 Markdown data exported to:\n %s\n", MdPath)
		}
	}
}

// exportToCsv serializes the findings into a comma-separated values (CSV) file.
func exportToCsv(findings []Finding) {
	file, err := os.Create(CsvPath)
	if err != nil {
		fmt.Printf("Error creating CSV file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"File", "Line", "Tag", "Author", "Content", "Date", "Status"}
	if err := writer.Write(header); err != nil {
		fmt.Printf("Error writing CSV header: %v\n", err)
		return
	}

	for _, finding := range findings {
		row := []string{
			finding.File,
			strconv.Itoa(finding.Line),
			finding.Tag,
			finding.Author,
			finding.Content,
			finding.When,
			finding.Status,
		}
		if err := writer.Write(row); err != nil {
			fmt.Printf("Error writing CSV row: %v\n", err)
			return
		}
	}

	if !IsQuiet {
		if absCsv, err := filepath.Abs(CsvPath); err == nil {
			fmt.Printf("📊 CSV data exported to:\n %s\n", absCsv)
		} else {
			fmt.Printf("📊 CSV data exported to:\n %s\n", CsvPath)
		}
	}
}

// generateHtmlReport creates a standalone HTML file containing the audit findings by injecting data into a template.
func generateHtmlReport(findings []Finding, config Config, outputPath string) error {
	jsonData, err := json.Marshal(findings)
	if err != nil {
		return fmt.Errorf("Failed to marshal findings to JSON: %w", err)
	}

	data := struct {
		FindingsJSON template.JS
		Config       Config
	}{
		FindingsJSON: template.JS(jsonData),
		Config:       config,
	}

	tmpl, err := template.ParseFS(content, "assets/template.html")
	if err != nil {
		return fmt.Errorf("Failed to parse embedded template: %w", err)
	}

	if err := extractAssets(); err != nil {
		return fmt.Errorf("Failed to extract assets: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("Failed to create report file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("Failed to execute template: %w", err)
	}

	return nil
}

// extractAssets copies the embedded assets/ folder to the local disk.
func extractAssets() error {
	return fs.WalkDir(content, "assets", func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if dir.IsDir() {
			return os.MkdirAll(path, 0755)
		}
		data, err := content.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(path, data, 0644)
	})
}
