package config

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	ErrUnknownName = errors.New("unknown name")
	ErrNotPresent  = errors.New("not present")
)

type Mode int

const (
	ModeDev = iota
	ModeProd
)

func (m Mode) String() string {
	switch m {
	case ModeDev:
		return "DEV"
	case ModeProd:
		return "PROD"
	default:
		return ""
	}
}

func (m Mode) MarshalText() ([]byte, error) {
	return []byte(m.String()), nil
}

func (m *Mode) UnmarshalText(text []byte) error {
	mode := string(text)
	switch strings.ToUpper(mode) {
	case "DEV":
		*m = ModeDev
	case "PROD":
		*m = ModeProd
	default:
		return fmt.Errorf("mode string \"%s\": %w", mode, ErrUnknownName)
	}

	return nil
}

type Config struct {
	Mode     Mode
	LogLevel slog.Level
	AppRoot  string

	ServerAddr string
	DSN        string
}

func Load() (Config, error) {
	cfg := Config{}

	var mode string
	flag.StringVar(&mode, "mode", "DEV", "determine application mode (DEV|PROD)")
	flag.Parse()

	if err := cfg.Mode.UnmarshalText([]byte(mode)); err != nil {
		return Config{}, fmt.Errorf("config: %w", err)
	}

	switch cfg.Mode {
	case ModeDev:
		if err := godotenv.Load(".env.local"); err != nil && !errors.Is(err, os.ErrNotExist) {
			return Config{}, fmt.Errorf("config: %w", err)
		}

		if err := godotenv.Load(".env"); err != nil {
			return Config{}, fmt.Errorf("config: %w", err)
		}
	case ModeProd:
		if err := validateEnv("LOG_LEVEL", "APP_ROOT", "SERVER_ADDR", "DSN"); err != nil {
			return Config{}, fmt.Errorf("config: %w", err)
		}
	}

	if err := cfg.LogLevel.UnmarshalText([]byte(os.Getenv("LOG_LEVEL"))); err != nil {
		return Config{}, err
	}

	cfg.AppRoot = os.Getenv("APP_ROOT")
	cfg.ServerAddr = os.Getenv("SERVER_ADDR")
	cfg.DSN = os.Getenv("DSN")

	return cfg, nil
}

func validateEnv(keys ...string) error {
	for _, key := range keys {
		_, ok := os.LookupEnv(key)
		if !ok {
			return fmt.Errorf("environment variable \"%s\": %w", key, ErrNotPresent)
		}
	}

	return nil
}
