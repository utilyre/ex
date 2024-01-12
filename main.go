package main

import (
	"github.com/utilyre/golang-backend-template/application"
	"github.com/utilyre/golang-backend-template/config"
)

func main() {
	cfg := config.Config{}
	cfg.Load()

	app := application.New(cfg)
	app.Init()
	app.Start()
}
