package main

import (
	"errors"
	"fmt"
	"fp/pkg/fp"
	"fp/pkg/repl"
	"github.com/chzyer/readline"
	"io"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	repl, welcome := repl.NewFP(fp.NewStdRuntime())
	_, _ = fmt.Fprintf(os.Stderr, welcome)

	// Create a readline instance with a static prompt
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          ">>> ",                 // Default prompt
		HistoryFile:     "/tmp/fp_repl_history", // Save command history
		InterruptPrompt: "^C",
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	// Channel to signal interrupts (Ctrl+C)
	interruptCh := make(chan struct{}, 1)

	// Channel for OS signals (SIGINT, SIGTERM)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	// Goroutine to listen for interrupts and notify REPL
	go func() {
		for sig := range signalCh {
			if sig == os.Interrupt {
				select {
				case interruptCh <- struct{}{}:
				default:
				}
			} else {
				os.Exit(0) // Exit cleanly on SIGTERM
			}
		}
	}()

	for {
		// Read input
		line, err := rl.Readline()
		if err != nil {
			if errors.Is(err, readline.ErrInterrupt) { // Handle Ctrl+C
				output := repl.ClearBuffer()
				_, _ = fmt.Fprint(os.Stderr, "    "+output)
				continue
			} else if err == io.EOF { // Handle Ctrl+D (exit)
				os.Exit(0)
			}
			panic(err)
		}

		// Process input in REPL
		output, executed := repl.ReplyInput(line, interruptCh)

		// Print REPL output
		if output != "" {
			_, _ = fmt.Fprint(os.Stderr, "    "+output)
		}

		// If executed is true, print prompt again
		if executed {
			// Reset the prompt to ">>> " when input is executed
			rl.SetPrompt(">>> ")
		} else {
			// Otherwise, indent continuation line (you can choose what to show)
			rl.SetPrompt("    ") // Or set it to "" for no prompt if not executed
		}
	}
}
