package cmd

import (
	"fmt"
	"io"
	"strings"
)

// read value of specified field and write it to stdout
func CmdGet(entryMap *EntryMap, path string, field string, stdout io.Writer, stderr io.Writer) int {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path // add missing "/"
	}

	values, ok := entryMap.GetValues(path)

	if !ok {
		fmt.Fprintf(stderr, "path '%s' does not exist\n", path)
		return 1 // failure
	}

	value, ok := values.GetValue(field)
	if !ok {
		fmt.Fprintf(stderr, "field '%s' does not exist in path '%s'\n", field, path)
		return 1 // failure
	}

	if ok {
		fmt.Fprintf(stdout, "%s", value) // stdout without newline!
	}

	return 0 // success
}
