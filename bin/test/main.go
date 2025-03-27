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

	r := fp.NewBasicRuntime().
		WithDebug(false).
		WithArithmeticExtension("div", func(nums ...fp.Object) fp.Object {
			if len(nums) != 2 {
				panic("runtime error")
			}
			return nums[0].(int) / nums[1].(int)
		}).
		WithArithmeticExtension("print", func(nums ...fp.Object) fp.Object {
			for _, num := range nums {
				fmt.Printf("%v ", num)
			}
			fmt.Println()
			return len(nums)
		}).
		WithArithmeticExtension("input", func(nums ...fp.Object) fp.Object {
			var v int
			_, err := fmt.Scanf("%d", &v)
			if err != nil {
				panic(err)
			}
			return v
		}).
		WithArithmeticExtension("make_list", func(nums ...fp.Object) fp.Object {
			var v []fp.Object
			for _, num := range nums {
				v = append(v, num)
			}
			return v
		}).
		WithArithmeticExtension("append_list", func(nums ...fp.Object) fp.Object {
			return append(nums[0].([]fp.Object), nums[1:]...)
		})
	for _, expr := range exprList {
		r.Step(expr)
	}
}
