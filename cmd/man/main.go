package main

import (
	"fmt"
	"fp/pkg/fp"
	"os"
	"sort"
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
		LoadExtension("div", fp.Extension{
			Exec: func(value ...fp.Object) (fp.Object, error) {
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
			},
			Man: "extension: division",
		})
	writeln("welcome to fp repl! type <function or module name> for help")
	var funcNameList []string
	for k := range r.Stack[0] {
		funcNameList = append(funcNameList, string(k))
	}
	sort.Strings(funcNameList)
	for _, name := range funcNameList {
		o, err := r.Step(fp.Name(name))
		if err != nil {
			panic(err)
		}
		writeln(">>>%s", name)
		writeln("%v", o)
	}
}
