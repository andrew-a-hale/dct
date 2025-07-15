package main

import (
	"context"
	"dct-mcp-server/server"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <dct-binary-path>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s ./dct\n", os.Args[0])
		os.Exit(1)
	}

	dctPath := os.Args[1]

	var err error
	var msg string
	if _, err = os.Stat(dctPath); os.IsNotExist(err) {
		msg = fmt.Sprintf("DCT binary not found at: %s", dctPath)
	}

	if _, err = exec.LookPath(dctPath); err != nil {
		msg = "DCT command not found. Is the command installed?"
	}

	if err != nil {
		log.Fatal(msg)
	}

	server := server.NewMCPServer(dctPath)

	log.Printf("Starting DCT MCP Server with DCT binary at: %s", dctPath)

	ctx := context.Background()
	if err := server.Start(ctx); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
