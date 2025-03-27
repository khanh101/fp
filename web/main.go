package main

import (
	"fmt"
	"fp/pkg/fp"
	"sort"
	"strings"
)

func repl() func(input string) (output string) {
	r := fp.NewStdRuntime()
	buffer := ""
	write := func(format string, a ...interface{}) {
		s := fmt.Sprintf(format, a...)
		strings.ReplaceAll(s, "\n", "<br>")
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
	return func(input string) (output string) {
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
		output = buffer
		buffer = ""
		return output
	}
}
