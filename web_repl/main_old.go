package main

import (
	"fmt"
	"syscall/js"
	"time"
)

func evaluate(this js.Value, p []js.Value) interface{} {
	if len(p) == 0 {
		return js.ValueOf("No input")
	}
	input := p[0].String()
	// Simple echo for now; replace with real evaluation logic
	output := fmt.Sprintf("You entered: %s", input)
	js.Global().Call("updateOutput", output)
	return nil
}

func sendOutputToWeb(this js.Value, p []js.Value) interface{} {
	// Conditionally send data to the web_repl, for example, without user input
	output := "This is data sent from Go without user input."
	js.Global().Call("updateOutput", output)
	return nil
}

func main() {
	// Expose Go functions to the global JavaScript context
	js.Global().Set("evaluate", js.FuncOf(evaluate))
	js.Global().Set("sendOutputToWeb", js.FuncOf(sendOutputToWeb))

	// Call sendOutputToWeb conditionally (for example, after a certain delay or event)
	go func() {
		// Example: Automatically send output after 5 seconds
		select {
		case <-time.After(5 * time.Second):
			js.Global().Call("sendOutputToWeb")
		}
	}()

	// Keep WebAssembly running
	select {}
}
