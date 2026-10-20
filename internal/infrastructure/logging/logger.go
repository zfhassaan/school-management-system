package logging

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
	"school-management/internal/platform/config"
)

// SetupLogger configures slog with both console and rotating file output.
func SetupLogger(cfg config.LoggingConfig) *slog.Logger {
	var level slog.Level
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var writers []io.Writer
	if cfg.Console {
		writers = append(writers, os.Stdout)
	}

	if cfg.FilePath != "" {
		rotator := &lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSizeMB,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAgeDays,
			Compress:   true,
		}
		writers = append(writers, rotator)
	}

	mw := io.MultiWriter(writers...)

	var handler slog.Handler
	if strings.ToLower(cfg.Format) == "json" {
		handler = slog.NewJSONHandler(mw, opts)
	} else {
		handler = slog.NewTextHandler(mw, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

