package cmd

import (
	"os"
	"strings"
	"testing"
)

// get value
func TestGetUserName(t *testing.T) {
	db := "test/test.kdbx"
	pw := "1234"

	args := []string{"get", "-d", db, "-p", pw, "-e", "/entry-1", "-f", "UserName"}
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	result := Run(args, &stdout, &stderr)
	if result != 0 {
		t.Errorf("secrets failed, result=%d", result)
		return
	}

	if stdout.String() != "admin1" {
		t.Errorf("value mismatch %s\n", stdout.String())
	}
}

// get value, password taken from environment variable
func TestGetUserNameEnv(t *testing.T) {
	db := "test/test.kdbx"

	args := []string{"get", "-d", db, "-e", "/entry-1", "-f", "UserName"}
	os.Setenv("KSPASSWORD", "1234")
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	result := Run(args, &stdout, &stderr)
	os.Setenv("KSPASSWORD", "")

	if result != 0 {
		t.Errorf("secrets failed, result=%d", result)
		return
	}

	if stdout.String() != "admin1" {
		t.Errorf("value mismatch %s\n", stdout.String())
	}
}

// get value inside sub group
func TestGetFromGroup(t *testing.T) {
	db := "test/test.kdbx"
	pw := "1234"

	args := []string{"get", "-d", db, "-p", pw, "-e", "/folder-b/entry-b1", "-f", "UserName"}
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	result := Run(args, &stdout, &stderr)
	if result != 0 {
		t.Errorf("secrets failed, result=%d", result)
		return
	}

	if stdout.String() != "admin-b1" {
		t.Errorf("value mismatch %s\n", stdout.String())
	}
}

// get value inside sub group, path without leading backslash
func TestGetFromGroupWithoutSlash(t *testing.T) {
	db := "test/test.kdbx"
	pw := "1234"

	args := []string{"get", "-d", db, "-p", pw, "-e", "folder-a/entry-a1", "-f", "UserName"}
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	result := Run(args, &stdout, &stderr)
	if result != 0 {
		t.Errorf("secrets failed, result=%d", result)
		return
	}

	if stdout.String() != "admin-a1" {
		t.Errorf("value mismatch %s\n", stdout.String())
	}
}

// get password value
func TestGetPassword(t *testing.T) {
	db := "test/test.kdbx"
	pw := "1234"

	args := []string{"get", "-d", db, "-p", pw, "-e", "/entry-1", "-f", "Password"}
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	result := Run(args, &stdout, &stderr)
	if result != 0 {
		t.Errorf("get failed, result=%d", result)
		return
	}

	if stdout.String() != "abcd" {
		t.Errorf("value mismatch %s\n", stdout.String())
	}
}

// get password value using POSIX option (--password instead of -p)
func TestGetPasswordPosix(t *testing.T) {
	db := "test/test.kdbx"
	pw := "1234"

	args := []string{"get", "--database", db, "--password", pw, "--entry", "/entry-1", "--field", "Password"}
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	result := Run(args, &stdout, &stderr)
	if result != 0 {
		t.Errorf("get failed, result=%d", result)
		return
	}

	if stdout.String() != "abcd" {
		t.Errorf("value mismatch %s\n", stdout.String())
	}
}

// get value of non existing entry
func TestGetInvalidPath(t *testing.T) {
	db := "test/test.kdbx"
	pw := "1234"

	args := []string{"get", "-d", db, "-p", pw, "-e", "/invalid", "-f", "Password"}
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	result := Run(args, &stdout, &stderr)
	if result != 1 {
		t.Errorf("get did not fail")
		return
	}

	expected := "path '/invalid' does not exist\n"
	actual := stderr.String()
	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
		return
	}
}

// get value of non existing field
func TestGetInvalidField(t *testing.T) {
	db := "test/test.kdbx"
	pw := "1234"

	args := []string{"get", "-d", db, "-p", pw, "-e", "/entry-1", "-f", "Invalid"}
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	result := Run(args, &stdout, &stderr)
	if result != 1 {
		t.Errorf("get did not fail")
		return
	}

	expected := "field 'Invalid' does not exist in path '/entry-1'\n"
	actual := stderr.String()
	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
		return
	}
}

// get value, database file is missing
func TestGetInvalidDatabase(t *testing.T) {
	db := "test/noexist.kdbx"
	pw := "1234"

	args := []string{"get", "-d", db, "-p", pw, "-e", "/entry-1", "-f", "Password"}
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	result := Run(args, &stdout, &stderr)
	if result == 0 {
		t.Errorf("run must fail")
		return
	}
}

// get value, incorrect password
func TestGetInvalidPassword(t *testing.T) {
	db := "test/test.kdbx"
	pw := "abcd"

	args := []string{"get", "-d", db, "-p", pw, "-e", "/entry-1", "-f", "Password"}
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	result := Run(args, &stdout, &stderr)
	if result == 0 {
		t.Errorf("run must fail")
		return
	}
}
