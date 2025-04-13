package fp

import (
	"errors"
	"slices"
)

// Expr : union of NameExpr, LambdaExpr
type Expr interface {
	String() string
	MustTypeExpr() // for type-safety every Expr must implement this
}

type NameExpr string

func (e NameExpr) String() string {
	return string(e)
}

func (e NameExpr) MustTypeExpr() {
}

// LambdaExpr : S-expression
type LambdaExpr struct {
	Name NameExpr
	Args []Expr
}

func (e LambdaExpr) String() string {
	s := ""
	s += "("
	s += e.Name.String()
	for _, arg := range e.Args {
		s += " " + arg.String()
	}
	s += ")"
	return s
}

func (e LambdaExpr) MustTypeExpr() {

}

func pop(tokenList []Token) ([]Token, Token, error) {
	if len(tokenList) == 0 {
		return nil, "", errors.New("empty token list")
	}
	return tokenList[1:], tokenList[0], nil
}

// ParseAll : parse a token list
func ParseAll(tokenList []Token) ([]Expr, []Token) {
	var expr Expr
	var exprList []Expr
	var err error
	for len(tokenList) > 0 {
		expr, tokenList, err = parseSingle(tokenList)
		if err != nil {
			panic(err)
		}
		exprList = append(exprList, expr)
	}
	return exprList, tokenList
}

type Parser struct {
	Buffer []Token
}

func (p *Parser) Clear() {
	p.Buffer = []Token{}
}

func (p *Parser) Input(tok Token) Expr {
	p.Buffer = append(p.Buffer, tok)
	// try parse single // TODO : do this for simplicity
	buffer := slices.Clone(p.Buffer)
	expr, buffer, err := parseSingle(buffer)
	if err != nil {
		// parse fail - don't do anything
		return nil
	} else {
		// parse ok - update buffer
		p.Buffer = buffer
		return expr
	}
}

func parseSingle(tokenList []Token) (Expr, []Token, error) {
	var parse func(tokenList []Token) (Expr, []Token, bool, error)
	parse = func(tokenList []Token) (Expr, []Token, bool, error) {
		if len(tokenList) == 0 {
			return nil, nil, false, errors.New("empty token list")
		}
		tokenList, head, err := pop(tokenList) // pop ( or [ or name
		if err != nil {
			return nil, nil, false, err
		}
		switch head {
		case "(": // start with Open
			tokenList, funcName, err := pop(tokenList)
			if err != nil {
				return nil, nil, false, err
			}
			if funcName == ")" { // empty
				return parse(tokenList)
			}
			var expr Expr
			var exprList []Expr
			var endWithClose bool
			for {
				expr, tokenList, endWithClose, err = parse(tokenList)
				if err != nil {
					return nil, nil, false, err
				}
				if endWithClose {
					// end with Close
					break
				}
				exprList = append(exprList, expr)
			}
			return LambdaExpr{
				Name: NameExpr(funcName),
				Args: exprList,
			}, tokenList, false, nil
		default:
			return NameExpr(head), tokenList, head == ")", nil
		}
	}

	expr, tokenList, endWithClose, err := parse(tokenList)
	if err != nil {
		return nil, nil, err
	}
	if endWithClose {
		return nil, nil, errors.New("parse error")
	}
	return expr, tokenList, nil
}
