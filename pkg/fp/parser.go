package fp

import "strings"

type Token = string

func pop(tokenList []Token) ([]Token, Token) {
	return tokenList[1:], tokenList[0]
}

func peak(tokenList []Token) Token {
	return tokenList[0]
}

func ParseFromString(str string) []Token {
	str = strings.ReplaceAll(str, "(", " ( ")
	str = strings.ReplaceAll(str, ")", " ) ")
	fields := strings.Fields(str)
	return fields
}

func ParseMany(tokenList []Token) ([]*Block, []Token) {
	var block *Block
	var blockList []*Block
	for {
		if len(tokenList) == 0 || peak(tokenList) == ")" {
			break
		}
		block, tokenList = parse(tokenList)
		blockList = append(blockList, block)
	}
	return blockList, tokenList
}
func parse(tokenList []Token) (*Block, []Token) {
	if len(tokenList) == 0 {
		return nil, nil
	}
	tokenList, head := pop(tokenList) // pop ( or name
	if head != "(" {
		return &Block{
			Type: BLOCKTYPE_LITERAL,
			Name: head,
		}, tokenList
	}
	tokenList, funcName := pop(tokenList)
	blockList, tokenList := ParseMany(tokenList)
	tokenList, tail := pop(tokenList) // pop )
	if tail != ")" {
		panic("parse error")
	}
	return &Block{
		Type: BLOCKTYPE_FUNCTION,
		Name: funcName,
		Args: blockList,
	}, tokenList
}
