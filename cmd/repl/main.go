package main

import (
	"bufio"
	"fmt"
	"fp/pkg/fp"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
)

func repl(web bool) (output string, repl func(input string) (output string), clear func() (output string)) {
	r := fp.NewStdRuntime()
	buffer := ""
	write := func(format string, a ...interface{}) {
		s := fmt.Sprintf(format, a...)
		if web {
			strings.ReplaceAll(s, "\n", "<br>")
		}
		buffer += s
	}
	writeln := func(s string) {
		write(s + "\n")
	}

	writeln("welcome to fp repl! type <function or module name> for help")
	write("loaded modules: ")
	var funcNameList []string
	for k := range r.Stack[0] {
		funcNameList = append(funcNameList, string(k))
	}
	sort.Strings(funcNameList)
	for _, name := range funcNameList {
		write("%s ", name)
	}
	writeln("")
	write(">>>")
	parser := &fp.Parser{}

	output, buffer = buffer, ""
	return output, func(input string) (output string) {
			tokenList := fp.Tokenize(input)
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
			output, buffer = buffer, ""
			return output
		}, func() (output string) {
			parser.Clear()
			writeln(">>> (Control + C) to clear buffer, (Control + D) to exit")
			writeln(">>>")
			output, buffer = buffer, ""
			return output
		}
}

func main() {
	output, repl, clearBuffer := repl(false)
	fmt.Printf(output)

	signCh := make(chan os.Signal, 1)
	signal.Notify(signCh, syscall.SIGINT, syscall.SIGTERM)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		select {
		case s := <-signCh:
			switch s {
			case syscall.SIGINT:
				output := clearBuffer()
				fmt.Printf(output)
			case syscall.SIGTERM:
				os.Exit(0)
			}
		default:
			input := scanner.Text()
			output := repl(input)
			fmt.Printf(output)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
