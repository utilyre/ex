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
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	application.New(cfg).Setup().Start()
}
