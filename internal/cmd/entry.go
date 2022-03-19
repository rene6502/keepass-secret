package cmd

type Entry struct {
	values map[string]string
}

func NewEntry() *Entry {
	return &Entry{values: make(map[string]string)}
}

func (entry *Entry) GetValue(name string) (string, bool) {
	value, ok := entry.values[name]
	return value, ok
}

func (entry *Entry) SetValue(name string, value string) {
	entry.values[name] = value
}

func (entry *Entry) GetNames() []string {
	names := make([]string, 0, len(entry.values))
	for name := range entry.values {
		names = append(names, name)
	}
	return names
}
