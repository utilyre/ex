package config

import (
	"log/slog"
	"os"
)

type Mode int

const (
	ModeDev = iota
	ModeProd
)

type Config struct {
	Root     string
	Mode     Mode
	LogLevel slog.Level
}

func (c *Config) Load() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if root, ok := os.LookupEnv("EX_ROOT"); ok {
		c.Root = root
	} else {
		c.Root = cwd
	}

	switch os.Getenv("EX_MODE") {
	case "dev":
		c.Mode = ModeDev
	case "prod":
		c.Mode = ModeProd
	}

	switch os.Getenv("EX_LOG_LEVEL") {
	case "debug":
		c.LogLevel = slog.LevelDebug
	case "info":
		c.LogLevel = slog.LevelInfo
	case "warn":
		c.LogLevel = slog.LevelWarn
	case "error":
		c.LogLevel = slog.LevelError
	}
}