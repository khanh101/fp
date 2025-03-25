package fp

const (
	BLOCKTYPE_LITERAL  = "literal"
	BLOCKTYPE_FUNCTION = "function"
)

type Block struct {
	Type string
	Name string
	Args []*Block
}
