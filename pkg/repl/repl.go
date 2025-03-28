package repl

import (
	"errors"
	"fmt"
	"fp/pkg/fp"
	"sort"
)

type REPL interface {
	ReplyInput(input string, interruptCh <-chan struct{}) (output string, executed bool)
	ClearBuffer() (output string)
}

type fpRepl struct {
	runtime *fp.Runtime
	parser  *fp.Parser
	buffer  string
}

func (r *fpRepl) ReplyInput(input string, interruptCh <-chan struct{}) (output string, executed bool) {
	tokenList := fp.Tokenize(input)
	executed = false
	if len(tokenList) == 0 {
		executed = true
	} else {
		for _, token := range tokenList {
			expr := r.parser.Input(token)
			if expr != nil {
				executed = true

				var copiedStack []fp.Frame
				for _, frame := range r.runtime.Stack {
					copiedStack = append(copiedStack, make(fp.Frame).Update(frame))
				}
				output, err := r.runtime.Step(expr, interruptCh)
				if err != nil {
					if errors.Is(err, fp.InterruptError) {
						// reset stack size
						r.runtime.Stack = copiedStack
						r.write("interrupted - stack was recovered")
					}
					r.writeln(err.Error())
					continue
				}
				r.write("%v\n", output)
			}
		}
	}
	return r.flush(), executed
}

func (r *fpRepl) ClearBuffer() (output string) {
	r.parser.Clear()
	r.writeln("(Control + C) to clear parser buffer, (Control + D) to exit")
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

func NewFP(runtime *fp.Runtime) (repl REPL, welcome string) {
	r := &fpRepl{
		runtime: runtime,
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
	return r, r.flush()
}
