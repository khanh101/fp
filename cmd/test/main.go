package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
)

func parseInput(input string) ([]string, error) {
	isSplitter := func(r rune) bool {
		return strings.ContainsRune("()[]{},", r)
	}

	appendToken := func(tokens []string, token string) []string {
		if token != "" {
			tokens = append(tokens, token)
		}
		return tokens
	}

	var tokens []string
	var current strings.Builder
	var quoteBuffer strings.Builder
	inQuotes := false
	escapeNext := false

	for _, r := range input {
		switch {
		case escapeNext:
			if inQuotes {
				quoteBuffer.WriteRune(r)
			} else {
				current.WriteRune(r)
			}
			escapeNext = false

		case r == '\\':
			escapeNext = true
			if inQuotes {
				quoteBuffer.WriteRune(r)
			}

		case r == '"':
			if inQuotes {
				quoteBuffer.WriteRune(r)
				inQuotes = false

				// Flush buffer before quoted string
				if current.Len() > 0 {
					tokens = appendToken(tokens, current.String())
					current.Reset()
				}

				// Decode quoted string
				quoted := quoteBuffer.String()
				var unquoted string
				err := json.Unmarshal([]byte(quoted), &unquoted)
				if err != nil {
					unquoted = quoted[1 : len(quoted)-1] // fallback
				}
				tokens = append(tokens, `"`+unquoted+`"`)
				quoteBuffer.Reset()

			} else {
				inQuotes = true
				if current.Len() > 0 {
					tokens = appendToken(tokens, current.String())
					current.Reset()
				}
				quoteBuffer.WriteRune(r)
			}

		case inQuotes:
			quoteBuffer.WriteRune(r)

		case unicode.IsSpace(r):
			if current.Len() > 0 {
				tokens = appendToken(tokens, current.String())
				current.Reset()
			}

		case isSplitter(r):
			// Flush any word before splitter
			if current.Len() > 0 {
				tokens = appendToken(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, string(r))

		default:
			current.WriteRune(r)
		}
	}

	if inQuotes {
		return nil, fmt.Errorf("unclosed quote in input at end")
	}
	if current.Len() > 0 {
		tokens = appendToken(tokens, current.String())
	}

	return tokens, nil
}

func main() {
	input := `hello("this is a test")12world"foo bar""\"quoted\"word"`
	tokens, err := parseInput(input)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for i, t := range tokens {
		fmt.Printf("Token %d: %v\n", i, t)
	}
}
