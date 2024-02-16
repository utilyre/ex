package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

var ErrInvalidMode = errors.New("invalid mode")

type Mode int

const (
	ModeDev = iota
	ModeProd
)

func (m Mode) MarshalText() ([]byte, error) {
	switch m {
	case ModeDev:
		return []byte("DEV"), nil
	case ModeProd:
		return []byte("PROD"), nil
	default:
		return nil, ErrInvalidMode
	}
}

func (m *Mode) UnmarshalText(text []byte) error {
	mode := string(text)
	switch strings.ToUpper(mode) {
	case "DEV":
		*m = ModeDev
	case "PROD":
		*m = ModeProd
	default:
		return fmt.Errorf("string \"%s\": %w", mode, ErrInvalidMode)
	}

	return nil
}

type Config struct {
	Mode     Mode       `env:"MODE,required"`
	LogLevel slog.Level `env:"LOG_LEVEL,required"`
	AppRoot  string     `env:"APP_ROOT,required"`

	ServerAddr string `env:"SERVER_ADDR,required"`
	DSN        string `env:"DSN,required"`
}

func Load() (Config, error) {
	cfg := Config{}

	if err := godotenv.Load(".env.local"); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("config: %w", err)
	}
	if err := godotenv.Load(".env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("config: %w", err)
	}

	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("config: %w", err)
	}

	return cfg, nil
}
