package main

import (
	"bufio"
	"fmt"
	"fp/pkg/fp"
	"os"
	"os/signal"
	"sort"
	"syscall"
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
			return a / b, nil
		})
	writeln("welcome to fp repl! ")
	write("loaded modules: ")
	var moduleNameList []string
	for k := range r.Module {
		moduleNameList = append(moduleNameList, string(k))
	}
	sort.Strings(moduleNameList)
	for _, name := range moduleNameList {
		write("%s ", name)
	}
	writeln("")

	signCh := make(chan os.Signal, 1)
	signal.Notify(signCh, syscall.SIGINT, syscall.SIGTERM)

	parser := &fp.Parser{}

	scanner := bufio.NewScanner(os.Stdin)
	write(">>>")
	for scanner.Scan() {
		select {
		case <-signCh:
			parser.Clear()
			writeln(">>> (Control + C) to clear buffer, (Control + D) to exit")
			writeln(">>>")
		default:
			line := scanner.Text()
			tokenList := fp.Tokenize(line)
			executed := false
			if len(tokenList) == 0 {
				executed = true
			} else {
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
			}
			if executed {
				write(">>>")
			}
		}

	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
