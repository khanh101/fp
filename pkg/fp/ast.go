package fp

const (
	BLOCKTYPE_LITERAL  = "literal"  // name
	BLOCKTYPE_FUNCTION = "function" // function_name + list of blocks
	BLOCKTYPE_LIST     = "list"     // list of names
)

type Block struct {
	Type string
	Name string
	Args []*Block
}
