package main

import (
	"fmt"
	"fp/pkg/fp"
	"sort"
	"strings"
	"syscall/js"
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

var reply func(input string) (output string)
var clear func() (output string)

func evaluate(this js.Value, p []js.Value) interface{} {
	if len(p) == 0 {
		return js.ValueOf("no input")
	}
	input := p[0].String()

	// repl here
	// end repl here

	// Simple echo for now; replace with real evaluation logic
	output := fmt.Sprintf("You entered: %s", input)
	js.Global().Call("updateOutput", output)
	return nil
}

func write(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	strings.ReplaceAll(s, "\n", "<br>")

}

func main() {
	// Expose Go functions to the global JavaScript context
	js.Global().Set("evaluate", js.FuncOf(evaluate))
	// js.Global().Set("sendOutputToWeb", js.FuncOf(sendOutputToWeb))

	//js.Global().Call("sendOutputToWeb")

	js.Global().Call("updateOutput", "this is weird")

	// Keep WebAssembly running
	select {}
}
