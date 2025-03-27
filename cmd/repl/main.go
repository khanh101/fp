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

func repl() (output string, reply func(input string) (output string), clear func() (output string)) {
	r := fp.NewStdRuntime()
	buffer := ""
	write := func(format string, a ...interface{}) {
		s := fmt.Sprintf(format, a...)
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
	output, reply, clearBuffer := repl()
	_, _ = fmt.Fprintf(os.Stderr, output)

	signCh := make(chan os.Signal, 1)
	signal.Notify(signCh, syscall.SIGINT, syscall.SIGTERM)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		select {
		case s := <-signCh:
			switch s {
			case syscall.SIGINT:
				output := clearBuffer()
				_, _ = fmt.Fprintf(os.Stderr, output)
			case syscall.SIGTERM:
				os.Exit(0)
			default:
				os.Exit(1)
			}
		default:
			input := scanner.Text()
			output := reply(input)
			_, _ = fmt.Fprintf(os.Stderr, output)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
