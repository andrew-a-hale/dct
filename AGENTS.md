# Agent Guide for DCT (Data Check Tool)

This document provides essential information for AI agents contributing to the DCT project. It outlines the build process, architectural patterns, and coding standards required to maintain consistency and quality.

## 🚀 Build, Lint, and Test

### Build
DCT is a Go application. Build the binary using the following command:
```bash
go build -o dct
```
Note: The integration tests in `test/test_dct.py` automatically attempt to build the binary if it is missing or needs an update, but manual builds are recommended after significant changes.

### Lint
The project follows standard Go linting practices. Although no `.golangci.yml` is present, it is recommended to use `golangci-lint`:
```bash
golangci-lint run
```

### Test
DCT employs a hybrid testing strategy:
1. **Integration Tests (Python)**: Located in `test/test_dct.py`. These tests execute the `dct` binary via `subprocess` and compare `stdout/stderr` against "golden files" in `test/expected/`.
2. **Unit Tests (Go)**: Used for internal logic, especially in `cmd/utils`.

**Running tests:**
```bash
# Run all integration tests (requires pytest)
pytest

# Run a single integration test by name
pytest test/test_dct.py::test_name

# Run tests matching a specific pattern
pytest -k "diff"

# Run Go unit tests
go test ./...
```

---

## 🛠️ Project Architecture

DCT is built as a modular CLI using the **Cobra** library. The core data processing engine is **DuckDB**, which allows for high-performance SQL execution on flat files.

- `main.go`: The entry point that simply calls `cmd.Execute()`.
- `cmd/`: Contains the CLI command hierarchy.
    - `root.go`: Defines the base `dct` command and registers all sub-commands.
    - `[command]/`: Each sub-command (e.g., `peek`, `diff`, `generator`, `profile`) resides in its own package.
- `cmd/utils/`: Shared logic and core integrations.
    - `duckdb.go`: Manages DuckDB connections and provides `Query` and `Execute` helpers.
    - `files.go`: Centralizes supported file types (CSV, JSON, NDJSON, Parquet) and extension validation.
    - `assert.go`: Helpers for validating CLI arguments.
    - `stack.go`: Generic stack implementation for complex parsing (e.g., in `flattify`).
- `test/`: Contains the integration test suite and sample data resources.

---

## 🎨 Code Style Guidelines

### Language & Formatting
- **Go Version**: 1.24 (as specified in `go.mod`).
- **Formatting**: Strictly follow `gofmt`.
- **Indentation**: Use **Tabs** for Go source files and **4 Spaces** for Python test files.

### Naming Conventions
- **Cobra Commands**: Exported command variables must use `PascalCase` with a `Cmd` suffix (e.g., `PeekCmd`, `InferCmd`).
- **Variables**: Use `camelCase` for local and unexported package-level variables.
- **Functions/Types**: Use `PascalCase` for exported entities and `camelCase` for internal ones.
- **Packages**: Use short, lowercase, singular names that match their directory name.

### Import Grouping
Imports should be organized into three distinct blocks separated by newlines:
1. Standard library imports (e.g., `fmt`, `os`).
2. Local project imports (`dct/cmd/...`).
3. Third-party library imports (e.g., `github.com/spf13/cobra`).

Example from `cmd/peek/peek.go`:
```go
import (
	"fmt"
	"log"

	"dct/cmd/utils"

	"github.com/spf13/cobra"
)
```

### Error Handling
- **Wrapping**: Always provide context when returning errors: `fmt.Errorf("failed to process file %s: %w", path, err)`.
- **CLI Feedback**: 
    - In `Run` functions, use `log.Fatalf` for fatal errors to exit the process with a clear message.
    - Use `log.Printf` for warnings that do not stop execution.
- **Custom Errors**: Define domain-specific error types in `cmd/utils/` to allow for structured error checking (e.g., `UnsupportedFileTypeErr`).

---

## 🦆 DuckDB Integration

All data-heavy operations should be offloaded to DuckDB to ensure performance and consistency across file formats.
- **`utils.Query(query string)`**: Executes a SQL query and returns a `Result` struct. The `Result` struct includes methods like `Render()` for tables and `ToCsv()` for exports.
- **`utils.Execute(query string)`**: Executes a statement without returning rows.
- **Type Switches**: When handling query results, ensure that the type switch in `utils.Query` (found in `cmd/utils/duckdb.go`) accounts for the specific DuckDB types returned by your query. DuckDB often returns types like `int32` or `float64` that must be safely mapped to Go's `any` interface.

### SQL Snippet Patterns
When writing queries for DCT, follow these DuckDB-specific patterns:
- **File Access**: Always wrap file paths in single quotes: `SELECT * FROM 'file.csv'`.
- **JSON Handling**: Use the `read_json_auto` or similar functions when fine-grained control is needed, though DuckDB usually handles this automatically with the extension-based sniffing.
- **Aggregates**: Leverage DuckDB's advanced aggregates for `diff` operations (e.g., `approx_count_distinct` for large datasets if performance is an issue).

---

## 🛠️ Command-Specific Guidance

- **Generator**: Uses the `expr` library for DSL evaluations. When adding new generators, ensure they are registered in `cmd/generator/sources/`.
- **Flattify**: Uses a custom stack-based parser to unnest JSON. It supports generating both raw JSON output and DuckDB-compliant SQL `SELECT` statements with deep path accessors (e.g., `json->'$.a[0].b'`).
- **Profile**: Performs character-level analysis. If adding new profiling metrics, ensure they are efficient as they run over the entire dataset in Go memory after the initial DuckDB query.

---

## 📝 Contribution Workflow

1. **Scaffold**: Create a new directory in `cmd/` for your command.
2. **Implement**: Define the `cobra.Command` and its `Run` logic. Keep the command logic minimal; delegate complex processing to helpers or `cmd/utils`.
3. **Register**: Add your command to `rootCmd` in `cmd/root.go`.
4. **Integration Test**:
    - Add a test case to `test/test_dct.py`.
    - If your command produces output, create a corresponding "golden file" in `test/expected/`.
    - Use `subprocess.run` to verify the CLI behavior.
5. **Validation**: Run `go build` followed by `pytest` to ensure no regressions were introduced.

---

## 🤖 AI Context & Best Practices
- **Idiomatic Go**: Prefer standard library features unless a third-party library is already established in the project.
- **Minimal Dependencies**: Current dependencies include `cobra`, `go-duckdb`, `lipgloss` (for tables), and `expr` (for generator logic). Avoid adding new ones without strong justification.
- **CLI Documentation**: Ensure every command has a concise `Short` description and a detailed `Long` description explaining flags and usage.
- **Performance**: Always prefer DuckDB's native functions (e.g., `count(*)`, `unnest()`) over manual data processing in Go.
