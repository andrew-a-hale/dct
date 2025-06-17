package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type DCTExecutor struct {
	dctPath string
}

func NewDCTExecutor(dctPath string) *DCTExecutor {
	return &DCTExecutor{dctPath: dctPath}
}

type ExecutionResult struct {
	Success    bool   `json:"success"`
	Output     string `json:"output"`
	Error      string `json:"error,omitempty"`
	ExitCode   int    `json:"exit_code"`
	Duration   string `json:"duration"`
}

func (e *DCTExecutor) Execute(command string, args []string) (*ExecutionResult, error) {
	start := time.Now()
	
	// Build the full command
	fullArgs := append([]string{command}, args...)
	cmd := exec.Command(e.dctPath, fullArgs...)
	
	// Capture both stdout and stderr
	output, err := cmd.CombinedOutput()
	duration := time.Since(start)
	
	result := &ExecutionResult{
		Output:   string(output),
		Duration: duration.String(),
	}
	
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
		} else {
			result.ExitCode = -1
		}
	} else {
		result.Success = true
		result.ExitCode = 0
	}
	
	return result, nil
}

// Helper function to create temporary files for commands that need file input
func (e *DCTExecutor) createTempFile(content string, extension string) (string, error) {
	tmpFile, err := os.CreateTemp("", fmt.Sprintf("dct-mcp-*%s", extension))
	if err != nil {
		return "", err
	}
	
	if _, err := tmpFile.WriteString(content); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", err
	}
	
	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}
	
	return tmpFile.Name(), nil
}

// Clean up temporary files
func (e *DCTExecutor) cleanup(filePath string) {
	if strings.Contains(filePath, "dct-mcp-") {
		os.Remove(filePath)
	}
}

// Execute peek command
func (e *DCTExecutor) ExecutePeek(filePath string, lines int, outputFile string) (*ExecutionResult, error) {
	args := []string{filePath}
	
	if lines > 0 {
		args = append(args, "-n", fmt.Sprintf("%d", lines))
	}
	
	if outputFile != "" {
		args = append(args, "-o", outputFile)
	}
	
	return e.Execute("peek", args)
}

// Execute diff command  
func (e *DCTExecutor) ExecuteDiff(keys string, file1, file2 string, metrics string, showAll bool, outputFile string) (*ExecutionResult, error) {
	args := []string{keys, file1, file2}
	
	if metrics != "" {
		args = append(args, "-m", metrics)
	}
	
	if showAll {
		args = append(args, "-a")
	}
	
	if outputFile != "" {
		args = append(args, "-o", outputFile)
	}
	
	return e.Execute("diff", args)
}

// Execute chart command
func (e *DCTExecutor) ExecuteChart(filePath string, colIndex int, width int32) (*ExecutionResult, error) {
	args := []string{}
	
	if filePath != "" {
		args = append(args, filePath)
	}
	
	if colIndex >= 0 {
		args = append(args, fmt.Sprintf("%d", colIndex))
	}
	
	if width > 0 {
		args = append(args, "-w", fmt.Sprintf("%d", width))
	}
	
	return e.Execute("chart", args)
}

// Execute generator command
func (e *DCTExecutor) ExecuteGenerate(schema string, lines int, format string, outputFile string) (*ExecutionResult, error) {
	// Create temp file for schema if it's JSON content
	var schemaFile string
	var err error
	var cleanup bool
	
	if strings.HasPrefix(strings.TrimSpace(schema), "{") || strings.HasPrefix(strings.TrimSpace(schema), "[") {
		// It's JSON content, create temp file
		schemaFile, err = e.createTempFile(schema, ".json")
		if err != nil {
			return nil, fmt.Errorf("failed to create temp schema file: %v", err)
		}
		cleanup = true
	} else {
		// It's a file path
		schemaFile = schema
	}
	
	args := []string{schemaFile}
	
	if lines > 0 {
		args = append(args, "-n", fmt.Sprintf("%d", lines))
	}
	
	if format != "" {
		args = append(args, "-f", format)
	}
	
	if outputFile != "" {
		args = append(args, "-o", outputFile)
	}
	
	result, execErr := e.Execute("gen", args)
	
	if cleanup {
		e.cleanup(schemaFile)
	}
	
	return result, execErr
}

// Execute flattify command
func (e *DCTExecutor) ExecuteFlattify(input string, sql bool, outputFile string) (*ExecutionResult, error) {
	var inputFile string
	var err error
	var cleanup bool
	
	if strings.HasPrefix(strings.TrimSpace(input), "{") || strings.HasPrefix(strings.TrimSpace(input), "[") {
		// It's JSON content, create temp file
		inputFile, err = e.createTempFile(input, ".json")
		if err != nil {
			return nil, fmt.Errorf("failed to create temp input file: %v", err)
		}
		cleanup = true
	} else {
		// It's a file path
		inputFile = input
	}
	
	args := []string{inputFile}
	
	if sql {
		args = append(args, "-s")
	}
	
	if outputFile != "" {
		args = append(args, "-o", outputFile)
	}
	
	result, execErr := e.Execute("flattify", args)
	
	if cleanup {
		e.cleanup(inputFile)
	}
	
	return result, execErr
}

// Execute js2sql command
func (e *DCTExecutor) ExecuteJs2Sql(schemaFile string, tableName string, outputFile string) (*ExecutionResult, error) {
	args := []string{schemaFile}
	
	if tableName != "" {
		args = append(args, "-t", tableName)
	}
	
	if outputFile != "" {
		args = append(args, "-o", outputFile)
	}
	
	return e.Execute("js2sql", args)
}

// Execute profile command
func (e *DCTExecutor) ExecuteProfile(filePath string, outputFile string) (*ExecutionResult, error) {
	args := []string{filePath}
	
	if outputFile != "" {
		args = append(args, "-o", outputFile)
	}
	
	return e.Execute("prof", args)
}