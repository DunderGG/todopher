// config.go handles the initialization of the application settings.
// It parses command-line flags and merges them with default configuration values.
package main

import (
	"flag"
	"fmt"
	"strings"
)

// parseFlags defines and parses command-line flags for customizing the tool's behavior.
func parseFlags() {
	// Custom usage message
	flag.Usage = func() {
		printIntro()
		fmt.Println("Usage: ToDopher [options] [search_directory]")
		fmt.Println("\nToDopher is a high-speed technical debt scanner that extracts TODOs from source code.")
		fmt.Println("\nOptions:")
		fmt.Println("  -c, --csv string     Optional path to export findings as a CSV file")
		fmt.Println("  -e, --exts string    Comma-separated list of additional file extensions")
		fmt.Println("  -h, --help           Show help message and usage examples")
		fmt.Println("  -i, --ignore string  Comma-separated list of additional folders to ignore")
		fmt.Println("  -j, --json string    Optional path to export findings as a JSON file")
		fmt.Println("  -m, --md string      Optional path to export findings as a Markdown file")
		fmt.Println("  -o, --output string  Output filename for the report (default \"report.html\")")
		fmt.Println("  -q, --quiet          Quiet mode (suppress output)")
		fmt.Println("  -t, --tags string    Comma-separated list of additional search tags")
		fmt.Println("  -w, --workers int    Number of concurrent workers (default 20)")

		fmt.Println("\nExamples:")
		fmt.Println("   ToDopher .                               # Scan current directory")
		fmt.Println("   ToDopher -q 	                            # Quiet mode")
		fmt.Println("   ToDopher -o \"report.html\" -j \"data.json\" # Custom outputs")
		fmt.Println("   ToDopher -t \"BUG,URGENT\" -e \".js,.ts\"    # Custom tags and extensions")
		fmt.Println("   ToDopher -i \"node_modules,dist\"          # Custom ignore folders")
		fmt.Println("   ToDopher -w 40                           # More threads!")
	}

	// The -q or --quiet flag enables quiet mode, which suppresses non-essential output for a cleaner experience when the user just wants the report.
	flag.BoolVar(&IsQuiet, "q", false, "Quiet mode (suppress output)")
	flag.BoolVar(&IsQuiet, "quiet", false, "Quiet mode (suppress output)")
	// The -o or --output flag allows the user to specify a custom file path for the generated report, defaulting to "report.html" if not provided.
	flag.StringVar(&OutputPath, "o", DefaultOutputPath, "Output filename for the report")
	flag.StringVar(&OutputPath, "output", DefaultOutputPath, "Output filename for the report")
	// The -t or --tags flag allows the user to append custom search tags (e.g., -t "SECURITY,IMPORTANT").
	flag.StringVar(&CustomTags, "t", "", "Comma-separated list of additional search tags")
	flag.StringVar(&CustomTags, "tags", "", "Comma-separated list of additional search tags")
	// The -e or --exts flag allows the user to append additional file extensions (e.g., -e ".js,.ts").
	flag.StringVar(&CustomExtensions, "e", "", "Comma-separated list of additional file extensions")
	flag.StringVar(&CustomExtensions, "exts", "", "Comma-separated list of additional file extensions")
	// The -i or --ignore flag allows the user to append additional folders to ignore (e.g., -i "node_modules,dist").
	flag.StringVar(&CustomIgnore, "i", "", "Comma-separated list of additional folders to ignore")
	flag.StringVar(&CustomIgnore, "ignore", "", "Comma-separated list of additional folders to ignore")
	// The -j or --json flag allows the user to specify a path for a JSON export of the findings.
	flag.StringVar(&JsonPath, "j", "", "Optional path to export findings as a JSON file")
	flag.StringVar(&JsonPath, "json", "", "Optional path to export findings as a JSON file")
	// The -c or --csv flag allows the user to specify a path for a CSV export of the findings.
	flag.StringVar(&CsvPath, "c", "", "Optional path to export findings as a CSV file")
	flag.StringVar(&CsvPath, "csv", "", "Optional path to export findings as a CSV file")
	// The -m or --md flag allows the user to specify a path for a Markdown export of the findings.
	flag.StringVar(&MdPath, "m", "", "Optional path to export findings as a Markdown file")
	flag.StringVar(&MdPath, "md", "", "Optional path to export findings as a Markdown file")
	// The -w or --workers flag allows the user to customize the number of concurrent scanning workers.
	flag.IntVar(&WorkerCount, "w", 20, "Number of concurrent workers")
	flag.IntVar(&WorkerCount, "workers", 20, "Number of concurrent workers")
	flag.Parse()
}

// initializeConfig sets up and returns the Config struct with default values
// and incorporates any custom tags or extensions provided via command-line flags.
func initializeConfig() Config {
	// Initialize configuration with default search tags, ignored folders, and allowed extensions
	config := Config{
		SearchTags:        []string{"TODO", "FIXME", "HACK", "BUG", "SUGGESTION", "IDEA", "REWORK"},
		IgnoreFolders:     []string{"Intermediate", "ThirdParty", ".git", "Binaries", "Saved", "Plugins"},
		AllowedExtensions: []string{".h", ".cpp", ".html", ".go", ".java", ".py", ".ini", ".cs"},
		WorkerCount:       WorkerCount,
	}

	// Update WorkerCount if it's less than 1
	if config.WorkerCount < 1 {
		config.WorkerCount = 1
	}

	// Append custom tags if provided via CLI
	if CustomTags != "" {
		tags := strings.Split(CustomTags, ",")
		for _, tag := range tags {
			trimmed := strings.TrimSpace(tag)
			if trimmed != "" {
				config.SearchTags = append(config.SearchTags, trimmed)
			}
		}
	}

	// Append custom extensions if provided via CLI
	if CustomExtensions != "" {
		exts := strings.Split(CustomExtensions, ",")
		for _, ext := range exts {
			trimmed := strings.TrimSpace(ext)
			if trimmed != "" {
				// Ensure it starts with a dot
				if !strings.HasPrefix(trimmed, ".") {
					trimmed = "." + trimmed
				}
				config.AllowedExtensions = append(config.AllowedExtensions, trimmed)
			}
		}
	}

	// Append custom ignore folders if provided via CLI
	if CustomIgnore != "" {
		folders := strings.Split(CustomIgnore, ",")
		for _, folder := range folders {
			trimmed := strings.TrimSpace(folder)
			if trimmed != "" {
				config.IgnoreFolders = append(config.IgnoreFolders, trimmed)
			}
		}
	}
	return config
}
