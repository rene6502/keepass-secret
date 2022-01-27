package cmd

import (
	"strings"
	"testing"
)

// no arguments -> usage printed
func TestOptionsEmpty(t *testing.T) {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	args := []string{}
	result := Run(args, &stdout, &stderr)

	if result == 0 {
		t.Errorf("run must fail")
		return
	}

	options := NewOptions()
	expected := options.getUsage() + "\n"
	actual := stderr.String()
	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
		return
	}
}

// invalid command
func TestOptionsInvalidCmd(t *testing.T) {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	args := []string{"test", "-d", "test.kdbx", "-p", "1234"}
	result := Run(args, &stdout, &stderr)

	if result == 0 {
		t.Errorf("run must fail")
		return
	}

	expected := "unknown command test\n"
	actual := stderr.String()
	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
		return
	}
}

// invalid option
func TestOptionsInvalidOption(t *testing.T) {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	args := []string{"get", "--invalid"}
	result := Run(args, &stdout, &stderr)

	if result == 0 {
		t.Errorf("run must fail")
		return
	}

	expected := "unknown flag: --invalid\n"
	actual := stderr.String()
	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
		return
	}
}

// missing --database option
func TestOptionsNoDb(t *testing.T) {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	args := []string{"get"}
	result := Run(args, &stdout, &stderr)

	if result == 0 {
		t.Errorf("run must fail")
		return
	}

	expected := "missing -d/--database parameter\n"
	actual := stderr.String()
	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
		return
	}
}

// missing --password option
func TestOptionsNoPw(t *testing.T) {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	args := []string{"get", "-d", "test.kdbx"}
	result := Run(args, &stdout, &stderr)

	if result == 0 {
		t.Errorf("run must fail")
		return
	}

	expected := "missing -p/--password parameter\n"
	actual := stderr.String()
	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
		return
	}
}

// missing --path option
func TestOptionsGetCmdNoPath(t *testing.T) {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	args := []string{"get", "-d", "test.kdbx", "-p", "1234"}
	result := Run(args, &stdout, &stderr)

	if result == 0 {
		t.Errorf("run must fail")
		return
	}

	expected := "missing -e/--entry parameter\n"
	actual := stderr.String()
	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
		return
	}
}

// missing --field option
func TestOptionsGetCmdNoField(t *testing.T) {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	args := []string{"get", "-d", "test.kdbx", "-p", "1234", "-e", "/entry-1"}
	result := Run(args, &stdout, &stderr)

	if result == 0 {
		t.Errorf("run must fail")
		return
	}

	expected := "missing -f/--field parameter\n"
	actual := stderr.String()
	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
		return
	}
}

// missing --entry option
func TestOptionsSetCmdNoPath(t *testing.T) {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	args := []string{"set", "-d", "test.kdbx", "-p", "1234"}
	result := Run(args, &stdout, &stderr)

	if result == 0 {
		t.Errorf("run must fail")
		return
	}

	expected := "missing -e/--entry parameter\n"
	actual := stderr.String()
	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
		return
	}
}

// missing --field option
func TestOptionsSetCmdNoField(t *testing.T) {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	args := []string{"set", "-d", "test.kdbx", "-p", "1234", "-e", "/entry-1"}
	result := Run(args, &stdout, &stderr)

	if result == 0 {
		t.Errorf("run must fail")
		return
	}

	expected := "missing -f/--field parameter\n"
	actual := stderr.String()
	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
		return
	}
}

// missing --out option
func TestOptionsSecretCmdNoOut(t *testing.T) {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	args := []string{"secrets", "-d", "test.kdbx", "-p", "1234"}
	result := Run(args, &stdout, &stderr)

	if result == 0 {
		t.Errorf("run must fail")
		return
	}

	expected := "missing -o/--out parameter\n"
	actual := stderr.String()
	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
		return
	}
}

// missing --out option
func TestOptionsExportNoOut(t *testing.T) {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	args := []string{"export", "-d", "test.kdbx", "-p", "1234"}
	result := Run(args, &stdout, &stderr)

	if result == 0 {
		t.Errorf("run must fail")
		return
	}

	expected := "missing -o/--out parameter\n"
	actual := stderr.String()
	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
		return
	}
}

// missing --in option
func TestOptionsImportNoIn(t *testing.T) {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	args := []string{"import", "-d", "test.kdbx", "-p", "1234"}
	result := Run(args, &stdout, &stderr)

	if result == 0 {
		t.Errorf("run must fail")
		return
	}

	expected := "missing -i/--in parameter\n"
	actual := stderr.String()
	if expected != actual {
		t.Errorf("expected: %s", expected)
		t.Errorf("actual:   %s", actual)
		return
	}
}

// call arrayFlags.Type() method
func TestFieldType(t *testing.T) {
	options := NewOptions()
	stderr := strings.Builder{}
	options.parseOptions([]string{"-f", "name=value"}, &stderr)
	if len(options.fields) != 1 {
		t.Errorf("invalid field length %d", len(options.fields))
	}

	if options.fields.Type() != "string" {
		t.Errorf("invalid type %s", options.fields.Type())
	}
}
