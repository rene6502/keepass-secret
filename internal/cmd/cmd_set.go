package cmd

import (
	"io"

	"github.com/tobischo/gokeepasslib/v3"
)

// updates fields of KeePass entry
// create the entry if it does not already exist
func CmdSet(root *gokeepasslib.Group, path string, fields []string, stdout io.Writer, stderr io.Writer) bool {
	overwrite := true
	return createOrUpdateEntry(root, path, fields, overwrite, stdout, stderr)
}
