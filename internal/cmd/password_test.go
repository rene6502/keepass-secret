package cmd

import (
	"regexp"
	"strings"
	"testing"
)

// invalid pattern
func TestInvalidPattern(t *testing.T) {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	password := createPasswordFromPattern("ABC", &stdout, &stderr)
	r := "^[A-Za-z0-9]{32}$"
	match, _ := regexp.MatchString(r, password)
	if !match {
		t.Errorf("createPasswordFromPattern failed, expected=%s actual=%s", r, password)
		return
	}

	expected := "unknown password pattern ABC, fallback to A32\n"
	actual := stderr.String()
	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
	}
}

// lower-case hex characters
func TestHexLowerCase(t *testing.T) {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	actual := createPasswordFromPattern("h20", &stdout, &stderr)

	expected := "^[0-9a-f]{20}$"
	match, _ := regexp.MatchString(expected, actual)
	if !match {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
	}
}

// upper-case hex characters
func TestHexUpperCase(t *testing.T) {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	actual := createPasswordFromPattern("H20", &stdout, &stderr)

	expected := "^[0-9A-F]{20}$"
	match, _ := regexp.MatchString(expected, actual)
	if !match {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
	}
}

// letter
func TestLetter(t *testing.T) {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	actual := createPasswordFromPattern("L32", &stdout, &stderr)

	expected := "^[A-Za-z]{32}$"
	match, _ := regexp.MatchString(expected, actual)
	if !match {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
	}
}

// alphanumeric
func TestAlphanumeric(t *testing.T) {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	actual := createPasswordFromPattern("A40", &stdout, &stderr)

	expected := "^[A-Za-z0-9]{40}$"
	match, _ := regexp.MatchString(expected, actual)
	if !match {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
	}
}

// printable
func TestPrintable(t *testing.T) {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	actual := createPasswordFromPattern("S32", &stdout, &stderr)

	expected := "^[\\x21-\\x7e]{32}$"
	match, _ := regexp.MatchString(expected, actual)
	if !match {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
	}
}
