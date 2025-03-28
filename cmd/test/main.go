package main

import (
	"fmt"
	"strings"
	"unicode"
)

func splitWords(input string) []string {
	var result []string
	var word strings.Builder
	inQuotes := false
	i := 0
	for i < len(input) {
		char := input[i]

		if char == '\\' && i+1 < len(input) { // Handle escaped characters
			nextChar := input[i+1]
			switch nextChar {
			case 'n':
				word.WriteByte('\n')
			case 't':
				word.WriteByte('\t')
			case 'r':
				word.WriteByte('\r')
			case '"':
				word.WriteByte('"')
			case '\\':
				word.WriteByte('\\')
			default:
				word.WriteByte(nextChar)
			}
			i += 2
			continue
		}

		if char == '"' { // Handle quoted strings
			inQuotes = !inQuotes
			i++
			continue
		}

		if unicode.IsSpace(rune(char)) && !inQuotes { // Split on whitespace outside quotes
			if word.Len() > 0 {
				result = append(result, word.String())
				word.Reset()
			}
			i++
			continue
		}

		word.WriteByte(char)
		i++
	}

	if word.Len() > 0 {
		result = append(result, word.String())
	}

	return result
}

func main() {
	input := "hello world \"this is a \\\"quoted\\\" string\" test\nnew line\ttab"
	words := splitWords(input)
	for _, w := range words {
		fmt.Println(w)
	}

	// Output: ["hello" "world" "this is a \"quoted\" string" "test"]
}
