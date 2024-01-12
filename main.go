package main

import (
	"github.com/utilyre/ex/application"
	"github.com/utilyre/ex/config"
)

func main() {
	cfg := config.Config{}
	cfg.Load()

	app := application.New(cfg)
	app.Init()
	app.Start()
}
