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

	exprList, tokenList := fp.ParseAll(tokenList)
	if len(tokenList) > 0 {
		panic("parse error")
	}

	r := fp.NewRuntime().WithExtension("div", func(nums ...int) int {
		if len(nums) != 2 {
			panic("runtime error")
		}
		return nums[0] / nums[1]
	})
	for _, expr := range exprList {
		r.Step(expr)
	}
}
