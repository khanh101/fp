package main

import (
	"fmt"
	"fp/pkg/fp"
	"os"
)

func main() {
	fmt.Println("hello")

	buffer, err := os.ReadFile("fib.txt")
	if err != nil {
		panic(err)
	}
	str := string(buffer)
	tokenList := fp.ParseFromString(str)
	fmt.Println(tokenList)

	b, tokenList := fp.ParseMany(tokenList)
	fmt.Println(b)

	r := fp.Runtime{
		FuncMap:     make(map[string]fp.Func),
		VarMapStack: []map[string]int{make(map[string]int)},
	}
	for _, block := range b {
		r.Eval(block)
	}
}
