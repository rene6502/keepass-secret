package cmd

import "strings"

const prefix = "secret-"

// models the contents of the 'Notes' field as a key/value map
// each line in the 'Notes' field is treated as an key/value pair
// if it starts with the "secret-" prefix
// the stored key exclude the prefix
// e.g. the line "secret-postgresql-user=UserName" will be stored
// with key="postgresql-user" and value="UserName"
type Notes struct {
	keys    []string
	entries map[string]string
}

func (notes *Notes) Get(key string) string {
	return notes.entries[key]
}

func (notes *Notes) GetKeys() []string {
	result := make([]string, 0)

	for i := 0; i < len(notes.keys); i++ {
		key := notes.keys[i]
		if key != "type" && key != "tags" {
			result = append(result, key)
		}
	}

	return result
}

func NewNotes(values Entry) *Notes {
	notes := Notes{make([]string, 0), make(map[string]string)}

	if notesStr, ok := values.GetValue("Notes"); ok {
		lines := strings.Split(strings.ReplaceAll(notesStr, "\r", ""), "\n")
		for i := 0; i < len(lines); i++ {
			line := lines[i]
			if strings.HasPrefix(line, prefix) {
				pos := strings.Index(line, "=")
				if pos > len(prefix) {
					key := line[len(prefix):pos]
					value := line[pos+1:]
					notes.entries[key] = value
					notes.keys = append(notes.keys, key)
				}
			}
		}
	}

	return &notes
}
