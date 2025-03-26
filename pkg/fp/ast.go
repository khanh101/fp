package fp

// Expr : union of string, LambdaExpr
type Expr interface{}
type LambdaExpr struct {
	Name string
	Args []Expr
}
