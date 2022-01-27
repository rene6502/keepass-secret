package cmd

import (
	"os"
	"strings"
	"testing"
)

// create -> fill -> export -> import -> compare
func TestImport(t *testing.T) {
	db0 := "test/export.kdbx"
	db1 := "test/import.kdbx"
	out0 := "test/exported0.json"
	out1 := "test/exported1.json"
	pw := "a1b2c3d4"

	if !testCreateDatabase(db0, pw, t) {
		return
	}

	if !testFillDatabase(db0, pw, t) {
		return
	}

	if !testExportDatabase(db0, pw, out0, 0, t) {
		return
	}

	if !testCreateDatabase(db1, pw, t) {
		return
	}

	// import with dry-run
	expected := "A created\n1/A created\n2/N created\n"
	actual := testImportDatabase(db1, pw, out0, true /*dryRun*/, t)
	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
		return
	}

	if !testExportDatabase(db1, pw, out1, 0, t) {
		return
	}

	expected = "[]\n\n"
	actual = readFile(out1, t)
	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
		return
	}

	// import without dry-run
	expected = "A created\n1/A created\n2/N created\n"
	actual = testImportDatabase(db1, pw, out0, false /*dryRun*/, t)
	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
		return
	}

	// 2nd import must ignore values
	expected = ""
	actual = testImportDatabase(db1, pw, out0, false /*dryRun*/, t)
	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
		return
	}

	if !testExportDatabase(db1, pw, out1, 0, t) {
		return
	}

	compareFiles(out0, out1, t)

	testDeleteFile(out0, t)
	testDeleteFile(out1, t)
	testDeleteFile(db0, t)
	testDeleteFile(db1, t)
}

// import, file does not exist
func TestImportFileNoExisting(t *testing.T) {
	db := "test/test.kdbx"
	pw := "1234"
	in := "test/missing.json"

	stdout := strings.Builder{}
	stderr := strings.Builder{}
	args := []string{"import", "-d", db, "-p", pw, "-i", in}
	result := Run(args, &stdout, &stderr)

	if result == 0 {
		t.Errorf("run must fail")
		return
	}

	expected := "test/missing.json file does not exist\n"
	actual := stderr.String()
	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
		return
	}
}

// import, file is invalid
func TestImportFileInvalidJson(t *testing.T) {
	db := "test/test.kdbx"
	pw := "1234"
	in := "test/invalid.json"
	json := []string{"a"}

	stdout := strings.Builder{}
	stderr := strings.Builder{}
	writeFile(in, &json, &stderr)
	args := []string{"import", "-d", db, "-p", pw, "-i", in}
	result := Run(args, &stdout, &stderr)

	if result == 0 {
		t.Errorf("run must fail")
		return
	}

	if err := os.Remove(in); err != nil {
		t.Errorf("cannot delete %s", in)
		return
	}
}

// import, entry does not contain path
func TestImportFileMissingPath(t *testing.T) {
	db := "test/test.kdbx"
	pw := "1234"
	in := "test/invalid.json"
	json := []string{"[{\"UserName\": \"admin\", \"Password\": \"secret\"}]"}

	stdout := strings.Builder{}
	stderr := strings.Builder{}
	writeFile(in, &json, &stderr)
	args := []string{"import", "-d", db, "-p", pw, "-i", in, "--dry-run"}
	result := Run(args, &stdout, &stderr)

	if result == 0 {
		t.Errorf("run must fail")
		return
	}

	expected := "missing path in entry #0\n"
	actual := stderr.String()
	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
		return
	}

	if err := os.Remove(in); err != nil {
		t.Errorf("cannot delete %s", in)
		return
	}
}
