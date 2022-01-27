package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/tobischo/gokeepasslib/v3"
)

// main command processing
// - parse commandline and open database
// - read all entries into a EntryMap structure for easy access
// - delegate command to Cmd... structures
// - save database if entries have been modified (and --dry-run is not set)
func Run(args []string, stdout io.Writer, stderr io.Writer) int {
	options := NewOptions()

	if !options.Parse(args, stderr) {
		return 1
	}

	if options.GetCmd() == "init" {
		return CmdInit(options.GetDb(), options.GetPw(), stdout, stderr)
	}

	readFile, err := os.Open(options.GetDb())
	if err != nil {
		fmt.Fprintf(stderr, "%s\n", err)
		return 1
	}
	defer readFile.Close()

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(options.GetPw())
	err = gokeepasslib.NewDecoder(readFile).Decode(db)
	if err != nil {
		fmt.Fprintf(stderr, "%s\n", err)
		return 1
	}

	db.UnlockProtectedEntries()

	if len(db.Content.Root.Groups) < 1 {
		fmt.Fprintf(stderr, "missing root group\n")
		return 1
	}

	root := &db.Content.Root.Groups[0]
	recycleBin := db.Content.Meta.RecycleBinUUID

	modified := false
	result := 0
	switch options.GetCmd() {
	case "secrets":
		entryMap := NewEntryMap(root, recycleBin)
		return CmdSecrets(entryMap, options.GetOut(), options.GetTag(), stdout, stderr) // write secrets to yaml file
	case "get":
		entryMap := NewEntryMap(root, recycleBin)
		return CmdGet(entryMap, options.GetPath(), options.GetFields()[0], stdout, stderr) // returns value in stdout
	case "set":
		modified = CmdSet(root, options.GetPath(), options.GetFields(), stdout, stderr) // writes to existing file
	case "export":
		entryMap := NewEntryMap(root, recycleBin)
		return CmdExport(entryMap, options.GetOut(), stdout, stderr) // export to json file
	case "import":
		modified, result = CmdImport(root, options.GetIn(), stdout, stderr) // import from json file
	default:
		return 1
	}

	if result != 0 {
		return result
	}

	if modified && !options.IsDryRun() {
		db.LockProtectedEntries()

		writeFile, err := os.Create(options.GetDb())
		if err != nil {
			fmt.Fprintf(stderr, "cannot create file %s\n", err)
			return 1
		}
		defer writeFile.Close()

		keepassEncoder := gokeepasslib.NewEncoder(writeFile)
		if err := keepassEncoder.Encode(db); err != nil {
			fmt.Fprintf(stderr, "cannot save database %s\n", err)
			return 1
		}
	}

	return 0
}
