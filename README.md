# <img width="64" height="64" alt="appicon" src="https://github.com/user-attachments/assets/0b314a1d-b14d-4a07-9c8b-af04b7eeca8c" /> **ToDopher**

**ToDopher** is a lightning-fast source code auditor built in Go. It's designed to help engineering teams manage technical debt by extracting `TODO`, `FIXME`, and `BUG` comments from large codebases (like Unreal Engine projects) and presenting them in a beautiful, searchable dashboard.

###### ToDo + Gopher = ToDopher. Get it? It's a pun. You're welcome...

## 🚀 Features

- **Concurrent Scanning**: Uses a tunable Goroutine worker pool to audit thousands of files in seconds.
- **Binary File Safety**: Robustly detects and skips binary files (like `.uasset` or `.exe`) to prevent scanning junk data.
- **Smart Filtering**: Automatically ignores common "noise" directories like `Intermediate/`, `Binaries/`, and `.git/`.
- **Regex Extraction**: Captures not just the comment, but the line number and the optional **author** (e.g., `TODO-Dunder: fix this`).
- **Interactive Report**: Generates a standalone, dark-mode-first HTML report with a **Configuration Info Tray** for full audit transparency.
- **Multi-Format Export**: Export results to **CSV**, **JSON**, or clean **Markdown** tables for documentation.
- **UE-Ready**: Pre-configured with filters for header files, source code, and configuration files common in Unreal Engine.

## 📦 Getting Started

### Prerequisites
- [Go 1.20+](https://go.dev/dl/)

### How to Run
The easiest way to run **ToDopher** is using the `go run` command:
```powershell
go run ToDopher.go "C:\Path\To\Your\Project"
```
If no path is provided, it defaults to scanning the current directory.

Alternatively, you can build a standalone executable:
```powershell
go build -o ToDopher.exe ToDopher.go
```

## 📖 Usage

Run ToDopher from the command line, providing an optional path to scan. 

```powershell
./ToDopher.exe "C:\Path\To\Your\Project"
```

### Command Line Options

| Flag  | Description |
| :---:  | :--- |
| `--csv`, `-c` | Optional path to export findings as a spreadsheet-friendly CSV file. |
| `--exts`, `-e` | Comma-separated list of additional file extensions to scan (e.g., `.js,.ts`). |
| `--help`, `-h` | Show help message and usage examples. |
| `--ignore`, `-i` | Comma-separated list of additional folders to ignore (e.g., `node_modules,build`). |
| `--json`, `-j` | Optional path to export findings as a machine-readable JSON file. |
| `--md`, `-m` | Optional path to export findings as a clean Markdown table (perfect for GitHub/Notion). |
| `--output`, `-o` | Custom **file path** for the generated HTML report (defaults to `report.html`). Supports absolute or relative paths. |
| `--quiet`, `-q` | Quiet mode. Suppresses the ASCII intro, progress bar, and status messages. Useful for CI/CD. |
| `--tags`, `-t` | Comma-separated list of additional tags to search for (e.g., `IMPORTANT,SECURITY`). |
| `--workers`, `-w` | Number of concurrent worker goroutines (defaults to `20`). |

## 📊 The Report
After running the audit, ToDopher generates a `report.html` file in the project folder. 

1. **Dark Mode First**: Optimized for long coding sessions.
2. **Instant Search**: Find all `FIXME` items or everything assigned to a specific author instantly.
3. **No Server Required**: The report is a portable, standalone file.

## 🛠️ Configuration
Currently, ToDopher is configured via the `Config` struct in [ToDopher.go](ToDopher.go#L72-L76):
- **Search Tags**: `TODO`, `FIXME`, `HACK`, `BUG`, `SUGGESTION`, `IDEA`, `REWORK`.
- **Allowed Extensions**: `.h`, `.cpp`, `.cs`, `.py`, `.ini`, `.go`, `.java`, `.html`.

## 🛠️ Built With

- **[Go](https://go.dev/)**: The core scanning engine, leveraging Goroutines for concurrency.
- **[DataTables](https://datatables.net/)**: Powering the interactive dashboard with RowGroup for module organization.
- **[Chart.js](https://www.chartjs.org/)**: Visualizing technical debt distribution via a dynamic pie chart.
- **[jQuery](https://jquery.com/)**: Handling DOM manipulation and DataTables events.
- **[Git](https://git-scm.com/)**: Integrated via `git blame` to automatically identify the authors of untagged TODOs.

## 👤 Author

**David Bennehag** - [@dunder](https://github.com/dunder) - [dunder.gg](https://dunder.gg)

## 📄 License

This project is licensed under the GPL-3.0 - see the [LICENSE](LICENSE) file for details.
