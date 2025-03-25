package main

import (
	"fp/pkg/fp"
	"os"
)

func main() {

	buffer, err := os.ReadFile("fib.lisp")
	if err != nil {
		panic(err)
	}
	str := string(buffer)
	tokenList := fp.ParseFromString(str)

	b, tokenList := fp.ParseMany(tokenList)

	r := fp.Runtime{
		FuncMap:     make(map[string]fp.Func),
		VarMapStack: []map[string]int{make(map[string]int)},
	}
	for _, block := range b {
		r.Eval(block)
	}
}
