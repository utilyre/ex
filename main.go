package main

import (
	"flag"

	"github.com/joho/godotenv"
	"github.com/utilyre/ex/application"
	"github.com/utilyre/ex/config"
)

var (
	env string
)

func init() {
	flag.StringVar(&env, "env", "", "path to an optional .env file")
	flag.Parse()
}

func main() {
	if env != "" {
		if err := godotenv.Load(env); err != nil {
			panic(err)
		}
	}

	cfg := config.Config{}
	cfg.Load()

	app := application.New(cfg)
	app.Setup()
	app.Start()
}
