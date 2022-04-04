package cmd

import (
	"os"
	"testing"

	"github.com/tobischo/gokeepasslib/v3"
)

func TestEntry(t *testing.T) {
	in := "test/test.kdbx"
	pw := "1234"

	readFile, err := os.Open(in)
	if err != nil {
		t.Errorf("cannot open %s", in)
		return
	}
	defer readFile.Close()

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(pw)
	err = gokeepasslib.NewDecoder(readFile).Decode(db)
	if err != nil {
		t.Errorf("cannot decode %s", in)
		return
	}

	db.UnlockProtectedEntries()

	entryMap := NewEntryMap(db)

	entry, ok := entryMap.GetValues("/binary")
	if !ok {
		t.Errorf("cannot get /binary entry")
		return
	}

	if userName, _ := entry.GetValue("UserName"); userName != "admin0" {
		t.Errorf("invalid UserName %s", userName)
		return
	}

	if password, _ := entry.GetValue("Password"); password != "ABCD" {
		t.Errorf("invalid Password %s", password)
		return
	}

	binaries := entry.GetBinaries()
	if len(binaries) != 2 {
		t.Errorf("incorrect number of binaries %d", len(binaries))
		return
	}

	file1, ok1 := entry.GetBinary("file1.bin")
	if !ok1 {
		t.Errorf("missing binary file1.bin")
		return
	}

	if !compareBinary(file1, "test/file1.bin", t) {
		return
	}

	file2, ok2 := entry.GetBinary("file2.bin")
	if !ok2 {
		t.Errorf("missing binary file2.bin")
		return
	}

	if !compareBinary(file2, "test/file2.bin", t) {
		return
	}
}
