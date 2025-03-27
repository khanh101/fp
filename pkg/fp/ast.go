package fp

// Expr : union of Name, LambdaExpr
type Expr interface {
	String() string
}

type Name string

func (e Name) String() string {
	return string(e)
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
