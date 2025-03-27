package main

import (
	"fmt"
	"fp/pkg/fp"
	"sort"
	"strings"
	"syscall/js"
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

var replyFunc func(input string) (output string)
var clearFunc func() (output string)

func write(format string, a ...interface{}) {
	output := fmt.Sprintf(format, a...)
	output = strings.ReplaceAll(output, "\n", "<br>")
	js.Global().Call("updateOutput", output)
}

func evaluate(this js.Value, p []js.Value) interface{} {
	if len(p) == 0 {
		return js.ValueOf("no input")
	}
	input := p[0].String()

	// repl here
	output := replyFunc(input)
	// end repl here

	output = strings.ReplaceAll(output, "\n", "<br>")
	return output
}

func main() {
	// initialize
	output, reply, clearBuffer := repl()
	replyFunc = reply
	clearFunc = clearBuffer
	write(output)

	js.Global().Set("evaluate", js.FuncOf(evaluate))
	// Keep WebAssembly running
	select {}
}
