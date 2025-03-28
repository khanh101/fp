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

	// Create a readline instance
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          ">>>",                  // We manually print the prompt
		HistoryFile:     "/tmp/fp_repl_history", // Save command history
		InterruptPrompt: "^C",
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	interruptCh := make(chan struct{}, 1)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

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
		line, err := rl.Readline()
		if err != nil {
			if errors.Is(err, readline.ErrInterrupt) { // Handle Ctrl+C
				output := repl.ClearBuffer()
				_, _ = fmt.Fprint(os.Stderr, output)
				continue
			} else if err == io.EOF { // Handle Ctrl+D (exit)
				os.Exit(0)
			}
			panic(err)
		}

		output, executed := repl.ReplyInput(line, interruptCh)

		if output != "" {
			_, _ = fmt.Fprint(os.Stderr, output)
		}
		_ = executed

	}
}
