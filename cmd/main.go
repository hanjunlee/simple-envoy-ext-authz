package main

import (
	"os"

	sub "github.com/hanjunlee/simple-envoy-ext-authz/cmd/subcmds"
)

func main() {
	if err := sub.Execute(); err != nil {
		os.Exit(2)
	}
}
