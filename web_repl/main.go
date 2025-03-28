package main

import (
	"fmt"
	"fp/pkg/repl"
	"strings"
	"syscall/js"
)

var r repl.REPL

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
	output := r.ReplyInput(input)
	// end repl here

	output = strings.ReplaceAll(output, "\n", "<br>")
	return output
}

// Go function to handle buffer clearing
func clearBuffer(this js.Value, p []js.Value) interface{} {
	write(r.ClearBuffer())
	return nil
}

func main() {
	// initialize
	var welcome string
	r, welcome = repl.NewFP()
	write(welcome)

	js.Global().Set("evaluate", js.FuncOf(evaluate))
	js.Global().Set("clearBuffer", js.FuncOf(clearBuffer))
	// Keep WebAssembly running
	select {}
}
