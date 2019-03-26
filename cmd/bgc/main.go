package main

import (
	"os"
	"sidus.io/boogrocha/internal/cli"
)

func main() {
	if err := cli.BgcCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
