package main

import (
	"fmt"
	"unicode"
)

type Token = string

func main() {
	input := `hello this is an "example" of "quoted\n string \"field\""`
	fmt.Println(input)
	tokens := tokenize(input)
	for i, t := range tokens {
		fmt.Printf("token[%d]:", i)
		fmt.Println(t)
	}
}
