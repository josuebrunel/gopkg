// Package xlog provides a structured logger based on log/slog.
package xlog

import (
	"io"
	"log/slog"
	"os"
	"runtime"
	"strconv"
)

var logger *slog.Logger

// Format constants
const (
	FormatJSON = "json"
	FormatText = "text"
)

// Config represents the configuration for the logger.
type Config struct {
	Level         string    // Level is the log level: "DEBUG", "INFO", "WARN", "ERROR". Default is "INFO".
	Output        io.Writer // Output destination. Default is os.Stdout.
	Format        string    // Format is the output format: FormatJSON (default) or FormatText.
	DisableSource bool      // DisableSource disables source file/line logging. Default is false (source enabled).
	SourceDepth   int       // SourceDepth is the stack depth for SourceKey. Default is 7.
	Color         bool      // Color enables color-coded, human-readable console output, overriding Format.
}

func init() {
	// Initialize with defaults
	Setup(Config{
		Level:         "INFO",
		Output:        os.Stdout,
		Format:        FormatJSON, // Default format
		DisableSource: false,      // Source enabled by default
		SourceDepth:   7,
	})
}

// Setup configures the package-level logger.
func Setup(cfg Config) {
	var level slog.Level

	// Parse Level
	switch cfg.Level {
	case "DEBUG", "debug":
		level = slog.LevelDebug
	case "INFO", "info":
		level = slog.LevelInfo
	case "WARN", "warn":
		level = slog.LevelWarn
	case "ERROR", "error":
		level = slog.LevelError
	default:
		// Attempt to parse integer level
		if v, err := strconv.Atoi(cfg.Level); err == nil {
			level = slog.Level(v)
		} else {
			level = slog.LevelInfo // Default fallback
		}
	}

	// Default Output
	if cfg.Output == nil {
		cfg.Output = os.Stdout
	}

	// Default SourceDepth
	if cfg.SourceDepth == 0 {
		cfg.SourceDepth = 7
	}

	opts := &slog.HandlerOptions{
		AddSource: !cfg.DisableSource,
		Level:     level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if !cfg.DisableSource && a.Key == slog.SourceKey {
				pc, f, l, _ := runtime.Caller(cfg.SourceDepth)
				a.Value = slog.GroupValue(
					slog.Attr{
						Key:   "file",
						Value: slog.StringValue(f),
					},
					slog.Attr{
						Key:   "line",
						Value: slog.IntValue(l),
					},
					slog.Attr{
						Key:   "function",
						Value: slog.StringValue(runtime.FuncForPC(pc).Name()),
					},
				)
			}
			return a
		},
	}

	var handler slog.Handler
	switch {
	case cfg.Color:
		handler = NewColorHandler(cfg.Output, opts)
	case cfg.Format == FormatText:
		handler = slog.NewTextHandler(cfg.Output, opts)
	default:
		handler = slog.NewJSONHandler(cfg.Output, opts)
	}

	logger = slog.New(handler)
}

// Info logs a message at Info level.
func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

// Debug logs a message at Debug level.
func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

// Error logs a message at Error level.
func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

// Warn logs a message at Warn level.
func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}
