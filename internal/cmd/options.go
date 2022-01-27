package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

var version = "0.0.0" // application version (must be set in build)
var commit = "local"  // commit hash (must be set in build)

// arrayFlags collects multiple string options into array (used for option 'field')
type arrayFlags []string

func (arr *arrayFlags) String() string {
	return strings.Join(*arr, ",")
}

func (arr *arrayFlags) Set(value string) error {
	*arr = append(*arr, value)
	return nil
}

func (arr *arrayFlags) Type() string {
	return "string"
}

// stores all commandline options
type Options struct {
	flags  *flag.FlagSet
	cmd    string
	db     string
	pw     string
	path   string
	tag    string
	fields arrayFlags
	out    string
	in     string
	dryRun bool
}

func NewOptions() Options {
	options := Options{}
	options.flags = flag.NewFlagSet("", flag.ContinueOnError)
	return options
}

// parse command line and make options availabe via getters
func (options *Options) parseCmd(args []string) ([]string, error) {
	if len(args) < 1 {
		return make([]string, 0), errors.New(options.getUsage())
	}

	options.cmd = args[0]

	if options.cmd != "secrets" && options.cmd != "get" && options.cmd != "export" && options.cmd != "import" && options.cmd != "init" && options.cmd != "set" {
		return make([]string, 0), errors.New("unknown command " + options.cmd)
	}

	options.cmd = args[0]
	return args[1:], nil
}

func (options *Options) parseOptions(args []string, stderr io.Writer) bool {
	dbFlag := options.flags.StringP("database", "d", "", "keepass 2.30 file")
	pwFlag := options.flags.StringP("password", "p", "", "password")
	pathFlag := options.flags.StringP("entry", "e", "", "path of keepass entry")
	tagFlag := options.flags.StringP("tag", "t", "", "filter by tag")
	outFlag := options.flags.StringP("out", "o", "", "output filename")
	inFlag := options.flags.StringP("in", "i", "", "input filename")
	dryRunFlag := options.flags.BoolP("dry-run", "", false, "do not modify database")
	options.flags.VarP(&options.fields, "field", "f", "field name and value")

	err := options.flags.Parse(args)
	if err != nil {
		fmt.Fprintf(stderr, "%s\n", err)
		return false
	}

	options.db = *dbFlag
	options.pw = *pwFlag
	options.path = *pathFlag
	options.tag = *tagFlag
	options.out = *outFlag
	options.in = *inFlag
	options.dryRun = *dryRunFlag

	if options.pw == "" && os.Getenv("KSPASSWORD") != "" {
		options.pw = os.Getenv("KSPASSWORD")
	}

	return true
}

func (options *Options) Parse(args []string, stderr io.Writer) bool {
	optionArgs, err := options.parseCmd(args)
	if err != nil {
		fmt.Fprintf(stderr, "%s\n", err)
		return false
	}

	if !options.parseOptions(optionArgs, stderr) {
		return false
	}

	if !options.verify(stderr) {
		return false
	}

	return true
}

func (options *Options) getUsage() string {
	usage := strings.Builder{}

	usage.WriteString(fmt.Sprintf("keepass-secret %s (%s)\n", version, commit))
	usage.WriteString("usage: keepass-secret secrets -d keepass.kdbx -p 1234 -o secrets.yaml [--tag abc]\n")
	usage.WriteString("       keepass-secret get     -d keepass.kdbx -p 1234 -e /entry-1 -f Password\n")
	usage.WriteString("       keepass-secret set     -d keepass.kdbx -p 1234 -e /entry-1 -f Password=1234 -f UserName=admin\n")
	usage.WriteString("       keepass-secret export  -d keepass.kdbx -p 1234 -o export.json\n")
	usage.WriteString("       keepass-secret import  -d keepass.kdbx -p 1234 -i import.json [--dry-run]\n")
	usage.WriteString("       keepass-secret init    -d keepass.kdbx -p 1234\n")
	usage.WriteString("\n")
	usage.WriteString("The password can also be set via the environment variable 'KSPASSWORD'\n")

	return usage.String()
}

// check presence of common mandatory options
func (options *Options) verifyCommon(stderr io.Writer) bool {
	if options.db == "" {
		fmt.Fprintf(stderr, "missing -d/--database parameter\n")
		return false
	}

	if options.pw == "" {
		fmt.Fprintf(stderr, "missing -p/--password parameter\n")
		return false
	}

	return true
}

// check presence of mandatory options for export and secrets command
func (options *Options) verifyExportOrSecrets(stderr io.Writer) bool {
	if options.out == "" {
		fmt.Fprintf(stderr, "missing -o/--out parameter\n")
		return false
	}

	return true
}

// check presence of mandatory options for get command
func (options *Options) verifyGet(stderr io.Writer) bool {
	if options.path == "" {
		fmt.Fprintf(stderr, "missing -e/--entry parameter\n")
		return false
	}

	if len(options.fields) != 1 {
		fmt.Fprintf(stderr, "missing -f/--field parameter\n")
		return false
	}

	return true
}

// check presence of mandatory options for set command
func (options *Options) verifySet(stderr io.Writer) bool {
	if options.path == "" {
		fmt.Fprintf(stderr, "missing -e/--entry parameter\n")
		return false
	}

	if len(options.fields) == 0 {
		fmt.Fprintf(stderr, "missing -f/--field parameter\n")
		return false
	}

	return true
}

// check presence of mandatory options for import command
func (options *Options) verifyImport(stderr io.Writer) bool {
	if options.in == "" {
		fmt.Fprintf(stderr, "missing -i/--in parameter\n")
		return false
	}

	return true
}

// check plausibility of commandline options
func (options *Options) verify(stderr io.Writer) bool {
	if !options.verifyCommon(stderr) {
		return false
	}

	if options.cmd == "secrets" && !options.verifyExportOrSecrets(stderr) {
		return false
	}

	if options.cmd == "get" && !options.verifyGet(stderr) {
		return false
	}

	if options.cmd == "set" && !options.verifySet(stderr) {
		return false
	}

	if options.cmd == "export" && !options.verifyExportOrSecrets(stderr) {
		return false
	}

	if options.cmd == "import" && !options.verifyImport(stderr) {
		return false
	}

	return true
}

func (options *Options) GetCmd() string {
	return options.cmd
}

func (options *Options) GetDb() string {
	return options.db
}

func (options *Options) GetPw() string {
	return options.pw
}

func (options *Options) GetPath() string {
	return options.path
}

func (options *Options) GetOut() string {
	return options.out
}

func (options *Options) GetIn() string {
	return options.in
}

func (options *Options) IsDryRun() bool {
	return options.dryRun
}

func (options *Options) GetTag() string {
	return options.tag
}

func (options *Options) GetFields() []string {
	return options.fields
}
