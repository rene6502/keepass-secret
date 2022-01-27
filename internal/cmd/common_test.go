package cmd

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

// read file into string
func readFile(file string, t *testing.T) string {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		t.Errorf("%s file does not exist\n", file)
		return ""
	}

	bytesRead, _ := ioutil.ReadFile(file)
	str := string(bytesRead)
	return str
}

func readLines(file string, t *testing.T) []string {
	str := readFile(file, t)
	if str == "" {
		return nil
	}
	lines := strings.Split(str, "\n")
	return lines
}

// compare two string arrays and fail if they are not equal
func compareLines(lines1 []string, lines2 []string, t *testing.T) bool {
	if len(lines1) != len(lines2) {
		t.Errorf("line count mismatch %d/%d\n", len(lines1), len(lines2))
		return false
	}

	for i := 0; i < len(lines1); i++ {
		line1 := strings.TrimSuffix(lines1[i], "\r")
		line2 := strings.TrimSuffix(lines2[i], "\r")

		if line1 != line2 {
			t.Errorf("line mismatch1 %s\n", line1)
			t.Errorf("line mismatch2 %s\n", line2)
			return false
		}
	}

	return true
}

// compare two text files and fail if they are not equal
func compareFiles(file1 string, file2 string, t *testing.T) bool {
	lines1 := readLines(file1, t)
	if lines1 == nil {
		return false
	}

	lines2 := readLines(file2, t)
	if lines2 == nil {
		return false
	}

	return compareLines(lines1, lines2, t)
}

// create empty database
func testCreateDatabase(db string, pw string, t *testing.T) bool {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	args := []string{"init", "-d", db, "-p", pw}
	result := Run(args, &stdout, &stderr)

	if result != 0 {
		t.Errorf("create failed, result=%d", result)
		return false
	}

	return true
}

// fill database with test values
func testFillDatabase(db string, pw string, t *testing.T) bool {
	stdout0 := strings.Builder{}
	stderr0 := strings.Builder{}
	args0 := []string{"set", "-d", db, "-p", pw, "-e", "/1/A", "-f", "UserName=admin", "-f", "Password=secret0"}
	result0 := Run(args0, &stdout0, &stderr0)

	if result0 != 0 {
		t.Errorf("fill failed, result=%d", result0)
		return false
	}

	stdout1 := strings.Builder{}
	stderr1 := strings.Builder{}
	args1 := []string{"set", "-d", db, "-p", pw, "-e", "/2/N", "-f", "UserName=user", "-f", "Password=secret1"}
	result1 := Run(args1, &stdout1, &stderr1)

	if result1 != 0 {
		t.Errorf("fill failed, result=%d", result1)
		return false
	}

	stdout2 := strings.Builder{}
	stderr2 := strings.Builder{}
	args2 := []string{"set", "-d", db, "-p", pw, "-e", "/A", "-f", "=user", "-f", "Password=secret1", "-f", "Notes=secret-type=opaque\nsecret-password=Password\nsecret-username=UserName"}
	result2 := Run(args2, &stdout2, &stderr2)

	if result2 != 0 {
		t.Errorf("fill failed, result=%d", result2)
		return false
	}

	return true
}

// export database to JSON string
func testExportDatabase(db string, pw string, out string, expected int, t *testing.T) bool {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	args := []string{"export", "-d", db, "-p", pw, "-o", out}
	result := Run(args, &stdout, &stderr)

	if result != expected {
		t.Errorf("export failed, result=%d", result)
		return false
	}

	return true
}

// import database from JSON string
func testImportDatabase(db string, pw string, in string, dryRun bool, t *testing.T) string {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	args := []string{"import", "-d", db, "-p", pw, "-i", in}
	if dryRun {
		args = append(args, "--dry-run")
	}
	result := Run(args, &stdout, &stderr)

	if result != 0 {
		t.Errorf("import failed, result=%d", result)
		return ""
	}

	return stdout.String()
}

// delete file
func testDeleteFile(file string, t *testing.T) bool {
	if err := os.Remove(file); err != nil {
		t.Errorf("cannot delete %s", file)
		return false
	}

	return true
}
