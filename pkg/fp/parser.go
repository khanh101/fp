package fp

import (
	"regexp"
	"strings"
)

type Token = string

func pop(tokenList []Token) ([]Token, Token) {
	return tokenList[1:], tokenList[0]
}

func peak(tokenList []Token) Token {
	return tokenList[0]
}

func Tokenize(str string) []Token {
	// remove comment
	parts := strings.Split(str, "\n")
	newParts := []string{}
	for _, part := range parts {
		newParts = append(newParts, strings.Split(part, "//")[0])
	}

	str = strings.Join(newParts, "\n")

	str = strings.ReplaceAll(str, "(", " ( ")
	str = strings.ReplaceAll(str, ")", " ) ")

	// splitBySpaceExceptQuotes : chatgpt prompt: golang split a string by empty space, except those enclosed with double quotes
	splitBySpaceExceptQuotes := func(s string) []string {
		// Define a regular expression to match either a quoted string or non-quoted words
		re := regexp.MustCompile(`(?:[^\s"]+|"[^"]*")`)

		// Find all substrings that match the regular expression
		return re.FindAllString(s, -1)
	}
	fields := splitBySpaceExceptQuotes(str)
	return fields
}

func ParseAll(tokenList []Token) ([]Expr, []Token) {
	var expr Expr
	var exprList []Expr
	for {
		if len(tokenList) == 0 || peak(tokenList) == ")" {
			break
		}
		expr, tokenList = parse(tokenList)
		exprList = append(exprList, expr)
	}
	return exprList, tokenList
}

func parse(tokenList []Token) (Expr, []Token) {
	if len(tokenList) == 0 {
		return nil, nil
	}
	tokenList, head := pop(tokenList) // pop ( or [ or name
	switch head {
	case "(":
		tokenList, funcName := pop(tokenList)
		exprList, tokenList := ParseAll(tokenList)
		tokenList, tail := pop(tokenList) // pop )
		if tail != ")" {
			panic("parse error")
		}
		return LambdaExpr{
			Name: Name(funcName),
			Args: exprList,
		}, tokenList
	default:
		if head[0] == '"' && head[len(head)-1] == '"' {
			return String(head[1 : len(head)-1]), tokenList
		} else {
			return Name(head), tokenList
		}
	}
}
