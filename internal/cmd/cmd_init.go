package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tobischo/gokeepasslib/v3"
	w "github.com/tobischo/gokeepasslib/v3/wrappers"
)

// creates new KeePass database file
func CmdInit(db string, pw string, stdout io.Writer, stderr io.Writer) int {
	file, err := os.Create(db)
	if err != nil {
		fmt.Fprintf(stderr, "cannot create file %s\n", err)
		return 1
	}

	defer file.Close()

	// create root group (use filename without extension as group name)
	rootGroup := gokeepasslib.NewGroup()
	rootGroup.Times.ExpiryTime = &w.TimeWrapper{Time: time.Now()}
	rootGroup.Times.Expires = w.NewBoolWrapper(false)

	basename := filepath.Base(db)
	rootGroup.Name = strings.TrimSuffix(basename, filepath.Ext(basename))
	root := &gokeepasslib.RootData{Groups: []gokeepasslib.Group{rootGroup}}
	meta := gokeepasslib.NewMetaData()
	meta.Generator = "keepass-secret"
	meta.DatabaseName = rootGroup.Name
	meta.DatabaseNameChanged = &w.TimeWrapper{Time: time.Now()}
	meta.DatabaseDescriptionChanged = &w.TimeWrapper{Time: time.Now()}
	meta.DefaultUserNameChanged = &w.TimeWrapper{Time: time.Now()}
	meta.RecycleBinChanged = &w.TimeWrapper{Time: time.Now()}
	meta.EntryTemplatesGroupChanged = &w.TimeWrapper{Time: time.Now()}

	content := &gokeepasslib.DBContent{Meta: meta, Root: root}

	// now create the database containing the root group
	database := &gokeepasslib.Database{
		Header:      gokeepasslib.NewHeader(),
		Credentials: gokeepasslib.NewPasswordCredentials(pw),
		Content:     content,
	}

	// Lock entries using stream cipher
	database.LockProtectedEntries()

	// and encode it into the file
	keepassEncoder := gokeepasslib.NewEncoder(file)
	if err := keepassEncoder.Encode(database); err != nil {
		fmt.Fprintf(stderr, "cannot save database %s\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "successfully created kdbx file: %s\n", db)

	return 0 // success
}
