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

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		// receive SIGINT, SIGTERM when running code -> stop code
		for range signalCh {
			cancel()
		}
	}()

	for {
		line, err := rl.Readline()
		if err != nil {
			if errors.Is(err, readline.ErrInterrupt) {
				// receive SIGINT when typing -> clear buffer
				func() {
					replMtx.Lock()
					defer replMtx.Unlock()
					output := repl.ClearBuffer()
					if output != "" {
						_, _ = fmt.Fprint(os.Stderr, "    "+output)
					}
				}()
				continue
			} else if err == io.EOF {
				// receive SIGTERM when typing -> exit
				os.Exit(0)
			}
			panic(err)
		}
		ctx, cancel = context.WithCancel(context.Background())
		func() {
			defer cancel()
			replMtx.Lock()
			defer replMtx.Unlock()
			output, executed := repl.ReplyInput(ctx, line)
			if output != "" {
				_, _ = fmt.Fprint(os.Stderr, "    "+output)
			}
			contextCancelled := false
			select {
			case <-ctx.Done():
				contextCancelled = true
			default:
			}
			if executed || contextCancelled {
				rl.SetPrompt(">>> ") // reset prompt if command is executed
			} else {
				rl.SetPrompt("    ") // otherwise
			}
		}()

	}
}
