package main

import (
	"bufio"
	"fmt"
	"fp/pkg/fp"
	"os"
)

func tokenIter() <-chan fp.Token {
	outCh := make(chan fp.Token)
	fmt.Println("welcome to fp repl")
	go func(outCh chan fp.Token) {
		defer close(outCh)
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Printf(">>>")
		for scanner.Scan() {
			line := scanner.Text()

			for _, tok := range fp.Tokenize(line) {
				outCh <- tok
			}
			fmt.Printf(">>>")
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}(outCh)
	return outCh
}

func main() {
	r := fp.NewBasicRuntime().
		WithArithmeticExtension("print", func(nums ...fp.Object) fp.Object {
			for _, num := range nums {
				fmt.Printf("%v ", num)
			}
			fmt.Println()
			return len(nums)
		}).
		WithArithmeticExtension("div", func(nums ...fp.Object) fp.Object {
			if len(nums) != 2 {
				panic("runtime error")
			}
			return nums[0].(int) / nums[1].(int)
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

	for expr := range fp.ParseAllREPL(tokenIter()) {
		r.Step(expr)
	}
	return
}
