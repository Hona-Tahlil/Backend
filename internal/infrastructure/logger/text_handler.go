package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const timeFormat = "2006-01-02 15:04:05.000"

const (
	colorReset  = "\x1b[0m"
	colorGray   = "\x1b[90m"
	colorRed    = "\x1b[31m"
	colorGreen  = "\x1b[32m"
	colorYellow = "\x1b[33m"
	colorCyan   = "\x1b[36m"
	colorBlue   = "\x1b[34m"
)

type PrettyTextHandler struct {
	out    io.Writer
	opts   slog.HandlerOptions
	attrs  []slog.Attr
	groups []string
	color  bool
	mu     sync.Mutex
}

func NewPrettyTextHandler(out io.Writer, opts *slog.HandlerOptions) *PrettyTextHandler {
	handlerOpts := slog.HandlerOptions{}
	if opts != nil {
		handlerOpts = *opts
	}
	return &PrettyTextHandler{
		out:   out,
		opts: handlerOpts,
		color: isTerminalWriter(out) &&
			os.Getenv("TERM") != "dumb",
	}
}

func (h *PrettyTextHandler) Enabled(_ context.Context, level slog.Level) bool {
	if h.opts.Level == nil {
		return level >= slog.LevelInfo
	}
	return level >= h.opts.Level.Level()
}

func (h *PrettyTextHandler) Handle(ctx context.Context, record slog.Record) error {
	if !h.Enabled(ctx, record.Level) {
		return nil
	}

	timestamp := record.Time
	if timestamp.IsZero() {
		timestamp = time.Now()
	}
	timeStr := timestamp.Format(timeFormat)

	levelLabel := strings.ToUpper(record.Level.String())
	levelPadded := fmt.Sprintf("%-5s", levelLabel)
	if h.color {
		timeStr = colorize(colorGray, timeStr)
		levelPadded = colorize(levelColor(record.Level), levelPadded)
	}

	lines := make([]string, 0, 6)
	if record.Message != "" {
		lines = append(lines, fmt.Sprintf("%s %s %s", timeStr, levelPadded, record.Message))
	} else {
		lines = append(lines, fmt.Sprintf("%s %s", timeStr, levelPadded))
	}

	attrs := make([]slog.Attr, 0, len(h.attrs)+record.NumAttrs()+1)
	attrs = append(attrs, h.attrs...)
	record.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, attr)
		return true
	})

	if h.opts.AddSource && record.PC != 0 {
		frame, _ := runtime.CallersFrames([]uintptr{record.PC}).Next()
		source := fmt.Sprintf("%s:%d", filepath.Base(frame.File), frame.Line)
		attrs = append(attrs, slog.String("source", source))
	}

	prefix := strings.Join(h.groups, ".")
	for _, attr := range attrs {
		appendAttrLines(&lines, prefix, attr, h.color)
	}

	output := strings.Join(lines, "\n") + "\n"
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := io.WriteString(h.out, output)
	return err
}

func (h *PrettyTextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	next := h.clone()
	next.attrs = append(append([]slog.Attr{}, h.attrs...), attrs...)
	return next
}

func (h *PrettyTextHandler) WithGroup(name string) slog.Handler {
	next := h.clone()
	if name != "" {
		next.groups = append(append([]string{}, h.groups...), name)
	} else {
		next.groups = append([]string{}, h.groups...)
	}
	return next
}

func (h *PrettyTextHandler) clone() *PrettyTextHandler {
	return &PrettyTextHandler{
		out:    h.out,
		opts:   h.opts,
		attrs:  append([]slog.Attr{}, h.attrs...),
		groups: append([]string{}, h.groups...),
		color:  h.color,
	}
}

func appendAttrLines(lines *[]string, prefix string, attr slog.Attr, color bool) {
	attr.Value = attr.Value.Resolve()
	if attr.Value.Kind() == slog.KindGroup {
		groupPrefix := attr.Key
		if prefix != "" && groupPrefix != "" {
			groupPrefix = prefix + "." + groupPrefix
		} else if prefix != "" {
			groupPrefix = prefix
		}
		for _, groupAttr := range attr.Value.Group() {
			appendAttrLines(lines, groupPrefix, groupAttr, color)
		}
		return
	}

	key := attr.Key
	if prefix != "" {
		if key != "" {
			key = prefix + "." + key
		} else {
			key = prefix
		}
	}
	if key == "" {
		key = "attr"
	}

	value := formatValue(attr.Value)
	if color {
		key = colorize(colorBlue, key)
	}
	*lines = append(*lines, fmt.Sprintf("  %s: %s", key, value))
}

func formatValue(value slog.Value) string {
	value = value.Resolve()
	switch value.Kind() {
	case slog.KindString:
		str := value.String()
		if strings.ContainsAny(str, " \t") {
			return strconv.Quote(str)
		}
		return str
	case slog.KindInt64:
		return strconv.FormatInt(value.Int64(), 10)
	case slog.KindUint64:
		return strconv.FormatUint(value.Uint64(), 10)
	case slog.KindFloat64:
		return strconv.FormatFloat(value.Float64(), 'f', -1, 64)
	case slog.KindBool:
		return strconv.FormatBool(value.Bool())
	case slog.KindDuration:
		return value.Duration().String()
	case slog.KindTime:
		return value.Time().Format(timeFormat)
	default:
		return fmt.Sprint(value.Any())
	}
}

func levelColor(level slog.Level) string {
	switch {
	case level >= slog.LevelError:
		return colorRed
	case level >= slog.LevelWarn:
		return colorYellow
	case level >= slog.LevelInfo:
		return colorGreen
	default:
		return colorCyan
	}
}

func colorize(color, text string) string {
	return color + text + colorReset
}

func isTerminalWriter(out io.Writer) bool {
	file, ok := out.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}
