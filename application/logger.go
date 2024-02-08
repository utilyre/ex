package application

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/utilyre/ex/config"
)

func NewLogger(cfg config.Config) (*slog.Logger, error) {
	opts := &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}

	var handler slog.Handler
	switch cfg.Mode {
	case config.ModeDev:
		handler = slog.NewTextHandler(os.Stdout, opts)
	case config.ModeProd:
		f, err := os.OpenFile(filepath.Join(cfg.AppRoot, "server.log"),
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			return nil, fmt.Errorf("application: %w", err)
		}

		handler = slog.NewJSONHandler(io.MultiWriter(os.Stdout, f), opts)
	}

	return slog.New(handler), nil
}
