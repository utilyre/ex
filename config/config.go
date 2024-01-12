package config

import "log/slog"

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
