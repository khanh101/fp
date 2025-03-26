package fp

// Expr : union of Name, LambdaExpr, String
type Expr interface{}

type Name string
type LambdaExpr struct {
	Name Name
	Args []Expr
}

type String string
