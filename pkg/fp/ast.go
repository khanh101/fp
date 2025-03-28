package fp

// Expr : union of Name, LambdaExpr
type Expr interface {
	String() string
	AssertExpr() struct{} // for type-safety
}

type Name string

func (e Name) String() string {
	return string(e)
}

func (e Name) AssertExpr() struct{} {
	return struct{}{}
}

func (e Name) AssertObject() struct{} {
	return struct{}{}
}

// LambdaExpr : S-expression
type LambdaExpr struct {
	Name Name
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

func (e LambdaExpr) AssertExpr() struct{} {
	return struct{}{}
}
