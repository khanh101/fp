package fp

// Expr : union of Name, LambdaExpr
type Expr interface{}

type Name string
type LambdaExpr struct {
	Name Name
	Args []Expr
}
