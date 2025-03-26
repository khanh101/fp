package fp

const (
	BLOCKTYPE_NAME = "name" // name
	BLOCKTYPE_EXPR = "expr" // name + list of blocks
)

type Block struct {
	Type string
	Name string
	Args []*Block
}

// Expr : union of string, LambdaExpr
type Expr interface{}
type LambdaExpr struct {
	Name string
	Args []Expr
}
