package main

import (
	"context"
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
	r := fp.NewStdRuntime()
	writeln("welcome to fp repl! type function or module name for help")
	var funcNameList []string
	for k := range r.Stack[0] {
		funcNameList = append(funcNameList, string(k))
	}
	sort.Strings(funcNameList)
	for _, name := range funcNameList {
		o, err := r.Step(context.Background(), fp.NameExpr(name))
		if err != nil {
			panic(err)
		}
		writeln(">>>%s", name)
		writeln("%v", o)
	}
}
