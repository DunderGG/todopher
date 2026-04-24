# ToDopher — Product Roadmap

This document outlines planned improvements and new features for the ToDopher technical debt scanner. Items are organized by theme and relative priority. Items marked **[Done]** reflect features already shipped.

---

## 1. Core Scanner & Detection

| Status | Item | Description |
| :---: | :--- | :--- |
| ✅ Done | Goroutine worker pool | Concurrent scanning via 20 parallel workers using channels and `sync.WaitGroup`. |
| ✅ Done | Regex extraction | Captures tag, optional author (`TODO-Dunder`), and full comment content. |
| ✅ Done | Author attribution | Parses inline author names from tagged comments. |
| ✅ Done | Staleness classification | Classifies findings as `Fresh`, `Stale`, or `Ancient` based on git history. |
| ✅ Done | Surrounding context | Captures adjacent lines around each finding for better readability. |
| ✅ Done | Non-UTF8 file handling | Gracefully skip or flag files with invalid encodings instead of crashing or producing garbage output. |
| ✅ Done | Binary file detection | Detect and skip binary files (`.uasset`, `.exe`, `.pak`, etc.) before reading content. |
| ✅ Done | Configurable worker count | Expose `--workers` flag so users can tune concurrency for their machine (e.g., `--workers 40`). |
| 🔲 Todo | Symlink handling | Add a flag to control whether symbolic links are followed during directory traversal. |

---

## 2. Configuration

| Status | Item | Description |
| :---: | :--- | :--- |
| ✅ Done | CLI flags | Full flag surface: `--tags`, `--exts`, `--ignore`, `--output`, `--json`, `--quiet`. |
| 🔲 Todo | Config file support | Load settings from a `todopher.yaml` or `.todopherrc` file in the project root, eliminating long CLI strings in CI/CD pipelines. |
| 🔲 Todo | Per-project overrides | Allow a `.todopher.yaml` file at any directory level to override global settings, similar to `.editorconfig`. |
| 🔲 Todo | Severity tagging | Map tags to severity levels (e.g., `BUG` = Critical, `FIXME` = High, `TODO` = Medium, `IDEA` = Low) configurable by the user. |

---

## 3. Reporting & Export

| Status | Item | Description |
| :---: | :--- | :--- |
| ✅ Done | HTML dashboard | Standalone dark-mode report with DataTables sorting and filtering. |
| ✅ Done | JSON export | Machine-readable export via `--json` flag, suitable for downstream tooling. |
| ✅ Done | Markdown export | Export findings as a structured Markdown file (`--md`) for pasting into GitHub wikis, Notion, or Confluence. |
| ✅ Done | CSV export | Export findings as a `.csv` file (`--csv`) for import into spreadsheets or project management tools. |
| 🔲 Todo | Severity color coding | Color-code findings in the HTML report by severity (e.g., red for `BUG`, orange for `FIXME`, blue for `TODO`). |
| 🔲 Todo | Per-author summary | Add a summary section to the HTML report breaking down finding counts by author. |
| 🔲 Todo | Per-file summary | Add a collapsible file tree view in the report showing the debt density per file and directory. |
| 🔲 Todo | Trend reporting | Compare two JSON exports over time and generate a diff report showing which debts were resolved, added, or regressed. |

---

## 4. Git & VCS Integration

| Status | Item | Description |
| :---: | :--- | :--- |
| ✅ Done | `git blame` timestamp | Uses `git blame` to populate the `When` field, classifying finding age. |
| 🔲 Todo | Commit hash linking | Include the git commit hash for each finding so users can jump directly to the introducing commit. |
| 🔲 Todo | Branch-aware scanning | Add a `--branch` flag to scan a specific git branch without checking it out locally. |
| 🔲 Todo | `.gitignore` respect | Automatically honour `.gitignore` rules during file discovery as an opt-in flag (`--gitignore`). |

---

## 5. CI/CD & Automation

| Status | Item | Description |
| :---: | :--- | :--- |
| ✅ Done | Quiet mode | `--quiet` flag suppresses all output, making it suitable for scripted use. |
| 🔲 Todo | Exit code on findings | Return a non-zero exit code when findings exceed a configurable threshold (`--fail-on N`), enabling pipeline gating. |
| 🔲 Todo | GitHub Actions example | Ship a ready-to-use `todopher.yml` workflow file that scans on every push and uploads the report as an artifact. |
| 🔲 Todo | Pre-commit hook template | Provide a sample `pre-commit` hook that blocks commits when critical (`BUG`, `FIXME`) findings are introduced. |
| 🔲 Todo | GitHub Issues integration | Optionally auto-create GitHub Issues from new findings using the GitHub API, tagged with a `technical-debt` label. |

---

## 6. Notifications

| Status | Item | Description |
| :---: | :--- | :--- |
| 🔲 Todo | Slack webhook | Post a summary message to a Slack channel after a scan completes, including finding counts by severity. |
| 🔲 Todo | Discord webhook | Same as above, targeting a Discord channel via an incoming webhook URL. |
| 🔲 Todo | Email report | Optionally email the generated HTML report to a recipient list after a scan (useful for scheduled CI jobs). |

---

## 7. Language & Ecosystem Support

| Status | Item | Description |
| :---: | :--- | :--- |
| ✅ Done | C/C++, Go, Python, C#, Java, INI | Default extension and comment-style coverage for common languages. |
| 🔲 Todo | JavaScript & TypeScript | Add `.js`, `.ts`, `.tsx`, `.jsx` to defaults and ensure `//` and `/* */` patterns are covered. |
| 🔲 Todo | Lua & Blueprint scripts | Add `.lua` support for Unreal Engine Lua scripting and automation scripts. |
| 🔲 Todo | GLSL / HLSL shaders | Add `.glsl`, `.hlsl`, `.usf`, `.ush` support for shader file auditing. |
| 🔲 Todo | Shell & YAML | Add `.sh`, `.bash`, `.yml`, `.yaml` to catch debt in infrastructure-as-code and CI configuration files. |

---

## 8. Developer Experience

| Status | Item | Description |
| :---: | :--- | :--- |
| 🔲 Todo | Live web server mode | Add a `--serve` flag to spin up a local HTTP server that auto-refreshes the dashboard as files change, instead of writing a static file. |
| 🔲 Todo | VS Code extension | Build a companion extension that displays ToDopher findings inline in the editor gutter and shows a findings panel. |
| 🔲 Todo | Interactive CLI mode | Add an interactive TUI (terminal UI) using a library like `bubbletea` to browse and navigate findings without opening a browser. |
| 🔲 Todo | Watch mode | Add a `--watch` flag to re-scan automatically when source files change, keeping the report up to date during active development. |

---

## 9. Performance & Scalability

| Status | Item | Description |
| :---: | :--- | :--- |
| ✅ Done | Concurrent file scanning | Worker pool distributes I/O across goroutines. |
| 🔲 Todo | Large project benchmarks | Profile scanning against a full AAA Unreal Engine project (100k+ files) and establish baseline performance targets. |
| 🔲 Todo | Streaming JSON output | Write findings to the JSON file incrementally as they are discovered, reducing peak memory usage on very large codebases. |
| 🔲 Todo | Scan caching | Cache file hashes between runs and skip unchanged files, significantly reducing re-scan time in watch or CI contexts. |
