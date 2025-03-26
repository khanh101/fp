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

	b, tokenList := fp.ParseMany(tokenList)

	r := fp.NewRuntime()
	for _, block := range b {
		r.Eval(block)
	}
}
