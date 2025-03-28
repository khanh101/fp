package fp

// Expr : union of Name, LambdaExpr
type Expr interface {
	String() string
	MustTypeExpr() struct{} // for type-safety every Expr must implement this
}

type Name string

func (e Name) String() string {
	return string(e)
}

func (e Name) MustTypeExpr() struct{} {
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

func (e LambdaExpr) MustTypeExpr() struct{} {
	return struct{}{}
}
