package cmd

type Entry struct {
	values   map[string]string
	binaries map[string][]byte
}

func NewEntry() *Entry {
	return &Entry{values: make(map[string]string), binaries: make(map[string][]byte)}
}

func (entry *Entry) SetValue(name string, value string) {
	entry.values[name] = value
}

func (entry *Entry) GetValue(name string) (string, bool) {
	value, ok := entry.values[name]
	return value, ok
}

func (entry *Entry) SetBinary(name string, value []byte) {
	entry.binaries[name] = value
}

func (entry *Entry) GetBinary(name string) ([]byte, bool) {
	value, ok := entry.binaries[name]
	return value, ok
}

func (entry *Entry) GetNames() []string {
	names := make([]string, 0, len(entry.values))
	for name := range entry.values {
		names = append(names, name)
	}
	return names
}

func (entry *Entry) GetBinaries() []string {
	names := make([]string, 0, len(entry.binaries))
	for name := range entry.binaries {
		names = append(names, name)
	}
	return names
}
