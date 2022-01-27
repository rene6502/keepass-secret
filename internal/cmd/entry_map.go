package cmd

import (
	"github.com/tobischo/gokeepasslib/v3"
)

// model complete KeePass database as flat list of entries
// each entry is defined by its path an a key/value map of the entry fields
type EntryMap struct {
	paths      []string                     // all full qualified paths (used for iteration)
	entries    map[string]map[string]string // map path to entry (key/value map)
	recycleBin gokeepasslib.UUID
}

// recursively process group (folder of entries) and store entries in map
func (entryMap *EntryMap) processGroup(group *gokeepasslib.Group, path string) {
	if entryMap.recycleBin.Compare(group.UUID) {
		return // ignore entries from recycle bin
	}

	for i := 0; i < len(group.Entries); i++ {
		entry := &group.Entries[i]

		var values map[string]string = make(map[string]string)
		for j := 0; j < len(entry.Values); j++ {
			value := &entry.Values[j]
			values[value.Key] = value.Value.Content
		}

		key := path + entry.GetTitle()
		entryMap.entries[key] = values
		entryMap.paths = append(entryMap.paths, key)
	}

	for i := 0; i < len(group.Groups); i++ {
		childGroup := &group.Groups[i]
		entryMap.processGroup(childGroup, path+childGroup.Name+"/")
	}
}

func (entryMap *EntryMap) GetPaths() []string {
	return entryMap.paths
}

func (entryMap *EntryMap) GetValues(path string) (map[string]string, bool) {
	values, ok := entryMap.entries[path]
	return values, ok
}

func NewEntryMap(group *gokeepasslib.Group, recycleBin gokeepasslib.UUID) *EntryMap {
	entryMap := EntryMap{make([]string, 0), make(map[string]map[string]string), recycleBin}

	entryMap.processGroup(group, "/")

	return &entryMap
}
