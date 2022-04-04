package cmd

import (
	"github.com/tobischo/gokeepasslib/v3"
)

// model complete KeePass database as flat list of entries
// each entry is defined by its path an a key/value map of the entry fields
type EntryMap struct {
	paths   []string         // all full qualified paths (used for iteration)
	entries map[string]Entry // map path to entry (key/value map)
}

// recursively process group (folder of entries) and store entries in map
func (entryMap *EntryMap) processGroup(group *gokeepasslib.Group, path string, recycleBin gokeepasslib.UUID, binaries gokeepasslib.Binaries) {
	if recycleBin.Compare(group.UUID) {
		return // ignore entries from recycle bin
	}

	for i := 0; i < len(group.Entries); i++ {
		entry := &group.Entries[i]

		values := NewEntry()
		for j := 0; j < len(entry.Values); j++ {
			value := &entry.Values[j]
			values.SetValue(value.Key, value.Value.Content)
		}

		for j := 0; j < len(entry.Binaries); j++ {
			value := &entry.Binaries[j]
			binary := binaries.Find(value.Value.ID)
			if binary != nil {
				content, err := binary.GetContent()
				if err == nil {
					values.SetBinary(value.Name, []byte(content))
				}
			}
		}

		key := path + entry.GetTitle()
		entryMap.entries[key] = *values

		entryMap.paths = append(entryMap.paths, key)
	}

	for i := 0; i < len(group.Groups); i++ {
		childGroup := &group.Groups[i]
		entryMap.processGroup(childGroup, path+childGroup.Name+"/", recycleBin, binaries)
	}
}

func (entryMap *EntryMap) GetPaths() []string {
	return entryMap.paths
}

func (entryMap *EntryMap) GetValues(path string) (Entry, bool) {
	values, ok := entryMap.entries[path]
	return values, ok
}

func NewEntryMap(db *gokeepasslib.Database) *EntryMap {
	entryMap := EntryMap{paths: make([]string, 0), entries: make(map[string]Entry)}

	root := &db.Content.Root.Groups[0]
	recycleBin := db.Content.Meta.RecycleBinUUID
	binaries := db.Content.Meta.Binaries
	entryMap.processGroup(root, "/", recycleBin, binaries)

	return &entryMap
}
