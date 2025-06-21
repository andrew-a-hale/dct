package main

import (
	"context"
	"dct-mcp-server/server"
	"fmt"
	"log"
	"os"
)

const TOOLS_PATH = "./mcp-server/tools.yaml"

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <dct-binary-path>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s ./dct\n", os.Args[0])
		os.Exit(1)
	}

	dctPath := os.Args[1]

	if _, err := os.Stat(dctPath); os.IsNotExist(err) {
		log.Fatalf("DCT binary not found at: %s", dctPath)
	}

	server := server.NewMCPServer(dctPath, TOOLS_PATH)

	log.Printf("Starting DCT MCP Server with DCT binary at: %s", dctPath)

	ctx := context.Background()
	if err := server.Start(ctx); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
