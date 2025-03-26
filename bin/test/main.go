package main

import (
	"fp/pkg/fp"
	"os"
)

func main() {

	buffer, err := os.ReadFile("example.lisp")
	if err != nil {
		panic(err)
	}
	str := string(buffer)
	tokenList := fp.Tokenize(str)

	blockList, tokenList := fp.ParseMany(tokenList)
	if len(tokenList) > 0 {
		panic("parse error")
	}

	r := fp.NewRuntime()
	for _, block := range blockList {
		r.Eval(block)
	}
}
