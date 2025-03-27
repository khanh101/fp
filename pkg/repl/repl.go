package repl

import (
	"fmt"
	"fp/pkg/fp"
	"sort"
)

type REPL interface {
	ReplyInput(input string) (output string)
	ClearBuffer() (output string)
}

type fpRepl struct {
	runtime *fp.Runtime
	parser  *fp.Parser
	buffer  string
}

func (r *fpRepl) ReplyInput(input string) (output string) {
	tokenList := fp.Tokenize(input)
	executed := false
	if len(tokenList) == 0 {
		executed = true
	} else {
		for _, token := range tokenList {
			expr := r.parser.Input(token)
			if expr != nil {
				executed = true
				output, err := r.runtime.Step(expr)
				if err != nil {
					r.writeln(err.Error())
					continue
				}
				r.write("%v\n", output)
			}
		}
	}
	if executed {
		r.write(">>>")
	}
	return r.flush()
}

func (r *fpRepl) ClearBuffer() (output string) {
	r.parser.Clear()
	r.writeln(">>> (Control + C) to clear buffer, (Control + D) to exit")
	r.writeln(">>>")
	return r.flush()
}

func (r *fpRepl) flush() (output string) {
	output, r.buffer = r.buffer, ""
	return output
}

func (r *fpRepl) write(format string, a ...interface{}) {
	r.buffer += fmt.Sprintf(format, a...)
}
func (r *fpRepl) writeln(format string, a ...interface{}) {
	r.write(format+"\n", a...)
}

func NewFP() (repl REPL, welcome string) {
	r := &fpRepl{
		runtime: fp.NewStdRuntime(),
		parser:  &fp.Parser{},
		buffer:  "",
	}
	r.writeln("welcome to fp repl! type <function or module name> for help")
	r.write("loaded modules: ")
	var funcNameList []string
	for k := range r.runtime.Stack[0] {
		funcNameList = append(funcNameList, string(k))
	}
	sort.Strings(funcNameList)
	for _, name := range funcNameList {
		r.write("%s ", name)
	}
	r.writeln("")
	r.write(">>>")
	return r, r.flush()
}
