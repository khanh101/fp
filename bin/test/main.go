package main

import (
	"fmt"
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

	r := fp.NewRuntime().WithExtension("div", func(nums ...fp.Value) fp.Value {
		if len(nums) != 2 {
			panic("runtime error")
		}
		return nums[0].(int) / nums[1].(int)
	}).WithExtension("output", func(nums ...fp.Value) fp.Value {
		for _, num := range nums {
			fmt.Printf("%v ", num)
		}
		fmt.Println()
		return len(nums)
	}).WithExtension("input", func(nums ...fp.Value) fp.Value {
		var v int
		_, err := fmt.Scanf("%d", &v)
		if err != nil {
			panic(err)
		}
		return v
	})
	for _, expr := range exprList {
		r.Step(expr)
	}
}
