// Package main provides the entry point for Ory Talos.
// See talos/AGENTS.md for development guidelines.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ory/x/cmdx"

	"github.com/ory/talos/cmd"
)

func main() {
	// Run in a separate function so deferred cleanup (such as stop()) runs
	// before os.Exit terminates the process.
	os.Exit(run())
}

func run() int {
	// Cancel the command context on SIGINT/SIGTERM so blocking commands
	// (such as 'migrate status --block') unwind gracefully on pod deletion.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Create and execute root command
	// Version info comes from internal/version package (set by build flags)
	rootCmd := cmd.NewRoot()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		// ErrNoPrintButFail means the command already communicated the failure
		// to the user via stderr; suppress the error message and just exit non-zero.
		if !errors.Is(err, cmdx.ErrNoPrintButFail) {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		return 1
	}

	return 0
}
