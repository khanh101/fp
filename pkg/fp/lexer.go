package fp

import (
	"fmt"
	"strings"
	"unicode"
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

func Tokenize(str string) []Token {
	str = removeComments(str)

	const (
		STATE_OUTSTRING = iota
		STATE_INSTRING
		STATE_INSTRING_ESCAPE
	)

	var tokens []Token
	state := STATE_OUTSTRING
	buffer := ""
	flushBuffer := func() {
		if len(buffer) > 0 {
			tokens = append(tokens, buffer)
		}
		buffer = ""
	}
	for _, ch := range str {
		switch state {
		case STATE_OUTSTRING:
			if unicode.IsSpace(ch) {
				flushBuffer()
			} else if ch == '(' || ch == ')' || ch == '*' {
				flushBuffer()
				buffer += string(ch)
				flushBuffer()
			} else if ch == '"' {
				flushBuffer()
				buffer += string(ch)
				state = STATE_INSTRING
			} else {
				buffer += string(ch)
			}
		case STATE_INSTRING:
			if ch == '\\' {
				buffer += string(ch)
				state = STATE_INSTRING_ESCAPE
			} else if ch == '"' {
				buffer += string(ch)
				flushBuffer()
				state = STATE_OUTSTRING
			} else {
				buffer += string(ch)
			}
		case STATE_INSTRING_ESCAPE:
			buffer += string(ch)
			state = STATE_INSTRING
		default:
			panic(fmt.Sprintf("invalid state: %d", state))
		}
	}
	flushBuffer()
	return tokens
}
