package main

import (
	"os"

	"github.com/wasilibs/go-ryl/internal/runner"
)

func main() {
	os.Exit(runner.Run("ryl", os.Args[1:], os.Stdin, os.Stdout, os.Stderr, "."))
}
