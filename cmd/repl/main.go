package main

import (
	"context"
	"errors"
	"fmt"
	"fp/pkg/fp"
	"fp/pkg/repl"
	"github.com/chzyer/readline"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	replMtx := &sync.Mutex{}
	repl, welcome := repl.NewFP(fp.NewStdRuntime())
	_, _ = fmt.Fprintf(os.Stderr, welcome)

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          ">>> ",                 // Default prompt
		HistoryFile:     "/tmp/fp_repl_history", // Save command history
		InterruptPrompt: "^C",
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	var ctx context.Context
	var cancel context.CancelFunc = func() {}

	// handle syscall.SIGINT, syscall.SIGTERM when running code
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range signalCh {
			cancel()
			switch sig {
			case syscall.SIGINT:
				replMtx.Lock()
				output := repl.ClearBuffer()
				replMtx.Unlock()
				_, _ = fmt.Fprint(os.Stderr, "    "+output)
			case syscall.SIGTERM:
				os.Exit(0)
			}
		}
	}()

	for {
		line, err := rl.Readline()
		if err != nil {
			if errors.Is(err, readline.ErrInterrupt) { // handle syscall.SIGINT when receiving input
				replMtx.Lock()
				output := repl.ClearBuffer()
				replMtx.Unlock()
				_, _ = fmt.Fprint(os.Stderr, "    "+output)
				continue
			} else if err == io.EOF { // handle syscall.SIGTERM when receiving input
				os.Exit(0)
			}
			panic(err)
		}
		ctx, cancel = context.WithCancel(context.Background())

		replMtx.Lock()
		output, executed := repl.ReplyInput(ctx, line)
		replMtx.Unlock()

		if output != "" {
			_, _ = fmt.Fprint(os.Stderr, "    "+output)
		}

		if executed {
			rl.SetPrompt(">>> ") // reset prompt if command is executed
		} else {
			rl.SetPrompt("    ") // otherwise
		}
	}
}
