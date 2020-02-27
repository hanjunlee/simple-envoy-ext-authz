package main

import (
	"os"

	sub "github.com/hanjunlee/simple-envoy-ext-authz/cmd/subcommands"
)

func main() {
	if err := sub.Execute(); err != nil {
		os.Exit(2)
	}
}
