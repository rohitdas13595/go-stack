package main

import (
	"os"

	"github.com/rohitdas13595/go-stack/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
