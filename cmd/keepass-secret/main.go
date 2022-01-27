package main

import (
	"fmt"
	"keepass-secret/internal/cmd"
	"os"
	"strings"
)

// delegate to cmd.Run
// stdout/stderr are passed as strings.Builder to allow unit testing
func main() {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	result := cmd.Run(os.Args[1:], &stdout, &stderr)
	fmt.Fprint(os.Stderr, stderr.String())
	fmt.Fprint(os.Stdout, stdout.String())
	os.Exit(result)
}
