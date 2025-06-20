package main

import (
	"context"
	"fmt"
	"log"
	"mcp-server/server"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <dct-binary-path>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s ../dct\n", os.Args[0])
		os.Exit(1)
	}

	dctPath := os.Args[1]

	// Verify DCT binary exists
	if _, err := os.Stat(dctPath); os.IsNotExist(err) {
		log.Fatalf("DCT binary not found at: %s", dctPath)
	}

	ctx := context.Background()
	server := server.NewMCPServer(dctPath)

	log.Printf("Starting DCT MCP Server with DCT binary at: %s", dctPath)

	if err := server.Start(ctx); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
