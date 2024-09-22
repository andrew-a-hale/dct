package utils

import (
	"fmt"
	"io"
	"log"
	"os/exec"
)

func Run(sql string, writer io.Writer) {
	cmd := exec.Command("duckdb", "-c", sql)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		log.Fatalf("Error: failed to start cmd: %s\n", cmd)
	}

	errStream, err := io.ReadAll(stderr)
	if err != nil {
		log.Fatalf("Error: failed to read stderr of subprocess: %v\n", err)
	}
	fmt.Println(string(errStream))

	outStream, err := io.ReadAll(stdout)
	if err != nil {
		log.Fatalf("Error: failed to read stdout of subprocess: %v\n", err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatalf("Error: failed to wait cmd: %s\n", cmd)
	}

	fmt.Fprint(writer, string(outStream))
}
