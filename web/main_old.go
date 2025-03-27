package main

import (
	"fmt"
	"syscall/js"
)

func evaluate(this js.Value, p []js.Value) interface{} {
	if len(p) == 0 {
		return js.ValueOf("No input")
	}
	input := p[0].String()
	// Simple echo for now; replace with real evaluation logic
	output := fmt.Sprintf("You have <br> entered: %s", input)
	return js.ValueOf(output)
}

func main() {
	// Expose Go functions to the global JavaScript context
	js.Global().Set("evaluate", js.FuncOf(evaluate))
	// Keep WebAssembly running
	select {}
}
