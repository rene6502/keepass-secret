package cmd

import (
	"os"
	"strings"
	"testing"
)

// export all secrets
func TestSecrets(t *testing.T) {
	db := "test/test.kdbx"
	pw := "1234"
	out := "test/test.yaml"

	args := []string{"secrets", "-d", db, "-p", pw, "-o", out}
	stdout0 := strings.Builder{}
	stderr0 := strings.Builder{}
	result := Run(args, &stdout0, &stderr0)
	if result != 0 {
		t.Errorf("secrets failed, result=%d", result)
		return
	}

	if _, err := os.Stat(out); os.IsNotExist(err) {
		t.Errorf("%s file does not exit", out)
		return
	}

	if !compareFiles(out, "test/test-all.yaml", t) {
		t.Fail()
	}

	// remove file
	if err := os.Remove(out); err != nil {
		t.Errorf("cannot delete %s", out)
		return
	}
}

// export all secrets with tag "taga"
func TestSecretsTagA(t *testing.T) {
	db := "test/test.kdbx"
	pw := "1234"
	out := "test/test.yaml"

	args := []string{"secrets", "-d", db, "-p", pw, "-o", out, "-t", "taga"}
	stdout0 := strings.Builder{}
	stderr0 := strings.Builder{}
	result := Run(args, &stdout0, &stderr0)
	if result != 0 {
		t.Errorf("secrets failed, result=%d", result)
		return
	}

	if _, err := os.Stat(out); os.IsNotExist(err) {
		t.Errorf("%s file does not exit", out)
		return
	}

	if !compareFiles(out, "test/test-taga.yaml", t) {
		t.Fail()
	}

	// remove file
	if err := os.Remove(out); err != nil {
		t.Errorf("cannot delete %s", out)
		return
	}
}

// export all secrets with tag "tagb"
func TestSecretsTagB(t *testing.T) {
	db := "test/test.kdbx"
	pw := "1234"
	out := "test/test.yaml"

	args := []string{"secrets", "-d", db, "-p", pw, "-o", out, "-t", "tagb"}
	stdout0 := strings.Builder{}
	stderr0 := strings.Builder{}
	result := Run(args, &stdout0, &stderr0)
	if result != 0 {
		t.Errorf("secrets failed, result=%d", result)
		return
	}

	if _, err := os.Stat(out); os.IsNotExist(err) {
		t.Errorf("%s file does not exit", out)
		return
	}

	if !compareFiles(out, "test/test-tagb.yaml", t) {
		t.Fail()
	}

	// remove file
	if err := os.Remove(out); err != nil {
		t.Errorf("cannot delete %s", out)
		return
	}
}

// export all secrets with tag "tagc"
func TestSecretsTagC(t *testing.T) {
	db := "test/test.kdbx"
	pw := "1234"
	out := "test/test.yaml"

	args := []string{"secrets", "-d", db, "-p", pw, "-o", out, "-t", "tagc"}
	stdout0 := strings.Builder{}
	stderr0 := strings.Builder{}
	result := Run(args, &stdout0, &stderr0)
	if result != 0 {
		t.Errorf("secrets failed, result=%d", result)
		return
	}

	if _, err := os.Stat(out); os.IsNotExist(err) {
		t.Errorf("%s file does not exit", out)
		return
	}

	if !compareFiles(out, "test/test-tagc.yaml", t) {
		t.Fail()
	}

	// remove file
	if err := os.Remove(out); err != nil {
		t.Errorf("cannot delete %s", out)
		return
	}
}

// export secrets, docker secret with missing title
func TestSecretsDockerSecretMissingTitle(t *testing.T) {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	lines := make([]string, 0)
	values := NewEntry()
	createDockerSecret("e0", "", *values, &lines, &stdout, &stderr)

	expected := "missing title for entry 'e0'\n"
	actual := stderr.String()
	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
		return
	}
}

// export secrets, docker secret with missing username
func TestSecretsDockerSecretMissingUserName(t *testing.T) {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	lines := make([]string, 0)
	values := NewEntry()
	values.SetValue("Title", "Title")
	createDockerSecret("e0", "", *values, &lines, &stdout, &stderr)

	expected := "missing UserName for entry 'e0'\n"
	actual := stderr.String()
	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
		return
	}
}

// export secrets, docker secret with missing password
func TestSecretsDockerSecretMissingPassword(t *testing.T) {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	lines := make([]string, 0)
	values := NewEntry()
	values.SetValue("Title", "Title")
	values.SetValue("UserName", "UserName")
	createDockerSecret("e0", "", *values, &lines, &stdout, &stderr)

	expected := "missing Password for entry 'e0'\n"
	actual := stderr.String()
	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
		return
	}
}

// export secrets, docker secret with missing URL
func TestSecretsDockerSecretMissingURL(t *testing.T) {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	lines := make([]string, 0)
	values := NewEntry()
	values.SetValue("Title", "Title")
	values.SetValue("UserName", "UserName")
	values.SetValue("Password", "Password")
	createDockerSecret("e0", "", *values, &lines, &stdout, &stderr)

	expected := "missing URL for entry 'e0'\n"
	actual := stderr.String()
	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
		return
	}
}

// export secrets, opaque secret with missing title
func TestSecretsOpaqueSecretMissingTitle(t *testing.T) {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	lines := make([]string, 0)
	values := NewEntry()
	notes := NewNotes(*values)
	createOpaqueSecret("e0", "", notes, *values, &lines, &stdout, &stderr)

	expected := "missing title for entry 'e0'\n"
	actual := stderr.String()
	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
		return
	}
}

// export secrets, opaque secret with missing value
func TestSecretsOpaqueSecretMissingValue(t *testing.T) {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	lines := make([]string, 0)
	values := NewEntry()
	values.SetValue("Title", "Title")
	values.SetValue("Notes", "secret-password=Password")
	notes := NewNotes(*values)
	createOpaqueSecret("e1", "", notes, *values, &lines, &stdout, &stderr)

	expected := "entry 'e1' does not contain value 'Password'\n"
	actual := stderr.String()
	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
		return
	}
}
