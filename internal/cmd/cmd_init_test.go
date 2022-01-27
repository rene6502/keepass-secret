package cmd

import (
	"os"
	"strings"
	"testing"
)

// create empty databse
func TestCreate(t *testing.T) {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	db := "init_test.kdbx"
	pw := "a1b2c3d4"

	// create empty database
	args := []string{"init", "-d", db, "-p", pw}
	result := Run(args, &stdout, &stderr)

	if result != 0 {
		t.Errorf("init failed, result=%d", result)
		return
	}

	if _, err := os.Stat(db); os.IsNotExist(err) {
		t.Errorf("%s file does not exit", db)
		return
	}

	// remove file
	if err := os.Remove(db); err != nil {
		t.Errorf("cannot delete %s", db)
		return
	}
}
