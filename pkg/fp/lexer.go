package fp

import (
	"fmt"
	"strings"
)

type Token = string

func removeComments(str string) string {
	lines := strings.Split(str, "\n")
	var newLines []string
	for _, line := range lines {
		newLines = append(newLines, strings.Split(line, "//")[0])
	}
	return strings.Join(newLines, "\n")
}

func processSpecialChar(str string) string {
	specialChars := map[rune]struct{}{
		'(': {},
		')': {},
		'*': {}, // unwrap symbol
	}
	newStr := ""
	for _, ch := range str {
		if _, ok := specialChars[ch]; ok {
			newStr += fmt.Sprintf(" %c ", ch)
		} else {
			newStr += string(ch)
		}
	}
	return newStr
}

// Tokenize : TODO - process raw string with double quote using json
func Tokenize(str string) []Token {
	str = removeComments(str)
	str = processSpecialChar(str)
	// tokenize
	return strings.Fields(str)
}
