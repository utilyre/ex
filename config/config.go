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
	Mode     Mode
	LogLevel slog.Level
	Root     string

	ServerAddr string
}

func (c *Config) Load() {
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

	if root, ok := os.LookupEnv("EX_ROOT"); ok {
		c.Root = root
	} else {
		c.Root = "."
	}

	if addr, ok := os.LookupEnv("EX_SERVER_ADDR"); ok {
		c.ServerAddr = addr
	} else {
		c.ServerAddr = "127.0.0.1:3000"
	}
}
