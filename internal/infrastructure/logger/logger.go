package logger

import (
	"context"
	"fmt"
	"hona/backend/bootstrap"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

var (
	initOnce      sync.Once
	defaultLogger *slog.Logger
	logFile       *os.File
)

func Init() *slog.Logger {
	initOnce.Do(func() {
		cfg := bootstrap.Run().Env.Logger
		handler := buildHandler(cfg)

		logger := slog.New(handler)
		logger = withStaticAttrs(logger, cfg)

		slog.SetDefault(logger)
		defaultLogger = logger
		redirectStdLog(logger)
	})

	if defaultLogger == nil {
		return slog.Default()
	}
	return defaultLogger
}

func Close() error {
	if logFile == nil {
		return nil
	}
	return logFile.Close()
}

func withStaticAttrs(logger *slog.Logger, cfg bootstrap.LoggerConfig) *slog.Logger {
	attrs := make([]any, 0, 4)
	if cfg.ServiceName != "" {
		attrs = append(attrs, "service", cfg.ServiceName)
	}
	if cfg.Environment != "" {
		attrs = append(attrs, "env", cfg.Environment)
	}
	if len(attrs) == 0 {
		return logger
	}
	return logger.With(attrs...)
}

func redirectStdLog(logger *slog.Logger) {
	log.SetFlags(0)
	log.SetPrefix("")
	log.SetOutput(slog.NewLogLogger(logger.Handler(), slog.LevelInfo).Writer())
}

func buildHandler(cfg bootstrap.LoggerConfig) slog.Handler {
	handlers := make([]slog.Handler, 0, 3)
	opts := &slog.HandlerOptions{
		Level:     cfg.Level,
		AddSource: cfg.AddSource,
	}

	if cfg.TextStdout {
		handlers = append(handlers, NewPrettyTextHandler(os.Stdout, opts))
	}
	if cfg.JSONStdout {
		handlers = append(handlers, slog.NewJSONHandler(os.Stdout, opts))
	}
	if cfg.JSONFilePath != "" {
		file, err := openLogFile(cfg.JSONFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "logger: unable to open log file %s: %v\n", cfg.JSONFilePath, err)
		} else {
			logFile = file
			handlers = append(handlers, slog.NewJSONHandler(file, opts))
		}
	}

	if len(handlers) == 0 {
		handlers = append(handlers, slog.NewTextHandler(os.Stdout, opts))
	}

	if len(handlers) == 1 {
		return handlers[0]
	}
	return &multiHandler{handlers: handlers}
}

func openLogFile(path string) (*os.File, error) {
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}
	return os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
}

type multiHandler struct {
	handlers []slog.Handler
}

func (h *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *multiHandler) Handle(ctx context.Context, record slog.Record) error {
	var firstErr error
	handled := false
	for _, handler := range h.handlers {
		if !handler.Enabled(ctx, record.Level) {
			continue
		}
		next := record
		if handled {
			next = record.Clone()
		}
		handled = true
		if err := handler.Handle(ctx, next); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (h *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	next := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		next[i] = handler.WithAttrs(attrs)
	}
	return &multiHandler{handlers: next}
}

func (h *multiHandler) WithGroup(name string) slog.Handler {
	next := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		next[i] = handler.WithGroup(name)
	}
	return &multiHandler{handlers: next}
}
