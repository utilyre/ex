package config

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Mode int

const (
	ModeDev = iota
	ModeProd
)

type Config struct {
	Mode     Mode
	Root     string
	LogLevel slog.Level

	ServerAddr string
}

func Load() Config {
	var mode string
	flag.StringVar(&mode, "mode", "dev", "determine application mode (dev|prod)")
	flag.Parse()

	cfg := Config{}

	switch mode {
	case "dev":
		cfg.Mode = ModeDev
		if err := godotenv.Load(".env.local", ".env"); err != nil && !errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "godotenv: %v\n", err)
			os.Exit(1)
		}
	case "prod":
		cfg.Mode = ModeProd
	default:
		fmt.Fprintf(os.Stderr, "invalid argument '%s' for '-mode'\n", mode)
		os.Exit(1)
	}

	if root, ok := os.LookupEnv("APP_ROOT"); ok {
		cfg.Root = root
	} else {
		cfg.Root = "."
	}

	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		cfg.LogLevel = slog.LevelDebug
	case "info":
		cfg.LogLevel = slog.LevelInfo
	case "warn":
		cfg.LogLevel = slog.LevelWarn
	case "error":
		cfg.LogLevel = slog.LevelError
	}

	if addr, ok := os.LookupEnv("SERVER_ADDR"); ok {
		cfg.ServerAddr = addr
	} else {
		cfg.ServerAddr = "127.0.0.1:3000"
	}

	return cfg
}
