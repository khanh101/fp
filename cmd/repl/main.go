package main

import (
	"bufio"
	"fmt"
	"fp/pkg/repl"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	repl, welcome := repl.NewFP()
	_, _ = fmt.Fprintf(os.Stderr, welcome)

	signCh := make(chan os.Signal, 1)
	signal.Notify(signCh, syscall.SIGINT, syscall.SIGTERM)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		select {
		case s := <-signCh:
			switch s {
			case syscall.SIGINT:
				output := repl.ClearBuffer()
				_, _ = fmt.Fprintf(os.Stderr, output)
			case syscall.SIGTERM:
				os.Exit(0)
			default:
				os.Exit(1)
			}
		default:
			input := scanner.Text()
			output := repl.ReplyInput(input)
			_, _ = fmt.Fprintf(os.Stderr, output)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
