package cmd

import (
	"encoding/json"
	"io"
	"strings"
)

// export complete database as JSON
func CmdExport(entryMap *EntryMap, out string, stdout io.Writer, stderr io.Writer) int {
	list := make([]map[string]string, 0)
	paths := entryMap.GetPaths()
	for i := 0; i < len(paths); i++ {
		path := paths[i]

		entry := make(map[string]string)
		entry["path"] = path

		values, ok := entryMap.GetValues(path)
		if ok {
			for key := range values {
				entry[key] = values[key]
			}
		}

		list = append(list, entry)
	}

	str := strings.Builder{}
	enc := json.NewEncoder(&str)
	enc.Encode(list)
	lines := strings.Split(str.String(), "\n")

	return writeFile(out, &lines, stderr)
}
