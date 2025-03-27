package fp

import (
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

	str = strings.ReplaceAll(str, "\n", " ")
	str = strings.ReplaceAll(str, "(", " ( ")
	str = strings.ReplaceAll(str, ")", " ) ")

	return strings.Fields(str)
}

func ParseAll(tokenList []Token) ([]Expr, []Token) {
	var expr Expr
	var exprList []Expr
	var endWithClose bool
	for len(tokenList) > 0 {
		expr, tokenList, endWithClose = parse(tokenList)
		if endWithClose {
			panicError("parse error")
		}
		exprList = append(exprList, expr)
	}
	return exprList, tokenList
}

func ParseAllREPL(tokenCh <-chan Token) <-chan Expr {
	outCh := make(chan Expr)
	go func(outCh chan Expr) {
		defer close(outCh)
		for {
			expr, endWithClose, eof := parseREPL(tokenCh)
			if eof {
				break
			}
			if endWithClose {
				panicError("parse error")
			}
			outCh <- expr
		}
	}(outCh)
	return outCh
}

func parseREPL(tokenCh <-chan Token) (expr Expr, endWithClose bool, eof bool) {
	for head := range tokenCh {
		switch head {
		case "(":
			funcName := <-tokenCh
			var exprList []Expr
			for {
				expr, endWithClose, eof = parseREPL(tokenCh)
				if eof {
					panicError("parse error")
				}
				if endWithClose {
					break
				}
				exprList = append(exprList, expr)
			}
			return LambdaExpr{
				Name: Name(funcName),
				Args: exprList,
			}, false, false
		default:
			return Name(head), head == ")", false
		}
	}
	return nil, false, true
}

func parse(tokenList []Token) (Expr, []Token, bool) {
	if len(tokenList) == 0 {
		return nil, nil, false
	}
	tokenList, head := pop(tokenList) // pop ( or [ or name
	switch head {
	case "(":
		tokenList, funcName := pop(tokenList)
		var expr Expr
		var exprList []Expr
		var endWithClose bool
		for {
			expr, tokenList, endWithClose = parse(tokenList)
			if endWithClose {
				break
			}
			exprList = append(exprList, expr)
		}
		return LambdaExpr{
			Name: Name(funcName),
			Args: exprList,
		}, tokenList, false
	default:
		return Name(head), tokenList, head == ")"
	}
}
