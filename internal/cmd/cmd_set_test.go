package cmd

import (
	"os"
	"strings"
	"testing"
)

// set value, create -> set -> export -> compare
func TestSet(t *testing.T) {
	db := "set_test.kdbx"
	pw := "a1b2c3d4"
	out := "test/test.json"

	// create empty database
	args := []string{"init", "-d", db, "-p", pw}
	stdout0 := strings.Builder{}
	stderr0 := strings.Builder{}
	result := Run(args, &stdout0, &stderr0)
	if result != 0 {
		t.Errorf("init failed, result=%d", result)
		return
	}

	if _, err := os.Stat(db); os.IsNotExist(err) {
		t.Errorf("%s file does not exit", db)
		return
	}

	// create value /1/A
	args = []string{"set", "-d", db, "-p", pw, "-e", "/1/A", "-f", "UserName=admin", "-f", "Password=secret1"}
	stdout1 := strings.Builder{}
	stderr1 := strings.Builder{}
	result = Run(args, &stdout1, &stderr1)
	if result != 0 {
		t.Errorf("set failed, result=%d", result)
		return
	}

	// overwrite value /1/A
	args = []string{"set", "-d", db, "-p", pw, "-e", "/1/A", "-f", "UserName=admin", "-f", "Password=secret2"}
	stdout2 := strings.Builder{}
	stderr2 := strings.Builder{}
	result = Run(args, &stdout2, &stderr2)
	if result != 0 {
		t.Errorf("set failed, result=%d", result)
		return
	}

	// export as json
	args = []string{"export", "-d", db, "-p", pw, "-o", out}
	stdout3 := strings.Builder{}
	stderr3 := strings.Builder{}
	result = Run(args, &stdout3, &stderr3)
	if result != 0 {
		t.Errorf("export failed, result=%d", result)
		return
	}

	if _, err := os.Stat(out); os.IsNotExist(err) {
		t.Errorf("%s file does not exit", out)
		return
	}

	expected := "[{\"Password\":\"secret2\",\"Title\":\"A\",\"UserName\":\"admin\",\"path\":\"/1/A\"}]\n\n"
	actual := readFile(out, t)

	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
		return
	}

	// remove file
	if err := os.Remove(out); err != nil {
		t.Errorf("cannot delete %s", out)
		return
	}

	// overwrite value /1/A with random password
	args = []string{"set", "-d", db, "-p", pw, "-e", "/1/A", "-f", "UserName=admin", "-f", "Password={A32}"}
	stdout4 := strings.Builder{}
	stderr4 := strings.Builder{}
	result = Run(args, &stdout4, &stderr4)
	if result != 0 {
		t.Errorf("set failed, result=%d", result)
		return
	}

	args = []string{"get", "-d", db, "-p", pw, "-e", "/1/A", "-f", "Password"}
	stdout5 := strings.Builder{}
	stderr5 := strings.Builder{}
	result = Run(args, &stdout5, &stderr5)
	if result != 0 {
		t.Errorf("get failed, result=%d", result)
		return
	}

	if len(stdout5.String()) != 32 {
		t.Errorf("value mismatch %s\n", stdout5.String())
	}

	// remove file
	if err := os.Remove(db); err != nil {
		t.Errorf("cannot delete %s", db)
		return
	}
}
