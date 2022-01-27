package cmd

import (
	"testing"
)

// create -> fill -> export
func TestExport(t *testing.T) {
	db := "test/export.kdbx"
	out := "test/exported.json"
	pw := "a1b2c3d4"

	if !testCreateDatabase(db, pw, t) {
		return
	}

	if !testFillDatabase(db, pw, t) {
		return
	}

	if !testExportDatabase(db, pw, out, 0, t) {
		return
	}

	testDeleteFile(out, t)
	testDeleteFile(db, t)
}
