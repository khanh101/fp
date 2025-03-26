package fp

// Expr : union of string, LambdaExpr
type Expr interface{}

type Name string
type LambdaExpr struct {
	Name Name
	Args []Expr
}
