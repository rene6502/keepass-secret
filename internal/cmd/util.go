package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/tobischo/gokeepasslib/v3"
	w "github.com/tobischo/gokeepasslib/v3/wrappers"
)

// find group with specified name, return nil if not found
func findGroup(node *gokeepasslib.Group, name string) *gokeepasslib.Group {
	for i := 0; i < len(node.Groups); i++ {
		group := &node.Groups[i]
		if group.Name == name {
			return group
		}
	}

	return nil
}

// returns true if entry with specified title exist in group
func entryExists(group *gokeepasslib.Group, title string) bool {
	for i := 0; i < len(group.Entries); i++ {
		entry := group.Entries[i]
		for j := 0; j < len(entry.Values); j++ {
			value := entry.Values[j]
			if value.Key == "Title" && value.Value.Content == title {
				return true // entry exists
			}
		}
	}

	return false // entry does not exist
}

// delete entry with specified title
// do nothing if entry cannot be found
func deleteEntry(group *gokeepasslib.Group, title string) {
	for i := 0; i < len(group.Entries); i++ {
		entry := group.Entries[i]
		for j := 0; j < len(entry.Values); j++ {
			value := entry.Values[j]
			if value.Key == "Title" && value.Value.Content == title {
				// remove entry
				copy(group.Entries[i:], group.Entries[i+1:])
				group.Entries = group.Entries[:len(group.Entries)-1]
				return
			}
		}
	}
}

// create entry with specified title and fill in fields
// caller is responsible for saving it to database
func createEntry(title string, fields []string, stdout io.Writer, stderr io.Writer) *gokeepasslib.Entry {
	entry := gokeepasslib.NewEntry()
	entry.Times.ExpiryTime = &w.TimeWrapper{Time: time.Now()}
	entry.Times.Expires = w.NewBoolWrapper(false)

	entry.Values = append(entry.Values, gokeepasslib.ValueData{Key: "Title", Value: gokeepasslib.V{Content: title}})

	for j := 0; j < len(fields); j++ {
		pos := strings.Index(fields[j], "=")
		if pos != -1 {
			key := fields[j][:pos]
			value := fields[j][pos+1:]

			if key == "Password" {
				if strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
					pattern := value[1 : len(value)-1]
					value = createPasswordFromPattern(pattern, stdout, stderr)
				}
				entry.Values = append(entry.Values, gokeepasslib.ValueData{Key: key, Value: gokeepasslib.V{Content: value, Protected: w.NewBoolWrapper(true)}})
			} else if key == "Notes" {
				value = strings.ReplaceAll(value, "\\n", "\n")
				entry.Values = append(entry.Values, gokeepasslib.ValueData{Key: key, Value: gokeepasslib.V{Content: value}})
			} else {
				entry.Values = append(entry.Values, gokeepasslib.ValueData{Key: key, Value: gokeepasslib.V{Content: value}})
			}
		}
	}

	return &entry
}

// creates and entry if it does not exist
// if it exists create it after deletion = update
// "xyz created" or "abc updated" is written to stdout
// to update an existing entry overwrite must be true, otherwise the changes will be ignored
// (command import will not overwrite, command set will overwrite)
func createOrUpdateEntry(root *gokeepasslib.Group, path string, fields []string, overwrite bool, stdout io.Writer, stderr io.Writer) bool {
	group := root
	path = strings.TrimPrefix(path, "/") // remove leading "/"
	groupNames := make([]string, 0)
	title := path

	if strings.Contains(path, "/") {
		groupNames = strings.Split(path, "/")
		title = groupNames[len(groupNames)-1]
		groupNames = groupNames[:len(groupNames)-1]
	}

	for i := 0; i < len(groupNames); i++ {
		groupName := groupNames[i]
		subGroup := findGroup(group, groupName)
		if subGroup != nil {
			group = subGroup
		} else {
			newGroup := gokeepasslib.NewGroup()
			newGroup.Name = groupName
			group.Groups = append(group.Groups, newGroup)
			group = findGroup(group, groupName)
		}
	}

	if group != nil {
		exists := entryExists(group, title)
		if exists && !overwrite {
			return false // entry exists already and overwrite is not allowed
		}

		if exists {
			deleteEntry(group, title)
		}

		entry := createEntry(title, fields, stdout, stderr)
		group.Entries = append(group.Entries, *entry)

		if exists {
			fmt.Fprintf(stdout, "%s updated\n", path)
		} else {
			fmt.Fprintf(stdout, "%s created\n", path)
		}

		return true // modified
	}

	return false // not modified
}

// write string to file
// on success return result=0
// on error a message is written to stderr and result=1 is returned
func writeFile(out string, pLines *[]string, stderr io.Writer) int {
	lines := *pLines

	file, err := os.Create(out)
	if err != nil {
		fmt.Fprintf(stderr, "%s\n", err)
		return 1 // failure
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for i := 0; i < len(lines); i++ {
		_, err = fmt.Fprintln(writer, lines[i])
		if err != nil {
			fmt.Fprintf(stderr, "%s\n", err)
			file.Close()
			return 1 // failure
		}
	}

	err = writer.Flush()
	if err != nil {
		fmt.Fprintf(stderr, "%s\n", err)
		return 1 // failure
	}

	err = file.Close()
	if err != nil {
		fmt.Fprintf(stderr, "%s\n", err)
		return 1 // failure
	}

	return 0 // success
}
