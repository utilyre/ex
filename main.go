package main

import (
	"fmt"
	"os"

	"github.com/utilyre/ex/application"
	"github.com/utilyre/ex/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ex: %v\n", err)
		os.Exit(1)
	}

	logger, err := application.NewLogger(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ex: %v\n", err)
		os.Exit(1)
	}

	application.
		New(cfg, logger).
		Setup().
		Start()
}
