package main

import (
	"bufio"
	"fmt"
	"fp/pkg/fp"
	"os"
)

func write(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
	_ = os.Stderr.Sync() // flush
}

func writeln(format string, args ...interface{}) {
	write(format+"\n", args...)
}

func main() {
	r := fp.NewDebugRuntime().
		WithArithmeticExtension("print", func(nums ...fp.Object) (fp.Object, error) {
			for _, num := range nums {
				fmt.Printf("%v ", num)
			}
			fmt.Println()
			return len(nums), nil
		}).
		WithArithmeticExtension("div", func(value ...fp.Object) (fp.Object, error) {
			if len(value) != 2 {
				return nil, fmt.Errorf("subtract requires 2 arguments")
			}
			a, ok := value[0].(int)
			if !ok {
				return nil, fmt.Errorf("subtract non-integer value")
			}
			b, ok := value[0].(int)
			if !ok {
				return nil, fmt.Errorf("subtract non-integer value")
			}
			if b == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return a - b, nil
		})
	writeln("welcome to fp repl! ")
	write("loaded modules: ")
	for k, _ := range r.Module {
		write("%v ", k)
	}
	writeln("")

	parser := &fp.Parser{}

	scanner := bufio.NewScanner(os.Stdin)
	write(">>>")
	for scanner.Scan() {
		line := scanner.Text()
		tokenList := fp.Tokenize(line)
		executed := false
		for _, token := range tokenList {
			expr := parser.Input(token)
			if expr != nil {
				executed = true
				output, err := r.Step(expr)
				if err != nil {
					writeln(err.Error())
					continue
				}
				write("%v\n", output)
			}
		}
		if executed {
			write(">>>")
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
