package main

import (
	"github.com/utilyre/ex/application"
	"github.com/utilyre/ex/config"
)

func main() {
	cfg := config.Load()
	application.New(cfg).Setup().Start()
}
