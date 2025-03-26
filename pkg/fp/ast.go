package fp

const (
	BLOCKTYPE_LITERAL  = "literal"  // name
	BLOCKTYPE_FUNCTION = "function" // name + list of blocks
)

type Block struct {
	Type string
	Name string
	Args []*Block
}
