# Agent Guide for DCT (Data Check Tool)

This document provides essential information for AI agents contributing to the DCT project.

## 🚀 Build, Lint, and Test

### Build
DCT is a Go application. You can build the binary using:
```bash
go build -o dct
```

### Lint
The project follows standard Go linting practices. Although no specific configuration file is present, it is recommended to use `golangci-lint`:
```bash
golangci-lint run
```

### Test
The project uses Python-based integration tests located in the `test/` directory.

**Requirements:**
- Python 3.x
- `pytest`
- Built `dct` binary in the root directory

**Running all tests:**
```bash
# First build the binary
go build -o dct && chmod +x ./dct
# Run tests
pytest
```

**Running a single test:**
```bash
pytest test/test_dct.py::test_name
```

**Unit tests:**
Standard Go unit tests can be run using the standard toolchain:
```bash
go test ./...
```

---

## 🛠️ Project Architecture

DCT is structured as a Cobra-based CLI tool:
- `main.go`: Entry point.
- `cmd/`: Contains all CLI commands.
    - `root.go`: Root command definition and sub-command registration.
    - `[command]/`: Each sub-command has its own package (e.g., `peek`, `diff`, `generator`).
- `cmd/utils/`: Shared utilities, including DuckDB integration, file handling, and common error types.
- `test/`: Integration tests and test resources.

The core data processing engine is **DuckDB**, which allows DCT to perform high-performance SQL queries directly on CSV, JSON, NDJSON, and Parquet files.

---

## 🎨 Code Style Guidelines

### Language & Formatting
- **Go Version**: 1.24
- **Formatting**: Always use `gofmt` (or `goimports`).
- **Indentation**: Tabs (standard Go convention).

### Naming Conventions
- **Commands**: Use `PascalCase` with a `Cmd` suffix for exported Cobra commands (e.g., `PeekCmd`, `GenCmd`).
- **Variables**: Use `camelCase` for local variables and unexported package-level variables.
- **Functions/Types**: Use `PascalCase` for exported entities and `camelCase` for unexported ones.
- **Packages**: Use short, lowercase, singular names that match the directory name.

### Imports
Group imports into three blocks separated by newlines:
1. Standard library imports.
2. Local project imports (`dct/...`).
3. Third-party library imports.

Example:
```go
import (
    "fmt"
    "io"

    "dct/cmd/utils"

    "github.com/spf13/cobra"
)
```

### Error Handling
- **Wrapping**: Use `fmt.Errorf("context: %w", err)` to provide context while preserving the original error.
- **CLI Output**: In the `Run` function of Cobra commands, use `log.Fatalf` for fatal errors and `log.Printf` for warnings to provide immediate feedback to the user.
- **Custom Errors**: Define custom error types in `cmd/utils/` for domain-specific errors (e.g., `UnsupportedFileTypeErr`).
- **Validation**: Perform argument validation early in the `Run` or `Args` phase of Cobra commands.

### DuckDB Integration
- Use `cmd/utils/duckdb.go` for executing queries.
- Prefer the `utils.Query` and `utils.Execute` helpers which manage connection lifecycle.
- When adding new data types, ensure they are handled in the type switch within `utils.Query`.

### Dependencies
- **CLI**: `github.com/spf13/cobra`
- **Database**: `github.com/marcboeker/go-duckdb`
- **UI/Tables**: `github.com/charmbracelet/lipgloss`
- **Expressions**: `github.com/expr-lang/expr` (used in the generator for derived fields)

---

## 📝 Contribution Workflow

1. **Implement**: Add functionality in a new package under `cmd/` or update existing utilities.
2. **Register**: Add new commands to `cmd/root.go`.
3. **Test**: Add a corresponding test case in `test/test_dct.py` or create a new Python test file.
4. **Verify**: Run `go build` followed by `pytest` to ensure no regressions.

---

## 🤖 AI Context
When working on this codebase, prioritize readability and consistency with existing patterns. Avoid introducing new dependencies unless absolutely necessary. When modifying the CLI, ensure the `--help` documentation is updated accordingly by providing clear `Short` and `Long` descriptions in the `cobra.Command` definition.
