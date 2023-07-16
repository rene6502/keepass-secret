package cmd

import (
	"fmt"
	"io"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var charsHexLower = []rune("0123456789abcdef")
var charsHexUpper = []rune("0123456789ABCDEF")
var charsLetter = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
var charsAlphanumeric = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
var charsPrintable = []rune("!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~")
var regexPattern = regexp.MustCompile("^[h|H|A|L|S][0-9]{1,3}$")

func createPassword(chars []rune, length int) string {

	source := rand.NewSource(time.Now().UnixNano())
	randSource := rand.New(source)

	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[randSource.Intn(len(chars))])
	}

	return b.String()
}

// creates password from specified pattern
// e.g. pattern=A32 creates a password with 32 random alphanumeric characters
func createPasswordFromPattern(pattern string, stdout io.Writer, stderr io.Writer) string {
	if !regexPattern.MatchString(pattern) {
		fmt.Fprintf(stderr, "unknown password pattern %s, fallback to A32\n", pattern)
		pattern = "A32"
	}

	length, _ := strconv.Atoi(pattern[1:])
	c := pattern[0:1]
	switch c {
	case "h": // lower-case hex characters
		return createPassword(charsHexLower, length)

	case "H": // upper-case hex characters
		return createPassword(charsHexUpper, length)

	case "L": // letter
		return createPassword(charsLetter, length)

	case "S": // printable
		return createPassword(charsPrintable, length)

	default: // alphanumeric
		return createPassword(charsAlphanumeric, length)
	}
}
