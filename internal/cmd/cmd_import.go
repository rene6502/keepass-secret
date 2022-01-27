package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/tobischo/gokeepasslib/v3"
)

// import JSON to KeePass database
func CmdImport(root *gokeepasslib.Group, in string, stdout io.Writer, stderr io.Writer) (bool, int) {
	if _, err := os.Stat(in); os.IsNotExist(err) {
		fmt.Fprintf(stderr, "%s file does not exist\n", in)
		return false, 1
	}

	readFile, err := os.Open(in)
	if err != nil {
		fmt.Fprintf(stderr, "%s\n", err)
		return false, 1
	}
	defer readFile.Close()

	reader := bufio.NewReader(readFile)
	dec := json.NewDecoder(reader)
	list := make([]map[string]string, 0)
	err = dec.Decode(&list)
	if err != nil {
		fmt.Fprintf(stderr, "%s\n", err)
		return false, 1
	}

	modified := false
	for i := 0; i < len(list); i++ {
		entry := list[i]
		path, ok := entry["path"]
		if !ok {
			fmt.Fprintf(stderr, "missing path in entry #%d\n", i)
			return false, 1
		}

		fields := make([]string, 0)
		for name, value := range entry {
			if name != "path" {
				fields = append(fields, name+"="+value)
			}
		}

		overwrite := false
		if createOrUpdateEntry(root, path, fields, overwrite, stdout, stderr) {
			modified = true
		}
	}

	return modified, 0
}
