package cmd

import (
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

// export all marked entries as Kubernetes secrets YAML
// supports  opaque (regular) and docker secrets
func CmdSecrets(entryMap *EntryMap, out string, tag string, stdout io.Writer, stderr io.Writer) int {
	paths := entryMap.GetPaths()
	lines := make([]string, 0)
	for i := 0; i < len(paths); i++ {
		path := paths[i]
		if values, ok := entryMap.GetValues(path); ok {
			notes := NewNotes(values)

			tags := strings.Split(notes.Get("tags"), ",")
			namespaces := strings.Split(notes.Get("namespace"), ",")

			for j := 0; j < len(namespaces); j++ {
				namespace := namespaces[j]
				if include(tags, tag) {
					secretType := notes.Get("type")
					switch secretType {
					case "opaque":
						createOpaqueSecret(path, namespace, notes, values, &lines, stdout, stderr)
					case "docker":
						createDockerSecret(path, namespace, values, &lines, stdout, stderr)
					case "tls":
						createTlsSecret(path, namespace, values, &lines, stdout, stderr)
					}
				}
			}
		}
	}

	return writeFile(out, &lines, stderr)
}

// check if entry should be included
// apply tag filter
func include(tags []string, tag string) bool {
	if tag == "" {
		return true
	}

	// tag specified -> entry must contain tag
	for i := 0; i < len(tags); i++ {
		if tags[i] == tag {
			return true
		}
	}

	return false
}

// create opaque (regular) secret
func createOpaqueSecret(path string, namespace string, notes *Notes, values Entry, lines *[]string, stdout io.Writer, stderr io.Writer) {
	title, _ := values.GetValue("Title")
	if title == "" {
		fmt.Fprintf(stderr, "missing title for entry '%s'\n", path)
		return
	}

	if len(*lines) > 0 {
		*lines = append(*lines, "")
		*lines = append(*lines, "---")
	}

	*lines = append(*lines, "apiVersion: v1")
	*lines = append(*lines, "kind: Secret")
	*lines = append(*lines, "metadata:")
	*lines = append(*lines, "  name: \""+title+"\"")
	if len(namespace) > 0 {
		*lines = append(*lines, "  namespace: \""+namespace+"\"")
	}
	*lines = append(*lines, "type: Opaque")
	*lines = append(*lines, "data:")

	secretKeys := notes.GetKeys()

	fmt.Fprintf(stdout, "secret opaque name=%s fields=%s\n", title, strings.Join(secretKeys, ","))

	for i := 0; i < len(secretKeys); i++ {
		secretKey := secretKeys[i]
		valuesKey := notes.Get(secretKey)
		value, ok := values.GetValue(valuesKey)
		if ok {
			value64 := base64.StdEncoding.EncodeToString([]byte(value))
			secretKey = strings.TrimPrefix(secretKey, ":")
			*lines = append(*lines, "  "+secretKey+": \""+value64+"\"")
		} else {
			fmt.Fprintf(stderr, "entry '%s' does not contain value '%s'\n", path, valuesKey)
		}
	}
}

// create docker secret
func createDockerSecret(path string, namespace string, values Entry, lines *[]string, stdout io.Writer, stderr io.Writer) {

	title, _ := values.GetValue("Title")
	if title == "" {
		fmt.Fprintf(stderr, "missing title for entry '%s'\n", path)
		return
	}

	username, _ := values.GetValue("UserName")
	if username == "" {
		fmt.Fprintf(stderr, "missing UserName for entry '%s'\n", path)
		return
	}

	password, _ := values.GetValue("Password")
	if password == "" {
		fmt.Fprintf(stderr, "missing Password for entry '%s'\n", path)
		return
	}

	url, _ := values.GetValue("URL")
	if url == "" {
		fmt.Fprintf(stderr, "missing URL for entry '%s'\n", path)
		return
	}

	fmt.Fprintf(stdout, "secret docker name=%s url=%s username=%s\n", title, url, username)

	if len(*lines) > 0 {
		*lines = append(*lines, "")
		*lines = append(*lines, "---")
	}

	email := "mail@example.de"
	auth := username + ":" + password
	auth64 := base64.StdEncoding.EncodeToString([]byte(auth))

	secretJson := "{\"auths\":{\"" + url +
		"\":{\"username\":\"" + username +
		"\",\"password\":\"" + password +
		"\",\"email\":\"" + email +
		"\",\"auth\":\"" + auth64 + "\"}}}"
	secret64 := base64.StdEncoding.EncodeToString([]byte(secretJson))

	*lines = append(*lines, "apiVersion: v1")
	*lines = append(*lines, "kind: Secret")
	*lines = append(*lines, "metadata:")
	*lines = append(*lines, "  name: \""+title+"\"")
	if len(namespace) > 0 {
		*lines = append(*lines, "  namespace: \""+namespace+"\"")
	}

	*lines = append(*lines, "type: kubernetes.io/dockerconfigjson")
	*lines = append(*lines, "data:")
	*lines = append(*lines, "  .dockerconfigjson: "+secret64)
}

// create docker secret
func createTlsSecret(path string, namespace string, values Entry, lines *[]string, stdout io.Writer, stderr io.Writer) {

	title, _ := values.GetValue("Title")
	if title == "" {
		fmt.Fprintf(stderr, "missing title for entry '%s'\n", path)
		return
	}

	username, _ := values.GetValue("UserName")
	if username == "" {
		fmt.Fprintf(stderr, "missing UserName for entry '%s'\n", path)
		return
	}

	password, _ := values.GetValue("Password")
	if password == "" {
		fmt.Fprintf(stderr, "missing Password for entry '%s'\n", path)
		return
	}

	fmt.Fprintf(stdout, "secret tls name=%s crt=%s key=%s\n", title, username[:27], password[:27])

	if len(*lines) > 0 {
		*lines = append(*lines, "")
		*lines = append(*lines, "---")
	}

	crt := base64.StdEncoding.EncodeToString([]byte(username))
	key := base64.StdEncoding.EncodeToString([]byte(password))

	*lines = append(*lines, "apiVersion: v1")
	*lines = append(*lines, "kind: Secret")
	*lines = append(*lines, "metadata:")
	*lines = append(*lines, "  name: \""+title+"\"")
	if len(namespace) > 0 {
		*lines = append(*lines, "  namespace: \""+namespace+"\"")
	}
	*lines = append(*lines, "type: kubernetes.io/tls")
	*lines = append(*lines, "data:")
	*lines = append(*lines, "  tls.crt: \""+crt+"\"")
	*lines = append(*lines, "  tls.key: \""+key+"\"")
}
