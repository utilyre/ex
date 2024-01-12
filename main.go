package main

import (
	"log/slog"

	"github.com/utilyre/golang-backend-template/application"
	"github.com/utilyre/golang-backend-template/config"
)

func main() {
	cfg := config.Config{
		Root:     ".",
		Mode:     config.ModeProd,
		LogLevel: slog.LevelDebug,
	}

	// TODO: get cfg from env
	// TODO: then override env by flag

	app := application.New(cfg)
	app.Init()
	app.Start()
}
