package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

type LoggingMiddleware struct {
	logger *slog.Logger
}

func NewLoggingMiddleware() *LoggingMiddleware {
	return &LoggingMiddleware{logger: slog.Default()}
}

func (lm *LoggingMiddleware) LogRequests(ctx *gin.Context) {
	start := time.Now()
	path := ctx.Request.URL.Path
	rawQuery := ctx.Request.URL.RawQuery
	route := ctx.FullPath()
	if route == "" {
		route = path
	}

	ctx.Next()

	status := ctx.Writer.Status()
	size := ctx.Writer.Size()
	if size < 0 {
		size = 0
	}

	attrs := []slog.Attr{
		slog.String("method", ctx.Request.Method),
		slog.String("path", path),
		slog.String("route", route),
		slog.Int("status", status),
		slog.Int("size", size),
		slog.String("ip", ctx.ClientIP()),
		slog.String("user_agent", ctx.Request.UserAgent()),
		slog.Duration("duration", time.Since(start)),
	}

	if rawQuery != "" {
		attrs = append(attrs, slog.String("query", rawQuery))
	}
	if requestID := ctx.GetHeader("X-Request-ID"); requestID != "" {
		attrs = append(attrs, slog.String("request_id", requestID))
	}
	if len(ctx.Errors) > 0 {
		attrs = append(attrs, slog.String("errors", ctx.Errors.String()))
	}

	level := slog.LevelInfo
	if status >= 500 {
		level = slog.LevelError
	} else if status >= 400 {
		level = slog.LevelWarn
	}

	lm.logger.LogAttrs(ctx.Request.Context(), level, "http_request", attrs...)
}
